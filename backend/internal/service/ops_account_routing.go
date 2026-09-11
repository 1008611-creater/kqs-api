package service

import (
	"context"
	"errors"
	"math"
	"time"
)

type OpsAccountRoutingStatsResponse struct {
	Window        string                    `json:"window"`
	StartTime     time.Time                 `json:"start_time"`
	EndTime       time.Time                 `json:"end_time"`
	Platform      string                    `json:"platform"`
	GroupID       *int64                    `json:"group_id"`
	TotalRequests int64                     `json:"total_requests"`
	Accounts      []*OpsAccountRoutingStats `json:"accounts"`
	GeneratedAt   time.Time                 `json:"generated_at"`
}

type OpsAccountRoutingStats struct {
	AccountID       int64    `json:"account_id"`
	AccountName     string   `json:"account_name"`
	Platform        string   `json:"platform"`
	GroupID         *int64   `json:"group_id"`
	GroupName       string   `json:"group_name"`
	Priority        int      `json:"priority"`
	AccountPriority int      `json:"account_priority"`
	GroupPriority   *int     `json:"group_priority"`
	LoadFactor      int      `json:"load_factor"`
	Schedulable     bool     `json:"schedulable"`
	Status          string   `json:"status"`
	RequestCount    int64    `json:"request_count"`
	SuccessCount    int64    `json:"success_count"`
	ErrorCount      int64    `json:"error_count"`
	RequestShare    float64  `json:"request_share"`
	SuccessRate     *float64 `json:"success_rate"`
	AvgLatencyMs    *int     `json:"avg_latency_ms"`
	AvgTTFTMs       *int     `json:"avg_ttft_ms"`
	SelectionRole   string   `json:"selection_role"`
}

type opsAccountRoutingStatsReader interface {
	GetAccountRoutingStats(ctx context.Context, startTime, endTime time.Time, platform string, groupID *int64) ([]*OpsAccountRoutingStats, error)
}

func (s *OpsService) GetAccountRoutingStats(ctx context.Context, window time.Duration, windowLabel, platform string, groupID *int64) (*OpsAccountRoutingStatsResponse, error) {
	if s == nil || s.opsRepo == nil {
		return nil, errors.New("ops service is not available")
	}
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
	}
	reader, ok := s.opsRepo.(opsAccountRoutingStatsReader)
	if !ok {
		return nil, errors.New("account routing stats are not available")
	}

	end := time.Now().UTC()
	start := end.Add(-window)
	items, err := reader.GetAccountRoutingStats(ctx, start, end, platform, groupID)
	if err != nil {
		return nil, err
	}

	total := finalizeAccountRoutingStats(items)

	return &OpsAccountRoutingStatsResponse{
		Window:        windowLabel,
		StartTime:     start,
		EndTime:       end,
		Platform:      platform,
		GroupID:       groupID,
		TotalRequests: total,
		Accounts:      items,
		GeneratedAt:   end,
	}, nil
}

func finalizeAccountRoutingStats(items []*OpsAccountRoutingStats) int64 {
	var total int64
	for _, item := range items {
		if item != nil && item.RequestCount > 0 {
			total += item.RequestCount
		}
	}

	preferredPriority := math.MaxInt
	for _, item := range items {
		if item == nil || !item.Schedulable || item.Status != StatusActive {
			continue
		}
		if item.Priority < preferredPriority {
			preferredPriority = item.Priority
		}
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		if total > 0 {
			item.RequestShare = math.Round(float64(item.RequestCount)/float64(total)*10000) / 100
		}
		if item.RequestCount > 0 {
			rate := float64(item.SuccessCount) / float64(item.RequestCount) * 100
			item.SuccessRate = &rate
		}
		switch {
		case !item.Schedulable || item.Status != StatusActive:
			item.SelectionRole = "unavailable"
		case item.Priority == preferredPriority:
			item.SelectionRole = "preferred"
		default:
			item.SelectionRole = "fallback"
		}
	}
	return total
}
