package steam

import (
	"fmt"
	"time"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

const DefaultAccumulatorCap = 10000.0

type AccumulatorService struct {
	epoch *store.EpochStore
	cap   float64
	audit *audit.Recorder
}

func NewAccumulatorService(epochStore *store.EpochStore, recorder *audit.Recorder) *AccumulatorService {
	return &AccumulatorService{
		epoch: epochStore,
		cap:   DefaultAccumulatorCap,
		audit: recorder,
	}
}

func (a *AccumulatorService) Cap() float64 {
	return a.cap
}

func (a *AccumulatorService) Accumulate(furnaceID string, reading float64) (float64, int64, error) {
	last, err := a.epoch.Load(furnaceID)
	lastEpoch := int64(0)
	lastReading := 0.0
	if err == nil {
		lastEpoch = last.Epoch
		lastReading = last.Reading
	}
	epoch := lastEpoch
	if reading < lastReading {
		epoch = lastEpoch + 1
	}
	current := reading + float64(epoch)*a.cap
	previous := lastReading + float64(lastEpoch)*a.cap
	delta := current - previous
	record := store.EpochRecord{
		FurnaceID: furnaceID,
		Epoch:     epoch,
		Reading:   reading,
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := a.epoch.Save(record); err != nil {
		return 0, epoch, err
	}
	message := fmt.Sprintf("steam accumulated %.2f into epoch %d", delta, epoch)
	if err := a.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeSteamAccumulated,
		Message:   message,
	}); err != nil {
		return 0, epoch, err
	}
	return delta, epoch, nil
}

func (a *AccumulatorService) Output(furnaceID string) (float64, int64, error) {
	last, err := a.epoch.Load(furnaceID)
	if err != nil {
		return 0, 0, err
	}
	return last.Reading + float64(last.Epoch)*a.cap, last.Epoch, nil
}
