package ash

import (
	"fmt"

	"wastegen/internal/audit"
)

type ScraperState struct {
	running map[string]bool
}

func NewScraperState() *ScraperState {
	return &ScraperState{
		running: make(map[string]bool),
	}
}

func (s *ScraperState) SetRunning(furnaceID string, running bool) {
	s.running[furnaceID] = running
}

func (s *ScraperState) IsRunning(furnaceID string) bool {
	return s.running[furnaceID]
}

type ScraperService struct {
	state *ScraperState
	audit *audit.Recorder
}

func NewScraperService(state *ScraperState, recorder *audit.Recorder) *ScraperService {
	return &ScraperService{
		state: state,
		audit: recorder,
	}
}

func (s *ScraperService) Start(furnaceID string) error {
	s.state.SetRunning(furnaceID, true)
	return s.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeScraperStarted,
		Message:   fmt.Sprintf("ash scraper for %s started", furnaceID),
	})
}

func (s *ScraperService) Stop(furnaceID string) error {
	s.state.SetRunning(furnaceID, false)
	return s.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeAshValveClosed,
		Message:   fmt.Sprintf("ash scraper for %s stopped", furnaceID),
	})
}

func (s *ScraperService) IsRunning(furnaceID string) bool {
	return s.state.IsRunning(furnaceID)
}
