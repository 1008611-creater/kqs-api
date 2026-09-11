package service

import "testing"

func TestFinalizeAccountRoutingStatsUsesRealCountsAndPriorityRoles(t *testing.T) {
	items := []*OpsAccountRoutingStats{
		{AccountID: 1, Priority: 2, Schedulable: true, Status: StatusActive, RequestCount: 40, SuccessCount: 38},
		{AccountID: 2, Priority: 1, Schedulable: true, Status: StatusActive, RequestCount: 50, SuccessCount: 50},
		{AccountID: 3, Priority: 1, Schedulable: false, Status: StatusError, RequestCount: 10},
	}

	if got := finalizeAccountRoutingStats(items); got != 100 {
		t.Fatalf("total requests = %d, want 100", got)
	}
	if items[0].RequestShare != 40 || items[1].RequestShare != 50 || items[2].RequestShare != 10 {
		t.Fatalf("unexpected request shares: %.2f %.2f %.2f", items[0].RequestShare, items[1].RequestShare, items[2].RequestShare)
	}
	if items[0].SelectionRole != "fallback" || items[1].SelectionRole != "preferred" || items[2].SelectionRole != "unavailable" {
		t.Fatalf("unexpected roles: %q %q %q", items[0].SelectionRole, items[1].SelectionRole, items[2].SelectionRole)
	}
	if items[0].SuccessRate == nil || *items[0].SuccessRate != 95 {
		t.Fatalf("success rate = %v, want 95", items[0].SuccessRate)
	}
}

func TestFinalizeAccountRoutingStatsWithoutSamplesLeavesShareAndRateEmpty(t *testing.T) {
	items := []*OpsAccountRoutingStats{{AccountID: 1, Priority: 1, Schedulable: true, Status: StatusActive}}
	if got := finalizeAccountRoutingStats(items); got != 0 {
		t.Fatalf("total requests = %d, want 0", got)
	}
	if items[0].RequestShare != 0 || items[0].SuccessRate != nil || items[0].SelectionRole != "preferred" {
		t.Fatalf("unexpected empty-window result: %+v", items[0])
	}
}
