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
	record, err := q.quotaStore.Load(furnaceID)
	if err != nil {
		return record, err
	}
	if record.Date != date {
		record.Date = date
		record.Used = 0
	}
	record.Used -= kg
	if record.Used < 0 {
		record.Used = 0
	}
	err = q.quotaStore.Save(record)
	return record, err
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
