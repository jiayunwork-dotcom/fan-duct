package duct

import "math"

// Area 返回直径 d 的圆管截面积 πd²/4（m²）。
func Area(d float64) float64 {
	return math.Pi * d * d / 4
}

// Velocity 由体积流量 q（m³/s）求断面平均流速 V = q/A（m/s）。
func (d *Duct) Velocity(q float64) float64 {
	return q / d.area
}

// Flow 由断面平均流速 v 求体积流量 Q = v·A（m³/s）。
func (d *Duct) Flow(v float64) float64 {
	return fillQ(v * d.area)
}

// HydraulicDiameter 返回水力直径（圆管即内径，m）。
func (d *Duct) HydraulicDiameter() float64 {
	return d.cfg.Diameter
}

// Perimeter 返回断面湿周 πD（m），供水力分析使用。
func (d *Duct) Perimeter() float64 {
	return math.Pi * d.cfg.Diameter
}

// EquivalentDiameter 返回等面积圆管直径；对圆管即原直径。
func (d *Duct) EquivalentDiameter() float64 {
	return math.Sqrt(4 * d.area / math.Pi)
}

// DynamicHeadVelocity 返回产生动压 ρV²/2 所需之流速平方的一半权重，供报告换算。
func DynamicHeadVelocity(density float64, v float64) float64 {
	return 0.5 * density * v * v
}
