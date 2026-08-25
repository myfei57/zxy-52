package steam

import (
	"fmt"

	"wastegen/internal/audit"
	"wastegen/internal/boiler"
)

type ConnectService struct {
	valve     *boiler.MainValveService
	turbine   *TurbineService
	audit     *audit.Recorder
	connected map[string]bool
}

func NewConnectService(valve *boiler.MainValveService, turbine *TurbineService, recorder *audit.Recorder) *ConnectService {
	return &ConnectService{
		valve:     valve,
		turbine:   turbine,
		audit:     recorder,
		connected: make(map[string]bool),
	}
}

func (c *ConnectService) Connect(furnaceID string) error {
	setpoint, err := c.valve.RefreshSetpoint(furnaceID)
	if err != nil {
		return err
	}
	if err := c.valve.SetPosition(furnaceID, setpoint); err != nil {
		return err
	}
	c.turbine.ConnectGrid(furnaceID)
	c.connected[furnaceID] = true
	if c.valve.Position(furnaceID) != c.valve.Setpoint(furnaceID) {
		c.turbine.Trip(furnaceID)
		return nil
	}
	return c.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeGridConnected,
		Message:   fmt.Sprintf("turbine for %s connected to grid at setpoint %.1f", furnaceID, setpoint),
	})
}

func (c *ConnectService) IsConnected(furnaceID string) bool {
	return c.connected[furnaceID]
}

func (c *ConnectService) SurgeCount(furnaceID string) int {
	return c.turbine.TripCount(furnaceID)
}

func (c *ConnectService) Reset(furnaceID string) {
	c.connected[furnaceID] = false
	c.turbine.Reset(furnaceID)
}
