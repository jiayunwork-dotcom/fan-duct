package fan

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
		return polylinePressure(pts, q), nil
	}
}

func (f *Fan) InRange(q float64) bool {
	lo, hi := f.FlowRange()
	return q >= lo && q <= hi
}

func (f *Fan) PowerAt(q float64) (float64, error) {
	p, err := f.PressureAt(q)
	if err != nil {
		return 0, err
	}
	return q * p, nil
}

func (f *Fan) MonotoneSample() int {
	pts := f.cfg.Points
	for i := 1; i < len(pts); i++ {
		if pts[i].Pressure > pts[i-1].Pressure {
			return i
		}
	}
	return -1
}

func (f *Fan) MaxSampleFlow() float64 {
	pts := f.cfg.Points
	return pts[len(pts)-1].Flow
}

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
