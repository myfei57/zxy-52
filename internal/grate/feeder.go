package grate

import (
	"fmt"

	"wastegen/internal/audit"
)

type FeederState struct {
	dropped map[string]float64
}

func NewFeederState() *FeederState {
	return &FeederState{
		dropped: make(map[string]float64),
	}
}

func (f *FeederState) Add(furnaceID string, kg float64) {
	f.dropped[furnaceID] += kg
}

func (f *FeederState) Total(furnaceID string) float64 {
	return f.dropped[furnaceID]
}

type FeederService struct {
	state *FeederState
	audit *audit.Recorder
}

func NewFeederService(state *FeederState, recorder *audit.Recorder) *FeederService {
	return &FeederService{
		state: state,
		audit: recorder,
	}
}

func (f *FeederService) Drop(furnaceID string, kg float64) error {
	f.state.Add(furnaceID, kg)
	return f.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeFeedDropped,
		Message:   fmt.Sprintf("feeder dropped %.1f kg for %s", kg, furnaceID),
	})
}

func (f *FeederService) Total(furnaceID string) float64 {
	return f.state.Total(furnaceID)
}
