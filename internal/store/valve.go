package store

import (
	"path/filepath"
	"time"
)

type ValveRecord struct {
	FurnaceID string
	Setpoint  float64
	Position  float64
	UpdatedAt string
}

type ValveStore struct {
	store *Store
}

func NewValveStore(store *Store) *ValveStore {
	return &ValveStore{store: store}
}

func (v *ValveStore) Path(furnaceID string) string {
	return filepath.Join(v.store.Dir("valve"), furnaceID+".json")
}

func (v *ValveStore) Save(record ValveRecord) error {
	return writeJSON(v.Path(record.FurnaceID), record)
}

func (v *ValveStore) Load(furnaceID string) (ValveRecord, error) {
	var record ValveRecord
	err := readJSON(v.Path(furnaceID), &record)
	return record, err
}

func (v *ValveStore) Touch(furnaceID string, setpoint float64, position float64) error {
	record := ValveRecord{
		FurnaceID: furnaceID,
		Setpoint:  setpoint,
		Position:  position,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	return v.Save(record)
}
