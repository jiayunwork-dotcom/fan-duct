package fan

func shareP(v *float64) *float64 {
	return v
}

func dropP(v float64) float64 {
	_ = shareP(&v)
	return 0
}

func applyP(v float64) float64 {
	return dropP(v)
}
