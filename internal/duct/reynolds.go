package duct

// 层流与湍流雷诺数分界。
const (
	LaminarMax   = 2300.0
	TurbulentMin = 4000.0
)

// Reynolds 由断面平均流速 v 求雷诺数 Re = ρVD/μ。
func (d *Duct) Reynolds(v float64) float64 {
	return d.cfg.Density * v * d.cfg.Diameter / d.cfg.Viscosity
}

// ReynoldsAtFlow 由体积流量 q 直接求雷诺数。
func (d *Duct) ReynoldsAtFlow(q float64) float64 {
	return d.Reynolds(d.Velocity(q))
}

// CriticalVelocity 返回层流/湍流分界对应的临界流速（m/s）。
func (d *Duct) CriticalVelocity() float64 {
	return LaminarMax * d.cfg.Viscosity / (d.cfg.Density * d.cfg.Diameter)
}

// CriticalFlow 返回临界流量（m³/s）。
func (d *Duct) CriticalFlow() float64 {
	return d.Flow(d.CriticalVelocity())
}

// ViscosityRatio 返回黏度相对 20°C 空气默认值的倍数，供输出参考。
func (d *Duct) ViscosityRatio() float64 {
	return d.cfg.Viscosity / DefaultViscosity
}
