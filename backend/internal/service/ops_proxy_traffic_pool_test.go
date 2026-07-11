package service

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

func TestBuildCurrentMihomoTrafficPoolItem(t *testing.T) {
	body := `
proxies:
  - { name: '39.61 GB | 150 GB', type: anytls }
  - { name: 'Traffic Reset: 28 Days Left', type: anytls }
  - { name: 'Expire Date: 2026-08-01', type: anytls }
`
	usage := parseSSRDOGSubscriptionBody(body)
	if usage == nil {
		t.Fatal("expected SSRDOG usage from mihomo config")
	}

	now := time.Date(2026, 6, 3, 12, 0, 0, 0, time.UTC)
	item := buildCurrentMihomoTrafficPoolItem(usage, now)

	if item.ID != "mihomo-active-config" {
		t.Fatalf("unexpected id: %s", item.ID)
	}
	if !item.Active {
		t.Fatal("expected current mihomo item to be active")
	}
	if item.Source != "mihomo-config" {
		t.Fatalf("unexpected source: %s", item.Source)
	}
	if item.Status != "ok" {
		t.Fatalf("unexpected status: %s", item.Status)
	}
	if item.UsedBytes == nil || item.TotalBytes == nil || item.RemainingBytes == nil || item.UsedPercent == nil {
		t.Fatal("expected usage fields to be populated")
	}
	if *item.UsedBytes <= 0 || *item.TotalBytes <= *item.UsedBytes {
		t.Fatalf("unexpected usage bytes: used=%v total=%v", *item.UsedBytes, *item.TotalBytes)
	}
	if item.ResetAfterDays == nil || *item.ResetAfterDays != 28 {
		t.Fatalf("unexpected reset days: %#v", item.ResetAfterDays)
	}
	if item.ExpiresAt == nil {
		t.Fatal("expected expiry")
	}
}

func TestProxyTrafficPoolMihomoStateConfiguredDoesNotImplyActive(t *testing.T) {
	subURL := "https://example.com/sub?id=abc#fragment"
	normalized := normalizeProxySubscriptionURL(subURL)
	state := proxyTrafficPoolMihomoState{
		subscriptions: map[string]struct{}{normalized: {}},
	}

	if !state.hasSubscription(" https://example.com/sub?id=abc ") {
		t.Fatal("expected subscription list match after normalization")
	}
	if state.referencesSubscription(subURL) {
		t.Fatal("subscription list membership alone must not mark a pool active")
	}

	state.content = "proxy-providers:\n  airport:\n    url: https://example.com/sub?id=abc#fragment\n"
	if !state.referencesSubscription(subURL) {
		t.Fatal("expected active match when current mihomo config references the subscription URL")
	}
}

func TestReadMihomoSubscriptionsNormalizesURLFragments(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "subscriptions-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	path := f.Name()
	if _, err := f.WriteString("\nhttps://example.com/sub?id=abc#secret-fragment\nnot a url\n"); err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	subs := readMihomoSubscriptions(path)
	if _, ok := subs["https://example.com/sub?id=abc"]; !ok {
		t.Fatalf("expected normalized URL without fragment, got %#v", subs)
	}
	if _, ok := subs["not a url"]; !ok {
		t.Fatal("expected non-URL lines to be preserved for exact matching")
	}
}

func TestProxyTrafficPoolSafeErrorMessageRedactsSubscriptionURL(t *testing.T) {
	err := fmt.Errorf("subscription returned HTTP 403 for https://example.com/sub?token=secret")
	msg := proxyTrafficPoolSafeErrorMessage(err, true)

	if !strings.Contains(msg, "HTTP 403") {
		t.Fatalf("expected HTTP 403 context, got %q", msg)
	}
	if strings.Contains(msg, "example.com") || strings.Contains(msg, "secret") || strings.Contains(msg, "token=") {
		t.Fatalf("message leaked subscription detail: %q", msg)
	}
	if !strings.Contains(msg, "Mihomo") {
		t.Fatalf("expected Mihomo fallback hint, got %q", msg)
	}
}
