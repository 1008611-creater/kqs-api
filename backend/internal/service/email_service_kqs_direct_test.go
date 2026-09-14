//go:build unit

package service

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldUseKqsDirectVerifyDelivery(t *testing.T) {
	t.Run("uses fallback for KQS address when SMTP is missing", func(t *testing.T) {
		require.True(t, shouldUseKqsDirectVerifyDelivery("2023308250104@cau.edu.cn", ErrEmailNotConfigured))
	})

	t.Run("does not use fallback for non-KQS address", func(t *testing.T) {
		require.False(t, shouldUseKqsDirectVerifyDelivery("user@example.com", ErrEmailNotConfigured))
	})

	t.Run("does not use fallback for unrelated SMTP failure", func(t *testing.T) {
		require.False(t, shouldUseKqsDirectVerifyDelivery("student@cau.edu.cn", errors.New("smtp auth failed")))
	})
}

func TestShouldContinueVerifyCodeLegacyDelivery(t *testing.T) {
	t.Run("continues for template or config fallback", func(t *testing.T) {
		require.True(t, shouldContinueVerifyCodeLegacyDelivery("user@example.com", notificationEmailConfigErr(errors.New("template unavailable"))))
	})

	t.Run("continues for KQS direct fallback when notification delivery reports missing SMTP", func(t *testing.T) {
		require.True(t, shouldContinueVerifyCodeLegacyDelivery("2023308250104@cau.edu.cn", notificationEmailDeliveryErr(ErrEmailNotConfigured)))
	})

	t.Run("does not continue for non KQS notification delivery failure", func(t *testing.T) {
		require.False(t, shouldContinueVerifyCodeLegacyDelivery("user@example.com", notificationEmailDeliveryErr(ErrEmailNotConfigured)))
	})
}

func TestEmailService_SendKqsDirectVerifyEmailRejectsNonCau(t *testing.T) {
	svc := NewEmailService(nil, nil)

	err := svc.SendKqsDirectVerifyEmail(context.Background(), "user@example.com", "subject", "body", "site")

	require.ErrorIs(t, err, ErrEmailNotConfigured)
}

func TestDirectMXMessageHelpers(t *testing.T) {
	msg := directMXEmail{
		FromAddress: kqsDirectMailSender,
		FromName:    "矿泉水API",
		ToAddress:   "2023308250104@cau.edu.cn",
		Subject:     "[矿泉水API] 验证码",
		HTMLBody:    "<p>123456</p>",
		HELOName:    kqsDirectMailHELO,
	}

	raw := buildDirectMXMessage(msg, msg.ToAddress)

	require.Contains(t, raw, "Content-Transfer-Encoding: base64")
	require.Contains(t, raw, "Message-ID: <")
	require.Contains(t, raw, "Subject: =?UTF-8?")
	require.Contains(t, raw, "From: =?UTF-8?")
	require.Equal(t, "cau.edu.cn", emailDomain(msg.ToAddress))
}
