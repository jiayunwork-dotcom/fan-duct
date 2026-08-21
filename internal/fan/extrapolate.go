package fan

// extrapolate 处理流量 q 超出样本范围 [lo, hi] 时的求值。
//
// ExtrapolateLinear：低端用前两点斜率、高端用后两点斜率线性外推；
// ExtrapolateQuadratic：沿用整条曲线的二次多项式；
// ExtrapolateError：返回 *OutOfRangeError。
func (f *Fan) extrapolate(q, lo, hi float64) (float64, error) {
	switch f.cfg.Extrapolate {
	case ExtrapolateLinear:
		return f.extrapolateLinear(q), nil
	case ExtrapolateQuadratic:
		return f.a*q*q + f.b*q + f.c, nil
	default:
		return 0, &OutOfRangeError{Flow: q, Lo: lo, Hi: hi}
	}
}

// extrapolateLinear 按边界段斜率线性外推。
func (f *Fan) extrapolateLinear(q float64) float64 {
	pts := f.cfg.Points
	if q < pts[0].Flow {
		p0, p1 := pts[0], pts[1]
		slope, _ := LinearlyFit(p0, p1)
		return p0.Pressure + slope*(q-p0.Flow)
	}
	n := len(pts)
	p0, p1 := pts[n-2], pts[n-1]
	slope, _ := LinearlyFit(p0, p1)
	return p1.Pressure + slope*(q-p1.Flow)
}

// ExtrapolatedPressure 在显式外推策略下计算 q 处压升，并在 q < 0 时报错。
func (f *Fan) ExtrapolatedPressure(q float64) (float64, error) {
	if f.cfg.Extrapolate == ExtrapolateError {
		lo, hi := f.FlowRange()
		return 0, &OutOfRangeError{Flow: q, Lo: lo, Hi: hi}
	}
	return f.extrapolateLinear(q), nil
}

// ZeroCrossing 返回分段线性曲线上压升首次降到给定目标值 p 的流量。
// 若压升始终大于目标或样本不足，返回 false。
func (f *Fan) ZeroCrossing(target float64) (float64, bool) {
	pts := f.cfg.Points
	if len(pts) < 2 || pts[0].Pressure <= target {
		return 0, false
	}
	for i := 1; i < len(pts); i++ {
		if pts[i].Pressure <= target {
			p0, p1 := pts[i-1], pts[i]
			t := (target - p0.Pressure) / (p1.Pressure - p0.Pressure)
			return p0.Flow + t*(p1.Flow-p0.Flow), true
		}
	}
	return 0, false
}
