package furnace

import (
	"sync"

	"wastegen/internal/ns"
)

type ThermocoupleService struct {
	mu       sync.Mutex
	readings map[string]float64
}

func NewThermocoupleService() *ThermocoupleService {
	return &ThermocoupleService{
		readings: make(map[string]float64),
	}
}

func (t *ThermocoupleService) SetReading(thermocoupleID string, temp float64) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.readings[thermocoupleID] = temp
}

func (t *ThermocoupleService) Reading(thermocoupleID string) float64 {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.readings[thermocoupleID]
}

func (t *ThermocoupleService) ZoneTemp(zone ns.Zone) float64 {
	return t.Reading(zone.ThermocoupleID)
}

func (t *ThermocoupleService) Average(furnaceID string, registry *ns.Registry) float64 {
	zones := registry.Zones(furnaceID)
	if len(zones) == 0 {
		return 0
	}
	var total float64
	for _, zone := range zones {
		total += t.ZoneTemp(zone)
	}
	return total / float64(len(zones))
}

func (t *ThermocoupleService) Max(furnaceID string, registry *ns.Registry) float64 {
	zones := registry.Zones(furnaceID)
	var maxTemp float64
	for _, zone := range zones {
		temp := t.ZoneTemp(zone)
		if temp > maxTemp {
			maxTemp = temp
		}
	}
	return maxTemp
}
