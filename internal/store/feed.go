package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type FeedRecord struct {
	FurnaceID string
	FeedKG    float64
	Seq       int64
	CreatedAt string
}

type FeedStore struct {
	store *Store
}

func NewFeedStore(store *Store) *FeedStore {
	return &FeedStore{store: store}
}

func (f *FeedStore) Path(furnaceID string) string {
	return filepath.Join(f.store.Dir("feed"), furnaceID+".json")
}

func (f *FeedStore) Save(record FeedRecord) error {
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(f.Path(record.FurnaceID)), 0o755); err != nil {
		return err
	}
	return os.WriteFile(f.Path(record.FurnaceID), data, 0o644)
}

func (f *FeedStore) Load(furnaceID string) (FeedRecord, error) {
	var record FeedRecord
	err := readJSON(f.Path(furnaceID), &record)
	return record, err
}

func (f *FeedStore) NextSeq(furnaceID string) int64 {
	record, err := f.Load(furnaceID)
	if err != nil {
		return 1
	}
	return record.Seq + 1
}

func (f *FeedStore) Touch(furnaceID string, feedKG float64) (FeedRecord, error) {
	record := FeedRecord{
		FurnaceID: furnaceID,
		FeedKG:    feedKG,
		Seq:       f.NextSeq(furnaceID),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	err := f.Save(record)
	return record, err
}

func (f *FeedStore) Exists(furnaceID string) bool {
	return fileExists(f.Path(furnaceID))
}
