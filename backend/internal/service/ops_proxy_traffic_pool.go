package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/imroc/req/v3"
)

const (
	defaultProxyTrafficPoolTimeout = 12 * time.Second
	defaultMihomoConfigPath        = "/etc/mihomo/config-sub2api.yaml"
	defaultMihomoStatusPath        = "/etc/mihomo/sub2api-traffic-status.txt"
	defaultMihomoSubscriptionsPath = "/etc/mihomo/subscriptions.txt"
	defaultMihomoProxyURL          = "http://127.0.0.1:7890"
)

type OpsProxyTrafficPoolConfigItem struct {
	ID               string `json:"id"`
	Name             string `json:"name"`
	SubscriptionURL  string `json:"subscription_url"`
	Provider         string `json:"provider,omitempty"`
	MihomoConfigPath string `json:"mihomo_config_path,omitempty"`
	Enabled          bool   `json:"enabled,omitempty"`
}

type OpsProxyTrafficPoolItem struct {
	ID                 string     `json:"id"`
	Name               string     `json:"name"`
	Provider           string     `json:"provider,omitempty"`
	Enabled            bool       `json:"enabled"`
	Active             bool       `json:"active"`
	ConfiguredInMihomo bool       `json:"configured_in_mihomo"`
	UsedBytes          *int64     `json:"used_bytes,omitempty"`
	TotalBytes         *int64     `json:"total_bytes,omitempty"`
	RemainingBytes     *int64     `json:"remaining_bytes,omitempty"`
	UsedPercent        *float64   `json:"used_percent,omitempty"`
	ExpiresAt          *time.Time `json:"expires_at,omitempty"`
	ResetAfterDays     *int       `json:"reset_after_days,omitempty"`
	Status             string     `json:"status"`
	Message            string     `json:"message,omitempty"`
	LastCheckedAt      time.Time  `json:"last_checked_at"`
	Source             string     `json:"source"`
	SubscriptionKey    string     `json:"subscription_key"`
}

type OpsProxyTrafficPoolResponse struct {
	Items        []*OpsProxyTrafficPoolItem `json:"items"`
	TotalBytes   int64                      `json:"total_bytes"`
	UsedBytes    int64                      `json:"used_bytes"`
	ActiveItemID string                     `json:"active_item_id,omitempty"`
	CheckedAt    time.Time                  `json:"checked_at"`
}

func (s *OpsService) GetProxyTrafficPools(ctx context.Context) (*OpsProxyTrafficPoolResponse, error) {
	if err := s.RequireMonitoringEnabled(ctx); err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	items := s.loadProxyTrafficPoolConfig(ctx)
	mihomoState := readMihomoTrafficPoolState(items)

	resp := &OpsProxyTrafficPoolResponse{Items: []*OpsProxyTrafficPoolItem{}, CheckedAt: now}
	hasMeasuredItem := false
	for _, cfg := range items {
		item := s.inspectProxyTrafficPool(ctx, cfg, mihomoState, now)
		resp.Items = append(resp.Items, item)
		if item.UsedBytes != nil && item.TotalBytes != nil {
			hasMeasuredItem = true
		}
		if item.TotalBytes != nil {
			resp.TotalBytes += *item.TotalBytes
		}
		if item.UsedBytes != nil {
			resp.UsedBytes += *item.UsedBytes
		}
		if item.Active {
			resp.ActiveItemID = item.ID
		}
	}
	if !hasMeasuredItem && mihomoState.usage != nil {
		item := buildCurrentMihomoTrafficPoolItem(mihomoState.usage, now)
		resp.Items = append(resp.Items, item)
		resp.ActiveItemID = item.ID
		if item.TotalBytes != nil {
			resp.TotalBytes += *item.TotalBytes
		}
		if item.UsedBytes != nil {
			resp.UsedBytes += *item.UsedBytes
		}
	}
	return resp, nil
}

func (s *OpsService) loadProxyTrafficPoolConfig(ctx context.Context) []*OpsProxyTrafficPoolConfigItem {
	if s == nil || s.settingRepo == nil {
		return nil
	}
	raw, err := s.settingRepo.GetValue(ctx, SettingKeyOpsProxyTrafficPools)
	if err != nil || strings.TrimSpace(raw) == "" {
		return nil
	}
	var items []*OpsProxyTrafficPoolConfigItem
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return nil
	}
	out := make([]*OpsProxyTrafficPoolConfigItem, 0, len(items))
	for _, item := range items {
		if item == nil || strings.TrimSpace(item.SubscriptionURL) == "" {
			continue
		}
		cp := *item
		cp.ID = strings.TrimSpace(cp.ID)
		cp.Name = strings.TrimSpace(cp.Name)
		cp.Provider = strings.TrimSpace(cp.Provider)
		cp.SubscriptionURL = strings.TrimSpace(cp.SubscriptionURL)
		cp.MihomoConfigPath = strings.TrimSpace(cp.MihomoConfigPath)
		if cp.ID == "" {
			cp.ID = "pool_" + shortURLHash(cp.SubscriptionURL)
		}
		if cp.Name == "" {
			cp.Name = cp.ID
		}
		if cp.MihomoConfigPath == "" {
			cp.MihomoConfigPath = defaultMihomoConfigPath
		}
		out = append(out, &cp)
	}
	return out
}

func (s *OpsService) inspectProxyTrafficPool(ctx context.Context, cfg *OpsProxyTrafficPoolConfigItem, mihomoState proxyTrafficPoolMihomoState, now time.Time) *OpsProxyTrafficPoolItem {
	registeredInMihomo := mihomoState.hasSubscription(cfg.SubscriptionURL)
	item := &OpsProxyTrafficPoolItem{
		ID:                 cfg.ID,
		Name:               cfg.Name,
		Provider:           cfg.Provider,
		Enabled:            cfg.Enabled,
		Active:             mihomoState.referencesSubscription(cfg.SubscriptionURL),
		ConfiguredInMihomo: registeredInMihomo,
		Status:             "unknown",
		LastCheckedAt:      now,
		Source:             "subscription",
		SubscriptionKey:    shortURLHash(cfg.SubscriptionURL),
	}
	body, headers, err := fetchProxySubscription(ctx, cfg.SubscriptionURL, pickProxyTrafficPoolProxyURL(s))
	if err != nil {
		item.Status = "error"
		item.Message = proxyTrafficPoolSafeErrorMessage(err, registeredInMihomo)
		return item
	}
	if usage := parseSubscriptionUserInfo(headers.Get("subscription-userinfo")); usage != nil {
		applyProxyTrafficUsage(item, usage)
		item.Source = "subscription-userinfo"
	}
	if usage := parseSSRDOGSubscriptionBody(body); usage != nil {
		applyProxyTrafficUsage(item, usage)
		item.Source = "subscription-body"
	}
	if item.UsedBytes == nil && item.TotalBytes == nil {
		item.Message = proxyTrafficPoolNoUsageMessage(registeredInMihomo)
	}
	classifyProxyTrafficPoolStatus(item)
	return item
}

type proxyTrafficUsage struct {
	UsedBytes      *int64
	TotalBytes     *int64
	ExpiresAt      *time.Time
	ResetAfterDays *int
}

func applyProxyTrafficUsage(item *OpsProxyTrafficPoolItem, usage *proxyTrafficUsage) {
	if usage.UsedBytes != nil {
		item.UsedBytes = usage.UsedBytes
	}
	if usage.TotalBytes != nil {
		item.TotalBytes = usage.TotalBytes
	}
	if usage.ExpiresAt != nil {
		item.ExpiresAt = usage.ExpiresAt
	}
	if usage.ResetAfterDays != nil {
		item.ResetAfterDays = usage.ResetAfterDays
	}
	if item.UsedBytes != nil && item.TotalBytes != nil && *item.TotalBytes > 0 {
		remaining := *item.TotalBytes - *item.UsedBytes
		if remaining < 0 {
			remaining = 0
		}
		percent := float64(*item.UsedBytes) / float64(*item.TotalBytes) * 100
		item.RemainingBytes = &remaining
		item.UsedPercent = &percent
	}
}

func classifyProxyTrafficPoolStatus(item *OpsProxyTrafficPoolItem) {
	item.Status = "unknown"
	if item.UsedBytes != nil && item.TotalBytes != nil && *item.TotalBytes > 0 {
		switch {
		case *item.UsedBytes >= *item.TotalBytes:
			item.Status = "exhausted"
		case item.UsedPercent != nil && *item.UsedPercent >= 90:
			item.Status = "warning"
		default:
			item.Status = "ok"
		}
	}
	if item.ExpiresAt != nil && time.Now().UTC().After(*item.ExpiresAt) {
		item.Status = "expired"
	}
}

func buildCurrentMihomoTrafficPoolItem(usage *proxyTrafficUsage, now time.Time) *OpsProxyTrafficPoolItem {
	item := &OpsProxyTrafficPoolItem{
		ID:                 "mihomo-active-config",
		Name:               "当前正在使用的机场配置",
		Provider:           "Mihomo",
		Enabled:            true,
		Active:             true,
		ConfiguredInMihomo: true,
		Status:             "unknown",
		LastCheckedAt:      now,
		Source:             "mihomo-config",
		SubscriptionKey:    "local-config",
		Message:            "订阅站点未返回可读流量元数据时，使用服务器本机 Mihomo 配置/状态文件里的当前流量。具体订阅链接已隐藏。",
	}
	applyProxyTrafficUsage(item, usage)
	classifyProxyTrafficPoolStatus(item)
	return item
}

func fetchProxySubscription(ctx context.Context, subURL, proxyURL string) (string, http.Header, error) {
	client := req.C().SetTimeout(defaultProxyTrafficPoolTimeout)
	if strings.TrimSpace(proxyURL) != "" {
		client.SetProxyURL(proxyURL)
	}
	resp, err := client.R().SetContext(ctx).Get(subURL)
	if err != nil && strings.TrimSpace(proxyURL) != "" {
		resp, err = req.C().SetTimeout(defaultProxyTrafficPoolTimeout).R().SetContext(ctx).Get(subURL)
	}
	if err != nil {
		return "", nil, fmt.Errorf("fetch subscription: %w", err)
	}
	if resp.GetStatusCode() < 200 || resp.GetStatusCode() >= 300 {
		return "", resp.Response.Header, fmt.Errorf("subscription returned HTTP %d", resp.GetStatusCode())
	}
	return resp.String(), resp.Response.Header, nil
}

func proxyTrafficPoolSafeErrorMessage(err error, configuredInMihomo bool) string {
	msg := strings.ToLower(strings.TrimSpace(err.Error()))
	switch {
	case strings.Contains(msg, "http 403"):
		if configuredInMihomo {
			return "订阅站点拒绝直接读取流量信息（HTTP 403）；该订阅已写入服务器 Mihomo，当前用量请看“当前正在使用的机场配置”。"
		}
		return "订阅站点拒绝直接读取流量信息（HTTP 403）；未在服务器 Mihomo 订阅清单中确认到该池。"
	case strings.Contains(msg, "http 401"):
		return "订阅站点需要认证，无法直接读取流量信息。"
	case strings.Contains(msg, "timeout") || strings.Contains(msg, "deadline"):
		return "读取订阅流量信息超时；服务器会继续展示本机 Mihomo 可用的状态。"
	}
	if configuredInMihomo {
		return "无法直接读取订阅流量信息；该订阅已写入服务器 Mihomo，当前用量请看本机 Mihomo 状态。"
	}
	return "无法直接读取订阅流量信息；请检查订阅地址是否可用。"
}

func proxyTrafficPoolNoUsageMessage(configuredInMihomo bool) string {
	if configuredInMihomo {
		return "订阅返回成功，但没有可解析的流量字段；该订阅已写入服务器 Mihomo。"
	}
	return "订阅返回成功，但没有可解析的流量字段。"
}

func parseSubscriptionUserInfo(raw string) *proxyTrafficUsage {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	values := map[string]int64{}
	for _, part := range strings.Split(raw, ";") {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		n, err := strconv.ParseInt(strings.TrimSpace(kv[1]), 10, 64)
		if err != nil {
			continue
		}
		values[strings.ToLower(strings.TrimSpace(kv[0]))] = n
	}
	upload, hasUp := values["upload"]
	download, hasDown := values["download"]
	total, hasTotal := values["total"]
	if !hasTotal {
		return nil
	}
	used := int64(0)
	if hasUp {
		used += upload
	}
	if hasDown {
		used += download
	}
	out := &proxyTrafficUsage{UsedBytes: &used, TotalBytes: &total}
	if expire, ok := values["expire"]; ok && expire > 0 {
		t := time.Unix(expire, 0).UTC()
		out.ExpiresAt = &t
	}
	return out
}

var (
	ssrdogTrafficRe = regexp.MustCompile(`(?i)(\d+(?:\.\d+)?)\s*([KMGT]?B)\s*\|\s*(\d+(?:\.\d+)?)\s*([KMGT]?B)`)
	ssrdogExpireRe  = regexp.MustCompile(`(?i)Expire Date:\s*(\d{4}-\d{2}-\d{2})`)
	ssrdogResetRe   = regexp.MustCompile(`(?i)Traffic Reset:\s*(\d+)\s*Days?\s*Left`)
)

func parseSSRDOGSubscriptionBody(body string) *proxyTrafficUsage {
	body = strings.TrimSpace(body)
	if body == "" {
		return nil
	}
	out := &proxyTrafficUsage{}
	if m := ssrdogTrafficRe.FindStringSubmatch(body); len(m) == 5 {
		used := parseTrafficBytes(m[1], m[2])
		total := parseTrafficBytes(m[3], m[4])
		if total > 0 {
			out.UsedBytes = &used
			out.TotalBytes = &total
		}
	}
	if m := ssrdogExpireRe.FindStringSubmatch(body); len(m) == 2 {
		if t, err := time.Parse("2006-01-02", m[1]); err == nil {
			t = t.Add(24*time.Hour - time.Second).UTC()
			out.ExpiresAt = &t
		}
	}
	if m := ssrdogResetRe.FindStringSubmatch(body); len(m) == 2 {
		if days, err := strconv.Atoi(m[1]); err == nil {
			out.ResetAfterDays = &days
		}
	}
	if out.UsedBytes == nil && out.TotalBytes == nil && out.ExpiresAt == nil && out.ResetAfterDays == nil {
		return nil
	}
	return out
}

func parseTrafficBytes(numRaw, unitRaw string) int64 {
	v, err := strconv.ParseFloat(strings.TrimSpace(numRaw), 64)
	if err != nil {
		return 0
	}
	switch strings.ToUpper(strings.TrimSpace(unitRaw)) {
	case "TB":
		v *= 1024
		fallthrough
	case "GB":
		v *= 1024
		fallthrough
	case "MB":
		v *= 1024
		fallthrough
	case "KB":
		v *= 1024
	}
	return int64(v)
}

func pickProxyTrafficPoolProxyURL(s *OpsService) string {
	if s != nil && s.cfg != nil && strings.TrimSpace(s.cfg.Update.ProxyURL) != "" {
		return strings.TrimSpace(s.cfg.Update.ProxyURL)
	}
	if v := strings.TrimSpace(os.Getenv("UPDATE_PROXY_URL")); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("HTTPS_PROXY")); v != "" {
		return v
	}
	if v := strings.TrimSpace(os.Getenv("HTTP_PROXY")); v != "" {
		return v
	}
	return defaultMihomoProxyURL
}

type proxyTrafficPoolMihomoState struct {
	content       string
	subscriptions map[string]struct{}
	usage         *proxyTrafficUsage
}

func (s proxyTrafficPoolMihomoState) hasSubscription(raw string) bool {
	key := normalizeProxySubscriptionURL(raw)
	if key == "" {
		return false
	}
	_, ok := s.subscriptions[key]
	return ok
}

func (s proxyTrafficPoolMihomoState) referencesSubscription(raw string) bool {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return false
	}
	return s.content != "" && strings.Contains(s.content, raw)
}

func readMihomoTrafficPoolState(items []*OpsProxyTrafficPoolConfigItem) proxyTrafficPoolMihomoState {
	content := readMihomoConfigForTrafficPool(items)
	return proxyTrafficPoolMihomoState{
		content:       content,
		subscriptions: readMihomoSubscriptions(defaultMihomoSubscriptionsPath),
		usage:         parseSSRDOGSubscriptionBody(content),
	}
}

func readMihomoConfigForTrafficPool(items []*OpsProxyTrafficPoolConfigItem) string {
	paths := []string{defaultMihomoStatusPath, defaultMihomoConfigPath}
	for _, item := range items {
		if item != nil && strings.TrimSpace(item.MihomoConfigPath) != "" {
			paths = append(paths, strings.TrimSpace(item.MihomoConfigPath))
		}
	}
	seen := map[string]struct{}{}
	var b strings.Builder
	for _, p := range paths {
		if _, ok := seen[p]; ok {
			continue
		}
		seen[p] = struct{}{}
		data, err := os.ReadFile(p)
		if err == nil {
			b.Write(data)
			b.WriteByte('\n')
		}
	}
	return b.String()
}

func readMihomoSubscriptions(path string) map[string]struct{} {
	out := map[string]struct{}{}
	data, err := os.ReadFile(path)
	if err != nil {
		return out
	}
	for _, line := range strings.Split(string(data), "\n") {
		key := normalizeProxySubscriptionURL(line)
		if key != "" {
			out[key] = struct{}{}
		}
	}
	return out
}

func normalizeProxySubscriptionURL(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Host == "" {
		return raw
	}
	parsed.Fragment = ""
	return parsed.String()
}

func shortURLHash(raw string) string {
	parsed, err := url.Parse(raw)
	key := raw
	if err == nil && parsed.Host != "" {
		key = parsed.Host + parsed.Path
	}
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])[:12]
}
