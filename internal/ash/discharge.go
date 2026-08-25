package ash

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/grate"
)

type DischargeService struct {
	scraper *ScraperService
	valve   *grate.AshValveService
	audit   *audit.Recorder
}

func NewDischargeService(scraper *ScraperService, valve *grate.AshValveService, recorder *audit.Recorder) *DischargeService {
	return &DischargeService{
		scraper: scraper,
		valve:   valve,
		audit:   recorder,
	}
}

func (d *DischargeService) Discharge(furnaceID string) error {
	if err := d.scraper.Start(furnaceID); err != nil {
		return err
	}
	if err := d.valve.Open(furnaceID, d.scraper.IsRunning(furnaceID)); err != nil {
		return err
	}
	return d.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeFeedDropped,
		Message:   fmt.Sprintf("ash discharge completed for %s", furnaceID),
	})
}

func (d *DischargeService) CanDischarge(furnaceID string) bool {
	return d.scraper.IsRunning(furnaceID) && d.valve.IsOpen(furnaceID)
}
