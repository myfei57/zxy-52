package grate

import (
	"fmt"
	"time"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
	"wastegen/internal/store"
)

type RampState struct {
	speeds map[string]float64
}

func NewRampState() *RampState {
	return &RampState{
		speeds: make(map[string]float64),
	}
}

func (r *RampState) SetSpeed(furnaceID string, speed float64) {
	r.speeds[furnaceID] = speed
}

func (r *RampState) Speed(furnaceID string) float64 {
	return r.speeds[furnaceID]
}

type RampService struct {
	state   *RampState
	feed    *store.FeedStore
	combust *furnace.CombustService
	audit   *audit.Recorder
}

func NewRampService(state *RampState, feed *store.FeedStore, combust *furnace.CombustService, recorder *audit.Recorder) *RampService {
	return &RampService{
		state:   state,
		feed:    feed,
		combust: combust,
		audit:   recorder,
	}
}

func (r *RampService) Ramp(furnaceID string, speed float64, feedKG float64) (float64, error) {
	record := store.FeedRecord{
		FurnaceID: furnaceID,
		FeedKG:    feedKG,
		Seq:       r.feed.NextSeq(furnaceID),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
	}
	if err := r.feed.Save(record); err != nil {
		return 0, err
	}
	r.state.SetSpeed(furnaceID, speed)
	air := r.combust.AdjustAir(record)
	message := fmt.Sprintf("feed %.1f kg persisted before ramp to %.1f", feedKG, speed)
	if err := r.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeFeedPersisted,
		Message:   message,
	}); err != nil {
		return 0, err
	}
	return air, nil
}

func (r *RampService) Recover(furnaceID string) (float64, error) {
	record, err := r.feed.Load(furnaceID)
	if err != nil {
		return 0, err
	}
	return r.combust.AdjustAir(record), nil
}

func (r *RampService) Speed(furnaceID string) float64 {
	return r.state.Speed(furnaceID)
}

func (r *RampService) Feed(furnaceID string) (store.FeedRecord, error) {
	return r.feed.Load(furnaceID)
}
