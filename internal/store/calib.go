package store

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type CalibrationRecord struct {
	FurnaceID  string
	BaselineO2 float64
	AdjustedAt string
	Seq        int64
}

type CalibStore struct {
	store *Store
}

func NewCalibStore(store *Store) *CalibStore {
	return &CalibStore{store: store}
}

func (c *CalibStore) Path(furnaceID string) string {
	return filepath.Join(c.store.Dir("calib"), furnaceID+".jsonl")
}

func (c *CalibStore) Save(record CalibrationRecord) error {
	return appendJSONL(c.Path(record.FurnaceID), record)
}

func (c *CalibStore) Load(furnaceID string) (CalibrationRecord, error) {
	var records []CalibrationRecord
	err := readJSONLLines(c.Path(furnaceID), func(raw []byte) error {
		var record CalibrationRecord
		if err := json.Unmarshal(raw, &record); err != nil {
			return err
		}
		records = append(records, record)
		return nil
	})
	if err != nil {
		return CalibrationRecord{}, err
	}
	if len(records) == 0 {
		return CalibrationRecord{}, os.ErrNotExist
	}
	return records[0], nil
}

func (c *CalibStore) NextSeq(furnaceID string) int64 {
	record, err := c.Load(furnaceID)
	if err != nil {
		return 1
	}
	return record.Seq + 1
}

func (c *CalibStore) Reset(furnaceID string) error {
	record := CalibrationRecord{
		FurnaceID:  furnaceID,
		BaselineO2: 20.9,
		AdjustedAt: time.Now().UTC().Format(time.RFC3339),
		Seq:        c.NextSeq(furnaceID),
	}
	return c.Save(record)
}
