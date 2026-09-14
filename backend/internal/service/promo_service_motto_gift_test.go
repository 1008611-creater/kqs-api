package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplyPromoCodeBlocksMottoGiftForNonCauEmail(t *testing.T) {
	repo := &kqsMottoGiftUserRepoStub{user: &User{ID: 7, Email: "user@example.com"}}
	svc := &PromoService{userRepo: repo}

	err := svc.ApplyPromoCode(context.Background(), 7, "JMSZDJYTXZYC")

	require.True(t, errors.Is(err, ErrKqsMottoGiftEmailRequired), "got %v", err)
}

func TestApplyPromoCodeBlocksMottoGiftAliasSpellings(t *testing.T) {
	for _, code := range []string{"jmszdjytxzyc", "jmszdj ytxzyc", " jmszdj-ytxzyc "} {
		repo := &kqsMottoGiftUserRepoStub{user: &User{ID: 7, Email: "user@example.com"}}
		svc := &PromoService{userRepo: repo}

		err := svc.ApplyPromoCode(context.Background(), 7, code)

		require.True(t, errors.Is(err, ErrKqsMottoGiftEmailRequired), "code %q got %v", code, err)
	}
}

func TestApplyPromoCodeMottoGiftRequiresUserLookup(t *testing.T) {
	svc := &PromoService{}

	err := svc.ApplyPromoCode(context.Background(), 7, "JMSZDJYTXZYC")

	require.True(t, errors.Is(err, ErrKqsMottoGiftEmailRequired), "got %v", err)
}
