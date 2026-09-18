package service

import (
	"context"
	"testing"
)

func TestResolveSubscriptionRedeemRebateBaseAcceptsSingularDay(t *testing.T) {
	ctx := context.Background()
	client := newPaymentConfigServiceTestClient(t)
	groupID := int64(8)

	client.SubscriptionPlan.Create().
		SetGroupID(groupID).
		SetName("周卡135刀").
		SetPrice(70).
		SetValidityDays(7).
		SetValidityUnit("day").
		SetForSale(true).
		SetSortOrder(1).
		SaveX(ctx)

	svc := &RedeemService{entClient: client}
	got := svc.resolveSubscriptionRedeemRebateBase(ctx, groupID, 7)
	if got != 70 {
		t.Fatalf("resolveSubscriptionRedeemRebateBase() = %.2f, want 70.00", got)
	}
}
