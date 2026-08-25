package furnace

import (
	"fmt"
	"sort"

	"wastegen/internal/ns"
	"wastegen/internal/store"
)

type Furnace struct {
	ID         string
	Name       string
	CapacityMW float64
	State      string
}

type FurnaceService struct {
	registry *ns.Registry
	snapshot *store.SnapshotStore
	furnaces map[string]Furnace
}

func NewFurnaceService(registry *ns.Registry, snapshot *store.SnapshotStore) *FurnaceService {
	return &FurnaceService{
		registry: registry,
		snapshot: snapshot,
		furnaces: make(map[string]Furnace),
	}
}

func (f *FurnaceService) Register(id string, name string, capacityMW float64) (Furnace, error) {
	furnace := Furnace{
		ID:         id,
		Name:       name,
		CapacityMW: capacityMW,
		State:      "idle",
	}
	if err := f.snapshot.Save("furnaces", id, furnace); err != nil {
		return Furnace{}, err
	}
	f.furnaces[id] = furnace
	if err := f.registry.RegisterFurnace(ns.FurnaceRef{
		ID:     id,
		Name:   name,
		Region: "region-" + id,
	}); err != nil {
		return Furnace{}, err
	}
	return furnace, nil
}

func (f *FurnaceService) Restore(id string) error {
	var furnace Furnace
	if err := f.snapshot.Load("furnaces", id, &furnace); err != nil {
		return err
	}
	f.furnaces[id] = furnace
	return f.registry.RegisterFurnace(ns.FurnaceRef{
		ID:     furnace.ID,
		Name:   furnace.Name,
		Region: "region-" + furnace.ID,
	})
}

func (f *FurnaceService) Get(id string) (Furnace, bool) {
	furnace, ok := f.furnaces[id]
	return furnace, ok
}

func (f *FurnaceService) List() []Furnace {
	furnaces := make([]Furnace, 0, len(f.furnaces))
	for _, furnace := range f.furnaces {
		furnaces = append(furnaces, furnace)
	}
	sort.Slice(furnaces, func(i, j int) bool {
		return furnaces[i].ID < furnaces[j].ID
	})
	return furnaces
}

func (f *FurnaceService) SetState(id string, state string) (Furnace, error) {
	furnace, ok := f.furnaces[id]
	if !ok {
		return Furnace{}, fmt.Errorf("furnace %s not found", id)
	}
	furnace.State = state
	f.furnaces[id] = furnace
	if err := f.snapshot.Save("furnaces", id, furnace); err != nil {
		return Furnace{}, err
	}
	return furnace, nil
}

func (f *FurnaceService) Count() int {
	return len(f.furnaces)
}
