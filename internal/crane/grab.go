package crane

import (
	"fmt"

	"wastegen/internal/audit"
)

type GrabService struct {
	loads  map[string]float64
	states map[string]string
	audit  *audit.Recorder
}

func NewGrabService(recorder *audit.Recorder) *GrabService {
	return &GrabService{
		loads:  make(map[string]float64),
		states: make(map[string]string),
		audit:  recorder,
	}
}

func (g *GrabService) Grab(grabID string, furnaceID string, loadKG float64) error {
	g.loads[grabID] = loadKG
	g.states[grabID] = "loaded"
	return g.audit.Record(audit.Event{
		FurnaceID: furnaceID,
		EventType: audit.TypeGrabLoaded,
		Message:   fmt.Sprintf("grab %s loaded %.1f kg for %s", grabID, loadKG, furnaceID),
	})
}

func (g *GrabService) Hoist(grabID string) {
	g.states[grabID] = "hoisted"
}

func (g *GrabService) State(grabID string) string {
	return g.states[grabID]
}

func (g *GrabService) Load(grabID string) float64 {
	return g.loads[grabID]
}

func (g *GrabService) Clear(grabID string) {
	delete(g.loads, grabID)
	delete(g.states, grabID)
}
