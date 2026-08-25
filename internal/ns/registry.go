package ns

import (
	"fmt"
	"sort"

	"wastegen/internal/store"
)

type Registry struct {
	furnaces  map[string]FurnaceRef
	zoneMap   map[string][]Zone
	byTC      map[string]Zone
	zoneStore *store.ZoneStore
}

func NewRegistry(zoneStore *store.ZoneStore) *Registry {
	return &Registry{
		furnaces:  make(map[string]FurnaceRef),
		zoneMap:   make(map[string][]Zone),
		byTC:      make(map[string]Zone),
		zoneStore: zoneStore,
	}
}

func (r *Registry) RegisterFurnace(ref FurnaceRef) error {
	r.furnaces[ref.ID] = ref
	return r.persistZones(ref.ID)
}

func (r *Registry) Furnaces() []FurnaceRef {
	refs := make([]FurnaceRef, 0, len(r.furnaces))
	for _, ref := range r.furnaces {
		refs = append(refs, ref)
	}
	sort.Slice(refs, func(i, j int) bool {
		return refs[i].ID < refs[j].ID
	})
	return refs
}

func (r *Registry) Furnace(id string) (FurnaceRef, bool) {
	ref, ok := r.furnaces[id]
	return ref, ok
}

func (r *Registry) AddZone(zone Zone) error {
	zones := r.zoneMap[zone.FurnaceID]
	replaced := false
	for i, existing := range zones {
		if existing.ThermocoupleID == zone.ThermocoupleID {
			zones[i] = zone
			replaced = true
			break
		}
	}
	if !replaced {
		zones = append(zones, zone)
	}
	sort.Slice(zones, func(i, j int) bool {
		return zones[i].Index < zones[j].Index
	})
	r.zoneMap[zone.FurnaceID] = zones
	r.rebuildByTC(zone.FurnaceID)
	return r.persistZones(zone.FurnaceID)
}

func (r *Registry) Zones(furnaceID string) []Zone {
	zones := r.zoneMap[furnaceID]
	copied := append([]Zone(nil), zones...)
	sort.Slice(copied, func(i, j int) bool {
		return copied[i].Index < copied[j].Index
	})
	return copied
}

func (r *Registry) ZoneForThermocouple(thermocoupleID string) (Zone, bool) {
	zone, ok := r.byTC[thermocoupleID]
	return zone, ok
}

func (r *Registry) RebindZones(furnaceID string, zones []Zone) error {
	sort.Slice(zones, func(i, j int) bool {
		return zones[i].Index < zones[j].Index
	})
	r.zoneMap[furnaceID] = append([]Zone(nil), zones...)
	r.rebuildByTC(furnaceID)
	return r.persistZones(furnaceID)
}

func (r *Registry) Restore(furnaceID string) error {
	snapshot, err := r.zoneStore.Load(furnaceID)
	if err != nil {
		return err
	}
	zones := make([]Zone, 0, len(snapshot.Bindings))
	for _, binding := range snapshot.Bindings {
		zones = append(zones, Zone{
			ID:             binding.ZoneID,
			FurnaceID:      furnaceID,
			Index:          binding.Index,
			ThermocoupleID: binding.ThermocoupleID,
		})
	}
	r.zoneMap[furnaceID] = zones
	r.rebuildByTC(furnaceID)
	return nil
}

func (r *Registry) rebuildByTC(furnaceID string) {
	for tc, zone := range r.byTC {
		if zone.FurnaceID == furnaceID {
			delete(r.byTC, tc)
		}
	}
	for _, zone := range r.zoneMap[furnaceID] {
		r.byTC[zone.ThermocoupleID] = zone
	}
}

func (r *Registry) persistZones(furnaceID string) error {
	zones := r.zoneMap[furnaceID]
	bindings := make([]store.ZoneBinding, 0, len(zones))
	for _, zone := range zones {
		bindings = append(bindings, store.ZoneBinding{
			ZoneID:         zone.ID,
			ThermocoupleID: zone.ThermocoupleID,
			Index:          zone.Index,
		})
	}
	return r.zoneStore.Touch(furnaceID, bindings)
}

func (r *Registry) Describe(furnaceID string) string {
	return fmt.Sprintf("furnace %s zones=%d", furnaceID, len(r.zoneMap[furnaceID]))
}
