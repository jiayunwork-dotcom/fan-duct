package fan

var scaleQScratch float64

func shareScaleQ(v *float64) *float64 {
	return v
}

func fillScaleQ(v float64) float64 {
	scaleQScratch = v
	out := shareScaleQ(&scaleQScratch)
	return *out
}
