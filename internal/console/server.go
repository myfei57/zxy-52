package console

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"wastegen/internal/ash"
	"wastegen/internal/audit"
	"wastegen/internal/boiler"
	"wastegen/internal/crane"
	"wastegen/internal/flue"
	"wastegen/internal/furnace"
	"wastegen/internal/grate"
	"wastegen/internal/ns"
	"wastegen/internal/quota"
	"wastegen/internal/steam"
	"wastegen/internal/store"
)

type Server struct {
	store     *store.Store
	registry  *ns.Registry
	furnaces  *furnace.FurnaceService
	zones     *furnace.ZoneService
	combust   *furnace.CombustService
	thermo    *furnace.ThermocoupleService
	baselines *furnace.BaselineService
	ramp      *grate.RampService
	hopper    *grate.HopperService
	feeder    *grate.FeederService
	ashValve  *grate.AshValveService
	burnout   *grate.BurnoutService
	pressure  *boiler.PressureService
	valve     *boiler.MainValveService
	connect   *steam.ConnectService
	accum     *steam.AccumulatorService
	turbine   *steam.TurbineService
	o2        *flue.O2Service
	dose      *flue.DoseService
	dedup     *flue.DedupService
	monitor   *flue.MonitorService
	scraper   *ash.ScraperService
	discharge *ash.DischargeService
	grab      *crane.GrabService
	tip       *crane.TipService
	quota     *quota.QuotaService
	recorder  *audit.Recorder
	config    *store.ConfigStore
	snapshot  *store.SnapshotStore
	zoneStore *store.ZoneStore
	calib     *store.CalibStore
	epoch     *store.EpochStore
	valveSt   *store.ValveStore
	quotaSt   *store.QuotaStore
	baselineSt *store.BaselineStore
	ledger    *store.LedgerStore
	feedSt    *store.FeedStore
	limits    quota.Limits
	router    chi.Router
}

func NewServer(st *store.Store) *Server {
	auditStore := store.NewAuditStore(st)
	recorder := audit.NewRecorder(auditStore)
	zoneStore := store.NewZoneStore(st)
	registry := ns.NewRegistry(zoneStore)
	snapshot := store.NewSnapshotStore(st)
	furnaces := furnace.NewFurnaceService(registry, snapshot)
	combust := furnace.NewCombustService()
	thermo := furnace.NewThermocoupleService()
	baselineStore := store.NewBaselineStore(st)
	baselines := furnace.NewBaselineService(baselineStore, recorder)
	zones := furnace.NewZoneService(registry, recorder)
	feedStore := store.NewFeedStore(st)
	rampState := grate.NewRampState()
	ramp := grate.NewRampService(rampState, feedStore, combust, recorder)
	hopperState := grate.NewHopperState()
	hopper := grate.NewHopperService(hopperState, recorder)
	feederState := grate.NewFeederState()
	feeder := grate.NewFeederService(feederState, recorder)
	ashValveState := grate.NewAshValveState()
	ashValve := grate.NewAshValveService(ashValveState, recorder)
	burnout := grate.NewBurnoutService(baselines, recorder)
	valveStore := store.NewValveStore(st)
	mainValve := boiler.NewMainValveService(valveStore, recorder)
	pressure := boiler.NewPressureService(combust, mainValve, recorder)
	turbine := steam.NewTurbineService()
	connect := steam.NewConnectService(mainValve, turbine, recorder)
	epochStore := store.NewEpochStore(st)
	accum := steam.NewAccumulatorService(epochStore, recorder)
	calibStore := store.NewCalibStore(st)
	o2 := flue.NewO2Service(calibStore, recorder)
	dose := flue.NewDoseService(o2, zones, recorder)
	ledger := store.NewLedgerStore(st)
	dedup := flue.NewDedupService(ledger, recorder)
	monitor := flue.NewMonitorService(dedup, recorder)
	scraperState := ash.NewScraperState()
	scraper := ash.NewScraperService(scraperState, recorder)
	discharge := ash.NewDischargeService(scraper, ashValve, recorder)
	grab := crane.NewGrabService(recorder)
	tip := crane.NewTipService(hopper, recorder)
	quotaStore := store.NewQuotaStore(st)
	quotaSvc := quota.NewQuotaService(quotaStore, recorder)
	config := store.NewConfigStore(st)

	s := &Server{
		store:      st,
		registry:   registry,
		furnaces:   furnaces,
		zones:      zones,
		combust:    combust,
		thermo:     thermo,
		baselines:  baselines,
		ramp:       ramp,
		hopper:     hopper,
		feeder:     feeder,
		ashValve:   ashValve,
		burnout:    burnout,
		pressure:   pressure,
		valve:      mainValve,
		connect:    connect,
		accum:      accum,
		turbine:    turbine,
		o2:         o2,
		dose:       dose,
		dedup:      dedup,
		monitor:    monitor,
		scraper:    scraper,
		discharge:  discharge,
		grab:       grab,
		tip:        tip,
		quota:      quotaSvc,
		recorder:   recorder,
		config:     config,
		snapshot:   snapshot,
		zoneStore:  zoneStore,
		calib:      calibStore,
		epoch:      epochStore,
		valveSt:    valveStore,
		quotaSt:    quotaStore,
		baselineSt: baselineStore,
		ledger:     ledger,
		feedSt:     feedStore,
		limits:     quota.DefaultLimits(),
	}
	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	registerRoutes(s)
	return s
}

func (s *Server) Router() chi.Router {
	return s.router
}

func (s *Server) Start(addr string) error {
	log.Printf("wastegen DCS console listening on %s", addr)
	return http.ListenAndServe(addr, s.router)
}
