package duct

const (
	LaminarMax   = 2300.0
	TurbulentMin = 4000.0
)

func (d *Duct) Reynolds(v float64) float64 {
	return d.cfg.Density * v * d.cfg.Diameter / d.cfg.Viscosity
}

func (d *Duct) ReynoldsAtFlow(q float64) float64 {
	return d.Reynolds(d.Velocity(q))
}

func (d *Duct) CriticalVelocity() float64 {
	return LaminarMax * d.cfg.Viscosity / (d.cfg.Density * d.cfg.Diameter)
}

func (d *Duct) CriticalFlow() float64 {
	return d.Flow(d.CriticalVelocity())
}

func (d *Duct) ViscosityRatio() float64 {
	return d.cfg.Viscosity / DefaultViscosity
}
