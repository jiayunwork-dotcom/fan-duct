package duct

// Sample 是管阻曲线上的一个离散点：流量与对应压降。
type Sample struct {
	Flow     float64 // m³/s
	Pressure float64 // Pa
}

// SampleCurve 在 [0, qMax] 上取 n+1 个等距流量点，计算管阻曲线并返回。
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
	return pts, nil
}

// SlopeAt 用中心差分估计管阻曲线在流量 q 处的斜率 d(Δp)/dQ。
// 步长取流量的 0.1% 且不小于 1e-9，中心差分误差为 O(h²)。
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

// PressureScaling 返回从流量 q0 到 q1 的压降倍数，用于"流量加倍压升约 4 倍"类检验。
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
