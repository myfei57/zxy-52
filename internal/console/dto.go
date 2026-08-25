package console

import (
	"encoding/json"
	"net/http"

	"wastegen/internal/store"
)

type FurnaceRequest struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	CapacityMW float64 `json:"capacity_mw"`
}

type StateRequest struct {
	State string `json:"state"`
}

type ZoneRequest struct {
	FurnaceID      string `json:"furnace_id"`
	ZoneID         string `json:"zone_id"`
	ThermocoupleID string `json:"thermocouple_id"`
	Index          int    `json:"index"`
}

type PartitionRequest struct {
	Bindings []store.ZoneBinding `json:"bindings"`
}

type RampRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Speed     float64 `json:"speed"`
	FeedKG    float64 `json:"feed_kg"`
}

type FeedRequest struct {
	FurnaceID string  `json:"furnace_id"`
	FeedKG    float64 `json:"feed_kg"`
}

type HopperRequest struct {
	HopperID string  `json:"hopper_id"`
	GrabID   string  `json:"grab_id"`
	LoadKG   float64 `json:"load_kg"`
}

type AshValveRequest struct {
	FurnaceID string `json:"furnace_id"`
}

type BurnoutRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Rate      float64 `json:"rate"`
}

type PressureRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Pressure  float64 `json:"pressure"`
	Target    float64 `json:"target"`
	Delta     float64 `json:"delta"`
}

type ValveSetpointRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Setpoint  float64 `json:"setpoint"`
}

type ValvePositionRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Position  float64 `json:"position"`
}

type AccumulateRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Reading   float64 `json:"reading"`
}

type CalibrateRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Baseline  float64 `json:"baseline"`
}

type CorrectRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Reading   float64 `json:"reading"`
}

type DoseRequest struct {
	FurnaceID      string  `json:"furnace_id"`
	ThermocoupleID string  `json:"thermocouple_id"`
	NOx            float64 `json:"nox"`
	O2             float64 `json:"o2"`
}

type ZoneDoseRequest struct {
	ZoneIndex int     `json:"zone_index"`
	NOx       float64 `json:"nox"`
}

type SampleRequest struct {
	FurnaceID  string  `json:"furnace_id"`
	Timestamp  string  `json:"timestamp"`
	ZoneID     string  `json:"zone_id"`
	O2         float64 `json:"o2"`
	NOx        float64 `json:"nox"`
	ThermocoupleID string `json:"thermocouple_id"`
}

type LedgerRewriteRequest struct {
	FurnaceID string                   `json:"furnace_id"`
	Records   []store.EmissionRecord `json:"records"`
}

type GrabRequest struct {
	GrabID    string  `json:"grab_id"`
	FurnaceID string  `json:"furnace_id"`
	LoadKG    float64 `json:"load_kg"`
}

type TipRequest struct {
	FurnaceID string  `json:"furnace_id"`
	GrabID    string  `json:"grab_id"`
	LoadKG    float64 `json:"load_kg"`
}

type QuotaRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Limit     float64 `json:"limit"`
	KG        float64 `json:"kg"`
	Date      string  `json:"date"`
}

type ConfigRequest struct {
	PlantName string `json:"plant_name"`
	Region    string `json:"region"`
	Timezone  string `json:"timezone"`
}

type ThermoRequest struct {
	ThermocoupleID string  `json:"thermocouple_id"`
	Temp           float64 `json:"temp"`
}

type RetrofitRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Threshold float64 `json:"threshold"`
	Model     string  `json:"model"`
}

type FlowRequest struct {
	FurnaceID string  `json:"furnace_id"`
	Value     float64 `json:"value"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func decodeJSON(r *http.Request, value any) error {
	return json.NewDecoder(r.Body).Decode(value)
}
