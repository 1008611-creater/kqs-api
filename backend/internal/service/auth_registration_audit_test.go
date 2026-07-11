//go:build unit

package service

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistrationAuditEmailDomainAndHash(t *testing.T) {
	require.Equal(t, "kqs.edu.cn", registrationAuditEmailDomain(" Student@KQS.EDU.CN "))
	require.Empty(t, registrationAuditEmailDomain("not-an-email"))

	hash1 := registrationAuditEmailHash(" Student@KQS.EDU.CN ")
	hash2 := registrationAuditEmailHash("student@kqs.edu.cn")
	require.Len(t, hash1, 64)
	require.Equal(t, hash1, hash2)
	require.Empty(t, registrationAuditEmailHash(" "))
}

func TestNormalizeRegistrationAuditInput(t *testing.T) {
	input := normalizeRegistrationAuditInput(RegistrationAuditInput{
		Action:        "unknown",
		Outcome:       "weird",
		FailureReason: strings.Repeat("x", registrationAuditMaxReasonLen+10),
		ClientIP:      strings.Repeat("1", registrationAuditMaxClientIPLen+10),
		UserAgent:     strings.Repeat("u", registrationAuditMaxUserAgentLen+10),
		RequestPath:   strings.Repeat("p", registrationAuditMaxPathLen+10),
	})

	require.Equal(t, RegistrationAuditActionRegister, input.Action)
	require.Equal(t, RegistrationAuditOutcomeFailed, input.Outcome)
	require.Len(t, input.FailureReason, registrationAuditMaxReasonLen)
	require.Len(t, input.ClientIP, registrationAuditMaxClientIPLen)
	require.Len(t, input.UserAgent, registrationAuditMaxUserAgentLen)
	require.Len(t, input.RequestPath, registrationAuditMaxPathLen)
}
