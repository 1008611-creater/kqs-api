package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAccount_IsAnthropicAPIKeyPassthroughEnabled(t *testing.T) {
	t.Run("Anthropic API Key 开启", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"anthropic_passthrough": true,
			},
		}
		require.True(t, account.IsAnthropicAPIKeyPassthroughEnabled())
	})

	t.Run("Anthropic API Key 关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"anthropic_passthrough": false,
			},
		}
		require.False(t, account.IsAnthropicAPIKeyPassthroughEnabled())
	})

	t.Run("字段类型非法默认关闭", func(t *testing.T) {
		account := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"anthropic_passthrough": "true",
			},
		}
		require.False(t, account.IsAnthropicAPIKeyPassthroughEnabled())
	})

	t.Run("非 Anthropic API Key 账号始终关闭", func(t *testing.T) {
		oauth := &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeOAuth,
			Extra: map[string]any{
				"anthropic_passthrough": true,
			},
		}
		require.False(t, oauth.IsAnthropicAPIKeyPassthroughEnabled())

		openai := &Account{
			Platform: PlatformOpenAI,
			Type:     AccountTypeAPIKey,
			Extra: map[string]any{
				"anthropic_passthrough": true,
			},
		}
		require.False(t, openai.IsAnthropicAPIKeyPassthroughEnabled())
	})
}

<<<<<<< HEAD
func TestAccount_GetAnthropicAPIKeyAuthScheme(t *testing.T) {
	tests := []struct {
		name    string
		account *Account
		want    string
	}{
		{
			name: "missing extra defaults to x-api-key",
			account: &Account{
				Platform: PlatformAnthropic,
				Type:     AccountTypeAPIKey,
			},
			want: AnthropicAPIKeyAuthSchemeXAPIKey,
		},
		{
			name: "explicit bearer",
			account: &Account{
				Platform: PlatformAnthropic,
				Type:     AccountTypeAPIKey,
				Extra: map[string]any{
					"anthropic_apikey_auth_scheme": AnthropicAPIKeyAuthSchemeAuthorizationBearer,
				},
			},
			want: AnthropicAPIKeyAuthSchemeAuthorizationBearer,
		},
		{
			name: "invalid value defaults to x-api-key",
			account: &Account{
				Platform: PlatformAnthropic,
				Type:     AccountTypeAPIKey,
				Extra: map[string]any{
					"anthropic_apikey_auth_scheme": "bearer",
				},
			},
			want: AnthropicAPIKeyAuthSchemeXAPIKey,
		},
		{
			name: "non Anthropic API key defaults to x-api-key",
			account: &Account{
				Platform: PlatformOpenAI,
				Type:     AccountTypeAPIKey,
				Extra: map[string]any{
					"anthropic_apikey_auth_scheme": AnthropicAPIKeyAuthSchemeAuthorizationBearer,
				},
			},
			want: AnthropicAPIKeyAuthSchemeXAPIKey,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.account.GetAnthropicAPIKeyAuthScheme())
		})
	}
=======
func TestAccount_AnthropicAPIKeyUsesBearerAuth(t *testing.T) {
	base := func(extra map[string]any) *Account {
		return &Account{
			Platform: PlatformAnthropic,
			Type:     AccountTypeAPIKey,
			Extra:    extra,
		}
	}

	require.False(t, base(nil).AnthropicAPIKeyUsesBearerAuth())
	require.False(t, base(map[string]any{"anthropic_auth_header": "x-api-key"}).AnthropicAPIKeyUsesBearerAuth())
	require.True(t, base(map[string]any{"anthropic_auth_header": "bearer"}).AnthropicAPIKeyUsesBearerAuth())
	require.True(t, base(map[string]any{"anthropic_auth_header": "Authorization"}).AnthropicAPIKeyUsesBearerAuth())

	require.False(t, (&Account{
		Platform: PlatformOpenAI,
		Type:     AccountTypeAPIKey,
		Extra:    map[string]any{"anthropic_auth_header": "bearer"},
	}).AnthropicAPIKeyUsesBearerAuth())
>>>>>>> 1157e4cfe6271a81c9dc96df84af6c6a3b22f831
}
