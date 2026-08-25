package store

import (
	"encoding/json"
	"path/filepath"
	"sync"
)

type AuditRecord struct {
	ID        string
	FurnaceID string
	EventType string
	Message   string
	At        string
}

type AuditStore struct {
	store *Store
	mu    sync.Mutex
}

func NewAuditStore(store *Store) *AuditStore {
	return &AuditStore{store: store}
}

func (a *AuditStore) Path() string {
	return filepath.Join(a.store.Dir("audit"), "events.jsonl")
}

func (a *AuditStore) Append(record AuditRecord) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	return appendJSONL(a.Path(), record)
}

func (a *AuditStore) List() ([]AuditRecord, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var records []AuditRecord
	err := readJSONLLines(a.Path(), func(raw []byte) error {
		var record AuditRecord
		if err := json.Unmarshal(raw, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	return records, err
}

func (a *AuditStore) Recent(limit int) ([]AuditRecord, error) {
	records, err := a.List()
	if err != nil {
		return nil, err
	}
	if limit <= 0 || limit >= len(records) {
		return records, nil
	}
	return records[len(records)-limit:], nil
}

func (a *AuditStore) ForFurnace(furnaceID string) ([]AuditRecord, error) {
	records, err := a.List()
	if err != nil {
		return nil, err
	}
	var filtered []AuditRecord
	for _, record := range records {
		if record.FurnaceID == furnaceID {
			filtered = append(filtered, record)
		}
	}
	return filtered, nil
}
