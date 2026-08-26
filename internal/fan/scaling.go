package fan

func RPMRatio(n1, n2 float64) (float64, error) {
	if n1 <= 0 {
		return 0, configErr("speed N1", "must be > 0")
	}
	if n2 <= 0 {
		return 0, configErr("speed N2", "must be > 0")
	}
	return n2 / n1, nil
}

func (f *Fan) ScaledByRatio(r float64) (*Fan, error) {
	if r <= 0 {
		return nil, configErr("speed ratio", "must be > 0")
	}
	pts := make([]Point, len(f.cfg.Points))
	for i, p := range f.cfg.Points {
		pts[i] = Point{Flow: p.Flow * r, Pressure: p.Pressure * r * r}
	}
	cfg := f.cfg
	cfg.Points = pts
	return New(cfg)
}

func ScalePoint(q, dp, r float64) (q2, dp2, p2 float64) {
	q2 = q * r
	dp2 = dp * r * r
	p2 = q2 * dp2
	return q2, dp2, p2
}

func AffinityChecks(q1, dp1, q2, dp2 float64) (ratioDP, ratioP float64) {
	ratioDP = dp2 / dp1
	ratioP = (q2 * dp2) / (q1 * dp1)
	return ratioDP, ratioP
}
