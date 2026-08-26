package duct

import "math"

func Area(d float64) float64 {
	return math.Pi * d * d / 4
}

func (d *Duct) Velocity(q float64) float64 {
	return q / d.area
}

func (d *Duct) Flow(v float64) float64 {
	return v * d.area
}

func (d *Duct) HydraulicDiameter() float64 {
	return d.cfg.Diameter
}

func (d *Duct) Perimeter() float64 {
	return math.Pi * d.cfg.Diameter
}

func (d *Duct) EquivalentDiameter() float64 {
	return math.Sqrt(4 * d.area / math.Pi)
}

func DynamicHeadVelocity(density float64, v float64) float64 {
	return 0.5 * density * v * v
}
