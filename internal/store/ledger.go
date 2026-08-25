package store

import (
	"encoding/json"
	"path/filepath"
	"sync"
)

type EmissionRecord struct {
	FurnaceID  string
	Timestamp  string
	O2         float64
	NOx        float64
	Ammonia    float64
	ZoneID     string
	RecordedAt string
	ID         string
}

type LedgerStore struct {
	store *Store
	mu    sync.Mutex
}

func NewLedgerStore(store *Store) *LedgerStore {
	return &LedgerStore{store: store}
}

func (l *LedgerStore) Path(furnaceID string) string {
	return filepath.Join(l.store.Dir("ledger"), furnaceID+".jsonl")
}

func (l *LedgerStore) Append(record EmissionRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	return appendJSONL(l.Path(record.FurnaceID), record)
}

func (l *LedgerStore) List(furnaceID string) ([]EmissionRecord, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	var records []EmissionRecord
	err := readJSONLLines(l.Path(furnaceID), func(raw []byte) error {
		var record EmissionRecord
		if err := json.Unmarshal(raw, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	return records, err
}

func (l *LedgerStore) HasTimestamp(furnaceID string, timestamp string) (bool, error) {
	records, err := l.List(furnaceID)
	if err != nil {
		return false, err
	}
	for _, record := range records {
		if record.ID == timestamp {
			return true, nil
		}
	}
	return false, nil
}

func (l *LedgerStore) Count(furnaceID string) (int, error) {
	records, err := l.List(furnaceID)
	if err != nil {
		return 0, err
	}
	return len(records), nil
}

func (l *LedgerStore) WriteAll(furnaceID string, records []EmissionRecord) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	path := l.Path(furnaceID)
	if err := mkdirAllFor(path); err != nil {
		return err
	}
	lines := make([]byte, 0, 256*len(records))
	for _, record := range records {
		line, err := json.Marshal(record)
		if err != nil {
			return err
		}
		lines = append(lines, line...)
		lines = append(lines, '\n')
	}
	return writeFileBytes(path, lines)
}
