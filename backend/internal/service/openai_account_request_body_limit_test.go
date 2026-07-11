package service

import "testing"

func TestAccountAllowsOpenAIRequestBodySize(t *testing.T) {
	account := &Account{Extra: map[string]any{"max_request_body_mb": 10}}

	if !accountAllowsOpenAIRequestBodySize(account, 10*1024*1024) {
		t.Fatal("expected request at the account body limit to be allowed")
	}
	if accountAllowsOpenAIRequestBodySize(account, 10*1024*1024+1) {
		t.Fatal("expected request above the account body limit to be skipped")
	}
	if !accountAllowsOpenAIRequestBodySize(&Account{Extra: map[string]any{}}, 256*1024*1024) {
		t.Fatal("expected accounts without a body limit to be allowed")
	}
	if !accountAllowsOpenAIRequestBodySize(account, 0) {
		t.Fatal("expected unknown body size to be allowed")
	}
}
