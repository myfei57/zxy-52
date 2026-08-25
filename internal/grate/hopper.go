package grate

import (
	"fmt"

	"wastegen/internal/audit"
)

type HopperState struct {
	aligned map[string]bool
	spills  map[string]int
	loads   map[string]float64
}

func NewHopperState() *HopperState {
	return &HopperState{
		aligned: make(map[string]bool),
		spills:  make(map[string]int),
		loads:   make(map[string]float64),
	}
}

func (h *HopperState) MarkAligned(hopperID string) {
	h.aligned[hopperID] = true
}

func (h *HopperState) Aligned(hopperID string) bool {
	return h.aligned[hopperID]
}

func (h *HopperState) AddSpill(hopperID string) {
	h.spills[hopperID]++
}

func (h *HopperState) SpillCount(hopperID string) int {
	return h.spills[hopperID]
}

func (h *HopperState) AddLoad(hopperID string, kg float64) {
	h.loads[hopperID] += kg
}

func (h *HopperState) Loaded(hopperID string) float64 {
	return h.loads[hopperID]
}

type HopperService struct {
	state *HopperState
	audit *audit.Recorder
}

func NewHopperService(state *HopperState, recorder *audit.Recorder) *HopperService {
	return &HopperService{
		state: state,
		audit: recorder,
	}
}

func (h *HopperService) Align(hopperID string) error {
	h.state.MarkAligned(hopperID)
	return h.audit.Record(audit.Event{
		FurnaceID: hopperID,
		EventType: audit.TypeHopperAligned,
		Message:   fmt.Sprintf("hopper %s aligned", hopperID),
	})
}

func (h *HopperService) Release(grabID string, hopperID string, loadKG float64) error {
	if !h.state.Aligned(hopperID) {
		h.state.AddSpill(hopperID)
		return fmt.Errorf("hopper %s not aligned, grab %s cannot release", hopperID, grabID)
	}
	h.state.AddLoad(hopperID, loadKG)
	return h.audit.Record(audit.Event{
		FurnaceID: hopperID,
		EventType: audit.TypeGrabReleased,
		Message:   fmt.Sprintf("grab %s released %.1f kg into hopper %s", grabID, loadKG, hopperID),
	})
}

func (h *HopperService) SpillCount(hopperID string) int {
	return h.state.SpillCount(hopperID)
}

func (h *HopperService) Loaded(hopperID string) float64 {
	return h.state.Loaded(hopperID)
}

func (h *HopperService) Aligned(hopperID string) bool {
	return h.state.Aligned(hopperID)
}
