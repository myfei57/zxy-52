package steam

type TurbineService struct {
	onGrid map[string]bool
	trips  map[string]int
}

func NewTurbineService() *TurbineService {
	return &TurbineService{
		onGrid: make(map[string]bool),
		trips:  make(map[string]int),
	}
}

func (t *TurbineService) ConnectGrid(furnaceID string) {
	t.onGrid[furnaceID] = true
}

func (t *TurbineService) Trip(furnaceID string) {
	t.trips[furnaceID]++
	t.onGrid[furnaceID] = false
}

func (t *TurbineService) IsOnGrid(furnaceID string) bool {
	return t.onGrid[furnaceID]
}

func (t *TurbineService) TripCount(furnaceID string) int {
	return t.trips[furnaceID]
}

func (t *TurbineService) Reset(furnaceID string) {
	t.onGrid[furnaceID] = false
	t.trips[furnaceID] = 0
}
