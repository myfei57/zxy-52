package crane

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/grate"
)

type TipService struct {
	hopper *grate.HopperService
	audit  *audit.Recorder
}

func NewTipService(hopper *grate.HopperService, recorder *audit.Recorder) *TipService {
	return &TipService{
		hopper: hopper,
		audit:  recorder,
	}
}

func HopperID(furnaceID string) string {
	return "hopper-" + furnaceID
}

func (t *TipService) Tip(furnaceID string, grabID string, loadKG float64) error {
	hopperID := HopperID(furnaceID)
	if err := t.hopper.Align(hopperID); err != nil {
		return err
	}
	if err := t.hopper.Release(grabID, hopperID, loadKG); err != nil {
		return err
	}
	return t.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeCraneTipped,
		Message:   fmt.Sprintf("crane tipped %.1f kg into %s", loadKG, hopperID),
	})
}

func (t *TipService) SpillCount(furnaceID string) int {
	return t.hopper.SpillCount(HopperID(furnaceID))
}
