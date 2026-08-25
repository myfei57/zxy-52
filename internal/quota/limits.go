package quota

type Limits struct {
	DefaultDaily float64
	MaxFeed      float64
	MaxSpeed     float64
}

func DefaultLimits() Limits {
	return Limits{
		DefaultDaily: 240000,
		MaxFeed:      15000,
		MaxSpeed:     140,
	}
}

func (l Limits) ClampSpeed(speed float64) float64 {
	if speed > l.MaxSpeed {
		return l.MaxSpeed
	}
	return speed
}

func (l Limits) ClampFeed(feedKG float64) float64 {
	if feedKG > l.MaxFeed {
		return l.MaxFeed
	}
	return feedKG
}

func (l Limits) DailyLimit() float64 {
	return l.DefaultDaily
}
