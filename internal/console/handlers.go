package console

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"wastegen/internal/ns"
	"wastegen/internal/store"
)

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok", "service": "wastegen"})
}

func (s *Server) handleStoreRoot(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"root": s.store.Root()})
}

func (s *Server) handleConfigGet(w http.ResponseWriter, r *http.Request) {
	config, err := s.config.Load()
	if err != nil {
		config = s.config.Default()
	}
	writeJSON(w, http.StatusOK, config)
}

func (s *Server) handleConfigPut(w http.ResponseWriter, r *http.Request) {
	var req ConfigRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	config := store.PlantConfig{
		PlantName: req.PlantName,
		Region:    req.Region,
		Timezone:  req.Timezone,
	}
	if err := s.config.Save(config); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, config)
}

func (s *Server) handleListFurnaces(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.furnaces.List())
}

func (s *Server) handleCreateFurnace(w http.ResponseWriter, r *http.Request) {
	var req FurnaceRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	furnace, err := s.furnaces.Register(req.ID, req.Name, req.CapacityMW)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, furnace)
}

func (s *Server) handleGetFurnace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	furnace, ok := s.furnaces.Get(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "furnace not found"})
		return
	}
	writeJSON(w, http.StatusOK, furnace)
}

func (s *Server) handleFurnaceExists(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"exists": s.snapshot.Has("furnaces", id),
		"total":  s.furnaces.Count(),
	})
}

func (s *Server) handleSetFurnaceState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req StateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	furnace, err := s.furnaces.SetState(id, req.State)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, furnace)
}

func (s *Server) handleNSFurnaces(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, s.registry.Furnaces())
}

func (s *Server) handleNSDescribe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if _, ok := s.registry.Furnace(id); !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "furnace not in namespace"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"description": s.registry.Describe(id)})
}

func (s *Server) handleNSFurnace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ref, ok := s.registry.Furnace(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "furnace not in namespace"})
		return
	}
	writeJSON(w, http.StatusOK, ref)
}

func (s *Server) handleAddZone(w http.ResponseWriter, r *http.Request) {
	var req ZoneRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	zone := ns.Zone{
		ID:             req.ZoneID,
		FurnaceID:      req.FurnaceID,
		Index:          req.Index,
		ThermocoupleID: req.ThermocoupleID,
	}
	if err := s.zones.AddZone(zone); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, zone)
}

func (s *Server) handleZones(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, s.zones.Zones(id))
}

func (s *Server) handleZoneBindings(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	bindings, err := s.zoneStore.Bindings(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, bindings)
}

func (s *Server) handleRePartition(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req PartitionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	zones, err := s.zones.RePartition(id, req.Bindings)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, zones)
}

func (s *Server) handleZoneRestore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.zones.Restore(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s.zones.Zones(id))
}

func (s *Server) handleZoneForThermocouple(w http.ResponseWriter, r *http.Request) {
	tc := chi.URLParam(r, "tc")
	zone, ok := s.zones.ZoneForThermocouple(tc)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no zone mapping"})
		return
	}
	writeJSON(w, http.StatusOK, zone)
}

func (s *Server) handleCombustionSnapshot(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, s.combust.Snapshot(id))
}

func (s *Server) handleSetAirFlow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req FlowRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.combust.SetAirFlow(id, req.Value)
	writeJSON(w, http.StatusOK, s.combust.Snapshot(id))
}

func (s *Server) handleSetFuelFlow(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req FlowRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.combust.SetFuelFlow(id, req.Value)
	writeJSON(w, http.StatusOK, s.combust.Snapshot(id))
}

func (s *Server) handleCombustionReset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s.combust.Reset(id)
	writeJSON(w, http.StatusOK, s.combust.Snapshot(id))
}

func (s *Server) handleSetThermoReading(w http.ResponseWriter, r *http.Request) {
	var req ThermoRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.thermo.SetReading(req.ThermocoupleID, req.Temp)
	writeJSON(w, http.StatusOK, map[string]float64{"temp": s.thermo.Reading(req.ThermocoupleID)})
}

func (s *Server) handleThermoAverage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{"average": s.thermo.Average(id, s.registry)})
}

func (s *Server) handleThermoMax(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{"max": s.thermo.Max(id, s.registry)})
}

func (s *Server) handleSetBaseline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req RetrofitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.baselines.SetBaseline(id, req.Threshold)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleGetBaseline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, s.baselines.CurrentBaseline(id))
}

func (s *Server) handleRetrofit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req RetrofitRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.baselines.Retrofit(id, req.Threshold, req.Model)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleBaselineRestore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.baselines.Restore(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, s.baselines.CurrentBaseline(id))
}

func (s *Server) handleBaselineDefault(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := s.baselineSt.Default(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleRamp(w http.ResponseWriter, r *http.Request) {
	var req RampRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	speed := s.limits.ClampSpeed(req.Speed)
	feed := s.limits.ClampFeed(req.FeedKG)
	air, err := s.ramp.Ramp(req.FurnaceID, speed, feed)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"air_flow": air, "speed": speed, "feed_kg": feed})
}

func (s *Server) handleRampSpeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{"speed": s.ramp.Speed(id)})
}

func (s *Server) handleRampRecover(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	air, err := s.ramp.Recover(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"air_flow": air})
}

func (s *Server) handleFeedTouch(w http.ResponseWriter, r *http.Request) {
	var req FeedRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	feed := s.limits.ClampFeed(req.FeedKG)
	record, err := s.feedSt.Touch(req.FurnaceID, feed)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	if err := s.feeder.Drop(req.FurnaceID, feed); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"record":        record,
		"dropped_total": s.feeder.Total(req.FurnaceID),
	})
}

func (s *Server) handleFeedGet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, err := s.ramp.Feed(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"record": record,
		"exists": s.feedSt.Exists(id),
	})
}

func (s *Server) handleHopperAlign(w http.ResponseWriter, r *http.Request) {
	var req HopperRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.hopper.Align(req.HopperID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"aligned": s.hopper.Aligned(req.HopperID)})
}

func (s *Server) handleHopperRelease(w http.ResponseWriter, r *http.Request) {
	var req HopperRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.hopper.Release(req.GrabID, req.HopperID, req.LoadKG); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"loaded": s.hopper.Loaded(req.HopperID),
		"spills": s.hopper.SpillCount(req.HopperID),
	})
}

func (s *Server) handleHopperState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"aligned": s.hopper.Aligned(id),
		"loaded":  s.hopper.Loaded(id),
		"spills":  s.hopper.SpillCount(id),
	})
}

func (s *Server) handleAshValveOpen(w http.ResponseWriter, r *http.Request) {
	var req AshValveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.ashValve.Open(req.FurnaceID, s.scraper.IsRunning(req.FurnaceID)); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"open": s.ashValve.IsOpen(req.FurnaceID)})
}

func (s *Server) handleAshValveClose(w http.ResponseWriter, r *http.Request) {
	var req AshValveRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.ashValve.Close(req.FurnaceID); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"open": s.ashValve.IsOpen(req.FurnaceID)})
}

func (s *Server) handleAshValveState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]bool{"open": s.ashValve.IsOpen(id)})
}

func (s *Server) handleBurnoutVerdict(w http.ResponseWriter, r *http.Request) {
	var req BurnoutRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	pass, baseline, err := s.burnout.Verdict(req.FurnaceID, req.Rate)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"pass":      pass,
		"threshold": baseline.Threshold,
		"model":     baseline.Model,
	})
}

func (s *Server) handleBurnoutThreshold(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{"threshold": s.burnout.Threshold(id)})
}

func (s *Server) handlePressureSet(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req PressureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.pressure.SetPressure(id, req.Pressure)
	writeJSON(w, http.StatusOK, map[string]float64{"pressure": s.pressure.Pressure(id)})
}

func (s *Server) handlePressureTarget(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req PressureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	s.pressure.SetTarget(id, req.Target)
	writeJSON(w, http.StatusOK, map[string]float64{"target": s.pressure.Target(id)})
}

func (s *Server) handlePressureRegulate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req PressureRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	pressure, err := s.pressure.Regulate(id, req.Delta)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"pressure": pressure,
		"lifts":    s.pressure.SafetyValveLiftCount(id),
	})
}

func (s *Server) handlePressureState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"pressure": s.pressure.Pressure(id),
		"target":   s.pressure.Target(id),
		"lifts":    s.pressure.SafetyValveLiftCount(id),
	})
}

func (s *Server) handlePressureReset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s.pressure.Reset(id)
	writeJSON(w, http.StatusOK, map[string]string{"reset": id})
}

func (s *Server) handleValveSetpoint(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ValveSetpointRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.valve.StoreSetpoint(id, req.Setpoint)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleValveRefresh(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	setpoint, err := s.valve.RefreshSetpoint(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"setpoint": setpoint})
}

func (s *Server) handleValvePosition(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ValvePositionRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.valve.SetPosition(id, req.Position); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"position": s.valve.Position(id)})
}

func (s *Server) handleValveState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{
		"position": s.valve.Position(id),
		"setpoint": s.valve.Setpoint(id),
	})
}

func (s *Server) handleValveRestore(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.valve.Restore(id); err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"setpoint": s.valve.Setpoint(id)})
}

func (s *Server) handleValveTouch(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ValveSetpointRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.valveSt.Touch(id, req.Setpoint, 0); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"touched": id})
}

func (s *Server) handleSteamConnect(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.connect.Connect(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"connected": s.connect.IsConnected(id),
		"surges":    s.connect.SurgeCount(id),
	})
}

func (s *Server) handleSteamState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"connected": s.connect.IsConnected(id),
		"on_grid":   s.turbine.IsOnGrid(id),
		"trips":     s.turbine.TripCount(id),
	})
}

func (s *Server) handleSteamReset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	s.connect.Reset(id)
	writeJSON(w, http.StatusOK, map[string]string{"reset": id})
}

func (s *Server) handleSteamAccumulate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req AccumulateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	delta, epoch, err := s.accum.Accumulate(id, req.Reading)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"delta": delta, "epoch": epoch})
}

func (s *Server) handleSteamOutput(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	output, epoch, err := s.accum.Output(id)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"output": output, "epoch": epoch})
}

func (s *Server) handleSteamCap(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]float64{"cap": s.accum.Cap()})
}

func (s *Server) handleO2Calibrate(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CalibrateRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	record, err := s.o2.Calibrate(id, req.Baseline)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleO2Correct(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req CorrectRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	corrected, err := s.o2.Correct(id, req.Reading)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"corrected": corrected})
}

func (s *Server) handleO2Baseline(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]float64{"baseline": s.o2.Baseline(id)})
}

func (s *Server) handleDose(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req DoseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	dose, zoneID, err := s.dose.Dose(id, req.ThermocoupleID, req.NOx, req.O2)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"dose": dose, "zone_id": zoneID})
}

func (s *Server) handleFlueSample(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req SampleRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	timestamp := req.Timestamp
	if timestamp == "" {
		timestamp = time.Now().UTC().Format(time.RFC3339)
	}
	accepted, record, err := s.monitor.Sample(id, timestamp, req.O2, req.NOx, req.ZoneID)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"accepted": accepted, "record": record})
}

func (s *Server) handleFlueLedger(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	records, err := s.dedup.List(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, records)
}

func (s *Server) handleFlueLatest(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	record, ok := s.monitor.Latest(id)
	if !ok {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": "no sample recorded yet"})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleFlueDedupCount(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	count, err := s.dedup.Count(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]int{"count": count})
}

func (s *Server) handleZoneDose(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ZoneDoseRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]float64{"dose": s.dose.ZoneDose(id, req.ZoneIndex, req.NOx)})
}

func (s *Server) handleLedgerRewrite(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req LedgerRewriteRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.ledger.WriteAll(id, req.Records); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"rewritten": id})
}

func (s *Server) handleScraperStart(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.scraper.Start(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"running": s.scraper.IsRunning(id)})
}

func (s *Server) handleScraperStop(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.scraper.Stop(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"running": s.scraper.IsRunning(id)})
}

func (s *Server) handleScraperState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]bool{"running": s.scraper.IsRunning(id)})
}

func (s *Server) handleDischarge(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.discharge.Discharge(id); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"scraper":  s.scraper.IsRunning(id),
		"valve":    s.ashValve.IsOpen(id),
		"complete": s.discharge.CanDischarge(id),
	})
}

func (s *Server) handleDischargeState(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]any{
		"scraper":  s.scraper.IsRunning(id),
		"valve":    s.ashValve.IsOpen(id),
		"complete": s.discharge.CanDischarge(id),
	})
}

func (s *Server) handleGrab(w http.ResponseWriter, r *http.Request) {
	var req GrabRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.grab.Grab(req.GrabID, req.FurnaceID, req.LoadKG); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusCreated, map[string]any{
		"grab_id": req.GrabID,
		"state":   s.grab.State(req.GrabID),
		"load":    s.grab.Load(req.GrabID),
	})
}

func (s *Server) handleGrabHoist(w http.ResponseWriter, r *http.Request) {
	grabID := chi.URLParam(r, "grabID")
	s.grab.Hoist(grabID)
	writeJSON(w, http.StatusOK, map[string]string{"state": s.grab.State(grabID)})
}

func (s *Server) handleGrabState(w http.ResponseWriter, r *http.Request) {
	grabID := chi.URLParam(r, "grabID")
	writeJSON(w, http.StatusOK, map[string]any{
		"state": s.grab.State(grabID),
		"load":  s.grab.Load(grabID),
	})
}

func (s *Server) handleGrabClear(w http.ResponseWriter, r *http.Request) {
	grabID := chi.URLParam(r, "grabID")
	s.grab.Clear(grabID)
	writeJSON(w, http.StatusOK, map[string]string{"cleared": grabID})
}

func (s *Server) handleTip(w http.ResponseWriter, r *http.Request) {
	var req TipRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	if err := s.tip.Tip(req.FurnaceID, req.GrabID, req.LoadKG); err != nil {
		writeJSON(w, http.StatusConflict, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"tipped": req.FurnaceID,
		"spills": s.tip.SpillCount(req.FurnaceID),
	})
}

func (s *Server) handleTipSpills(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	writeJSON(w, http.StatusOK, map[string]int{"spills": s.tip.SpillCount(id)})
}

func (s *Server) handleQuotaLimit(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req QuotaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	date := req.Date
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	record, err := s.quota.SetLimit(id, req.Limit, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleQuotaReserve(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req QuotaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	date := req.Date
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	ok, err := s.quota.Reserve(id, req.KG, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]bool{"reserved": ok})
}

func (s *Server) handleQuotaRelease(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req QuotaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	date := req.Date
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	record, err := s.quota.Release(id, req.KG, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleQuotaUsage(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	date := time.Now().UTC().Format("2006-01-02")
	record, err := s.quota.Usage(id, date)
	if err != nil {
		writeJSON(w, http.StatusNotFound, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleQuotaSeed(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req QuotaRequest
	if err := decodeJSON(r, &req); err != nil {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}
	limit := req.Limit
	if limit <= 0 {
		limit = s.limits.DailyLimit()
	}
	date := req.Date
	if date == "" {
		date = time.Now().UTC().Format("2006-01-02")
	}
	record, err := s.quotaSt.Seed(id, limit, date)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, record)
}

func (s *Server) handleAuditEvents(w http.ResponseWriter, r *http.Request) {
	events, err := s.recorder.Events()
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleAuditFurnace(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	events, err := s.recorder.ForFurnace(id)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleAuditRecent(w http.ResponseWriter, r *http.Request) {
	limit := 100
	raw := chi.URLParam(r, "limit")
	if parsed, err := strconv.Atoi(raw); err == nil && parsed > 0 {
		limit = parsed
	}
	events, err := s.recorder.Recent(limit)
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, events)
}

func (s *Server) handleAuditCount(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]int{"count": s.recorder.Count()})
}

func (s *Server) handleCalibReset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.calib.Reset(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"reset": id})
}

func (s *Server) handleEpochReset(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := s.epoch.Reset(id); err != nil {
		writeJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"reset": id})
}
