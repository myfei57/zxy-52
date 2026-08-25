package grate

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
	"wastegen/internal/store"
)

type BurnoutService struct {
	baseline *furnace.BaselineService
	audit    *audit.Recorder
}

func NewBurnoutService(baseline *furnace.BaselineService, recorder *audit.Recorder) *BurnoutService {
	return &BurnoutService{
		baseline: baseline,
		audit:    recorder,
	}
}

func (b *BurnoutService) Verdict(furnaceID string, rate float64) (bool, store.BaselineRecord, error) {
	baseline := b.baseline.CurrentBaseline(furnaceID)
	pass := rate >= baseline.Threshold
	if !pass {
		message := fmt.Sprintf("burnout rate %.3f below baseline %.3f", rate, baseline.Threshold)
		if err := b.audit.Record(audit.Event{
			FurnaceID: furnaceID,
			EventType: audit.TypeBurnoutPoor,
			Message:   message,
		}); err != nil {
			return false, baseline, err
		}
	}
	return pass, baseline, nil
}

func (b *BurnoutService) Threshold(furnaceID string) float64 {
	return b.baseline.CurrentBaseline(furnaceID).Threshold
}
