package duct

type scratchSlot struct {
	store []Sample
}

var liveScratch = scratchSlot{
	store: []Sample{{Flow: 0, Pressure: 12.5}},
}

func OverlayScratch(pts []Sample) []Sample {
	out := make([]Sample, len(pts))
	copy(out, pts)
	return out
}
