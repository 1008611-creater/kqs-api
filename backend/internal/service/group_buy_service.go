package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math"
	"strings"
	"time"

	infraerrors "github.com/Wei-Shaw/sub2api/internal/pkg/errors"
)

const (
	GroupBuyRoomStatusActive    = "active"
	GroupBuyRoomStatusCompleted = "completed"
	GroupBuyRoomStatusExpired   = "expired"
)

var (
	ErrGroupBuyPlanNotFound         = infraerrors.NotFound("GROUP_BUY_PLAN_NOT_FOUND", "subscription plan not found")
	ErrGroupBuyRoomNotFound         = infraerrors.NotFound("GROUP_BUY_ROOM_NOT_FOUND", "group-buy room not found")
	ErrGroupBuyInvalidTarget        = infraerrors.BadRequest("GROUP_BUY_INVALID_TARGET", "target count must be 3, 5, or 10")
	ErrGroupBuyNoActiveSub          = infraerrors.Forbidden("GROUP_BUY_NO_ACTIVE_SUBSCRIPTION", "active subscription for this plan is required")
	ErrGroupBuyAlreadyJoined        = infraerrors.Conflict("GROUP_BUY_ALREADY_JOINED", "you have already joined this room")
	ErrGroupBuyRoomClosed           = infraerrors.Conflict("GROUP_BUY_ROOM_CLOSED", "group-buy room is no longer active")
	ErrGroupBuySubscriptionMismatch = infraerrors.BadRequest("GROUP_BUY_SUBSCRIPTION_MISMATCH", "subscription does not match room plan")
)

type GroupBuyService struct {
	db           *sql.DB
	subRepo      UserSubscriptionRepository
	billingCache *BillingCacheService
}

type GroupBuyPlanSnapshot struct {
	ID             int64   `json:"id"`
	GroupID        int64   `json:"group_id"`
	Name           string  `json:"name"`
	GroupName      string  `json:"group_name"`
	Price          float64 `json:"price"`
	ValidityDays   int     `json:"validity_days"`
	DailyLimitUSD  float64 `json:"daily_limit_usd"`
	WeeklyLimitUSD float64 `json:"weekly_limit_usd"`
}

type GroupBuyRoom struct {
	ID                     int64      `json:"id"`
	PlanID                 int64      `json:"plan_id"`
	PlanName               string     `json:"plan_name"`
	GroupID                int64      `json:"group_id"`
	GroupName              string     `json:"group_name"`
	OwnerUserID            int64      `json:"owner_user_id"`
	TargetCount            int        `json:"target_count"`
	MemberCount            int        `json:"member_count"`
	BonusMultiplier        float64    `json:"bonus_multiplier"`
	BaseDailyLimitUSD      float64    `json:"base_daily_limit_usd"`
	UpgradedDailyLimitUSD  float64    `json:"upgraded_daily_limit_usd"`
	BaseWeeklyLimitUSD     float64    `json:"base_weekly_limit_usd"`
	UpgradedWeeklyLimitUSD float64    `json:"upgraded_weekly_limit_usd"`
	Status                 string     `json:"status"`
	ExpiresAt              time.Time  `json:"expires_at"`
	CompletedAt            *time.Time `json:"completed_at,omitempty"`
	CreatedAt              time.Time  `json:"created_at"`
	Joined                 bool       `json:"joined"`
	CanJoin                bool       `json:"can_join"`
}

type GroupBuyHall struct {
	Rooms   []GroupBuyRoom `json:"rooms"`
	MyRooms []GroupBuyRoom `json:"my_rooms"`
}

type groupBuyPlan struct {
	ID             int64
	GroupID        int64
	Name           string
	GroupName      string
	Price          float64
	ValidityDays   int
	DailyLimitUSD  float64
	WeeklyLimitUSD float64
}

type groupBuyInvalidation struct {
	UserID  int64
	GroupID int64
}

func NewGroupBuyService(db *sql.DB, subRepo UserSubscriptionRepository, billingCache *BillingCacheService) *GroupBuyService {
	return &GroupBuyService{
		db:           db,
		subRepo:      subRepo,
		billingCache: billingCache,
	}
}

func (s *GroupBuyService) ListHall(ctx context.Context, userID int64) (*GroupBuyHall, error) {
	if err := s.expireStaleRooms(ctx); err != nil {
		return nil, err
	}
	rooms, err := s.queryRooms(ctx, userID, false)
	if err != nil {
		return nil, err
	}
	myRooms, err := s.queryRooms(ctx, userID, true)
	if err != nil {
		return nil, err
	}
	return &GroupBuyHall{Rooms: rooms, MyRooms: myRooms}, nil
}

func (s *GroupBuyService) CreateRoom(ctx context.Context, userID, planID int64, targetCount int) (*GroupBuyRoom, error) {
	multiplier, ok := groupBuyBonusMultiplier(targetCount)
	if !ok {
		return nil, ErrGroupBuyInvalidTarget
	}
	plan, err := s.getPlan(ctx, planID)
	if err != nil {
		return nil, err
	}
	sub, err := s.subRepo.GetActiveByUserIDAndGroupID(ctx, userID, plan.GroupID)
	if err != nil {
		return nil, ErrGroupBuyNoActiveSub
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer rollbackUnlessCommitted(tx)

	var roomID int64
	err = tx.QueryRowContext(ctx, `
		INSERT INTO group_buy_rooms (
			plan_id, group_id, owner_user_id, target_count, bonus_multiplier,
			status, expires_at, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, 'active', NOW() + INTERVAL '24 hours', NOW(), NOW())
		RETURNING id
	`, plan.ID, plan.GroupID, userID, targetCount, multiplier).Scan(&roomID)
	if err != nil {
		return nil, err
	}
	if err := s.insertRoomMember(ctx, tx, roomID, userID, sub, plan); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}

	return s.GetRoom(ctx, userID, roomID)
}

func (s *GroupBuyService) JoinRoom(ctx context.Context, userID, roomID int64) (*GroupBuyRoom, error) {
	if err := s.expireStaleRooms(ctx); err != nil {
		return nil, err
	}
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer rollbackUnlessCommitted(tx)

	plan, targetCount, status, expiresAt, err := s.getRoomPlanForUpdate(ctx, tx, roomID)
	if err != nil {
		return nil, err
	}
	if status != GroupBuyRoomStatusActive || !expiresAt.After(time.Now()) {
		return nil, ErrGroupBuyRoomClosed
	}
	sub, err := s.subRepo.GetActiveByUserIDAndGroupID(ctx, userID, plan.GroupID)
	if err != nil {
		return nil, ErrGroupBuyNoActiveSub
	}
	if sub.GroupID != plan.GroupID {
		return nil, ErrGroupBuySubscriptionMismatch
	}
	if err := s.insertRoomMember(ctx, tx, roomID, userID, sub, plan); err != nil {
		return nil, err
	}

	memberCount, err := countRoomMembers(ctx, tx, roomID)
	if err != nil {
		return nil, err
	}
	var invalidations []groupBuyInvalidation
	if memberCount >= targetCount {
		invalidations, err = s.completeRoom(ctx, tx, roomID, plan)
		if err != nil {
			return nil, err
		}
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	s.invalidateSubscriptions(ctx, invalidations)

	return s.GetRoom(ctx, userID, roomID)
}

func (s *GroupBuyService) GetRoom(ctx context.Context, userID, roomID int64) (*GroupBuyRoom, error) {
	rooms, err := s.queryRoomsByID(ctx, userID, roomID)
	if err != nil {
		return nil, err
	}
	if len(rooms) == 0 {
		return nil, ErrGroupBuyRoomNotFound
	}
	return &rooms[0], nil
}

func (s *GroupBuyService) getPlan(ctx context.Context, planID int64) (*groupBuyPlan, error) {
	var plan groupBuyPlan
	err := s.db.QueryRowContext(ctx, `
		SELECT sp.id, sp.group_id, sp.name, g.name, sp.price, sp.validity_days,
		       COALESCE(g.daily_limit_usd, 0), COALESCE(g.weekly_limit_usd, 0)
		FROM subscription_plans sp
		JOIN groups g ON g.id = sp.group_id AND g.deleted_at IS NULL
		WHERE sp.id = $1 AND sp.for_sale = true
	`, planID).Scan(&plan.ID, &plan.GroupID, &plan.Name, &plan.GroupName, &plan.Price, &plan.ValidityDays, &plan.DailyLimitUSD, &plan.WeeklyLimitUSD)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrGroupBuyPlanNotFound
	}
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (s *GroupBuyService) getRoomPlanForUpdate(ctx context.Context, tx *sql.Tx, roomID int64) (*groupBuyPlan, int, string, time.Time, error) {
	var plan groupBuyPlan
	var targetCount int
	var status string
	var expiresAt time.Time
	err := tx.QueryRowContext(ctx, `
		SELECT r.target_count, r.status, r.expires_at,
		       sp.id, sp.group_id, sp.name, g.name, sp.price, sp.validity_days,
		       COALESCE(g.daily_limit_usd, 0), COALESCE(g.weekly_limit_usd, 0)
		FROM group_buy_rooms r
		JOIN subscription_plans sp ON sp.id = r.plan_id
		JOIN groups g ON g.id = r.group_id AND g.deleted_at IS NULL
		WHERE r.id = $1 AND r.deleted_at IS NULL
		FOR UPDATE OF r
	`, roomID).Scan(&targetCount, &status, &expiresAt, &plan.ID, &plan.GroupID, &plan.Name, &plan.GroupName, &plan.Price, &plan.ValidityDays, &plan.DailyLimitUSD, &plan.WeeklyLimitUSD)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, 0, "", time.Time{}, ErrGroupBuyRoomNotFound
	}
	if err != nil {
		return nil, 0, "", time.Time{}, err
	}
	return &plan, targetCount, status, expiresAt, nil
}

func (s *GroupBuyService) insertRoomMember(ctx context.Context, tx *sql.Tx, roomID, userID int64, sub *UserSubscription, plan *groupBuyPlan) error {
	_, err := tx.ExecContext(ctx, `
		INSERT INTO group_buy_room_members (
			room_id, user_id, user_subscription_id, joined_at,
			daily_limit_before, weekly_limit_before
		)
		VALUES ($1, $2, $3, NOW(), $4, $5)
	`, roomID, userID, sub.ID, nullablePositiveLimit(sub.EffectiveDailyLimitUSD(&Group{DailyLimitUSD: positiveLimitPtr(plan.DailyLimitUSD)})), nullablePositiveLimit(sub.EffectiveWeeklyLimitUSD(&Group{WeeklyLimitUSD: positiveLimitPtr(plan.WeeklyLimitUSD)})))
	if err != nil {
		if isUniqueViolation(err) {
			return ErrGroupBuyAlreadyJoined
		}
		return err
	}
	return nil
}

func (s *GroupBuyService) completeRoom(ctx context.Context, tx *sql.Tx, roomID int64, plan *groupBuyPlan) ([]groupBuyInvalidation, error) {
	var multiplier float64
	err := tx.QueryRowContext(ctx, `
		UPDATE group_buy_rooms
		SET status = 'completed', completed_at = NOW(), updated_at = NOW()
		WHERE id = $1 AND status = 'active'
		RETURNING bonus_multiplier
	`, roomID).Scan(&multiplier)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	upgradedDaily := roundQuota(plan.DailyLimitUSD * multiplier)
	upgradedWeekly := roundQuota(plan.WeeklyLimitUSD * multiplier)
	rows, err := tx.QueryContext(ctx, `
		UPDATE user_subscriptions us
		SET daily_limit_override_usd = CASE
				WHEN $2::numeric > COALESCE(us.daily_limit_override_usd, 0) THEN $2::numeric
				ELSE us.daily_limit_override_usd
			END,
			weekly_limit_override_usd = CASE
				WHEN $3::numeric > COALESCE(us.weekly_limit_override_usd, 0) THEN $3::numeric
				ELSE us.weekly_limit_override_usd
			END,
			quota_bonus_multiplier = GREATEST(COALESCE(us.quota_bonus_multiplier, 1), $4::numeric),
			quota_bonus_source = 'group_buy:' || $1::text,
			updated_at = NOW()
		FROM group_buy_room_members m
		WHERE m.room_id = $1
		  AND m.user_subscription_id = us.id
		  AND us.deleted_at IS NULL
		RETURNING us.user_id, us.group_id, us.id, us.daily_limit_override_usd, us.weekly_limit_override_usd
	`, roomID, upgradedDaily, upgradedWeekly, multiplier)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	invalidations := make([]groupBuyInvalidation, 0)
	for rows.Next() {
		var userID, groupID, subID int64
		var dailyAfter, weeklyAfter sql.NullFloat64
		if err := rows.Scan(&userID, &groupID, &subID, &dailyAfter, &weeklyAfter); err != nil {
			return nil, err
		}
		invalidations = append(invalidations, groupBuyInvalidation{UserID: userID, GroupID: groupID})
		_, err := tx.ExecContext(ctx, `
			UPDATE group_buy_room_members
			SET daily_limit_after = $3, weekly_limit_after = $4, quota_applied_at = NOW()
			WHERE room_id = $1 AND user_subscription_id = $2
		`, roomID, subID, nullableFloat(dailyAfter), nullableFloat(weeklyAfter))
		if err != nil {
			return nil, err
		}
	}
	return invalidations, rows.Err()
}

func (s *GroupBuyService) queryRooms(ctx context.Context, userID int64, onlyMine bool) ([]GroupBuyRoom, error) {
	whereMine := ""
	if onlyMine {
		whereMine = "AND EXISTS (SELECT 1 FROM group_buy_room_members mine WHERE mine.room_id = r.id AND mine.user_id = $1)"
	}
	return s.scanRooms(ctx, `
		SELECT r.id, r.plan_id, sp.name, r.group_id, g.name, r.owner_user_id,
		       r.target_count, COALESCE(mc.member_count, 0), r.bonus_multiplier,
		       COALESCE(g.daily_limit_usd, 0), COALESCE(g.weekly_limit_usd, 0),
		       r.status, r.expires_at, r.completed_at, r.created_at,
		       EXISTS (SELECT 1 FROM group_buy_room_members m WHERE m.room_id = r.id AND m.user_id = $1) AS joined,
		       EXISTS (
		         SELECT 1 FROM user_subscriptions us
		         WHERE us.user_id = $1 AND us.group_id = r.group_id
		           AND us.status = 'active' AND us.expires_at > NOW() AND us.deleted_at IS NULL
		       ) AS can_join
		FROM group_buy_rooms r
		JOIN subscription_plans sp ON sp.id = r.plan_id
		JOIN groups g ON g.id = r.group_id AND g.deleted_at IS NULL
		LEFT JOIN (
		  SELECT room_id, COUNT(*)::int AS member_count
		  FROM group_buy_room_members
		  GROUP BY room_id
		) mc ON mc.room_id = r.id
		WHERE r.deleted_at IS NULL
		  AND (r.status = 'active' OR (r.status = 'completed' AND r.completed_at > NOW() - INTERVAL '72 hours'))
		`+whereMine+`
		ORDER BY CASE WHEN r.status = 'active' THEN 0 ELSE 1 END, r.created_at DESC
		LIMIT 80
	`, userID)
}

func (s *GroupBuyService) queryRoomsByID(ctx context.Context, userID, roomID int64) ([]GroupBuyRoom, error) {
	return s.scanRooms(ctx, `
		SELECT r.id, r.plan_id, sp.name, r.group_id, g.name, r.owner_user_id,
		       r.target_count, COALESCE(mc.member_count, 0), r.bonus_multiplier,
		       COALESCE(g.daily_limit_usd, 0), COALESCE(g.weekly_limit_usd, 0),
		       r.status, r.expires_at, r.completed_at, r.created_at,
		       EXISTS (SELECT 1 FROM group_buy_room_members m WHERE m.room_id = r.id AND m.user_id = $1) AS joined,
		       EXISTS (
		         SELECT 1 FROM user_subscriptions us
		         WHERE us.user_id = $1 AND us.group_id = r.group_id
		           AND us.status = 'active' AND us.expires_at > NOW() AND us.deleted_at IS NULL
		       ) AS can_join
		FROM group_buy_rooms r
		JOIN subscription_plans sp ON sp.id = r.plan_id
		JOIN groups g ON g.id = r.group_id AND g.deleted_at IS NULL
		LEFT JOIN (
		  SELECT room_id, COUNT(*)::int AS member_count
		  FROM group_buy_room_members
		  GROUP BY room_id
		) mc ON mc.room_id = r.id
		WHERE r.deleted_at IS NULL AND r.id = $2
	`, userID, roomID)
}

func (s *GroupBuyService) scanRooms(ctx context.Context, query string, args ...any) ([]GroupBuyRoom, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	rooms := make([]GroupBuyRoom, 0)
	for rows.Next() {
		var room GroupBuyRoom
		if err := rows.Scan(
			&room.ID,
			&room.PlanID,
			&room.PlanName,
			&room.GroupID,
			&room.GroupName,
			&room.OwnerUserID,
			&room.TargetCount,
			&room.MemberCount,
			&room.BonusMultiplier,
			&room.BaseDailyLimitUSD,
			&room.BaseWeeklyLimitUSD,
			&room.Status,
			&room.ExpiresAt,
			&room.CompletedAt,
			&room.CreatedAt,
			&room.Joined,
			&room.CanJoin,
		); err != nil {
			return nil, err
		}
		room.UpgradedDailyLimitUSD = roundQuota(room.BaseDailyLimitUSD * room.BonusMultiplier)
		room.UpgradedWeeklyLimitUSD = roundQuota(room.BaseWeeklyLimitUSD * room.BonusMultiplier)
		if room.Joined || room.Status != GroupBuyRoomStatusActive || !room.ExpiresAt.After(time.Now()) {
			room.CanJoin = false
		}
		rooms = append(rooms, room)
	}
	return rooms, rows.Err()
}

func (s *GroupBuyService) expireStaleRooms(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, `
		UPDATE group_buy_rooms
		SET status = 'expired', updated_at = NOW()
		WHERE status = 'active'
		  AND expires_at <= NOW()
		  AND deleted_at IS NULL
	`)
	return err
}

func (s *GroupBuyService) invalidateSubscriptions(ctx context.Context, items []groupBuyInvalidation) {
	if s.billingCache == nil || len(items) == 0 {
		return
	}
	cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	for _, item := range items {
		_ = s.billingCache.InvalidateSubscription(cacheCtx, item.UserID, item.GroupID)
	}
	_ = ctx
}

func groupBuyBonusMultiplier(targetCount int) (float64, bool) {
	switch targetCount {
	case 3:
		return 1.10, true
	case 5:
		return 1.17, true
	case 10:
		return 1.34, true
	default:
		return 0, false
	}
}

func countRoomMembers(ctx context.Context, tx *sql.Tx, roomID int64) (int, error) {
	var count int
	err := tx.QueryRowContext(ctx, `SELECT COUNT(*)::int FROM group_buy_room_members WHERE room_id = $1`, roomID).Scan(&count)
	return count, err
}

func rollbackUnlessCommitted(tx *sql.Tx) {
	_ = tx.Rollback()
}

func nullablePositiveLimit(v *float64) any {
	if v == nil || *v <= 0 {
		return nil
	}
	return *v
}

func positiveLimitPtr(v float64) *float64 {
	if v <= 0 {
		return nil
	}
	return &v
}

func nullableFloat(v sql.NullFloat64) any {
	if !v.Valid {
		return nil
	}
	return v.Float64
}

func roundQuota(v float64) float64 {
	return math.Round(v*100) / 100
}

func isUniqueViolation(err error) bool {
	if err == nil {
		return false
	}
	text := strings.ToLower(fmt.Sprint(err))
	return strings.Contains(text, "duplicate key") || strings.Contains(text, "unique constraint")
}
