package ns

type FurnaceRef struct {
	ID     string
	Name   string
	Region string
}

type Zone struct {
	ID             string
	FurnaceID      string
	Index          int
	ThermocoupleID string
}
