package flue

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/furnace"
	"wastegen/internal/ns"
)

type DoseService struct {
	o2    *O2Service
	zones *furnace.ZoneService
	audit *audit.Recorder
	byTC  map[string]ns.Zone
}

func NewDoseService(o2 *O2Service, zones *furnace.ZoneService, recorder *audit.Recorder) *DoseService {
	return &DoseService{
		o2:    o2,
		zones: zones,
		audit: recorder,
		byTC:  make(map[string]ns.Zone),
	}
}

func (d *DoseService) Dose(furnaceID string, thermocoupleID string, noxReading float64, o2Reading float64) (float64, string, error) {
	zone, ok := d.byTC[thermocoupleID]
	if !ok {
		zone, ok = d.zones.ZoneForThermocouple(thermocoupleID)
		if ok {
			d.byTC[thermocoupleID] = zone
		}
	}
	if !ok {
		return 0, "", fmt.Errorf("thermocouple %s has no zone mapping", thermocoupleID)
	}
	corrected, err := d.o2.Correct(furnaceID, o2Reading)
	if err != nil {
		return 0, zone.ID, err
	}
	dose := noxReading*0.8 + corrected*0.2 + float64(zone.Index)*0.5
	if err := d.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeAmmoniaDosed,
		Message:   fmt.Sprintf("ammonia dose %.2f kg/h for zone %s", dose, zone.ID),
	}); err != nil {
		return 0, zone.ID, err
	}
	return dose, zone.ID, nil
}

func (d *DoseService) ZoneDose(furnaceID string, zoneIndex int, noxReading float64) float64 {
	return noxReading*0.7 + float64(zoneIndex)*0.4
}
