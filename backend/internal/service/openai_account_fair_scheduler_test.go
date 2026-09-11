package service

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAllOpenAICandidatePrioritiesEqual(t *testing.T) {
	tests := []struct {
		name     string
		pool     []openAIAccountCandidateScore
		expected bool
	}{
		{
			name: "same priority",
			pool: []openAIAccountCandidateScore{
				{effectivePriority: 1},
				{effectivePriority: 1},
			},
			expected: true,
		},
		{
			name: "different priority",
			pool: []openAIAccountCandidateScore{
				{effectivePriority: 1},
				{effectivePriority: 2},
			},
			expected: false,
		},
		{
			name: "zero is legacy unset priority",
			pool: []openAIAccountCandidateScore{
				{effectivePriority: 0},
				{effectivePriority: 0},
			},
			expected: false,
		},
		{
			name:     "single account is not a peer pool",
			pool:     []openAIAccountCandidateScore{{effectivePriority: 1}},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, allOpenAICandidatePrioritiesEqual(tt.pool))
		})
	}
}

func TestBuildOpenAIFairSelectionOrderRotatesStableAccountIDs(t *testing.T) {
	scheduler := &defaultOpenAIAccountScheduler{}
	groupID := int64(40)
	pool := []openAIAccountCandidateScore{
		{account: &Account{ID: 30}, effectivePriority: 1},
		{account: &Account{ID: 10}, effectivePriority: 1},
		{account: &Account{ID: 20}, effectivePriority: 1},
	}

	starts := make([]int64, 0, 3)
	for i := 0; i < 3; i++ {
		order := scheduler.buildOpenAIFairSelectionOrder(pool, &groupID)
		require.Len(t, order, 3)
		starts = append(starts, order[0].account.ID)
	}

	require.Equal(t, []int64{10, 20, 30}, starts)
}

func TestBuildOpenAISelectionOrderUsesFairPeersBeforeScoreModel(t *testing.T) {
	scheduler := &defaultOpenAIAccountScheduler{}
	groupID := int64(40)
	plan := openAIAccountLoadPlan{
		candidates: []openAIAccountCandidateScore{
			{account: &Account{ID: 30}, effectivePriority: 1, score: 99},
			{account: &Account{ID: 10}, effectivePriority: 1, score: 1},
			{account: &Account{ID: 20}, effectivePriority: 1, score: 50},
		},
		topK: 1,
	}

	order := scheduler.buildOpenAISelectionOrder(OpenAIAccountScheduleRequest{GroupID: &groupID}, plan)
	require.Equal(t, []int64{10, 20, 30}, []int64{order[0].account.ID, order[1].account.ID, order[2].account.ID})
}
