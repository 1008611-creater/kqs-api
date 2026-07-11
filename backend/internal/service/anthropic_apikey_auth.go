package service

import (
	"net/http"
	"strings"
)

const anthropicAPIKeyAuthHeaderModeBearer = "bearer"

// AnthropicAPIKeyUsesBearerAuth returns whether this Anthropic API-key account
// should authenticate with Authorization: Bearer instead of x-api-key.
func (a *Account) AnthropicAPIKeyUsesBearerAuth() bool {
	if a == nil || a.Platform != PlatformAnthropic || a.Type != AccountTypeAPIKey {
		return false
	}
	mode := strings.ToLower(strings.TrimSpace(a.getExtraString("anthropic_auth_header")))
	switch mode {
	case anthropicAPIKeyAuthHeaderModeBearer, "authorization", "authorization_bearer":
		return true
	default:
		return false
	}
}

func (a *Account) applyAnthropicAPIKeyAuthHeader(h http.Header, token string) {
	if h == nil {
		return
	}
	h.Del("authorization")
	h.Del("x-api-key")
	h.Del("x-goog-api-key")
	h.Del("cookie")
	if a.AnthropicAPIKeyUsesBearerAuth() {
		setHeaderRaw(h, "authorization", "Bearer "+token)
		return
	}
	setHeaderRaw(h, "x-api-key", token)
}
