package furnace

import (
	"sync"

	"wastegen/internal/store"
)

type CombustionState struct {
	AirFlow  float64
	FuelFlow float64
}

type CombustService struct {
	mu     sync.Mutex
	states map[string]CombustionState
}

func NewCombustService() *CombustService {
	return &CombustService{
		states: make(map[string]CombustionState),
	}
}

func (c *CombustService) AdjustAir(record store.FeedRecord) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	state := c.states[record.FurnaceID]
	air := record.FeedKG * 9.8
	state.AirFlow = air
	c.states[record.FurnaceID] = state
	return air
}

func (c *CombustService) AdjustFuel(furnaceID string, delta float64) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	state := c.states[furnaceID]
	fuel := state.FuelFlow + delta
	if fuel < 0 {
		fuel = 0
	}
	state.FuelFlow = fuel
	c.states[furnaceID] = state
	return fuel
}

func (c *CombustService) SetAirFlow(furnaceID string, air float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	state := c.states[furnaceID]
	state.AirFlow = air
	c.states[furnaceID] = state
}

func (c *CombustService) SetFuelFlow(furnaceID string, fuel float64) {
	c.mu.Lock()
	defer c.mu.Unlock()
	state := c.states[furnaceID]
	state.FuelFlow = fuel
	c.states[furnaceID] = state
}

func (c *CombustService) AirFlow(furnaceID string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.states[furnaceID].AirFlow
}

func (c *CombustService) FuelFlow(furnaceID string) float64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.states[furnaceID].FuelFlow
}

func (c *CombustService) Snapshot(furnaceID string) CombustionState {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.states[furnaceID]
}

func (c *CombustService) Reset(furnaceID string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.states[furnaceID] = CombustionState{}
}
