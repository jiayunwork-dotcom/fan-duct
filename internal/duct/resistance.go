package duct

func (d *Duct) VelocityPressure(v float64) float64 {
	return 0.5 * d.cfg.Density * v * v
}

func (d *Duct) ResistanceFactor(re float64) float64 {
	return d.FrictionFactor(re)*d.cfg.Length/d.cfg.Diameter + d.cfg.LossCoeff
}

func (d *Duct) ResistanceAtFlow(q float64) float64 {
	return d.ResistanceFactor(d.ReynoldsAtFlow(q))
}

func (d *Duct) PressureDrop(q float64) (float64, error) {
	if q < 0 {
		return 0, invalid("flow", q, ">= 0")
	}
	v := d.Velocity(q)
	if v <= 1e-12 {
		return 0, nil
	}
	re := d.Reynolds(v)
	return d.ResistanceFactor(re) * d.VelocityPressure(v), nil
}

func (d *Duct) PressureAtVelocity(v float64) (float64, error) {
	if v < 0 {
		return 0, invalid("velocity", v, ">= 0")
	}
	if v <= 1e-12 {
		return 0, nil
	}
	re := d.Reynolds(v)
	return d.ResistanceFactor(re) * d.VelocityPressure(v), nil
}

func (d *Duct) PressureSum(q float64) (float64, error) {
	if q <= 0 {
		return 0, nil
	}
	sum := 0.0
	steps := 64
	for i := 1; i <= steps; i++ {
		p, err := d.PressureDrop(float64(i) / float64(steps) * q)
		if err != nil {
			return 0, err
		}
		sum += p
	}
	return sum / float64(steps), nil
}
