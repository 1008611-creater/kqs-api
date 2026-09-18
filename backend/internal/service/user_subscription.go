package service

import "time"

const quotaComparisonEpsilon = 1e-9

type UserSubscription struct {
	ID      int64
	UserID  int64
	GroupID int64

	StartsAt     time.Time
	ExpiresAt    time.Time
	Status       string
	ValidityDays int

	DailyWindowStart   *time.Time
	WeeklyWindowStart  *time.Time
	MonthlyWindowStart *time.Time

	DailyUsageUSD    float64
	DailyRolloverUSD float64
	WeeklyUsageUSD   float64
	MonthlyUsageUSD  float64

	DailyLimitOverrideUSD   *float64
	WeeklyLimitOverrideUSD  *float64
	MonthlyLimitOverrideUSD *float64
	QuotaBonusMultiplier    float64
	QuotaBonusSource        string

	AssignedBy *int64
	AssignedAt time.Time
	Notes      string

	CreatedAt time.Time
	UpdatedAt time.Time

	User           *User
	Group          *Group
	AssignedByUser *User
}

func (s *UserSubscription) IsActive() bool {
	return s.Status == SubscriptionStatusActive && time.Now().Before(s.ExpiresAt)
}

func (s *UserSubscription) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

func (s *UserSubscription) DaysRemaining() int {
	if s.IsExpired() {
		return 0
	}
	return int(time.Until(s.ExpiresAt).Hours() / 24)
}

func (s *UserSubscription) IsWindowActivated() bool {
	return s.DailyWindowStart != nil || s.WeeklyWindowStart != nil || s.MonthlyWindowStart != nil
}

func (s *UserSubscription) HasOneTimeDailyQuota() bool {
	if s == nil || s.StartsAt.IsZero() || s.ExpiresAt.IsZero() {
		return false
	}
	return !s.ExpiresAt.After(s.StartsAt.AddDate(0, 0, 1))
}

func (s *UserSubscription) NeedsDailyReset() bool {
	return s.NeedsDailyResetAt(time.Now())
}

func (s *UserSubscription) NeedsDailyResetAt(now time.Time) bool {
	if s.DailyWindowStart == nil {
		return false
	}
	return !now.Before(nextDailyResetBoundary(*s.DailyWindowStart))
}

func (s *UserSubscription) NeedsWeeklyReset() bool {
	if s.WeeklyWindowStart == nil {
		return false
	}
	return time.Since(*s.WeeklyWindowStart) >= 7*24*time.Hour
}

func (s *UserSubscription) NeedsMonthlyReset() bool {
	if s.MonthlyWindowStart == nil {
		return false
	}
	return time.Since(*s.MonthlyWindowStart) >= 30*24*time.Hour
}

func (s *UserSubscription) DailyResetTime() *time.Time {
	if s.DailyWindowStart == nil {
		return nil
	}
	t := nextDailyResetBoundary(*s.DailyWindowStart)
	return &t
}

func (s *UserSubscription) WeeklyResetTime() *time.Time {
	if s.WeeklyWindowStart == nil {
		return nil
	}
	t := s.WeeklyWindowStart.Add(7 * 24 * time.Hour)
	return &t
}

func (s *UserSubscription) MonthlyResetTime() *time.Time {
	if s.MonthlyWindowStart == nil {
		return nil
	}
	t := s.MonthlyWindowStart.Add(30 * 24 * time.Hour)
	return &t
}

func (s *UserSubscription) CheckDailyLimit(group *Group, additionalCost float64) bool {
	available := s.EffectiveDailyAvailableUSD(group)
	if available == nil {
		return true
	}
	return s.DailyUsageUSD+additionalCost <= *available
}

func (s *UserSubscription) CheckWeeklyLimit(group *Group, additionalCost float64) bool {
	limit := s.EffectiveWeeklyLimitUSD(group)
	if limit == nil {
		return true
	}
	return s.WeeklyUsageUSD+additionalCost <= *limit
}

func (s *UserSubscription) CheckMonthlyLimit(group *Group, additionalCost float64) bool {
	limit := s.EffectiveMonthlyLimitUSD(group)
	if limit == nil {
		return true
	}
	return s.MonthlyUsageUSD+additionalCost <= *limit
}

func (s *UserSubscription) CheckAllLimits(group *Group, additionalCost float64) (daily, weekly, monthly bool) {
	daily = s.CheckDailyLimit(group, additionalCost)
	weekly = s.CheckWeeklyLimit(group, additionalCost)
	monthly = s.CheckMonthlyLimit(group, additionalCost)
	return
}

// CurrentWindowSnapshot returns a copy with expired usage windows reset in
// memory. It does not write to the database; callers can use it for routing
// and preflight checks before async maintenance persists the reset.
func (s *UserSubscription) CurrentWindowSnapshot(group *Group) UserSubscription {
	if s == nil {
		return UserSubscription{}
	}
	cp := *s
	if cp.NeedsDailyReset() {
		cp.DailyRolloverUSD = cp.CalculateNextDailyRolloverUSD(group)
		cp.DailyUsageUSD = 0
	}
	if cp.NeedsWeeklyReset() {
		cp.WeeklyUsageUSD = 0
	}
	if cp.NeedsMonthlyReset() {
		cp.MonthlyUsageUSD = 0
	}
	return cp
}

// HasBillableQuota reports whether this subscription has any remaining quota
// in every configured window. Unlimited windows are treated as available.
func (s *UserSubscription) HasBillableQuota(group *Group) bool {
	if s == nil {
		return false
	}
	if s.Status != SubscriptionStatusActive || s.IsExpired() {
		return false
	}
	cp := s.CurrentWindowSnapshot(group)

	if available := cp.EffectiveDailyAvailableUSD(group); available != nil && cp.DailyUsageUSD >= *available-quotaComparisonEpsilon {
		return false
	}
	if limit := cp.EffectiveWeeklyLimitUSD(group); limit != nil && cp.WeeklyUsageUSD >= *limit-quotaComparisonEpsilon {
		return false
	}
	if limit := cp.EffectiveMonthlyLimitUSD(group); limit != nil && cp.MonthlyUsageUSD >= *limit-quotaComparisonEpsilon {
		return false
	}
	return true
}

func (s *UserSubscription) EffectiveDailyLimitUSD(group *Group) *float64 {
	if s != nil && s.DailyLimitOverrideUSD != nil && *s.DailyLimitOverrideUSD > 0 {
		return s.DailyLimitOverrideUSD
	}
	if group != nil && group.HasDailyLimit() {
		return group.DailyLimitUSD
	}
	return nil
}

func (s *UserSubscription) EffectiveDailyAvailableUSD(group *Group) *float64 {
	limit := s.EffectiveDailyLimitUSD(group)
	if limit == nil {
		return nil
	}
	available := *limit
	if s != nil && s.DailyRolloverUSD > 0 {
		available += s.DailyRolloverUSD
	}
	return &available
}

func (s *UserSubscription) DailyRemainingUSD(group *Group) *float64 {
	available := s.EffectiveDailyAvailableUSD(group)
	if available == nil {
		return nil
	}
	remaining := *available - s.DailyUsageUSD
	if remaining < 0 {
		remaining = 0
	}
	return &remaining
}

func (s *UserSubscription) CalculateNextDailyRolloverUSD(group *Group) float64 {
	limit := s.EffectiveDailyLimitUSD(group)
	if limit == nil {
		return 0
	}
	previousRollover := s.DailyRolloverUSD
	if previousRollover < 0 {
		previousRollover = 0
	}
	usageFromTodayQuota := s.DailyUsageUSD - previousRollover
	if usageFromTodayQuota < 0 {
		usageFromTodayQuota = 0
	}
	rollover := *limit - usageFromTodayQuota
	if rollover < 0 {
		return 0
	}
	if rollover > *limit {
		return *limit
	}
	return rollover
}

func (s *UserSubscription) EffectiveWeeklyLimitUSD(group *Group) *float64 {
	if s != nil && s.WeeklyLimitOverrideUSD != nil && *s.WeeklyLimitOverrideUSD > 0 {
		return s.WeeklyLimitOverrideUSD
	}
	if group != nil && group.HasWeeklyLimit() {
		return group.WeeklyLimitUSD
	}
	return nil
}

func (s *UserSubscription) EffectiveMonthlyLimitUSD(group *Group) *float64 {
	if s != nil && s.MonthlyLimitOverrideUSD != nil && *s.MonthlyLimitOverrideUSD > 0 {
		return s.MonthlyLimitOverrideUSD
	}
	if group != nil && group.HasMonthlyLimit() {
		return group.MonthlyLimitUSD
	}
	return nil
}

func nextDailyResetBoundary(windowStart time.Time) time.Time {
	dayStart := time.Date(windowStart.Year(), windowStart.Month(), windowStart.Day(), 0, 0, 0, 0, windowStart.Location())
	return dayStart.AddDate(0, 0, 1)
}
