package repository

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestOpsRepositoryGetAccountRoutingStatsReturnsAccountsAndSamples(t *testing.T) {
	db, mock := newSQLMock(t)
	repo := &opsRepository{db: db}
	start := time.Date(2026, 8, 2, 10, 0, 0, 0, time.UTC)
	end := start.Add(time.Hour)
	groupID := int64(40)

	rows := mock.NewRows([]string{
		"id", "name", "platform", "account_priority", "group_priority", "group_id", "group_name",
		"load_factor", "schedulable", "status", "request_count", "success_count", "error_count",
		"avg_latency_ms", "avg_ttft_ms",
	}).AddRow(7, "route-a", "openai", 5, 2, groupID, "private", 8, true, "active", int64(42), int64(40), int64(2), 3200.4, 850.2)

	mock.ExpectQuery(`WITH request_samples`).
		WithArgs(start, end, "openai", groupID).
		WillReturnRows(rows)

	items, err := repo.GetAccountRoutingStats(context.Background(), start, end, "openai", &groupID)
	require.NoError(t, err)
	require.Len(t, items, 1)
	require.Equal(t, int64(7), items[0].AccountID)
	require.Equal(t, 2, items[0].Priority)
	require.NotNil(t, items[0].GroupPriority)
	require.Equal(t, 3200, *items[0].AvgLatencyMs)
	require.Equal(t, 850, *items[0].AvgTTFTMs)
	require.NoError(t, mock.ExpectationsWereMet())
}
