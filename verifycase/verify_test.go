package verifycase

import (
	"math"
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/flue"
	"wastegen/internal/store"
)

func TestWgFlueO2CorrectionFresh(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	calib := store.NewCalibStore(st)
	o2 := flue.NewO2Service(calib, recorder)
	if _, err := o2.Calibrate("f1", 20.5); err != nil {
		t.Fatal(err)
	}
	got, err := o2.Correct("f1", 21.0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-20.7) > 1e-6 {
		t.Fatalf("first correction used the wrong baseline: %v", got)
	}
	if _, err := o2.Calibrate("f1", 20.0); err != nil {
		t.Fatal(err)
	}
	got, err = o2.Correct("f1", 21.0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-20.4) > 1e-6 {
		t.Fatalf("correction lagged behind recalibration: %v", got)
	}
	fresh := flue.NewO2Service(calib, recorder)
	if got := fresh.Baseline("f1"); math.Abs(got-20.0) > 1e-6 {
		t.Fatalf("restart baseline stale: %v", got)
	}
	got, err = fresh.Correct("f1", 21.0)
	if err != nil {
		t.Fatal(err)
	}
	if math.Abs(got-20.4) > 1e-6 {
		t.Fatalf("restart correction stale: %v", got)
	}
}
