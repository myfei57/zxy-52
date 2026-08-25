package verifycase

import (
	"os"
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/flue"
	"wastegen/internal/store"
)

func TestWgFlueLedgerErrorPropagates(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	ledger := store.NewLedgerStore(st)
	dedup := flue.NewDedupService(ledger, recorder)
	monitor := flue.NewMonitorService(dedup, recorder)
	accepted, _, err := monitor.Sample("f1", "2026-08-25T03:00:00Z", 6.0, 95.0, "z1")
	if err != nil {
		t.Fatal(err)
	}
	if !accepted {
		t.Fatal("first emission sample should be accepted")
	}
	path := ledger.Path("f1")
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(path, 0o644)
	if _, _, err := monitor.Sample("f1", "2026-08-25T03:01:00Z", 6.1, 96.0, "z1"); err == nil {
		t.Fatal("emission ledger write failure swallowed by the sample flow")
	}
}
