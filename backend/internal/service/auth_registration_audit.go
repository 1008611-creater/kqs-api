package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"unicode/utf8"

	"github.com/Wei-Shaw/sub2api/internal/pkg/logger"
)

const (
	RegistrationAuditActionRegister       = "register"
	RegistrationAuditActionSendVerifyCode = "send_verify_code"

	RegistrationAuditOutcomeSuccess = "success"
	RegistrationAuditOutcomeFailed  = "failed"
)

const (
	registrationAuditMaxReasonLen    = 120
	registrationAuditMaxDomainLen    = 255
	registrationAuditMaxClientIPLen  = 80
	registrationAuditMaxUserAgentLen = 1024
	registrationAuditMaxPathLen      = 255
)

type RegistrationAuditInput struct {
	Action           string
	Outcome          string
	FailureReason    string
	UserID           int64
	Email            string
	ClientIP         string
	UserAgent        string
	RequestPath      string
	TurnstilePresent bool
}

func (s *AuthService) RecordRegistrationAudit(ctx context.Context, input RegistrationAuditInput) {
	if s == nil || s.entClient == nil {
		return
	}
	input = normalizeRegistrationAuditInput(input)

	var userID any
	if input.UserID > 0 {
		userID = input.UserID
	}

	_, err := s.entClient.ExecContext(ctx, `
INSERT INTO registration_audit_events (
	action,
	outcome,
	failure_reason,
	user_id,
	email_domain,
	email_sha256,
	client_ip,
	user_agent,
	request_path,
	turnstile_present
) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
`,
		input.Action,
		input.Outcome,
		input.FailureReason,
		userID,
		registrationAuditEmailDomain(input.Email),
		registrationAuditEmailHash(input.Email),
		input.ClientIP,
		input.UserAgent,
		input.RequestPath,
		input.TurnstilePresent,
	)
	if err != nil {
		logger.LegacyPrintf("service.auth", "[Auth] Failed to record registration audit event: %v", err)
	}
}

func normalizeRegistrationAuditInput(input RegistrationAuditInput) RegistrationAuditInput {
	switch strings.TrimSpace(input.Action) {
	case RegistrationAuditActionRegister, RegistrationAuditActionSendVerifyCode:
		input.Action = strings.TrimSpace(input.Action)
	default:
		input.Action = RegistrationAuditActionRegister
	}

	switch strings.TrimSpace(input.Outcome) {
	case RegistrationAuditOutcomeSuccess, RegistrationAuditOutcomeFailed:
		input.Outcome = strings.TrimSpace(input.Outcome)
	default:
		input.Outcome = RegistrationAuditOutcomeFailed
	}

	input.FailureReason = truncateAuditString(strings.TrimSpace(input.FailureReason), registrationAuditMaxReasonLen)
	input.ClientIP = truncateAuditString(strings.TrimSpace(input.ClientIP), registrationAuditMaxClientIPLen)
	input.UserAgent = truncateAuditString(strings.TrimSpace(input.UserAgent), registrationAuditMaxUserAgentLen)
	input.RequestPath = truncateAuditString(strings.TrimSpace(input.RequestPath), registrationAuditMaxPathLen)
	return input
}

func registrationAuditEmailDomain(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	if at := strings.LastIndex(email, "@"); at >= 0 && at+1 < len(email) {
		return truncateAuditString(email[at+1:], registrationAuditMaxDomainLen)
	}
	return ""
}

func registrationAuditEmailHash(email string) string {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return ""
	}
	sum := sha256.Sum256([]byte(email))
	return hex.EncodeToString(sum[:])
}

func truncateAuditString(value string, maxRunes int) string {
	if maxRunes <= 0 || value == "" || utf8.RuneCountInString(value) <= maxRunes {
		return value
	}
	out := make([]rune, 0, maxRunes)
	for _, r := range value {
		if len(out) >= maxRunes {
			break
		}
		out = append(out, r)
	}
	return string(out)
}
