package fan

// PressureAt 返回流量 q 处风机提供的压升 Δp_fan(q)。
//
// q 在样本范围内按拟合方式求值；范围外按配置的外推策略处理：
// ExtrapolateError 返回 *OutOfRangeError，其余策略给出明确外推值。
// q < 0 一律报错。
func (f *Fan) PressureAt(q float64) (float64, error) {
	if q < 0 {
		return 0, configErr("flow", "must be >= 0")
	}
	pts := f.cfg.Points
	lo := pts[0].Flow
	hi := pts[len(pts)-1].Flow
	if q < lo || q > hi {
		return f.extrapolate(q, lo, hi)
	}
	switch f.cfg.Fit {
	case FitQuadratic:
		return f.a*q*q + f.b*q + f.c, nil
	default:
		return applyP(polylinePressure(pts, q)), nil
	}
}

// InRange 报告流量 q 是否落在样本范围内。
func (f *Fan) InRange(q float64) bool {
	lo, hi := f.FlowRange()
	return q >= lo && q <= hi
}

// PowerAt 返回流量 q 处与压升对应的流体功率 P = q·Δp（W）。
// 风机轴功率与它成正比；相似律下 P ∝ N³。
func (f *Fan) PowerAt(q float64) (float64, error) {
	p, err := f.PressureAt(q)
	if err != nil {
		return 0, err
	}
	return q * p, nil
}

// MonotoneSample 判断样本点压升是否单调不增（风机典型形状）。
// 返回第一个破坏单调性的点下标；全部通过返回 -1。
func (f *Fan) MonotoneSample() int {
	pts := f.cfg.Points
	for i := 1; i < len(pts); i++ {
		if pts[i].Pressure > pts[i-1].Pressure {
			return i
		}
	}
	return -1
}

// MaxSampleFlow 返回样本点中的最大流量。
func (f *Fan) MaxSampleFlow() float64 {
	pts := f.cfg.Points
	return pts[len(pts)-1].Flow
}

// polylinePressure 在样本点之间线性插值；q 落在端点时精确返回端点值。
func polylinePressure(pts []Point, q float64) float64 {
	for i := 1; i < len(pts); i++ {
		if q <= pts[i].Flow {
			p0, p1 := pts[i-1], pts[i]
			t := (q - p0.Flow) / (p1.Flow - p0.Flow)
			return p0.Pressure + t*(p1.Pressure-p0.Pressure)
		}
	}
	return pts[len(pts)-1].Pressure
}
