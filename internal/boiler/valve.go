package boiler

import (
	"fmt"
	"sync"
	"time"

	"wastegen/internal/audit"
	"wastegen/internal/store"
)

type MainValveService struct {
	mu        sync.Mutex
	valveStore *store.ValveStore
	positions map[string]float64
	setpoints map[string]float64
	audit     *audit.Recorder
}

func NewMainValveService(valveStore *store.ValveStore, recorder *audit.Recorder) *MainValveService {
	return &MainValveService{
		valveStore: valveStore,
		positions:  make(map[string]float64),
		setpoints:  make(map[string]float64),
		audit:      recorder,
	}
}

func (v *MainValveService) StoreSetpoint(furnaceID string, setpoint float64) (store.ValveRecord, error) {
	v.mu.Lock()
	record := store.ValveRecord{
		FurnaceID: furnaceID,
		Setpoint:  setpoint,
		Position:  v.positions[furnaceID],
		UpdatedAt: nowRFC3339(),
	}
	v.mu.Unlock()
	if err := v.valveStore.Save(record); err != nil {
		return record, err
	}
	if err := v.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeValveSetpointStore,
		Message:   fmt.Sprintf("main steam valve setpoint stored at %.1f", setpoint),
	}); err != nil {
		return record, err
	}
	return record, nil
}

func (v *MainValveService) RefreshSetpoint(furnaceID string) (float64, error) {
	record, err := v.valveStore.Load(furnaceID)
	if err != nil {
		return 0, err
	}
	if err := v.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeValveSetpointFresh,
		Message:   fmt.Sprintf("main steam valve setpoint refreshed to %.1f", record.Setpoint),
	}); err != nil {
		return 0, err
	}
	return record.Setpoint, nil
}

func (v *MainValveService) SetPosition(furnaceID string, position float64) error {
	v.mu.Lock()
	v.positions[furnaceID] = position
	v.mu.Unlock()
	return v.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeValvePositionSet,
		Message:   fmt.Sprintf("main steam valve position set to %.1f", position),
	})
}

func (v *MainValveService) Move(furnaceID string, delta float64) error {
	v.mu.Lock()
	v.positions[furnaceID] += delta
	v.mu.Unlock()
	return nil
}

func (v *MainValveService) Position(furnaceID string) float64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.positions[furnaceID]
}

func (v *MainValveService) Setpoint(furnaceID string) float64 {
	v.mu.Lock()
	defer v.mu.Unlock()
	return v.setpoints[furnaceID]
}

func (v *MainValveService) Restore(furnaceID string) error {
	record, err := v.valveStore.Load(furnaceID)
	if err != nil {
		return err
	}
	v.mu.Lock()
	v.setpoints[furnaceID] = record.Setpoint
	v.positions[furnaceID] = record.Position
	v.mu.Unlock()
	return nil
}

func nowRFC3339() string {
	return time.Now().UTC().Format(time.RFC3339)
}
