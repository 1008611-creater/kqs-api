package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
)

// GetAccountRoutingStats returns all accounts in the selected scope and joins
// real success/error logs for the requested window. A missing sample remains a
// row with zero request_count so the admin can distinguish idle from missing data.
func (r *opsRepository) GetAccountRoutingStats(ctx context.Context, startTime, endTime time.Time, platform string, groupID *int64) ([]*service.OpsAccountRoutingStats, error) {
	if r == nil || r.db == nil {
		return nil, fmt.Errorf("nil ops repository")
	}
	const query = `
WITH request_samples AS (
  SELECT account_id, COUNT(*)::bigint AS request_count,
         COUNT(*)::bigint AS success_count, 0::bigint AS error_count,
         COALESCE(SUM(duration_ms), 0)::float8 AS latency_sum,
         COUNT(duration_ms)::float8 AS latency_count,
         COALESCE(SUM(first_token_ms), 0)::float8 AS ttft_sum,
         COUNT(first_token_ms)::float8 AS ttft_count
  FROM usage_logs
  WHERE account_id IS NOT NULL
    AND created_at >= $1 AND created_at < $2
    AND ($4::bigint IS NULL OR group_id = $4)
  GROUP BY account_id
  UNION ALL
  SELECT account_id, COUNT(*)::bigint AS request_count,
         0::bigint AS success_count, COUNT(*)::bigint AS error_count,
         COALESCE(SUM(response_latency_ms), 0)::float8 AS latency_sum,
         COUNT(response_latency_ms)::float8 AS latency_count,
         COALESCE(SUM(time_to_first_token_ms), 0)::float8 AS ttft_sum,
         COUNT(time_to_first_token_ms)::float8 AS ttft_count
  FROM ops_error_logs
  WHERE account_id IS NOT NULL
    AND created_at >= $1 AND created_at < $2
    AND ($3 = '' OR COALESCE(platform, '') = $3)
    AND ($4::bigint IS NULL OR group_id = $4)
  GROUP BY account_id
), aggregated AS (
  SELECT account_id, SUM(request_count)::bigint AS request_count,
         SUM(success_count)::bigint AS success_count,
         SUM(error_count)::bigint AS error_count,
         SUM(latency_sum)::float8 AS latency_sum,
         SUM(latency_count)::float8 AS latency_count,
         SUM(ttft_sum)::float8 AS ttft_sum,
         SUM(ttft_count)::float8 AS ttft_count
  FROM request_samples
  GROUP BY account_id
)
SELECT a.id, a.name, a.platform, a.priority, ag.priority,
       ag.group_id, COALESCE(g.name, ''),
       COALESCE(a.load_factor, NULLIF(a.concurrency, 0), 1),
       a.schedulable, a.status,
       COALESCE(aggregated.request_count, 0),
       COALESCE(aggregated.success_count, 0),
       COALESCE(aggregated.error_count, 0),
       CASE WHEN aggregated.latency_count > 0 THEN aggregated.latency_sum / aggregated.latency_count END,
       CASE WHEN aggregated.ttft_count > 0 THEN aggregated.ttft_sum / aggregated.ttft_count END
FROM accounts a
LEFT JOIN account_groups ag
  ON ag.account_id = a.id
 AND $4::bigint IS NOT NULL
 AND ag.group_id = $4
LEFT JOIN groups g ON g.id = ag.group_id
LEFT JOIN aggregated ON aggregated.account_id = a.id
WHERE ($3 = '' OR a.platform = $3)
  AND ($4::bigint IS NULL OR ag.group_id IS NOT NULL)
ORDER BY COALESCE(ag.priority, a.priority), a.id`

	var groupArg any
	if groupID != nil {
		groupArg = *groupID
	}
	rows, err := r.db.QueryContext(ctx, query, startTime.UTC(), endTime.UTC(), platform, groupArg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*service.OpsAccountRoutingStats, 0)
	for rows.Next() {
		var item service.OpsAccountRoutingStats
		var groupIDValue sql.NullInt64
		var groupPriority sql.NullInt64
		var latency, ttft sql.NullFloat64
		if err := rows.Scan(
			&item.AccountID, &item.AccountName, &item.Platform, &item.AccountPriority,
			&groupPriority, &groupIDValue, &item.GroupName, &item.LoadFactor,
			&item.Schedulable, &item.Status, &item.RequestCount, &item.SuccessCount,
			&item.ErrorCount, &latency, &ttft,
		); err != nil {
			return nil, err
		}
		item.Priority = item.AccountPriority
		if groupPriority.Valid {
			value := int(groupPriority.Int64)
			item.GroupPriority = &value
			item.Priority = value
		}
		if groupIDValue.Valid {
			value := groupIDValue.Int64
			item.GroupID = &value
		}
		if latency.Valid {
			value := int(latency.Float64 + 0.5)
			item.AvgLatencyMs = &value
		}
		if ttft.Valid {
			value := int(ttft.Float64 + 0.5)
			item.AvgTTFTMs = &value
		}
		items = append(items, &item)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
