package verifycase

import (
	"os"
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
	"wastegen/internal/grate"
	"wastegen/internal/store"
)

func TestWgGrateAfterFeedPersist(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	combust := furnace.NewCombustService()
	feed := store.NewFeedStore(st)
	ramp := grate.NewRampService(grate.NewRampState(), feed, combust, recorder)
	if _, err := ramp.Ramp("f1", 100, 500); err != nil {
		t.Fatal(err)
	}
	path := feed.Path("f1")
	if err := os.Chmod(path, 0o444); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(path, 0o644)
	if _, err := ramp.Ramp("f1", 120, 700); err == nil {
		t.Fatal("ramp should fail while the feed record cannot be persisted")
	}
	if got := ramp.Speed("f1"); got != 100 {
		t.Fatalf("grate speed ramped before the feed record was durable: speed = %v", got)
	}
	loaded, err := feed.Load("f1")
	if err != nil {
		t.Fatal(err)
	}
	if loaded.FeedKG != 500 {
		t.Fatalf("combustion feed should stay at the persisted 500 kg, got %v", loaded.FeedKG)
	}
}
