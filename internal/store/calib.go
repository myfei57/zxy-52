package store

import (
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
	return filepath.Join(c.store.Dir("calib"), furnaceID+".json")
}

func (c *CalibStore) Save(record CalibrationRecord) error {
	return writeJSON(c.Path(record.FurnaceID), record)
}

func (c *CalibStore) Load(furnaceID string) (CalibrationRecord, error) {
	var record CalibrationRecord
	err := readJSON(c.Path(furnaceID), &record)
	return record, err
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
