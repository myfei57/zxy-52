package store

import (
	"path/filepath"
	"time"
)

type QuotaRecord struct {
	FurnaceID  string
	DailyLimit float64
	Used       float64
	Date       string
	UpdatedAt  string
}

type QuotaStore struct {
	store *Store
	cache map[string]QuotaRecord
}

func NewQuotaStore(store *Store) *QuotaStore {
	return &QuotaStore{
		store: store,
		cache: make(map[string]QuotaRecord),
	}
}

func (q *QuotaStore) Path(furnaceID string) string {
	return filepath.Join(q.store.Dir("quota"), furnaceID+".json")
}

func (q *QuotaStore) Save(record QuotaRecord) error {
	if err := writeJSON(q.Path(record.FurnaceID), record); err != nil {
		return err
	}
	q.cache[record.FurnaceID] = record
	return nil
}

func (q *QuotaStore) Load(furnaceID string) (QuotaRecord, error) {
	return q.loadLocked(furnaceID)
}

func (q *QuotaStore) loadLocked(furnaceID string) (QuotaRecord, error) {
	if record, ok := q.cache[furnaceID]; ok {
		return record, nil
	}
	var record QuotaRecord
	if err := readJSON(q.Path(furnaceID), &record); err != nil {
		return record, err
	}
	q.cache[furnaceID] = record
	return record, nil
}

func (q *QuotaStore) Reserve(furnaceID string, kg float64, date string) (bool, error) {
	record, err := q.loadLocked(furnaceID)
	if err != nil {
		record = QuotaRecord{FurnaceID: furnaceID}
	}
	if record.Date != date {
		record.Date = date
		record.Used = 0
	}
	time.Sleep(2 * time.Millisecond)
	if record.Used+kg > record.DailyLimit {
		return false, nil
	}
	record.Used += kg
	record.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	if err := writeJSON(q.Path(furnaceID), record); err != nil {
		return false, err
	}
	q.cache[furnaceID] = record
	return true, nil
}

func (q *QuotaStore) Seed(furnaceID string, limit float64, date string) (QuotaRecord, error) {
	record := QuotaRecord{
		FurnaceID:  furnaceID,
		DailyLimit: limit,
		Used:       0,
		Date:       date,
		UpdatedAt:  time.Now().UTC().Format(time.RFC3339),
	}
	err := q.Save(record)
	return record, err
}
