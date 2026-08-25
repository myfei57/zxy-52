package furnace

import (
	"fmt"
	"time"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type BaselineService struct {
	baselineStore *store.BaselineStore
	baselines     map[string]store.BaselineRecord
	audit         *audit.Recorder
}

func NewBaselineService(baselineStore *store.BaselineStore, recorder *audit.Recorder) *BaselineService {
	return &BaselineService{
		baselineStore: baselineStore,
		baselines:     make(map[string]store.BaselineRecord),
		audit:         recorder,
	}
}

func (b *BaselineService) SetBaseline(furnaceID string, threshold float64) (store.BaselineRecord, error) {
	record := store.BaselineRecord{
		FurnaceID: furnaceID,
		Threshold: threshold,
		Model:     "standard",
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := b.baselineStore.Save(record); err != nil {
		return record, err
	}
	b.baselines[furnaceID] = record
	return record, nil
}

func (b *BaselineService) CurrentBaseline(furnaceID string) store.BaselineRecord {
	if record, ok := b.baselines[furnaceID]; ok {
		return record
	}
	record, err := b.baselineStore.Load(furnaceID)
	if err != nil {
		defaultRecord := store.BaselineRecord{
			FurnaceID: furnaceID,
			Threshold: 0.9,
			Model:     "standard",
		}
		b.baselines[furnaceID] = defaultRecord
		return defaultRecord
	}
	b.baselines[furnaceID] = record
	return record
}

func (b *BaselineService) Retrofit(furnaceID string, threshold float64, model string) (store.BaselineRecord, error) {
	record := store.BaselineRecord{
		FurnaceID: furnaceID,
		Threshold: threshold,
		Model:     model,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := b.baselineStore.Save(record); err != nil {
		return record, err
	}
	b.baselines[furnaceID] = record
	message := fmt.Sprintf("furnace retrofitted with baseline threshold %.3f model %s", threshold, model)
	if err := b.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeBaselineRetrofit,
		Message:   message,
	}); err != nil {
		return record, err
	}
	return record, nil
}

func (b *BaselineService) Restore(furnaceID string) error {
	record, err := b.baselineStore.Load(furnaceID)
	if err != nil {
		return err
	}
	b.baselines[furnaceID] = record
	return nil
}
