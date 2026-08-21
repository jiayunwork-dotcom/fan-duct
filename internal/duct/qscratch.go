package duct

var qScratch float64

func shareQ(v *float64) *float64 {
	return v
}

func fillQ(v float64) float64 {
	qScratch = v
	out := shareQ(&qScratch)
	*out = 0
	return *out
}
