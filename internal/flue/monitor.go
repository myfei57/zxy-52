package flue

import (
	"time"

	"github.com/google/uuid"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type MonitorService struct {
	dedup  *DedupService
	audit  *audit.Recorder
	latest map[string]store.EmissionRecord
}

func NewMonitorService(dedup *DedupService, recorder *audit.Recorder) *MonitorService {
	return &MonitorService{
		dedup:  dedup,
		audit:  recorder,
		latest: make(map[string]store.EmissionRecord),
	}
}

func (m *MonitorService) Sample(furnaceID string, timestamp string, o2 float64, nox float64, zoneID string) (bool, store.EmissionRecord, error) {
	record := store.EmissionRecord{
		FurnaceID:  furnaceID,
		Timestamp:  timestamp,
		O2:         o2,
		NOx:        nox,
		ZoneID:     zoneID,
		RecordedAt: time.Now().UTC().Format(time.RFC3339),
		ID:         uuid.NewString(),
	}
	accepted, err := m.dedup.Register(record)
	if err != nil {
		return false, record, err
	}
	if accepted {
		m.latest[furnaceID] = record
	}
	return accepted, record, nil
}

func (m *MonitorService) Latest(furnaceID string) (store.EmissionRecord, bool) {
	record, ok := m.latest[furnaceID]
	return record, ok
}
