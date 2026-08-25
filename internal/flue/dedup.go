package flue

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type DedupService struct {
	ledger *store.LedgerStore
	audit  *audit.Recorder
	seen   map[string]bool
}

func NewDedupService(ledger *store.LedgerStore, recorder *audit.Recorder) *DedupService {
	return &DedupService{
		ledger: ledger,
		audit:  recorder,
		seen:   make(map[string]bool),
	}
}

func (d *DedupService) Register(record store.EmissionRecord) (bool, error) {
	key := record.FurnaceID + "|" + record.Timestamp
	if d.seen[key] {
		if err := d.audit.Record(audit.Event{
			FurnaceID: record.FurnaceID,
			EventType: audit.TypeEmissionDuplicate,
			Message:   fmt.Sprintf("duplicate flue sample at %s skipped", record.Timestamp),
		}); err != nil {
			return false, err
		}
		return false, nil
	}
	if err := d.ledger.Append(record); err != nil {
		return false, err
	}
	d.seen[key] = true
	if err := d.audit.Record(audit.Event{
		FurnaceID: record.FurnaceID,
		EventType: audit.TypeEmissionRecorded,
		Message:   fmt.Sprintf("flue sample at %s recorded", record.Timestamp),
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (d *DedupService) Count(furnaceID string) (int, error) {
	return d.ledger.Count(furnaceID)
}

func (d *DedupService) List(furnaceID string) ([]store.EmissionRecord, error) {
	return d.ledger.List(furnaceID)
}
