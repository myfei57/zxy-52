package verifycase

import (
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/flue"
	"wastegen/internal/store"
)

func TestWgFlueDedupSingleWrite(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	ledger := store.NewLedgerStore(st)
	record := store.EmissionRecord{
		FurnaceID:  "f1",
		Timestamp:  "2026-08-25T03:00:00Z",
		O2:         6.0,
		NOx:        95.0,
		ZoneID:     "z1",
		RecordedAt: "2026-08-25T03:00:00Z",
		ID:         "rec-1",
	}
	collectorA := flue.NewDedupService(ledger, recorder)
	accepted, err := collectorA.Register(record)
	if err != nil {
		t.Fatal(err)
	}
	if !accepted {
		t.Fatal("first collector sample should be accepted")
	}
	collectorB := flue.NewDedupService(ledger, recorder)
	accepted, err = collectorB.Register(record)
	if err != nil {
		t.Fatal(err)
	}
	if accepted {
		t.Fatal("dedup mark missed the second collector write")
	}
	count, err := ledger.Count("f1")
	if err != nil {
		t.Fatal(err)
	}
	if count != 1 {
		t.Fatalf("ledger has %d records for one timestamp", count)
	}
}
