package service

import (
	"context"
	"sort"
	"strings"
	"time"
)

const (
	PublicChannelStatusModelID          = "gpt-5.5"
	PublicChannelStatusDisplayName      = "GPT-plus"
	PublicChannelStatusPro4DisplayName  = "GPT-pro"
	PublicChannelStatusGPT56DisplayName = "GPT-5.6"
	PublicChannelStatusGPT56ModelID     = "gpt-5.6"
	publicChannelStatusCronExpression   = "*/5 * * * *"
	publicChannelStatusMaxResults       = 240
)

type publicChannelStatusProfile struct {
	Key          string
	DisplayName  string
	ModelID      string
	GroupNames   []string
	AccountNames []string
}

var PublicChannelStatusGroupNames = []string{
	"codex-plus-2倍",
	"Codex Plus 2倍",
	"gptplus2.5倍",
	"gptplus 2.5倍",
	"GPTPlus2.5倍",
	"codex-plus-2.5倍",
	"Codex Plus 2.5倍",
	"codex-plus-2.8倍",
	"每日10刀订阅-2倍",
	"每日20刀订阅-2倍",
	"每日45刀订阅-2倍",
	"每日60刀订阅-2倍",
	"每日90刀订阅-2倍",
	"每日135刀订阅-2倍",
	"每日180刀订阅-2倍",
}

var PublicChannelStatusPro4GroupNames = []string{
	"codex-pro-四倍",
	"限时福利codex-pro-4倍",
	"pro4倍",
	"pro 4倍",
	"Pro4倍",
	"Pro 4倍",
	"PRO4倍",
	"gptpro4倍",
	"gpt pro4倍",
	"gptplus pro4倍",
	"GPTPlus Pro4倍",
}

var PublicChannelStatusPro4AccountNames = []string{
	"虾",
}

// The Grox trial is deliberately pinned to its only account. This keeps the
// public health signal aligned with the group users actually invoke.
var PublicChannelStatusGPT56GroupNames = []string{
	"GPT-5.6 3倍",
}

var PublicChannelStatusGPT56GroxAccountNames = []string{
	"grox",
}

var publicChannelStatusProfiles = []publicChannelStatusProfile{
	{
		Key:         "gptplus25",
		DisplayName: PublicChannelStatusDisplayName,
		ModelID:     PublicChannelStatusModelID,
		GroupNames:  PublicChannelStatusGroupNames,
	},
	{
		Key:          "pro4",
		DisplayName:  PublicChannelStatusPro4DisplayName,
		ModelID:      PublicChannelStatusModelID,
		GroupNames:   PublicChannelStatusPro4GroupNames,
		AccountNames: PublicChannelStatusPro4AccountNames,
	},
	{
		Key:          "gpt56",
		DisplayName:  PublicChannelStatusGPT56DisplayName,
		ModelID:      PublicChannelStatusGPT56ModelID,
		GroupNames:   PublicChannelStatusGPT56GroupNames,
		AccountNames: PublicChannelStatusGPT56GroxAccountNames,
	},
}

type PublicChannelStatusService struct {
	planRepo   ScheduledTestPlanRepository
	resultRepo ScheduledTestResultRepository
}

func NewPublicChannelStatusService(
	planRepo ScheduledTestPlanRepository,
	resultRepo ScheduledTestResultRepository,
) *PublicChannelStatusService {
	return &PublicChannelStatusService{
		planRepo:   planRepo,
		resultRepo: resultRepo,
	}
}

func (s *PublicChannelStatusService) EnsurePlans(ctx context.Context) (int64, error) {
	nextRun, err := computeNextRun(publicChannelStatusCronExpression, time.Now())
	if err != nil {
		return 0, err
	}
	var total int64
	for _, profile := range publicChannelStatusProfiles {
		created, err := s.planRepo.EnsurePublicStatusPlans(
			ctx,
			profile.GroupNames,
			profile.AccountNames,
			profile.ModelID,
			publicChannelStatusCronExpression,
			publicChannelStatusMaxResults,
			nextRun,
		)
		if err != nil {
			return total, err
		}
		total += created
	}
	return total, nil
}

func (s *PublicChannelStatusService) GetSnapshot(ctx context.Context) (*PublicChannelStatusSnapshot, error) {
	return s.GetSnapshotByProfile(ctx, "")
}

func (s *PublicChannelStatusService) GetSnapshotByProfile(ctx context.Context, rawProfile string) (*PublicChannelStatusSnapshot, error) {
	if _, err := s.EnsurePlans(ctx); err != nil {
		return nil, err
	}

	profile := publicChannelStatusProfileByKey(rawProfile)
	now := time.Now()
	stats, err := s.resultRepo.GetPublicChannelStatusStats(ctx, profile.GroupNames, profile.AccountNames, profile.ModelID, now, 160)
	if err != nil {
		return nil, err
	}
	if stats == nil {
		stats = &PublicChannelStatusStats{}
	}

	return &PublicChannelStatusSnapshot{
		GroupName:        profile.DisplayName,
		ModelID:          profile.ModelID,
		Status:           publicChannelStatus(stats, now),
		ActiveChannels:   stats.ActiveChannels,
		CheckedChannels:  stats.CheckedChannels,
		SuccessRate1h:    roundPercent(stats.SuccessRate1h),
		SuccessRate24h:   roundPercent(stats.SuccessRate24h),
		AverageLatencyMs: stats.AverageLatencyMs,
		P95LatencyMs:     stats.P95LatencyMs,
		LastCheckedAt:    stats.LastCheckedAt,
		Timeline:         buildPublicChannelTimeline(stats.RecentSamples, now),
		GeneratedAt:      now,
	}, nil
}

func publicChannelStatusProfileByKey(raw string) publicChannelStatusProfile {
	key := normalizePublicChannelStatusProfile(raw)
	for _, profile := range publicChannelStatusProfiles {
		if profile.Key == key {
			return profile
		}
	}
	return publicChannelStatusProfiles[0]
}

func normalizePublicChannelStatusProfile(raw string) string {
	value := strings.ToLower(strings.TrimSpace(raw))
	value = strings.ReplaceAll(value, " ", "")
	value = strings.ReplaceAll(value, "_", "")
	value = strings.ReplaceAll(value, "-", "")
	switch value {
	case "pro4", "pro4x", "pro4倍", "gptpro4", "gptpro4倍":
		return "pro4"
	case "gpt56", "gpt5.6", "gpt5_6", "gpt-5.6", "grox":
		return "gpt56"
	default:
		return "gptplus25"
	}
}

func publicChannelStatus(stats *PublicChannelStatusStats, now time.Time) string {
	if stats == nil || stats.ActiveChannels == 0 || stats.CheckedChannels == 0 || stats.LastCheckedAt == nil {
		return "unknown"
	}
	if now.Sub(*stats.LastCheckedAt) > 10*time.Minute {
		return "degraded"
	}
	if stats.SuccessRate1h >= 99 {
		return "operational"
	}
	if stats.SuccessRate1h > 0 {
		return "degraded"
	}
	return "failed"
}

func buildPublicChannelTimeline(samples []PublicChannelStatusSample, now time.Time) []PublicChannelStatusPoint {
	const bucketCount = 24
	const bucketSize = 5 * time.Minute

	start := now.Truncate(bucketSize).Add(-time.Duration(bucketCount-1) * bucketSize)
	type bucket struct {
		total   int
		success int
	}
	buckets := make([]bucket, bucketCount)

	for _, sample := range samples {
		if sample.CheckedAt.Before(start) || sample.CheckedAt.After(now.Add(bucketSize)) {
			continue
		}
		idx := int(sample.CheckedAt.Truncate(bucketSize).Sub(start) / bucketSize)
		if idx < 0 || idx >= bucketCount {
			continue
		}
		buckets[idx].total++
		if sample.Status == "success" {
			buckets[idx].success++
		}
	}

	points := make([]PublicChannelStatusPoint, 0, bucketCount)
	for i, b := range buckets {
		status := "unknown"
		if b.total > 0 {
			switch {
			case b.success == b.total:
				status = "operational"
			case b.success > 0:
				status = "degraded"
			default:
				status = "failed"
			}
		}
		points = append(points, PublicChannelStatusPoint{
			Status:    status,
			CheckedAt: start.Add(time.Duration(i) * bucketSize),
		})
	}
	return points
}

func roundPercent(v float64) float64 {
	return float64(int(v*10+0.5)) / 10
}

func init() {
	sort.Strings(PublicChannelStatusGroupNames)
	sort.Strings(PublicChannelStatusPro4GroupNames)
	sort.Strings(PublicChannelStatusPro4AccountNames)
	sort.Strings(PublicChannelStatusGPT56GroupNames)
	sort.Strings(PublicChannelStatusGPT56GroxAccountNames)
}
