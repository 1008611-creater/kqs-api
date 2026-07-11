package service

import (
	"context"
	"errors"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

func (s *OpsService) GetUpstreamTrafficStats(ctx context.Context, filter *OpsDashboardFilter, bucketSeconds int) (*OpsUpstreamTrafficStatsResponse, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
	}
	if s.opsRepo == nil {
		return nil, infraerrors.ServiceUnavailable("OPS_REPO_UNAVAILABLE", "Ops repository not available")
	}
	if filter == nil {
		return nil, infraerrors.BadRequest("OPS_FILTER_REQUIRED", "filter is required")
	}
	if filter.StartTime.IsZero() || filter.EndTime.IsZero() {
		return nil, infraerrors.BadRequest("OPS_TIME_RANGE_REQUIRED", "start_time/end_time are required")
	}
	if filter.StartTime.After(filter.EndTime) {
		return nil, infraerrors.BadRequest("OPS_TIME_RANGE_INVALID", "start_time must be <= end_time")
	}
	if bucketSeconds <= 0 {
		bucketSeconds = pickOpsTrafficBucketSeconds(filter.EndTime.Sub(filter.StartTime))
	}
	if bucketSeconds < 60 {
		bucketSeconds = 60
	}
	if bucketSeconds > 86400 {
		bucketSeconds = 86400
	}

	data, err := s.opsRepo.GetUpstreamTrafficStats(ctx, filter, bucketSeconds)
	if err != nil {
		if errors.Is(err, ErrOpsPreaggregatedNotPopulated) {
			return nil, infraerrors.Conflict("OPS_PREAGG_NOT_READY", "Pre-aggregated ops metrics are not populated yet")
		}
		return nil, err
	}
	return data, nil
}

func pickOpsTrafficBucketSeconds(window time.Duration) int {
	switch {
	case window <= 2*time.Hour:
		return 60
	case window <= 24*time.Hour:
		return 300
	default:
		return 3600
	}
}
