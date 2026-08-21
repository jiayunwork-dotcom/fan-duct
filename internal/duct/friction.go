package duct

import "math"

// BlasiusTurbulent 返回光滑管湍流摩阻系数 0.3164·Re^(-1/4)。
// 适用于 Re ≥ 4000，是 Blasius 湍流关联式的常用形式。
func BlasiusTurbulent(re float64) float64 {
	return 0.3164 * math.Pow(re, -0.25)
}

// Laminar 返回层流摩阻系数 64/Re，适用于 Re ≤ 2300。
func Laminar(re float64) float64 {
	return 64 / re
}

// AutomaticFriction 按雷诺数选择摩阻系数：
// Re ≤ 2300 用 64/Re；Re ≥ 4000 用 Blasius；中间过渡区在两端值间线性插值，
// 使函数连续无跳变。Re 过小（< 1）时返回层流下限值，仅供数值保护。
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

// FrictionFactor 返回给定 Re 下的摩阻系数：配置了固定 f 则恒用固定值，
// 配置了相对粗糙度则用 Swamee–Jain，否则按 AutomaticFriction（Blasius）。
func (d *Duct) FrictionFactor(re float64) float64 {
	if d.cfg.Friction > 0 {
		return d.cfg.Friction
	}
	if d.cfg.Roughness > 0 {
		return RoughnessFriction(re, d.cfg.Roughness)
	}
	return AutomaticFriction(re)
}

// FrictionAtFlow 返回给定流量下的摩阻系数（固定或自动）。
func (d *Duct) FrictionAtFlow(q float64) float64 {
	return d.FrictionFactor(d.ReynoldsAtFlow(q))
}

// FrictionLabels 返回摩阻取值的模式说明，供报告展示。
func (d *Duct) FrictionLabels() string {
	if d.cfg.Friction > 0 {
		return "fixed"
	}
	return "auto (Re-based)"
}
