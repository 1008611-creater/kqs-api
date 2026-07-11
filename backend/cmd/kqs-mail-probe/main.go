package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/Wei-Shaw/sub2api/internal/service"
	"github.com/redis/go-redis/v9"
)

type emptySettingRepo struct{}

func (emptySettingRepo) Get(context.Context, string) (*service.Setting, error) {
	return nil, service.ErrSettingNotFound
}

func (emptySettingRepo) GetValue(context.Context, string) (string, error) {
	return "", service.ErrSettingNotFound
}

func (emptySettingRepo) Set(context.Context, string, string) error {
	return nil
}

func (emptySettingRepo) GetMultiple(_ context.Context, keys []string) (map[string]string, error) {
	out := make(map[string]string, len(keys))
	for _, key := range keys {
		out[key] = ""
	}
	return out, nil
}

func (emptySettingRepo) SetMultiple(context.Context, map[string]string) error {
	return nil
}

func (emptySettingRepo) GetAll(context.Context) (map[string]string, error) {
	return map[string]string{}, nil
}

func (emptySettingRepo) Delete(context.Context, string) error {
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: kqs-mail-probe email")
		os.Exit(2)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	rdb := redis.NewClient(&redis.Options{
		Addr: "sub2api-gg-redis:6379",
	})
	defer func() { _ = rdb.Close() }()

	cache := probeEmailCache{rdb: rdb}
	emailService := service.NewEmailService(emptySettingRepo{}, cache)
	if err := emailService.SendVerifyCode(ctx, os.Args[1], "矿泉水API", "zh-CN"); err != nil {
		fmt.Fprintf(os.Stderr, "send failed: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("sent")
}

type probeEmailCache struct {
	rdb *redis.Client
}

func (c probeEmailCache) GetVerificationCode(ctx context.Context, email string) (*service.VerificationCodeData, error) {
	return nil, redis.Nil
}

func (c probeEmailCache) SetVerificationCode(ctx context.Context, email string, data *service.VerificationCodeData, ttl time.Duration) error {
	key := "verify_code:" + email
	value := fmt.Sprintf(`{"Code":"%s","Attempts":%d,"CreatedAt":"%s","ExpiresAt":"%s"}`,
		data.Code,
		data.Attempts,
		data.CreatedAt.Format(time.RFC3339Nano),
		data.ExpiresAt.Format(time.RFC3339Nano),
	)
	return c.rdb.Set(ctx, key, value, ttl).Err()
}

func (c probeEmailCache) DeleteVerificationCode(ctx context.Context, email string) error {
	return c.rdb.Del(ctx, "verify_code:"+email).Err()
}

func (c probeEmailCache) GetNotifyVerifyCode(context.Context, string) (*service.VerificationCodeData, error) {
	return nil, redis.Nil
}

func (c probeEmailCache) SetNotifyVerifyCode(context.Context, string, *service.VerificationCodeData, time.Duration) error {
	return nil
}

func (c probeEmailCache) DeleteNotifyVerifyCode(context.Context, string) error {
	return nil
}

func (c probeEmailCache) GetPasswordResetToken(context.Context, string) (*service.PasswordResetTokenData, error) {
	return nil, redis.Nil
}

func (c probeEmailCache) SetPasswordResetToken(context.Context, string, *service.PasswordResetTokenData, time.Duration) error {
	return nil
}

func (c probeEmailCache) DeletePasswordResetToken(context.Context, string) error {
	return nil
}

func (c probeEmailCache) IsPasswordResetEmailInCooldown(context.Context, string) bool {
	return false
}

func (c probeEmailCache) SetPasswordResetEmailCooldown(context.Context, string, time.Duration) error {
	return nil
}

func (c probeEmailCache) IncrNotifyCodeUserRate(context.Context, int64, time.Duration) (int64, error) {
	return 0, nil
}

func (c probeEmailCache) GetNotifyCodeUserRate(context.Context, int64) (int64, error) {
	return 0, redis.Nil
}
