package grate

import (
	"fmt"

	"wastegen/internal/audit"
)

type AshValveState struct {
	open map[string]bool
}

func NewAshValveState() *AshValveState {
	return &AshValveState{
		open: make(map[string]bool),
	}
}

func (v *AshValveState) SetOpen(furnaceID string, open bool) {
	v.open[furnaceID] = open
}

func (v *AshValveState) IsOpen(furnaceID string) bool {
	return v.open[furnaceID]
}

type AshValveService struct {
	state *AshValveState
	audit *audit.Recorder
}

func NewAshValveService(state *AshValveState, recorder *audit.Recorder) *AshValveService {
	return &AshValveService{
		state: state,
		audit: recorder,
	}
}

func (v *AshValveService) Open(furnaceID string, scraperRunning bool) error {
	if !scraperRunning {
		return fmt.Errorf("ash valve for %s cannot open before the scraper starts", furnaceID)
	}
	v.state.SetOpen(furnaceID, true)
	return v.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeAshValveOpened,
		Message:   fmt.Sprintf("ash valve for %s opened after scraper start", furnaceID),
	})
}

func (v *AshValveService) Close(furnaceID string) error {
	v.state.SetOpen(furnaceID, false)
	return v.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeAshValveClosed,
		Message:   fmt.Sprintf("ash valve for %s closed", furnaceID),
	})
}

func (v *AshValveService) IsOpen(furnaceID string) bool {
	return v.state.IsOpen(furnaceID)
}
