package quota

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type QuotaService struct {
	quotaStore *store.QuotaStore
	audit      *audit.Recorder
}

func NewQuotaService(quotaStore *store.QuotaStore, recorder *audit.Recorder) *QuotaService {
	return &QuotaService{
		quotaStore: quotaStore,
		audit:      recorder,
	}
}

func (q *QuotaService) SetLimit(furnaceID string, limit float64, date string) (store.QuotaRecord, error) {
	record := store.QuotaRecord{
		FurnaceID:  furnaceID,
		DailyLimit: limit,
		Date:       date,
	}
	if err := q.quotaStore.Save(record); err != nil {
		return record, err
	}
	return record, nil
}

func (q *QuotaService) Reserve(furnaceID string, kg float64, date string) (bool, error) {
	// The store performs the load-check-update-persist under one mutex so two
	// concurrent reservations against the same furnace cannot interleave and
	// lose a deduction.
	ok, err := q.quotaStore.Reserve(furnaceID, kg, date)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, nil
	}
	if err := q.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeQuotaReserved,
		Message:   fmt.Sprintf("reserved %.1f kg against daily quota", kg),
	}); err != nil {
		return false, err
	}
	return true, nil
}

func (q *QuotaService) Release(furnaceID string, kg float64, date string) (store.QuotaRecord, error) {
	// Delegated to the store so the load, roll-over, subtract and persist run
	// under one critical section — a concurrent Reserve cannot drop this
	// returned quota by interleaving between a separate Load and Save.
	record, err := q.quotaStore.Release(furnaceID, kg, date)
	if err != nil {
		return record, err
	}
	return record, nil
}

func (q *QuotaService) Usage(furnaceID string, date string) (store.QuotaRecord, error) {
	record, err := q.quotaStore.Load(furnaceID)
	if err != nil {
		return record, err
	}
	if record.Date != date {
		record.Date = date
		record.Used = 0
	}
	return record, nil
}
