package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/lib/pq"
)

// --- Plan Repository ---

type scheduledTestPlanRepository struct {
	db *sql.DB
}

func NewScheduledTestPlanRepository(db *sql.DB) service.ScheduledTestPlanRepository {
	return &scheduledTestPlanRepository{db: db}
}

func (r *scheduledTestPlanRepository) Create(ctx context.Context, plan *service.ScheduledTestPlan) (*service.ScheduledTestPlan, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO scheduled_test_plans (account_id, model_id, cron_expression, enabled, max_results, auto_recover, next_run_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, account_id, model_id, cron_expression, enabled, max_results, auto_recover, last_run_at, next_run_at, created_at, updated_at
	`, plan.AccountID, plan.ModelID, plan.CronExpression, plan.Enabled, plan.MaxResults, plan.AutoRecover, plan.NextRunAt)
	return scanPlan(row)
}

func (r *scheduledTestPlanRepository) GetByID(ctx context.Context, id int64) (*service.ScheduledTestPlan, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT id, account_id, model_id, cron_expression, enabled, max_results, auto_recover, last_run_at, next_run_at, created_at, updated_at
		FROM scheduled_test_plans WHERE id = $1
	`, id)
	return scanPlan(row)
}

func (r *scheduledTestPlanRepository) ListByAccountID(ctx context.Context, accountID int64) ([]*service.ScheduledTestPlan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, account_id, model_id, cron_expression, enabled, max_results, auto_recover, last_run_at, next_run_at, created_at, updated_at
		FROM scheduled_test_plans WHERE account_id = $1
		ORDER BY created_at DESC
	`, accountID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanPlans(rows)
}

func (r *scheduledTestPlanRepository) ListDue(ctx context.Context, now time.Time) ([]*service.ScheduledTestPlan, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, account_id, model_id, cron_expression, enabled, max_results, auto_recover, last_run_at, next_run_at, created_at, updated_at
		FROM scheduled_test_plans
		WHERE enabled = true AND next_run_at <= $1
		ORDER BY next_run_at ASC
	`, now)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()
	return scanPlans(rows)
}

func (r *scheduledTestPlanRepository) EnsurePublicStatusPlans(ctx context.Context, groupNames []string, accountNames []string, modelID string, cronExpression string, maxResults int, nextRunAt time.Time) (int64, error) {
	if len(groupNames) == 0 || modelID == "" {
		return 0, nil
	}
	if cronExpression == "" {
		cronExpression = "*/5 * * * *"
	}
	if maxResults <= 0 {
		maxResults = 120
	}
	result, err := r.db.ExecContext(ctx, `
		INSERT INTO scheduled_test_plans (account_id, model_id, cron_expression, enabled, max_results, auto_recover, next_run_at, created_at, updated_at)
		SELECT DISTINCT a.id, $2, $3, true, $4::int, true, $5::timestamptz, NOW(), NOW()
		FROM accounts a
		JOIN account_groups ag ON ag.account_id = a.id
		JOIN groups g ON g.id = ag.group_id
		WHERE g.name = ANY($1)
		  AND g.deleted_at IS NULL
		  AND a.deleted_at IS NULL
		  AND a.status = 'active'
		  AND a.schedulable = true
		  AND (a.auto_pause_on_expired = false OR a.expires_at IS NULL OR a.expires_at > NOW())
		  AND (a.overload_until IS NULL OR a.overload_until <= NOW())
		  AND (a.rate_limit_reset_at IS NULL OR a.rate_limit_reset_at <= NOW())
		  AND (a.temp_unschedulable_until IS NULL OR a.temp_unschedulable_until <= NOW())
		  AND (coalesce(cardinality($6::text[]), 0) = 0 OR a.name = ANY($6))
		  AND NOT EXISTS (
		    SELECT 1 FROM scheduled_test_plans p
		    WHERE p.account_id = a.id AND p.model_id = $2
		  )
	`, pq.Array(groupNames), modelID, cronExpression, maxResults, nextRunAt, pq.Array(accountNames))
	if err != nil {
		return 0, err
	}
	inserted, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	updated, err := r.db.ExecContext(ctx, `
		UPDATE scheduled_test_plans p
		SET cron_expression = $3,
		    enabled = true,
		    max_results = $4,
		    auto_recover = true,
		    next_run_at = CASE
		        WHEN p.enabled = false OR p.next_run_at IS NULL THEN $5
		        ELSE p.next_run_at
		    END,
		    updated_at = NOW()
		WHERE p.model_id = $2
		  AND EXISTS (
		    SELECT 1
		    FROM accounts a
		    JOIN account_groups ag ON ag.account_id = a.id
		    JOIN groups g ON g.id = ag.group_id
		    WHERE a.id = p.account_id
		      AND g.name = ANY($1)
		      AND g.deleted_at IS NULL
		      AND a.deleted_at IS NULL
		      AND a.status = 'active'
		      AND a.schedulable = true
		      AND (a.auto_pause_on_expired = false OR a.expires_at IS NULL OR a.expires_at > NOW())
		      AND (a.overload_until IS NULL OR a.overload_until <= NOW())
		      AND (a.rate_limit_reset_at IS NULL OR a.rate_limit_reset_at <= NOW())
		      AND (a.temp_unschedulable_until IS NULL OR a.temp_unschedulable_until <= NOW())
		      AND (coalesce(cardinality($6::text[]), 0) = 0 OR a.name = ANY($6))
		  )
	`, pq.Array(groupNames), modelID, cronExpression, maxResults, nextRunAt, pq.Array(accountNames))
	if err != nil {
		return inserted, err
	}
	updatedRows, err := updated.RowsAffected()
	if err != nil {
		return inserted, err
	}
	return inserted + updatedRows, nil
}

func (r *scheduledTestPlanRepository) Update(ctx context.Context, plan *service.ScheduledTestPlan) (*service.ScheduledTestPlan, error) {
	row := r.db.QueryRowContext(ctx, `
		UPDATE scheduled_test_plans
		SET model_id = $2, cron_expression = $3, enabled = $4, max_results = $5, auto_recover = $6, next_run_at = $7, updated_at = NOW()
		WHERE id = $1
		RETURNING id, account_id, model_id, cron_expression, enabled, max_results, auto_recover, last_run_at, next_run_at, created_at, updated_at
	`, plan.ID, plan.ModelID, plan.CronExpression, plan.Enabled, plan.MaxResults, plan.AutoRecover, plan.NextRunAt)
	return scanPlan(row)
}

func (r *scheduledTestPlanRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM scheduled_test_plans WHERE id = $1`, id)
	return err
}

func (r *scheduledTestPlanRepository) UpdateAfterRun(ctx context.Context, id int64, lastRunAt time.Time, nextRunAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE scheduled_test_plans SET last_run_at = $2, next_run_at = $3, updated_at = NOW() WHERE id = $1
	`, id, lastRunAt, nextRunAt)
	return err
}

// --- Result Repository ---

type scheduledTestResultRepository struct {
	db *sql.DB
}

func NewScheduledTestResultRepository(db *sql.DB) service.ScheduledTestResultRepository {
	return &scheduledTestResultRepository{db: db}
}

func (r *scheduledTestResultRepository) Create(ctx context.Context, result *service.ScheduledTestResult) (*service.ScheduledTestResult, error) {
	row := r.db.QueryRowContext(ctx, `
		INSERT INTO scheduled_test_results (plan_id, status, response_text, error_message, latency_ms, started_at, finished_at, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW())
		RETURNING id, plan_id, status, response_text, error_message, latency_ms, started_at, finished_at, created_at
	`, result.PlanID, result.Status, result.ResponseText, result.ErrorMessage, result.LatencyMs, result.StartedAt, result.FinishedAt)

	out := &service.ScheduledTestResult{}
	if err := row.Scan(
		&out.ID, &out.PlanID, &out.Status, &out.ResponseText, &out.ErrorMessage,
		&out.LatencyMs, &out.StartedAt, &out.FinishedAt, &out.CreatedAt,
	); err != nil {
		return nil, err
	}
	return out, nil
}

func (r *scheduledTestResultRepository) ListByPlanID(ctx context.Context, planID int64, limit int) ([]*service.ScheduledTestResult, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, plan_id, status, response_text, error_message, latency_ms, started_at, finished_at, created_at
		FROM scheduled_test_results
		WHERE plan_id = $1
		ORDER BY created_at DESC
		LIMIT $2
	`, planID, limit)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var results []*service.ScheduledTestResult
	for rows.Next() {
		r := &service.ScheduledTestResult{}
		if err := rows.Scan(
			&r.ID, &r.PlanID, &r.Status, &r.ResponseText, &r.ErrorMessage,
			&r.LatencyMs, &r.StartedAt, &r.FinishedAt, &r.CreatedAt,
		); err != nil {
			return nil, err
		}
		results = append(results, r)
	}
	return results, rows.Err()
}

func (r *scheduledTestResultRepository) GetPublicChannelStatusStats(ctx context.Context, groupNames []string, accountNames []string, modelID string, now time.Time, recentLimit int) (*service.PublicChannelStatusStats, error) {
	if len(groupNames) == 0 || modelID == "" {
		return &service.PublicChannelStatusStats{}, nil
	}
	if recentLimit <= 0 {
		recentLimit = 120
	}

	stats := &service.PublicChannelStatusStats{}
	row := r.db.QueryRowContext(ctx, `
		WITH eligible_accounts AS (
			SELECT DISTINCT a.id
			FROM accounts a
			JOIN account_groups ag ON ag.account_id = a.id
			JOIN groups g ON g.id = ag.group_id
			WHERE g.name = ANY($1)
			  AND g.deleted_at IS NULL
			  AND a.deleted_at IS NULL
			  AND a.status = 'active'
			  AND a.schedulable = true
			  AND (coalesce(cardinality($4::text[]), 0) = 0 OR a.name = ANY($4))
		),
		eligible_plans AS (
			SELECT p.id, p.account_id
			FROM scheduled_test_plans p
			JOIN eligible_accounts ea ON ea.id = p.account_id
			WHERE p.model_id = $2 AND p.enabled = true
		),
		latest_per_account AS (
			SELECT DISTINCT ON (ep.account_id)
				ep.account_id,
				r.status,
				r.created_at
			FROM eligible_plans ep
			JOIN scheduled_test_results r ON r.plan_id = ep.id
			ORDER BY ep.account_id, r.created_at DESC
		),
		recent_1h AS (
			SELECT r.status, r.latency_ms
			FROM eligible_plans ep
			JOIN scheduled_test_results r ON r.plan_id = ep.id
			WHERE r.created_at >= $3::timestamptz - INTERVAL '1 hour'
		),
		recent_24h AS (
			SELECT r.status
			FROM eligible_plans ep
			JOIN scheduled_test_results r ON r.plan_id = ep.id
			WHERE r.created_at >= $3::timestamptz - INTERVAL '24 hours'
		)
		SELECT
			(SELECT COUNT(*) FROM eligible_accounts) AS active_channels,
			(SELECT COUNT(*) FROM latest_per_account) AS checked_channels,
			COALESCE((SELECT 100.0 * COUNT(*) FILTER (WHERE status = 'success') / NULLIF(COUNT(*), 0) FROM recent_1h), 0) AS success_rate_1h,
			COALESCE((SELECT 100.0 * COUNT(*) FILTER (WHERE status = 'success') / NULLIF(COUNT(*), 0) FROM recent_24h), 0) AS success_rate_24h,
			(SELECT AVG(latency_ms)::int FROM recent_1h WHERE status = 'success' AND latency_ms > 0) AS avg_latency_ms,
			(SELECT percentile_cont(0.95) WITHIN GROUP (ORDER BY latency_ms)::int FROM recent_1h WHERE status = 'success' AND latency_ms > 0) AS p95_latency_ms,
			(SELECT MAX(created_at) FROM latest_per_account) AS last_checked_at
	`, pq.Array(groupNames), modelID, now, pq.Array(accountNames))

	var avgLatency sql.NullInt64
	var p95Latency sql.NullInt64
	var lastChecked sql.NullTime
	if err := row.Scan(
		&stats.ActiveChannels,
		&stats.CheckedChannels,
		&stats.SuccessRate1h,
		&stats.SuccessRate24h,
		&avgLatency,
		&p95Latency,
		&lastChecked,
	); err != nil {
		return nil, err
	}
	if avgLatency.Valid {
		v := int(avgLatency.Int64)
		stats.AverageLatencyMs = &v
	}
	if p95Latency.Valid {
		v := int(p95Latency.Int64)
		stats.P95LatencyMs = &v
	}
	if lastChecked.Valid {
		stats.LastCheckedAt = &lastChecked.Time
	}

	rows, err := r.db.QueryContext(ctx, `
		WITH eligible_accounts AS (
			SELECT DISTINCT a.id
			FROM accounts a
			JOIN account_groups ag ON ag.account_id = a.id
			JOIN groups g ON g.id = ag.group_id
			WHERE g.name = ANY($1)
			  AND g.deleted_at IS NULL
			  AND a.deleted_at IS NULL
			  AND a.status = 'active'
			  AND a.schedulable = true
			  AND (coalesce(cardinality($5::text[]), 0) = 0 OR a.name = ANY($5))
		),
		eligible_plans AS (
			SELECT p.id, p.account_id
			FROM scheduled_test_plans p
			JOIN eligible_accounts ea ON ea.id = p.account_id
			WHERE p.model_id = $2 AND p.enabled = true
		)
		SELECT ep.account_id, r.status, r.latency_ms, r.created_at
		FROM eligible_plans ep
		JOIN scheduled_test_results r ON r.plan_id = ep.id
		WHERE r.created_at >= $3::timestamptz - INTERVAL '2 hours'
		ORDER BY r.created_at DESC
		LIMIT $4
	`, pq.Array(groupNames), modelID, now, recentLimit, pq.Array(accountNames))
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var sample service.PublicChannelStatusSample
		if err := rows.Scan(&sample.AccountID, &sample.Status, &sample.LatencyMs, &sample.CheckedAt); err != nil {
			return nil, err
		}
		stats.RecentSamples = append(stats.RecentSamples, sample)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return stats, nil
}

func (r *scheduledTestResultRepository) PruneOldResults(ctx context.Context, planID int64, keepCount int) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM scheduled_test_results
		WHERE id IN (
			SELECT id FROM (
				SELECT id, ROW_NUMBER() OVER (PARTITION BY plan_id ORDER BY created_at DESC) AS rn
				FROM scheduled_test_results
				WHERE plan_id = $1
			) ranked
			WHERE rn > $2
		)
	`, planID, keepCount)
	return err
}

// --- scan helpers ---

type scannable interface {
	Scan(dest ...any) error
}

func scanPlan(row scannable) (*service.ScheduledTestPlan, error) {
	p := &service.ScheduledTestPlan{}
	if err := row.Scan(
		&p.ID, &p.AccountID, &p.ModelID, &p.CronExpression, &p.Enabled, &p.MaxResults, &p.AutoRecover,
		&p.LastRunAt, &p.NextRunAt, &p.CreatedAt, &p.UpdatedAt,
	); err != nil {
		return nil, err
	}
	return p, nil
}

func scanPlans(rows *sql.Rows) ([]*service.ScheduledTestPlan, error) {
	var plans []*service.ScheduledTestPlan
	for rows.Next() {
		p, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}
