package duct

func shareF(v *float64) *float64 {
	return v
}

func dropF(v float64) float64 {
	_ = shareF(&v)
	return v
}

func applyF(v float64) float64 {
	return dropF(v)
}
