package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/pkg/pagination"
	"github.com/stretchr/testify/require"
)

type kqsMottoGiftUserRepoStub struct {
	user                *User
	updateBalanceCalled bool
}

func (s *kqsMottoGiftUserRepoStub) Create(context.Context, *User) error { return nil }
func (s *kqsMottoGiftUserRepoStub) GetByID(context.Context, int64) (*User, error) {
	if s.user == nil {
		return nil, ErrUserNotFound
	}
	user := *s.user
	return &user, nil
}
func (s *kqsMottoGiftUserRepoStub) GetByEmail(context.Context, string) (*User, error) {
	return nil, ErrUserNotFound
}
func (s *kqsMottoGiftUserRepoStub) GetFirstAdmin(context.Context) (*User, error) {
	return nil, ErrUserNotFound
}
func (s *kqsMottoGiftUserRepoStub) Update(context.Context, *User) error { return nil }
func (s *kqsMottoGiftUserRepoStub) Delete(context.Context, int64) error { return nil }
func (s *kqsMottoGiftUserRepoStub) GetUserAvatar(context.Context, int64) (*UserAvatar, error) {
	return nil, nil
}
func (s *kqsMottoGiftUserRepoStub) UpsertUserAvatar(context.Context, int64, UpsertUserAvatarInput) (*UserAvatar, error) {
	return nil, nil
}
func (s *kqsMottoGiftUserRepoStub) DeleteUserAvatar(context.Context, int64) error { return nil }
func (s *kqsMottoGiftUserRepoStub) List(context.Context, pagination.PaginationParams) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (s *kqsMottoGiftUserRepoStub) ListWithFilters(context.Context, pagination.PaginationParams, UserListFilters) ([]User, *pagination.PaginationResult, error) {
	return nil, nil, nil
}
func (s *kqsMottoGiftUserRepoStub) GetLatestUsedAtByUserIDs(context.Context, []int64) (map[int64]*time.Time, error) {
	return map[int64]*time.Time{}, nil
}
func (s *kqsMottoGiftUserRepoStub) GetLatestUsedAtByUserID(context.Context, int64) (*time.Time, error) {
	return nil, nil
}
func (s *kqsMottoGiftUserRepoStub) UpdateUserLastActiveAt(context.Context, int64, time.Time) error {
	return nil
}
func (s *kqsMottoGiftUserRepoStub) UpdateBalance(context.Context, int64, float64) error {
	s.updateBalanceCalled = true
	return nil
}
func (s *kqsMottoGiftUserRepoStub) DeductBalance(context.Context, int64, float64) error { return nil }
func (s *kqsMottoGiftUserRepoStub) UpdateConcurrency(context.Context, int64, int) error { return nil }
func (s *kqsMottoGiftUserRepoStub) BatchSetConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
}
func (s *kqsMottoGiftUserRepoStub) BatchAddConcurrency(context.Context, []int64, int) (int, error) {
	return 0, nil
}
func (s *kqsMottoGiftUserRepoStub) ExistsByEmail(context.Context, string) (bool, error) {
	return false, nil
}
func (s *kqsMottoGiftUserRepoStub) RemoveGroupFromAllowedGroups(context.Context, int64) (int64, error) {
	return 0, nil
}
func (s *kqsMottoGiftUserRepoStub) AddGroupToAllowedGroups(context.Context, int64, int64) error {
	return nil
}
func (s *kqsMottoGiftUserRepoStub) RemoveGroupFromUserAllowedGroups(context.Context, int64, int64) error {
	return nil
}
func (s *kqsMottoGiftUserRepoStub) ListUserAuthIdentities(context.Context, int64) ([]UserAuthIdentityRecord, error) {
	return nil, nil
}
func (s *kqsMottoGiftUserRepoStub) UnbindUserAuthProvider(context.Context, int64, string) error {
	return nil
}
func (s *kqsMottoGiftUserRepoStub) UpdateTotpSecret(context.Context, int64, *string) error {
	return nil
}
func (s *kqsMottoGiftUserRepoStub) EnableTotp(context.Context, int64) error  { return nil }
func (s *kqsMottoGiftUserRepoStub) DisableTotp(context.Context, int64) error { return nil }

func TestKqsMottoGiftCodeNormalization(t *testing.T) {
	require.True(t, isKqsMottoGiftCode("JMSZDJYTXZYC"))
	require.True(t, isKqsMottoGiftCode("jmszdj ytxzyc"))
	require.True(t, isKqsMottoGiftCode("jmszdj-ytxzyc"))
	require.False(t, isKqsMottoGiftCode("JMSZDJ"))
}

func TestKqsMottoGiftIPHash(t *testing.T) {
	hash1 := hashKqsMottoGiftIP(" 203.0.113.42 ")
	hash2 := hashKqsMottoGiftIP("203.0.113.42")

	require.Len(t, hash1, 64)
	require.Equal(t, hash1, hash2)
	require.Empty(t, hashKqsMottoGiftIP(" "))
	require.NotContains(t, hash1, "203.0.113.42")
}

func TestKqsMottoGiftEligibleEmail(t *testing.T) {
	require.True(t, isKqsMottoGiftEligibleEmail("student@cau.edu.cn"))
	require.True(t, isKqsMottoGiftEligibleEmail(" STUDENT@CAU.EDU.CN "))

	require.False(t, isKqsMottoGiftEligibleEmail("student@gmail.com"))
	require.False(t, isKqsMottoGiftEligibleEmail("student@foo.cau.edu.cn"))
	require.False(t, isKqsMottoGiftEligibleEmail("studentcau.edu.cn"))
	require.False(t, isKqsMottoGiftEligibleEmail(""))
}

func TestKqsMottoGiftRequiresCauEmailBeforeClaim(t *testing.T) {
	repo := &kqsMottoGiftUserRepoStub{user: &User{ID: 42, Email: "student@gmail.com"}}
	svc := &RedeemService{userRepo: repo}

	_, err := svc.redeemKqsMottoGift(context.Background(), 42, "203.0.113.42")

	require.True(t, errors.Is(err, ErrKqsMottoGiftEmailRequired), "got %v", err)
	require.False(t, repo.updateBalanceCalled)
}
