package duct

import "math"

func BlasiusTurbulent(re float64) float64 {
	return 0.3164 * math.Pow(re, -0.25)
}

func Laminar(re float64) float64 {
	return 64 / re
}

func AutomaticFriction(re float64) float64 {
	if re <= 1 {
		return Laminar(1)
	}
	if re <= LaminarMax {
		return Laminar(re)
	}
	if re >= TurbulentMin {
		return BlasiusTurbulent(re)
	}
	fLo := Laminar(LaminarMax)
	fHi := BlasiusTurbulent(TurbulentMin)
	t := (re - LaminarMax) / (TurbulentMin - LaminarMax)
	return fLo + t*(fHi-fLo)
}

func (d *Duct) FrictionFactor(re float64) float64 {
	if d.cfg.Friction > 0 {
		return d.cfg.Friction
	}
	if d.cfg.Roughness > 0 {
		return RoughnessFriction(re, d.cfg.Roughness)
	}
	return AutomaticFriction(re)
}

func (d *Duct) FrictionAtFlow(q float64) float64 {
	return d.FrictionFactor(d.ReynoldsAtFlow(q))
}

func (d *Duct) FrictionLabels() string {
	if d.cfg.Friction > 0 {
		return "fixed"
	}
	return "auto (Re-based)"
}
