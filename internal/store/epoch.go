package store

import (
	"path/filepath"
	"time"
)

type EpochRecord struct {
	FurnaceID string
	Epoch     int64
	Reading   float64
	UpdatedAt string
}

type EpochStore struct {
	store *Store
}

func NewEpochStore(store *Store) *EpochStore {
	return &EpochStore{store: store}
}

func (e *EpochStore) Path(furnaceID string) string {
	return filepath.Join(e.store.Dir("epoch"), furnaceID+".json")
}

func (e *EpochStore) Save(record EpochRecord) error {
	return writeJSON(e.Path(record.FurnaceID), record)
}

func (e *EpochStore) Load(furnaceID string) (EpochRecord, error) {
	var record EpochRecord
	err := readJSON(e.Path(furnaceID), &record)
	return record, err
}

func (e *EpochStore) Reset(furnaceID string) error {
	record := EpochRecord{
		FurnaceID: furnaceID,
		Epoch:     0,
		Reading:   0,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return e.Save(record)
}
