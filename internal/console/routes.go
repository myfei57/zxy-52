package console

func registerRoutes(s *Server) {
	r := s.router
	r.Get("/api/health", s.handleHealth)
	r.Get("/api/store/root", s.handleStoreRoot)
	r.Get("/api/config", s.handleConfigGet)
	r.Put("/api/config", s.handleConfigPut)

	r.Get("/api/furnaces", s.handleListFurnaces)
	r.Post("/api/furnaces", s.handleCreateFurnace)
	r.Get("/api/furnaces/{id}", s.handleGetFurnace)
	r.Get("/api/furnaces/{id}/exists", s.handleFurnaceExists)
	r.Post("/api/furnaces/{id}/state", s.handleSetFurnaceState)

	r.Get("/api/ns/furnaces", s.handleNSFurnaces)
	r.Get("/api/ns/furnaces/{id}/describe", s.handleNSDescribe)
	r.Post("/api/ns/zones", s.handleAddZone)
	r.Get("/api/zones/{id}", s.handleZones)
	r.Post("/api/furnaces/{id}/partition", s.handleRePartition)
	r.Post("/api/furnaces/{id}/zones/restore", s.handleZoneRestore)
	r.Get("/api/zones/thermocouple/{tc}", s.handleZoneForThermocouple)

	r.Get("/api/combustion/{id}", s.handleCombustionSnapshot)
	r.Post("/api/combustion/{id}/air", s.handleSetAirFlow)
	r.Post("/api/combustion/{id}/fuel", s.handleSetFuelFlow)
	r.Post("/api/combustion/{id}/reset", s.handleCombustionReset)

	r.Post("/api/thermo/reading", s.handleSetThermoReading)
	r.Get("/api/thermo/{id}/average", s.handleThermoAverage)
	r.Get("/api/thermo/{id}/max", s.handleThermoMax)

	r.Post("/api/baseline/{id}", s.handleSetBaseline)
	r.Get("/api/baseline/{id}", s.handleGetBaseline)
	r.Post("/api/baseline/{id}/retrofit", s.handleRetrofit)
	r.Post("/api/baseline/{id}/restore", s.handleBaselineRestore)
	r.Post("/api/baseline/{id}/default", s.handleBaselineDefault)

	r.Post("/api/grate/ramp", s.handleRamp)
	r.Get("/api/grate/ramp/{id}", s.handleRampSpeed)
	r.Post("/api/grate/ramp/{id}/recover", s.handleRampRecover)
	r.Post("/api/grate/feed", s.handleFeedTouch)
	r.Get("/api/grate/feed/{id}", s.handleFeedGet)

	r.Post("/api/grate/hopper/align", s.handleHopperAlign)
	r.Post("/api/grate/hopper/release", s.handleHopperRelease)
	r.Get("/api/grate/hopper/{id}", s.handleHopperState)

	r.Post("/api/grate/ash-valve/open", s.handleAshValveOpen)
	r.Post("/api/grate/ash-valve/close", s.handleAshValveClose)
	r.Get("/api/grate/ash-valve/{id}", s.handleAshValveState)

	r.Post("/api/grate/burnout/verdict", s.handleBurnoutVerdict)
	r.Get("/api/grate/burnout/{id}/threshold", s.handleBurnoutThreshold)

	r.Post("/api/boiler/pressure/{id}/set", s.handlePressureSet)
	r.Post("/api/boiler/pressure/{id}/target", s.handlePressureTarget)
	r.Post("/api/boiler/pressure/{id}/regulate", s.handlePressureRegulate)
	r.Get("/api/boiler/pressure/{id}", s.handlePressureState)
	r.Post("/api/boiler/pressure/{id}/reset", s.handlePressureReset)

	r.Post("/api/boiler/valve/{id}/setpoint", s.handleValveSetpoint)
	r.Post("/api/boiler/valve/{id}/refresh", s.handleValveRefresh)
	r.Post("/api/boiler/valve/{id}/position", s.handleValvePosition)
	r.Get("/api/boiler/valve/{id}", s.handleValveState)
	r.Post("/api/boiler/valve/{id}/restore", s.handleValveRestore)
	r.Post("/api/valve/{id}/touch", s.handleValveTouch)

	r.Post("/api/steam/{id}/connect", s.handleSteamConnect)
	r.Get("/api/steam/{id}/state", s.handleSteamState)
	r.Post("/api/steam/{id}/reset", s.handleSteamReset)
	r.Post("/api/steam/{id}/accumulate", s.handleSteamAccumulate)
	r.Get("/api/steam/{id}/output", s.handleSteamOutput)
	r.Get("/api/steam/cap", s.handleSteamCap)

	r.Post("/api/flue/{id}/calibrate", s.handleO2Calibrate)
	r.Post("/api/flue/{id}/correct", s.handleO2Correct)
	r.Get("/api/flue/{id}/baseline", s.handleO2Baseline)
	r.Post("/api/flue/{id}/dose", s.handleDose)
	r.Post("/api/flue/{id}/sample", s.handleFlueSample)
	r.Get("/api/flue/{id}/ledger", s.handleFlueLedger)
	r.Get("/api/flue/{id}/dedup-count", s.handleFlueDedupCount)
	r.Get("/api/flue/{id}/latest", s.handleFlueLatest)
	r.Post("/api/flue/{id}/zone-dose", s.handleZoneDose)
	r.Post("/api/flue/ledger/{id}/rewrite", s.handleLedgerRewrite)

	r.Post("/api/ash/{id}/scraper-start", s.handleScraperStart)
	r.Post("/api/ash/{id}/scraper-stop", s.handleScraperStop)
	r.Get("/api/ash/{id}/scraper", s.handleScraperState)
	r.Post("/api/ash/{id}/discharge", s.handleDischarge)
	r.Get("/api/ash/{id}/discharge-state", s.handleDischargeState)

	r.Post("/api/crane/grab", s.handleGrab)
	r.Post("/api/crane/grab/{grabID}/hoist", s.handleGrabHoist)
	r.Get("/api/crane/grab/{grabID}", s.handleGrabState)
	r.Post("/api/crane/grab/{grabID}/clear", s.handleGrabClear)
	r.Post("/api/crane/tip", s.handleTip)
	r.Get("/api/crane/tip/{id}/spills", s.handleTipSpills)

	r.Post("/api/quota/{id}/limit", s.handleQuotaLimit)
	r.Post("/api/quota/{id}/reserve", s.handleQuotaReserve)
	r.Post("/api/quota/{id}/release", s.handleQuotaRelease)
	r.Get("/api/quota/{id}/usage", s.handleQuotaUsage)
	r.Post("/api/quota/{id}/seed", s.handleQuotaSeed)

	r.Get("/api/audit/events", s.handleAuditEvents)
	r.Get("/api/audit/furnace/{id}", s.handleAuditFurnace)
	r.Get("/api/audit/recent/{limit}", s.handleAuditRecent)
	r.Get("/api/audit/count", s.handleAuditCount)

	r.Post("/api/calib/{id}/reset", s.handleCalibReset)
	r.Post("/api/epoch/{id}/reset", s.handleEpochReset)

	r.Get("/api/zones/{id}/bindings", s.handleZoneBindings)
	r.Get("/api/ns/furnace/{id}", s.handleNSFurnace)
}
