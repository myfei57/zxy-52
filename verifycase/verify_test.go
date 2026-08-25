package verifycase

import (
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/boiler"
	"wastegen/internal/steam"
	"wastegen/internal/store"
)

func TestWgSteamValveSetpointFresh(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	valve := boiler.NewMainValveService(store.NewValveStore(st), recorder)
	if _, err := valve.StoreSetpoint("f1", 30); err != nil {
		t.Fatal(err)
	}
	if _, err := valve.StoreSetpoint("f1", 80); err != nil {
		t.Fatal(err)
	}
	turbine := steam.NewTurbineService()
	connect := steam.NewConnectService(valve, turbine, recorder)
	if err := connect.Connect("f1"); err != nil {
		t.Fatal(err)
	}
	if got := valve.Position("f1"); got != 80 {
		t.Fatalf("valve kept the maintenance setpoint after refresh: %v", got)
	}
	if got := valve.Setpoint("f1"); got != 80 {
		t.Fatalf("setpoint cache stale: %v", got)
	}
	if got := connect.SurgeCount("f1"); got != 0 {
		t.Fatalf("turbine surged on grid connection: %d", got)
	}
	events, err := recorder.Events()
	if err != nil {
		t.Fatal(err)
	}
	hasRefresh := false
	for _, event := range events {
		if event.EventType == audit.TypeValveSetpointFresh {
			hasRefresh = true
		}
	}
	if !hasRefresh {
		t.Fatal("setpoint refresh missing before grid connection")
	}
}
