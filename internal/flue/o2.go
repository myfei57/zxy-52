package flue

import (
	"fmt"
	"time"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type O2Service struct {
	calib    *store.CalibStore
	baseline float64
	audit    *audit.Recorder
}

func NewO2Service(calibStore *store.CalibStore, recorder *audit.Recorder) *O2Service {
	return &O2Service{
		calib:    calibStore,
		baseline: 20.9,
		audit:    recorder,
	}
}

func (o *O2Service) Calibrate(furnaceID string, baseline float64) (store.CalibrationRecord, error) {
	record := store.CalibrationRecord{
		FurnaceID:  furnaceID,
		BaselineO2: baseline,
		AdjustedAt: time.Now().UTC().Format(time.RFC3339),
		Seq:        o.calib.NextSeq(furnaceID),
	}
	if err := o.calib.Save(record); err != nil {
		return record, err
	}
	if err := o.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeO2Calibrated,
		Message:   fmt.Sprintf("oxygen calibration baseline set to %.2f", baseline),
	}); err != nil {
		return record, err
	}
	return record, nil
}

func (o *O2Service) Correct(furnaceID string, reading float64) (float64, error) {
	// Read the latest calibration for THIS furnace so the correction tracks
	// the most recent baseline. The previous implementation used o.baseline, a
	// single in-memory float shared across all furnaces that never refreshed
	// after a process restart or when calibration was written via the store
	// directly (e.g. the /api/calib/{id}/reset endpoint). After a load change
	// the correction kept trailing the old baseline, biasing the ammonia dose
	// and pushing outlet NOx past the emission limit.
	baseline := o.Baseline(furnaceID)
	corrected := reading + (baseline-reading)*0.6
	if err := o.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeO2Corrected,
		Message:   fmt.Sprintf("oxygen corrected to %.2f against baseline %.2f", corrected, baseline),
	}); err != nil {
		return 0, err
	}
	return corrected, nil
}

func (o *O2Service) Baseline(furnaceID string) float64 {
	record, err := o.calib.Load(furnaceID)
	if err != nil {
		return o.baseline
	}
	return record.BaselineO2
}
