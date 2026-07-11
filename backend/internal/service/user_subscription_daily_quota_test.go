package service

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

type dailyResetTrackingUserSubRepo struct {
	userSubRepoNoop

	resetDailyCalled bool
	lastRolloverUSD  float64
}

func (r *dailyResetTrackingUserSubRepo) ResetDailyUsage(_ context.Context, _ int64, _ time.Time, rolloverUSD float64) error {
	r.resetDailyCalled = true
	r.lastRolloverUSD = rolloverUSD
	return nil
}

func TestAssignOrExtendSubscription_ExpiredDailyCardStartsNewOneTimeQuota(t *testing.T) {
	groupRepo := &subscriptionGroupRepoStub{
		group: &Group{ID: 1, SubscriptionType: SubscriptionTypeSubscription},
	}
	subRepo := newSubscriptionUserSubRepoStub()
	oldStart := time.Now().AddDate(0, 0, -3)
	oldWindowStart := startOfDay(oldStart)
	subRepo.seed(&UserSubscription{
		ID:                 100,
		UserID:             200,
		GroupID:            1,
		StartsAt:           oldStart,
		ExpiresAt:          oldStart.AddDate(0, 0, 1),
		Status:             SubscriptionStatusExpired,
		DailyWindowStart:   &oldWindowStart,
		WeeklyWindowStart:  &oldWindowStart,
		MonthlyWindowStart: &oldWindowStart,
		DailyUsageUSD:      10,
		WeeklyUsageUSD:     20,
		MonthlyUsageUSD:    30,
		Notes:              "old",
	})
	svc := NewSubscriptionService(groupRepo, subRepo, nil, nil, nil)

	renewed, reused, err := svc.AssignOrExtendSubscription(context.Background(), &AssignSubscriptionInput{
		UserID:       200,
		GroupID:      1,
		ValidityDays: 1,
		Notes:        "new",
	})

	require.NoError(t, err)
	require.True(t, reused)
	require.True(t, renewed.HasOneTimeDailyQuota(), "过期后重新购买 1 日卡仍应被识别为一次性日额度")
	require.Equal(t, SubscriptionStatusActive, renewed.Status)
	require.True(t, renewed.StartsAt.After(oldStart), "重新购买过期订阅时应重置当前周期 StartsAt")
	require.False(t, renewed.ExpiresAt.After(renewed.StartsAt.AddDate(0, 0, 1)))
	require.NotNil(t, renewed.DailyWindowStart)
	require.Equal(t, startOfDay(renewed.StartsAt), *renewed.DailyWindowStart)
	require.Equal(t, 0.0, renewed.DailyUsageUSD)
	require.Equal(t, 0.0, renewed.WeeklyUsageUSD)
	require.Equal(t, 0.0, renewed.MonthlyUsageUSD)
	require.Equal(t, "old\nnew", renewed.Notes)
}

func TestUserSubscriptionNeedsDailyReset_DailyCardRefreshesWithRollover(t *testing.T) {
	start := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	dailyWindowStart := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	sub := &UserSubscription{
		StartsAt:         start,
		ExpiresAt:        start.Add(24 * time.Hour),
		DailyWindowStart: &dailyWindowStart,
		DailyUsageUSD:    10,
	}

	require.True(t, sub.HasOneTimeDailyQuota())
	require.True(t, sub.NeedsDailyResetAt(dailyWindowStart.Add(25*time.Hour)), "日卡跨过日窗口后应触发结转刷新")
}

func TestUserSubscriptionNeedsDailyReset_MultiDaySubscriptionStillRefreshes(t *testing.T) {
	start := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	dailyWindowStart := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	sub := &UserSubscription{
		StartsAt:         start,
		ExpiresAt:        start.AddDate(0, 0, 2),
		DailyWindowStart: &dailyWindowStart,
	}

	require.False(t, sub.HasOneTimeDailyQuota())
	require.True(t, sub.NeedsDailyResetAt(dailyWindowStart.Add(24*time.Hour)), "多日订阅仍应按 24 小时日窗口刷新")
}

func TestUserSubscriptionDailyResetTime_DailyCardReturnsDailyWindowEnd(t *testing.T) {
	start := time.Date(2026, 5, 18, 12, 0, 0, 0, time.UTC)
	dailyWindowStart := time.Date(2026, 5, 18, 0, 0, 0, 0, time.UTC)
	sub := &UserSubscription{
		StartsAt:         start,
		ExpiresAt:        start.Add(24 * time.Hour),
		DailyWindowStart: &dailyWindowStart,
	}

	resetAt := sub.DailyResetTime()
	require.NotNil(t, resetAt)
	require.Equal(t, dailyWindowStart.Add(24*time.Hour), *resetAt, "日卡展示的日额度刷新时间应为日窗口结束时间")
}

func TestUserSubscriptionDailyResetTime_UsesNextMidnightForNonMidnightWindowStart(t *testing.T) {
	dailyWindowStart := time.Date(2026, 5, 18, 15, 30, 0, 0, time.UTC)
	sub := &UserSubscription{DailyWindowStart: &dailyWindowStart}

	resetAt := sub.DailyResetTime()

	require.NotNil(t, resetAt)
	require.Equal(t, time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC), *resetAt)
}

func TestUserSubscriptionNeedsDailyReset_UsesNextMidnightForNonMidnightWindowStart(t *testing.T) {
	dailyWindowStart := time.Date(2026, 5, 18, 15, 30, 0, 0, time.UTC)
	sub := &UserSubscription{DailyWindowStart: &dailyWindowStart}

	require.False(t, sub.NeedsDailyResetAt(time.Date(2026, 5, 18, 23, 59, 59, 0, time.UTC)))
	require.True(t, sub.NeedsDailyResetAt(time.Date(2026, 5, 19, 0, 0, 0, 0, time.UTC)))
}

func TestCheckAndResetWindows_DailyCardResetsAndCarriesRollover(t *testing.T) {
	now := time.Now()
	startsAt := now.Add(-23 * time.Hour)
	dailyWindowStart := now.Add(-25 * time.Hour)
	repo := &dailyResetTrackingUserSubRepo{}
	svc := NewSubscriptionService(groupRepoNoop{}, repo, nil, nil, nil)
	dailyLimit := 15.0
	sub := &UserSubscription{
		ID:               1,
		UserID:           10,
		GroupID:          20,
		StartsAt:         startsAt,
		ExpiresAt:        startsAt.Add(24 * time.Hour),
		DailyUsageUSD:    10,
		Group:            &Group{DailyLimitUSD: &dailyLimit},
		DailyWindowStart: &dailyWindowStart,
	}

	err := svc.CheckAndResetWindows(context.Background(), sub)

	require.NoError(t, err)
	require.True(t, repo.resetDailyCalled, "日卡跨过日窗口后应重置 daily usage")
	require.Equal(t, 5.0, repo.lastRolloverUSD)
	require.Equal(t, 0.0, sub.DailyUsageUSD)
	require.Equal(t, 5.0, sub.DailyRolloverUSD)
}

func TestCheckAndResetWindows_MultiDaySubscriptionStillResetsDailyUsage(t *testing.T) {
	now := time.Now()
	startsAt := now.Add(-48 * time.Hour)
	dailyWindowStart := now.Add(-25 * time.Hour)
	repo := &dailyResetTrackingUserSubRepo{}
	svc := NewSubscriptionService(groupRepoNoop{}, repo, nil, nil, nil)
	sub := &UserSubscription{
		ID:               1,
		UserID:           10,
		GroupID:          20,
		StartsAt:         startsAt,
		ExpiresAt:        startsAt.AddDate(0, 0, 2),
		DailyUsageUSD:    10,
		DailyWindowStart: &dailyWindowStart,
	}

	err := svc.CheckAndResetWindows(context.Background(), sub)

	require.NoError(t, err)
	require.True(t, repo.resetDailyCalled, "多日订阅仍应重置过期 daily window")
	require.Equal(t, 0.0, sub.DailyUsageUSD)
}

func TestValidateAndCheckLimits_DailyCardCarriesRemainingQuotaAfterMidnight(t *testing.T) {
	start := time.Now().Add(-23 * time.Hour)
	dailyWindowStart := time.Now().Add(-25 * time.Hour)
	dailyLimit := 10.0
	sub := &UserSubscription{
		Status:           SubscriptionStatusActive,
		StartsAt:         start,
		ExpiresAt:        start.Add(24 * time.Hour),
		DailyWindowStart: &dailyWindowStart,
		DailyUsageUSD:    3,
	}
	group := &Group{
		SubscriptionType: SubscriptionTypeSubscription,
		DailyLimitUSD:    &dailyLimit,
	}
	svc := NewSubscriptionService(groupRepoNoop{}, userSubRepoNoop{}, nil, nil, nil)

	needsMaintenance, err := svc.ValidateAndCheckLimits(sub, group)

	require.True(t, needsMaintenance, "日卡跨过日窗口后应触发 daily reset 维护")
	require.NoError(t, err)
	require.Equal(t, 0.0, sub.DailyUsageUSD)
	require.Equal(t, 7.0, sub.DailyRolloverUSD)
}

func TestUserSubscriptionHasBillableQuotaRejectsExhaustedWindow(t *testing.T) {
	dailyLimit := 90.0
	group := &Group{DailyLimitUSD: &dailyLimit}
	sub := &UserSubscription{
		Status:        SubscriptionStatusActive,
		StartsAt:      time.Now().Add(-time.Hour),
		ExpiresAt:     time.Now().Add(23 * time.Hour),
		DailyUsageUSD: 90,
	}

	require.False(t, sub.HasBillableQuota(group))

	sub.DailyUsageUSD = 89.99
	require.True(t, sub.HasBillableQuota(group))
}

func TestCalculateNextDailyRollover_ExpiresPreviousRolloverAfterOneDay(t *testing.T) {
	dailyLimit := 10.0
	group := &Group{DailyLimitUSD: &dailyLimit}

	sub := &UserSubscription{
		DailyUsageUSD:    0,
		DailyRolloverUSD: 10,
	}
	require.Equal(t, 10.0, sub.CalculateNextDailyRolloverUSD(group), "昨天结转未用完时，不能继续叠加成 20")

	sub.DailyUsageUSD = 3
	sub.DailyRolloverUSD = 10
	require.Equal(t, 10.0, sub.CalculateNextDailyRolloverUSD(group), "先消费昨天结转，今天额度未动时只结转今天 10")

	sub.DailyUsageUSD = 12
	sub.DailyRolloverUSD = 10
	require.Equal(t, 8.0, sub.CalculateNextDailyRolloverUSD(group), "用完昨天结转后，今天额度剩多少才结转多少")

	sub.DailyUsageUSD = 25
	sub.DailyRolloverUSD = 10
	require.Equal(t, 0.0, sub.CalculateNextDailyRolloverUSD(group), "今天额度也用完后不再结转")
}
