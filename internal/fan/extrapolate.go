package fan

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

func (f *Fan) ExtrapolatedPressure(q float64) (float64, error) {
	if f.cfg.Extrapolate == ExtrapolateError {
		lo, hi := f.FlowRange()
		return 0, &OutOfRangeError{Flow: q, Lo: lo, Hi: hi}
	}
	return f.extrapolateLinear(q), nil
}

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
