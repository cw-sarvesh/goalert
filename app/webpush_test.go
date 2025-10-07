package app

import (
	"testing"

	"github.com/target/goalert/oncall"
)

func TestFilterPrimaryStepUserIDs(t *testing.T) {
	t.Run("only lowest step returned", func(t *testing.T) {
		users := []oncall.ServiceOnCallUser{
			{StepNumber: 2, UserID: "u2"},
			{StepNumber: 1, UserID: "u1"},
			{StepNumber: 3, UserID: "u3"},
			{StepNumber: 1, UserID: "u4"},
		}

		ids, step := filterPrimaryStepUserIDs(users)
		if len(ids) != 2 {
			t.Fatalf("expected 2 ids, got %d", len(ids))
		}
		if ids[0] != "u1" || ids[1] != "u4" {
			t.Fatalf("unexpected ids: %v", ids)
		}
		if step != 1 {
			t.Fatalf("expected step 1, got %d", step)
		}
	})

	t.Run("all steps filtered when equal", func(t *testing.T) {
		users := []oncall.ServiceOnCallUser{
			{StepNumber: 2, UserID: "u2"},
			{StepNumber: 2, UserID: "u3"},
		}

		ids, step := filterPrimaryStepUserIDs(users)
		if len(ids) != 2 {
			t.Fatalf("expected both ids, got %d", len(ids))
		}
		if step != 2 {
			t.Fatalf("expected step 2, got %d", step)
		}
	})

	t.Run("empty list", func(t *testing.T) {
		ids, step := filterPrimaryStepUserIDs(nil)
		if len(ids) != 0 {
			t.Fatalf("expected empty result, got %v", ids)
		}
		if step != -1 {
			t.Fatalf("expected step -1, got %d", step)
		}
	})
}
