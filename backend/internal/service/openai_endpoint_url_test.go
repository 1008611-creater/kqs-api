package service

import "testing"

func TestBuildOpenAIEndpointURL_ExplicitNoV1Marker(t *testing.T) {
	base := appendOpenAIBaseURLNoV1Marker("https://upstream.example/coding")

	if got, want := buildOpenAIResponsesURL(base), "https://upstream.example/coding/responses"; got != want {
		t.Fatalf("responses URL = %q, want %q", got, want)
	}
	if got, want := buildOpenAIChatCompletionsURL(base), "https://upstream.example/coding/chat/completions"; got != want {
		t.Fatalf("chat completions URL = %q, want %q", got, want)
	}
}

func TestBuildOpenAIEndpointURL_DefaultStillAddsV1(t *testing.T) {
	if got, want := buildOpenAIResponsesURL("https://upstream.example/coding"), "https://upstream.example/coding/v1/responses"; got != want {
		t.Fatalf("responses URL = %q, want %q", got, want)
	}
}
