package service

import (
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/timezone"
)

const subscriptionDayDuration = 24 * time.Hour

const quotaComparisonEpsilon = 1e-9

type UserSubscription struct {
	ID      int64
	UserID  int64
	GroupID int64

	StartsAt  time.Time
	ExpiresAt time.Time
	Status    string
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
	DeletedAt *time.Time

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
	return s.daysRemainingAt(time.Now())
}

func (s *UserSubscription) daysRemainingAt(now time.Time) int {
	remaining := s.ExpiresAt.Sub(now)
	if remaining <= 0 {
		return 0
	}

	days := int(remaining / subscriptionDayDuration)
	if remaining%subscriptionDayDuration != 0 {
		days++
	}
	return days
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
<<<<<<< HEAD
	_, ok := s.automaticDailyWindowStartAt(now)
	return ok
=======
	if s.DailyWindowStart == nil {
		return false
	}
	return !now.Before(nextDailyResetBoundary(*s.DailyWindowStart))
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
}

func (s *UserSubscription) NeedsWeeklyReset() bool {
	return s.NeedsWeeklyResetAt(time.Now())
}

func (s *UserSubscription) NeedsWeeklyResetAt(now time.Time) bool {
	if s.WeeklyWindowStart == nil {
		return false
	}
	return !now.Before(s.WeeklyWindowStart.Add(7 * 24 * time.Hour))
}

func (s *UserSubscription) NeedsMonthlyReset() bool {
	return s.NeedsMonthlyResetAt(time.Now())
}

func (s *UserSubscription) NeedsMonthlyResetAt(now time.Time) bool {
	if s.MonthlyWindowStart == nil {
		return false
	}
	return !now.Before(s.MonthlyWindowStart.Add(30 * 24 * time.Hour))
}

func (s *UserSubscription) canAutomaticallyResetDailyAt(now time.Time) bool {
	_, ok := s.automaticDailyWindowStartAt(now)
	return ok
}

// automaticDailyWindowStartAt 计算日窗口按“配置时区日历日”对齐后的当前窗口起点。
// 日额度固定在每天 0 点刷新（与周/月的期限对齐滚动窗口语义不同），因此只要持久化
// 的窗口起点落在更早的日历日，就允许推进到今天 0 点。手动重置、激活等写入的任何
// 非 0 点锚点都会在下一个 0 点被拉回日历日边界，不会永久漂移刷新时刻。
func (s *UserSubscription) automaticDailyWindowStartAt(now time.Time) (time.Time, bool) {
	if s.DailyWindowStart == nil {
		return time.Time{}, false
	}
	if s.HasOneTimeDailyQuota() {
		return time.Time{}, false
	}
	today := timezone.StartOfDay(now)
	if !today.After(timezone.StartOfDay(*s.DailyWindowStart)) {
		return time.Time{}, false
	}
	return today, true
}

func (s *UserSubscription) canAutomaticallyResetWeeklyAt(now time.Time) bool {
	_, ok := s.automaticWindowStartAt(s.WeeklyWindowStart, 7*24*time.Hour, now)
	return ok
}

func (s *UserSubscription) canAutomaticallyResetMonthlyAt(now time.Time) bool {
	_, ok := s.automaticWindowStartAt(s.MonthlyWindowStart, 30*24*time.Hour, now)
	return ok
}

// windowResetAnchor 返回周/月窗口实际推进所依据的锚点。
// 早期订阅把首个窗口初始化在开通日零点；只有这个初始值是无歧义的，之后出现的
// 零点锚点可能来自手动重置，必须保持权威。
// 自动推进（automaticWindowStartAt）与对外展示的重置时间（WeeklyResetTime/
// MonthlyResetTime）必须共用这一修正，否则仪表盘显示的重置时间会早于窗口实际
// 滚动的时间。
// 日窗口按日历日对齐（automaticDailyWindowStartAt），不走这里。
func (s *UserSubscription) windowResetAnchor(previous time.Time) time.Time {
	legacyAnchor := startOfDay(s.StartsAt)
	if legacyAnchor.Before(s.StartsAt) && previous.Equal(legacyAnchor) {
		return s.StartsAt
	}
	return previous
}

// automaticWindowStartAt 计算周/月窗口（期限对齐滚动窗口）的当前窗口起点。
// 窗口从锚点按整数个 period 步进，且不越过订阅到期时间，避免最后一个不完整
// 周期重复发放额度（issue #5051）。日窗口不走此函数，见 automaticDailyWindowStartAt。
func (s *UserSubscription) automaticWindowStartAt(previous *time.Time, period time.Duration, now time.Time) (time.Time, bool) {
	if previous == nil {
		return time.Time{}, false
	}

	anchor := s.windowResetAnchor(*previous)
	next := anchor.Add(period)
	if now.Before(next) || !next.Before(s.ExpiresAt) {
		return time.Time{}, false
	}

	periods := now.Sub(anchor) / period
	lastPeriodBeforeExpiry := (s.ExpiresAt.Sub(anchor) - 1) / period
	if periods > lastPeriodBeforeExpiry {
		periods = lastPeriodBeforeExpiry
	}
	return anchor.Add(periods * period), true
}

func (s *UserSubscription) DailyResetTime() *time.Time {
	if s.DailyWindowStart == nil {
		return nil
	}
<<<<<<< HEAD
	if s.HasOneTimeDailyQuota() {
		t := s.ExpiresAt
		return &t
	}
	// 日窗口按日历日对齐：下次刷新固定在窗口起点所在日的次日 0 点。
	t := timezone.StartOfDay(*s.DailyWindowStart).AddDate(0, 0, 1)
=======
	t := nextDailyResetBoundary(*s.DailyWindowStart)
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
	return &t
}

func (s *UserSubscription) WeeklyResetTime() *time.Time {
	if s.WeeklyWindowStart == nil {
		return nil
	}
	t := s.windowResetAnchor(*s.WeeklyWindowStart).Add(7 * 24 * time.Hour)
	return &t
}

func (s *UserSubscription) MonthlyResetTime() *time.Time {
	if s.MonthlyWindowStart == nil {
		return nil
	}
	t := s.windowResetAnchor(*s.MonthlyWindowStart).Add(30 * 24 * time.Hour)
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
