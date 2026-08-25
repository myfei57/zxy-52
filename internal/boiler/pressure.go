package boiler

import (
	"fmt"
	"sync"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
)

type PressureState struct {
	Pressure float64
	Target   float64
}

type PressureService struct {
	mu      sync.Mutex
	states  map[string]PressureState
	lifts   map[string]int
	combust *furnace.CombustService
	valve   *MainValveService
	audit   *audit.Recorder
}

func NewPressureService(combust *furnace.CombustService, valve *MainValveService, recorder *audit.Recorder) *PressureService {
	return &PressureService{
		states:  make(map[string]PressureState),
		lifts:   make(map[string]int),
		combust: combust,
		valve:   valve,
		audit:   recorder,
	}
}

func (p *PressureService) SetPressure(furnaceID string, pressure float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.states[furnaceID]
	state.Pressure = pressure
	p.states[furnaceID] = state
}

func (p *PressureService) SetTarget(furnaceID string, target float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	state := p.states[furnaceID]
	state.Target = target
	p.states[furnaceID] = state
}

func (p *PressureService) Pressure(furnaceID string) float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.states[furnaceID].Pressure
}

func (p *PressureService) Target(furnaceID string) float64 {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.states[furnaceID].Target
}

func (p *PressureService) Regulate(furnaceID string, pressureDelta float64) (float64, error) {
	p.mu.Lock()
	state := p.states[furnaceID]
	p.mu.Unlock()
	fuel := p.combust.AdjustFuel(furnaceID, pressureDelta*2.0)
	if err := p.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeCombustionAdjust,
		Message:   fmt.Sprintf("combustion fuel adjusted to %.1f", fuel),
	}); err != nil {
		return 0, err
	}
	if err := p.valve.Move(furnaceID, pressureDelta); err != nil {
		return 0, err
	}
	if err := p.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeSteamValveMove,
		Message:   fmt.Sprintf("main steam valve moved by %.1f", pressureDelta),
	}); err != nil {
		return 0, err
	}
	p.mu.Lock()
	state = p.states[furnaceID]
	state.Pressure = state.Pressure + pressureDelta*0.1
	if state.Pressure >= state.Target+3.0 {
		p.lifts[furnaceID]++
	}
	p.states[furnaceID] = state
	p.mu.Unlock()
	return state.Pressure, nil
}

func (p *PressureService) SafetyValveLiftCount(furnaceID string) int {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.lifts[furnaceID]
}

func (p *PressureService) Reset(furnaceID string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.states[furnaceID] = PressureState{}
	p.lifts[furnaceID] = 0
}
