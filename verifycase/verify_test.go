package verifycase

import (
	"math"
	"testing"

	"wastegen/internal/audit"
	"wastegen/internal/flue"
	"wastegen/internal/furnace"
	"wastegen/internal/ns"
	"wastegen/internal/store"
)

func TestWgFurnaceZoneMappingFresh(t *testing.T) {
	root := t.TempDir()
	st, err := store.Open(root)
	if err != nil {
		t.Fatal(err)
	}
	recorder := audit.NewRecorder(store.NewAuditStore(st))
	registry := ns.NewRegistry(store.NewZoneStore(st))
	zones := furnace.NewZoneService(registry, recorder)
	if err := zones.AddZone(ns.Zone{ID: "z1", FurnaceID: "f1", Index: 0, ThermocoupleID: "tc-a"}); err != nil {
		t.Fatal(err)
	}
	if err := zones.AddZone(ns.Zone{ID: "z2", FurnaceID: "f1", Index: 1, ThermocoupleID: "tc-b"}); err != nil {
		t.Fatal(err)
	}
	o2 := flue.NewO2Service(store.NewCalibStore(st), recorder)
	if _, err := o2.Calibrate("f1", 20.5); err != nil {
		t.Fatal(err)
	}
	dose := flue.NewDoseService(o2, zones, recorder)
	if _, _, err := dose.Dose("f1", "tc-a", 120, 21.0); err != nil {
		t.Fatal(err)
	}
	bindings := []store.ZoneBinding{
		{ZoneID: "z1", ThermocoupleID: "tc-b", Index: 0},
		{ZoneID: "z2", ThermocoupleID: "tc-a", Index: 1},
	}
	if _, err := zones.RePartition("f1", bindings); err != nil {
		t.Fatal(err)
	}
	zone, ok := zones.ZoneForThermocouple("tc-a")
	if !ok {
		t.Fatal("tc-a unmapped after re-partition")
	}
	if zone.ID != "z2" {
		t.Fatalf("zone mapping not refreshed: tc-a still %s", zone.ID)
	}
	got, zoneID, err := dose.Dose("f1", "tc-a", 120, 21.0)
	if err != nil {
		t.Fatal(err)
	}
	if zoneID != "z2" {
		t.Fatalf("ammonia dose targeted stale zone %s", zoneID)
	}
	expected := 120*0.8 + 20.7*0.2 + 1*0.5
	if math.Abs(got-expected) > 1e-6 {
		t.Fatalf("dose mismatch: %v", got)
	}
	fresh := flue.NewDoseService(o2, furnace.NewZoneService(registry, recorder), recorder)
	if _, zoneID, err := fresh.Dose("f1", "tc-a", 120, 21.0); err != nil || zoneID != "z2" {
		t.Fatalf("restart dose stale: zone=%s err=%v", zoneID, err)
	}
}
