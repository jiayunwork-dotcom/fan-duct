package duct

import "fmt"

// PressureSlopeAnalytic 返回管阻曲线在流量 q 处的解析斜率 d(Δp)/dQ。
//
// 固定摩阻系数时 Δp = R·ρ·Q²/(2A²)，其中 R = fL/D + ΣK 与 Q 无关，
// 因此 d(Δp)/dQ = R·ρ·Q/A²。自动摩阻时 R 依赖 Re，退化到数值差分。
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
	// Δp = (fL/D+K)·ρ·V²/2，V = Q/A ⇒ dΔp/dQ = (fL/D+K)·ρ·Q/A²。
	return r * d.cfg.Density * q / (d.area * d.area), nil
}

// NewtonSolve 用牛顿法在起点 q0 附近求管阻曲线等于 target 压降的流量。
// 固定摩阻时使用解析斜率，自动摩阻时使用数值斜率。
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
