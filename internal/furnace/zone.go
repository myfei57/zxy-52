package furnace

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/ns"
	"wastegen/internal/store"
)

type ZoneService struct {
	registry *ns.Registry
	audit    *audit.Recorder
}

func NewZoneService(registry *ns.Registry, recorder *audit.Recorder) *ZoneService {
	return &ZoneService{
		registry: registry,
		audit:    recorder,
	}
}

func (z *ZoneService) AddZone(zone ns.Zone) error {
	return z.registry.AddZone(zone)
}

func (z *ZoneService) Zones(furnaceID string) []ns.Zone {
	return z.registry.Zones(furnaceID)
}

func (z *ZoneService) RePartition(furnaceID string, bindings []store.ZoneBinding) ([]ns.Zone, error) {
	zones := make([]ns.Zone, 0, len(bindings))
	for _, binding := range bindings {
		zones = append(zones, ns.Zone{
			ID:             binding.ZoneID,
			FurnaceID:      furnaceID,
			Index:          binding.Index,
			ThermocoupleID: binding.ThermocoupleID,
		})
	}
	if err := z.registry.RebindZones(furnaceID, zones); err != nil {
		return nil, err
	}
	message := fmt.Sprintf("zone mapping updated to %d bindings", len(zones))
	if err := z.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeZonePartition,
		Message:   message,
	}); err != nil {
		return nil, err
	}
	return zones, nil
}

func (z *ZoneService) ZoneForThermocouple(thermocoupleID string) (ns.Zone, bool) {
	return z.registry.ZoneForThermocouple(thermocoupleID)
}

func (z *ZoneService) Restore(furnaceID string) error {
	return z.registry.Restore(furnaceID)
}
