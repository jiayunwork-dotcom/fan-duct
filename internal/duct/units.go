package duct

// FlowFromCFM 把立方英尺每分钟（CFM）换算为 m³/s。
func FlowFromCFM(cfm float64) float64 {
	return cfm * 0.000471947443
}

// FlowToCFM 把 m³/s 换算为立方英尺每分钟。
func FlowToCFM(q float64) float64 {
	return q / 0.000471947443
}

// PressureFromMMH2O 把毫米水柱（mmH2O）换算为 Pa。按 g=9.80665 m/s²、水密度 1000 kg/m³。
func PressureFromMMH2O(mm float64) float64 {
	return mm * 9.80665
}

// PressureToMMH2O 把 Pa 换算为毫米水柱。
func PressureToMMH2O(pa float64) float64 {
	return pa / 9.80665
}

// VelocityFromFPM 把英尺每分钟（fpm）换算为 m/s。
func VelocityFromFPM(fpm float64) float64 {
	return fpm * 0.00508
}

// VelocityToFPM 把 m/s 换算为英尺每分钟。
func VelocityToFPM(v float64) float64 {
	return v / 0.00508
}

// DiameterFromMM 把毫米换算为米。
func DiameterFromMM(mm float64) float64 {
	return mm * 0.001
}

// PressureFromInWG 把英寸水柱（inWC）换算为 Pa（1 inWC = 249.089 Pa）。
func PressureFromInWG(in float64) float64 {
	return in * 249.089
}

// PressureToInWG 把 Pa 换算为英寸水柱。
func PressureToInWG(pa float64) float64 {
	return pa / 249.089
}

// FlowFromLps 把升每秒（L/s）换算为 m³/s。
func FlowFromLps(lps float64) float64 {
	return lps * 0.001
}

// FlowToLps 把 m³/s 换算为升每秒。
func FlowToLps(q float64) float64 {
	return q * 1000
}
