package verifycase

import (
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/crane"
	"wastegen/internal/grate"
	"wastegen/internal/store"
)

func TestWgCraneGrabOrder(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	hopper := grate.NewHopperService(grate.NewHopperState(), recorder)
	tip := crane.NewTipService(hopper, recorder)
	if err := tip.Tip("f1", "g1", 3200); err != nil {
		t.Fatal(err)
	}
	if got := tip.SpillCount("f1"); got != 0 {
		t.Fatalf("waste spilled before the grab aligned: %d", got)
	}
	events, err := recorder.Events()
	if err != nil {
		t.Fatal(err)
	}
	idxAlign := -1
	idxRelease := -1
	for i, event := range events {
		if event.EventType == audit.TypeHopperAligned {
			idxAlign = i
		}
		if event.EventType == audit.TypeGrabReleased {
			idxRelease = i
		}
	}
	if idxAlign == -1 || idxRelease == -1 {
		t.Fatal("hopper events missing")
	}
	if idxAlign >= idxRelease {
		t.Fatalf("grab released before the hopper aligned: align=%d release=%d", idxAlign, idxRelease)
	}
	if got := hopper.Loaded("hopper-f1"); got != 3200 {
		t.Fatalf("hopper load not recorded: %v", got)
	}
}
