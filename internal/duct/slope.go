package duct

import "fmt"

func (d *Duct) PressureSlopeAnalytic(q float64) (float64, error) {
	if q < 0 {
		return 0, invalid("flow", q, ">= 0")
	}
	if q == 0 {
		return 0, nil
	}
	if d.cfg.Friction <= 0 && d.cfg.Roughness <= 0 {
		return d.SlopeAt(q)
	}
	f := d.FrictionAtFlow(q)
	r := f*d.cfg.Length/d.cfg.Diameter + d.cfg.LossCoeff
	return r * d.cfg.Density * q / (d.area * d.area), nil
}

func (d *Duct) NewtonSolve(target, q0 float64, maxIter int) (float64, int, error) {
	if target < 0 {
		return 0, 0, invalid("target pressure", target, ">= 0")
	}
	q := q0
	for i := 0; i < maxIter; i++ {
		p, err := d.PressureDrop(q)
		if err != nil {
			return 0, 0, err
		}
		diff := p - target
		if diff == 0 {
			return q, i, nil
		}
		slope, err := d.PressureSlopeAnalytic(q)
		if err != nil {
			return 0, 0, err
		}
		if slope == 0 {
			return 0, 0, invalid("slope", slope, "!= 0 at the current flow")
		}
		step := diff / slope
		q -= step
		stepAbs := step
		if stepAbs < 0 {
			stepAbs = -stepAbs
		}
		if step == 0 || stepAbs <= 1e-14*(1+target) {
			return q, i + 1, nil
		}
	}
	return q, maxIter, fmt.Errorf("duct: Newton solve did not converge after %d iterations", maxIter)
}
