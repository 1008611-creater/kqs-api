package service

import (
	"context"
	"time"
)

// ScheduledTestPlan represents a scheduled test plan domain model.
type ScheduledTestPlan struct {
	ID             int64      `json:"id"`
	AccountID      int64      `json:"account_id"`
	ModelID        string     `json:"model_id"`
	CronExpression string     `json:"cron_expression"`
	Enabled        bool       `json:"enabled"`
	MaxResults     int        `json:"max_results"`
	AutoRecover    bool       `json:"auto_recover"`
	LastRunAt      *time.Time `json:"last_run_at"`
	NextRunAt      *time.Time `json:"next_run_at"`
	CreatedAt      time.Time  `json:"created_at"`
	UpdatedAt      time.Time  `json:"updated_at"`
}

// ScheduledTestResult represents a single test execution result.
type ScheduledTestResult struct {
	ID           int64     `json:"id"`
	PlanID       int64     `json:"plan_id"`
	Status       string    `json:"status"`
	ResponseText string    `json:"response_text"`
	ErrorMessage string    `json:"error_message"`
	LatencyMs    int64     `json:"latency_ms"`
	StartedAt    time.Time `json:"started_at"`
	FinishedAt   time.Time `json:"finished_at"`
	CreatedAt    time.Time `json:"created_at"`
}

type PublicChannelStatusSnapshot struct {
	GroupName        string
	ModelID          string
	Status           string
	ActiveChannels   int
	CheckedChannels  int
	SuccessRate1h    float64
	SuccessRate24h   float64
	AverageLatencyMs *int
	P95LatencyMs     *int
	LastCheckedAt    *time.Time
	Timeline         []PublicChannelStatusPoint
	GeneratedAt      time.Time
}

type PublicChannelStatusPoint struct {
	Status    string
	CheckedAt time.Time
}

type PublicChannelStatusSample struct {
	AccountID int64
	Status    string
	LatencyMs int64
	CheckedAt time.Time
}

// ScheduledTestAccountTelemetry is a small, read-only probe summary used by
// the gateway scheduler after a process restart. It deliberately contains no
// upstream response body or credential material.
type ScheduledTestAccountTelemetry struct {
	AccountID        int64
	SampleCount      int
	SuccessRate      float64
	AverageLatencyMs float64
	LastCheckedAt    time.Time
}

type PublicChannelStatusStats struct {
	ActiveChannels   int
	CheckedChannels  int
	SuccessRate1h    float64
	SuccessRate24h   float64
	AverageLatencyMs *int
	P95LatencyMs     *int
	LastCheckedAt    *time.Time
	RecentSamples    []PublicChannelStatusSample
}

// ScheduledTestPlanRepository defines the data access interface for test plans.
type ScheduledTestPlanRepository interface {
	Create(ctx context.Context, plan *ScheduledTestPlan) (*ScheduledTestPlan, error)
	GetByID(ctx context.Context, id int64) (*ScheduledTestPlan, error)
	ListByAccountID(ctx context.Context, accountID int64) ([]*ScheduledTestPlan, error)
	ListDue(ctx context.Context, now time.Time) ([]*ScheduledTestPlan, error)
	EnsurePublicStatusPlans(ctx context.Context, groupNames []string, accountNames []string, modelID string, cronExpression string, maxResults int, nextRunAt time.Time) (int64, error)
	Update(ctx context.Context, plan *ScheduledTestPlan) (*ScheduledTestPlan, error)
	Delete(ctx context.Context, id int64) error
	UpdateAfterRun(ctx context.Context, id int64, lastRunAt time.Time, nextRunAt time.Time) error
}

// ScheduledTestResultRepository defines the data access interface for test results.
type ScheduledTestResultRepository interface {
	Create(ctx context.Context, result *ScheduledTestResult) (*ScheduledTestResult, error)
	ListByPlanID(ctx context.Context, planID int64, limit int) ([]*ScheduledTestResult, error)
	GetPublicChannelStatusStats(ctx context.Context, groupNames []string, accountNames []string, modelID string, now time.Time, recentLimit int) (*PublicChannelStatusStats, error)
	GetRecentAccountTelemetry(ctx context.Context, accountIDs []int64, now time.Time) (map[int64]ScheduledTestAccountTelemetry, error)
	PruneOldResults(ctx context.Context, planID int64, keepCount int) error
}
