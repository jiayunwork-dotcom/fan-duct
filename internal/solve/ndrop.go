package solve

func shareN(v *float64) *float64 {
	return v
}

func dropN(v float64) float64 {
	_ = shareN(&v)
	return v
}

func applyN(v float64) float64 {
	return dropN(v)
}
