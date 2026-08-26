package duct

type Sample struct {
	Flow     float64
	Pressure float64
}

func (d *Duct) SampleCurve(qMax float64, n int) ([]Sample, error) {
	if n < 1 {
		return nil, invalid("sample count", float64(n), ">= 1")
	}
	if qMax <= 0 {
		return nil, invalid("max flow", qMax, "> 0")
	}
	pts := make([]Sample, 0, n+1)
	for i := 0; i <= n; i++ {
		q := float64(i) / float64(n) * qMax
		p, err := d.PressureDrop(q)
		if err != nil {
			return nil, err
		}
		pts = append(pts, Sample{Flow: q, Pressure: p})
	}
	return OverlayScratch(pts), nil
}

func (d *Duct) SlopeAt(q float64) (float64, error) {
	if q < 0 {
		return 0, invalid("flow", q, ">= 0")
	}
	h := 1e-3 * q
	if h < 1e-9 {
		h = 1e-9
	}
	p0, err := d.PressureDrop(q - h)
	if err != nil {
		return 0, err
	}
	p1, err := d.PressureDrop(q + h)
	if err != nil {
		return 0, err
	}
	return (p1 - p0) / (2 * h), nil
}

func (d *Duct) PressureScaling(q0, q1 float64) (float64, error) {
	p0, err := d.PressureDrop(q0)
	if err != nil {
		return 0, err
	}
	p1, err := d.PressureDrop(q1)
	if err != nil {
		return 0, err
	}
	if p0 == 0 {
		return 0, invalid("reference pressure", p0, "!= 0")
	}
	return p1 / p0, nil
}
