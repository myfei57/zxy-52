package store

import (
	"path/filepath"
	"time"
)

type ZoneBinding struct {
	ZoneID         string
	ThermocoupleID string
	Index          int
}

type ZoneSnapshot struct {
	FurnaceID string
	Bindings  []ZoneBinding
	SavedAt   string
}

type ZoneStore struct {
	store *Store
}

func NewZoneStore(store *Store) *ZoneStore {
	return &ZoneStore{store: store}
}

func (z *ZoneStore) Path(furnaceID string) string {
	return filepath.Join(z.store.Dir("zones"), furnaceID+".json")
}

func (z *ZoneStore) Save(snapshot ZoneSnapshot) error {
	return writeJSON(z.Path(snapshot.FurnaceID), snapshot)
}

func (z *ZoneStore) Load(furnaceID string) (ZoneSnapshot, error) {
	var snapshot ZoneSnapshot
	err := readJSON(z.Path(furnaceID), &snapshot)
	return snapshot, err
}

func (z *ZoneStore) Bindings(furnaceID string) ([]ZoneBinding, error) {
	snapshot, err := z.Load(furnaceID)
	if err != nil {
		return nil, err
	}
	return snapshot.Bindings, nil
}

func (z *ZoneStore) Touch(furnaceID string, bindings []ZoneBinding) error {
	snapshot := ZoneSnapshot{
		FurnaceID: furnaceID,
		Bindings:  bindings,
		SavedAt:   time.Now().UTC().Format(time.RFC3339),
	}
	return z.Save(snapshot)
}
