package store

import (
	"path/filepath"
)

type SnapshotStore struct {
	store *Store
}

func NewSnapshotStore(store *Store) *SnapshotStore {
	return &SnapshotStore{store: store}
}

func (s *SnapshotStore) Path(namespace string, key string) string {
	return filepath.Join(s.store.Dir("snapshots"), namespace, key+".json")
}

func (s *SnapshotStore) Save(namespace string, key string, value any) error {
	return writeJSON(s.Path(namespace, key), value)
}

func (s *SnapshotStore) Load(namespace string, key string, value any) error {
	return readJSON(s.Path(namespace, key), value)
}

func (s *SnapshotStore) Has(namespace string, key string) bool {
	return fileExists(s.Path(namespace, key))
}
