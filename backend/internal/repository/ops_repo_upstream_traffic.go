package repository

import (
	"context"
	"database/sql"
	"fmt"
	"math"
	"strings"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

func (r *opsRepository) GetUpstreamTrafficStats(ctx context.Context, filter *service.OpsDashboardFilter, bucketSeconds int) (*service.OpsUpstreamTrafficStatsResponse, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil ops repository")
	}
	if filter == nil {
		return nil, fmt.Errorf("nil filter")
	}
	if filter.StartTime.IsZero() || filter.EndTime.IsZero() {
		return nil, fmt.Errorf("start_time/end_time required")
	}
	if bucketSeconds <= 0 {
		bucketSeconds = 300
	}

	start := filter.StartTime.UTC()
	end := filter.EndTime.UTC()
	join, where, args, next := buildUsageWhere(filter, start, end, 1)

	totalQuery := `
SELECT
  COALESCE(SUM(ul.upstream_request_bytes), 0),
  COALESCE(SUM(ul.upstream_response_bytes), 0),
  COALESCE(COUNT(*), 0)
FROM usage_logs ul
` + join + `
` + where

	var requestBytes, responseBytes, requestCount int64
	if err := r.db.QueryRowContext(ctx, totalQuery, args...).Scan(&requestBytes, &responseBytes, &requestCount); err != nil {
		return nil, err
	}

	trendQuery := `
SELECT
  to_timestamp(floor(extract(epoch from ul.created_at) / $` + fmt.Sprint(next) + `) * $` + fmt.Sprint(next) + `) AT TIME ZONE 'UTC' AS bucket_start,
  COALESCE(SUM(ul.upstream_request_bytes), 0),
  COALESCE(SUM(ul.upstream_response_bytes), 0),
  COALESCE(COUNT(*), 0)
FROM usage_logs ul
` + join + `
` + where + `
GROUP BY 1
ORDER BY 1 ASC`

	trendArgs := append(append([]any{}, args...), bucketSeconds)
	trendRows, err := r.db.QueryContext(ctx, trendQuery, trendArgs...)
	if err != nil {
		return nil, err
	}
	trend := make([]*service.OpsUpstreamTrafficTrendPoint, 0, 64)
	for trendRows.Next() {
		var p service.OpsUpstreamTrafficTrendPoint
		if err := trendRows.Scan(&p.BucketStart, &p.RequestBytes, &p.ResponseBytes, &p.RequestCount); err != nil {
			_ = trendRows.Close()
			return nil, err
		}
		p.BucketStart = p.BucketStart.UTC()
		p.TotalBytes = p.RequestBytes + p.ResponseBytes
		trend = append(trend, &p)
	}
	if err := trendRows.Close(); err != nil {
		return nil, err
	}
	if err := trendRows.Err(); err != nil {
		return nil, err
	}

	topJoin := join
	if !strings.Contains(topJoin, "accounts a") {
		topJoin += " LEFT JOIN accounts a ON a.id = ul.account_id"
	}
	topQuery := `
SELECT
  ul.account_id,
  COALESCE(NULLIF(a.name, ''), 'Account #' || ul.account_id::text) AS account_name,
  COALESCE(NULLIF(a.platform, ''), '') AS platform,
  COALESCE(SUM(ul.upstream_request_bytes), 0),
  COALESCE(SUM(ul.upstream_response_bytes), 0),
  COALESCE(COUNT(*), 0),
  AVG(ul.duration_ms) FILTER (WHERE ul.duration_ms IS NOT NULL),
  AVG(ul.first_token_ms) FILTER (WHERE ul.first_token_ms IS NOT NULL)
FROM usage_logs ul
` + topJoin + `
` + where + `
GROUP BY ul.account_id, a.name, a.platform
ORDER BY (COALESCE(SUM(ul.upstream_request_bytes), 0) + COALESCE(SUM(ul.upstream_response_bytes), 0)) DESC, COUNT(*) DESC
LIMIT 10`

	topRows, err := r.db.QueryContext(ctx, topQuery, args...)
	if err != nil {
		return nil, err
	}
	topAccounts := make([]*service.OpsUpstreamTrafficAccountItem, 0, 10)
	for topRows.Next() {
		var item service.OpsUpstreamTrafficAccountItem
		var avgDuration, avgTTFT sql.NullFloat64
		if err := topRows.Scan(
			&item.AccountID,
			&item.AccountName,
			&item.Platform,
			&item.RequestBytes,
			&item.ResponseBytes,
			&item.RequestCount,
			&avgDuration,
			&avgTTFT,
		); err != nil {
			_ = topRows.Close()
			return nil, err
		}
		item.TotalBytes = item.RequestBytes + item.ResponseBytes
		item.AvgDurationMs = nullableAvgMs(avgDuration)
		item.AvgTTFTMs = nullableAvgMs(avgTTFT)
		topAccounts = append(topAccounts, &item)
	}
	if err := topRows.Close(); err != nil {
		return nil, err
	}
	if err := topRows.Err(); err != nil {
		return nil, err
	}

	return &service.OpsUpstreamTrafficStatsResponse{
		StartTime:     start,
		EndTime:       end,
		Platform:      strings.TrimSpace(filter.Platform),
		GroupID:       filter.GroupID,
		Bucket:        bucketLabelFromSeconds(bucketSeconds),
		RequestBytes:  requestBytes,
		ResponseBytes: responseBytes,
		TotalBytes:    requestBytes + responseBytes,
		RequestCount:  requestCount,
		Trend:         trend,
		TopAccounts:   topAccounts,
	}, nil
}

func nullableAvgMs(v sql.NullFloat64) *int {
	if !v.Valid {
		return nil
	}
	out := int(math.Round(v.Float64))
	return &out
}

func bucketLabelFromSeconds(seconds int) string {
	switch seconds {
	case 60:
		return "1m"
	case 300:
		return "5m"
	case 3600:
		return "1h"
	case 86400:
		return "1d"
	default:
		return fmt.Sprintf("%ds", seconds)
	}
}
