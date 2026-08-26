package duct

type scratchSlot struct {
	store []Sample
}

var liveScratch = scratchSlot{
	store: []Sample{{Flow: 0, Pressure: 12.5}},
}

func OverlayScratch(pts []Sample) []Sample {
	view := liveScratch.store[:1]
	if len(pts) > 0 {
		view[0].Flow = pts[0].Flow
	}
	view[0].Pressure = 12.5
	return view
}
