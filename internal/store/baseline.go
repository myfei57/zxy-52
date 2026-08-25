package store

import (
	"path/filepath"
	"time"
)

type BaselineRecord struct {
	FurnaceID string
	Threshold float64
	Model     string
	UpdatedAt string
}

type BaselineStore struct {
	store *Store
}

func NewBaselineStore(store *Store) *BaselineStore {
	return &BaselineStore{store: store}
}

func (b *BaselineStore) Path(furnaceID string) string {
	return filepath.Join(b.store.Dir("baseline"), furnaceID+".json")
}

func (b *BaselineStore) Save(record BaselineRecord) error {
	return writeJSON(b.Path(record.FurnaceID), record)
}

func (b *BaselineStore) Load(furnaceID string) (BaselineRecord, error) {
	var record BaselineRecord
	err := readJSON(b.Path(furnaceID), &record)
	return record, err
}

func (b *BaselineStore) Default(furnaceID string) (BaselineRecord, error) {
	record := BaselineRecord{
		FurnaceID: furnaceID,
		Threshold: 0.9,
		Model:     "standard",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	err := b.Save(record)
	return record, err
}
