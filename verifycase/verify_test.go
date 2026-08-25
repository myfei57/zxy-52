package verifycase

import (
	"sync"
	"sync/atomic"
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/quota"
	"wastegen/internal/store"
)

func TestWgQuotaConcurrentReserve(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	quotaStore := store.NewQuotaStore(st)
	service := quota.NewQuotaService(quotaStore, recorder)
	if _, err := service.SetLimit("f1", 50, "2026-08-25"); err != nil {
		t.Fatal(err)
	}
	var group sync.WaitGroup
	var accepted int32
	for i := 0; i < 100; i++ {
		group.Add(1)
		go func() {
			defer group.Done()
			ok, err := service.Reserve("f1", 1, "2026-08-25")
			if err == nil && ok {
				atomic.AddInt32(&accepted, 1)
			}
		}()
	}
	group.Wait()
	if got := atomic.LoadInt32(&accepted); got != 50 {
		t.Fatalf("daily quota accounting broken under concurrent reservations: accepted=%d want 50", got)
	}
	record, err := quotaStore.Load("f1")
	if err != nil {
		t.Fatal(err)
	}
	if record.Used != 50 {
		t.Fatalf("quota used under-counted: used=%v want 50", record.Used)
	}
}
