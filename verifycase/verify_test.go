package verifycase

import (
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
	"wastegen/internal/grate"
	"wastegen/internal/store"
)

func TestWgBurnoutBaselineFresh(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	baselines := furnace.NewBaselineService(store.NewBaselineStore(st), recorder)
	if _, err := baselines.SetBaseline("f1", 0.90); err != nil {
		t.Fatal(err)
	}
	burnout := grate.NewBurnoutService(baselines, recorder)
	pass, baseline, err := burnout.Verdict("f1", 0.92)
	if err != nil {
		t.Fatal(err)
	}
	if !pass {
		t.Fatalf("baseline 0.90 should accept rate 0.92, got threshold %v", baseline.Threshold)
	}
	if _, err := baselines.Retrofit("f1", 0.95, "retrofit"); err != nil {
		t.Fatal(err)
	}
	pass, baseline, err = burnout.Verdict("f1", 0.92)
	if err != nil {
		t.Fatal(err)
	}
	if pass {
		t.Fatalf("burnout judged against the pre-retrofit baseline %v", baseline.Threshold)
	}
	if baseline.Threshold != 0.95 {
		t.Fatalf("baseline not refreshed: %v", baseline.Threshold)
	}
	fresh := grate.NewBurnoutService(baselines, recorder)
	if pass, _, err := fresh.Verdict("f1", 0.92); err != nil || pass {
		t.Fatalf("restart verdict stale: pass=%v err=%v", pass, err)
	}
}
