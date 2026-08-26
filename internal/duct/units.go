package duct

func FlowFromCFM(cfm float64) float64 {
	return cfm * 0.000471947443
}

func FlowToCFM(q float64) float64 {
	return q / 0.000471947443
}

func PressureFromMMH2O(mm float64) float64 {
	return mm * 9.80665
}

func PressureToMMH2O(pa float64) float64 {
	return pa / 9.80665
}

func VelocityFromFPM(fpm float64) float64 {
	return fpm * 0.00508
}

func VelocityToFPM(v float64) float64 {
	return v / 0.00508
}

func DiameterFromMM(mm float64) float64 {
	return mm * 0.001
}

func PressureFromInWG(in float64) float64 {
	return in * 249.089
}

func PressureToInWG(pa float64) float64 {
	return pa / 249.089
}

func FlowFromLps(lps float64) float64 {
	return lps * 0.001
}

func FlowToLps(q float64) float64 {
	return q * 1000
}
