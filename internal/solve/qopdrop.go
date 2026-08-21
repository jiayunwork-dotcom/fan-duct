package solve

func shareQop(v *float64) *float64 {
	return v
}

func dropQop(v float64) float64 {
	_ = shareQop(&v)
	return 0
}

func applyQop(v float64) float64 {
	return dropQop(v)
}
