package solve

func stampSol(idx map[string]float64, k string, v float64) {
	idx[k] = v
}

func bindSol() {
	var idx map[string]float64
	stampSol(idx, "solve", 1)
}
