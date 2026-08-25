package audit

import (
	"time"

	"github.com/google/uuid"
)

type Event struct {
	ID        string
	FurnaceID string
	EventType string
	Message   string
	At        string
}

func NewEvent(furnaceID string, eventType string, message string) Event {
	return Event{
		ID:        uuid.NewString(),
		FurnaceID: furnaceID,
		EventType: eventType,
		Message:   message,
		At:        time.Now().UTC().Format(time.RFC3339),
	}
}

const (
	TypeFeedPersisted      = "feed_persisted"
	TypeCombustionAdjust   = "combustion_adjust"
	TypeSteamValveMove     = "steam_valve_move"
	TypeHopperAligned      = "hopper_aligned"
	TypeGrabReleased       = "grab_released"
	TypeScraperStarted     = "scraper_started"
	TypeAshValveOpened     = "ash_valve_opened"
	TypeAshValveClosed     = "ash_valve_closed"
	TypeValveSetpointStore = "valve_setpoint_stored"
	TypeValveSetpointFresh = "valve_setpoint_refreshed"
	TypeValvePositionSet   = "valve_position_set"
	TypeGridConnected      = "grid_connected"
	TypeSteamAccumulated   = "steam_accumulated"
	TypeO2Calibrated       = "o2_calibrated"
	TypeO2Corrected        = "o2_corrected"
	TypeAmmoniaDosed       = "ammonia_dosed"
	TypeEmissionRecorded   = "emission_recorded"
	TypeEmissionDuplicate  = "emission_duplicate_skipped"
	TypeBurnoutPoor        = "burnout_poor"
	TypeFeedDropped        = "feed_dropped"
	TypeCraneTipped        = "crane_tipped"
	TypeGrabLoaded         = "grab_loaded"
	TypeQuotaReserved      = "quota_reserved"
	TypeZonePartition      = "zone_partition"
	TypeBaselineRetrofit   = "baseline_retrofit"
)
