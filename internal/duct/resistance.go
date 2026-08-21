package duct

// VelocityPressure 返回流速 v 对应的动压 ρv²/2（Pa）。
func (d *Duct) VelocityPressure(v float64) float64 {
	return 0.5 * d.cfg.Density * v * v
}

// ResistanceFactor 返回给定 Re 下的总阻力系数 R = f·L/D + ΣK。
// 沿程摩阻按 FrictionFactor 取值，因此固定 f 时 R 与 Re 无关，
// 自动 f 时 R 随 Re 变化。
func (d *Duct) ResistanceFactor(re float64) float64 {
	return d.FrictionFactor(re)*d.cfg.Length/d.cfg.Diameter + d.cfg.LossCoeff
}

// ResistanceAtFlow 返回给定流量下的总阻力系数。
func (d *Duct) ResistanceAtFlow(q float64) float64 {
	return d.ResistanceFactor(d.ReynoldsAtFlow(q))
}

// PressureDrop 由体积流量 q 求管阻压降
// Δp = (f·L/D + ΣK)·ρV²/2，其中 V = q/A。
//
// q = 0 时 V = 0，动压为零，Δp 返回 0；q < 0 报错。
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

// PressureAtVelocity 直接由流速求管阻压降，等价于把 q = v·A 代入 PressureDrop。
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

// PressureSum 返回给定流量区间 [0, q] 内的平均压降，供曲线面积类指标使用。
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
