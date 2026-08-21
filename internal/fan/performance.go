package fan

import "fmt"

// ShaftPower 估算风机在流量 q、效率 eta 下的轴功率 P = q·Δp/η（W）。
// eta 是 0~1 的总效率；流体功率 q·Δp 是理想功率。
func (f *Fan) ShaftPower(q, eta float64) (float64, error) {
	if eta <= 0 || eta > 1 {
		return 0, configErr("efficiency", "must be in (0, 1]")
	}
	p, err := f.PowerAt(q)
	if err != nil {
		return 0, err
	}
	return p / eta, nil
}

// EfficiencyAt 返回给定流量下的总效率，按效率样本线性插值。
// 未提供效率样本时返回默认 0.7，并返回 false 表示采用默认值。
func (f *Fan) EfficiencyAt(q float64) (float64, bool) {
	eff := f.cfg.Efficiency
	if len(eff) == 0 {
		return 0.7, false
	}
	pts := f.cfg.Points
	if q <= pts[0].Flow {
		return eff[0], true
	}
	for i := 1; i < len(pts); i++ {
		if q <= pts[i].Flow {
			t := (q - pts[i-1].Flow) / (pts[i].Flow - pts[i-1].Flow)
			return eff[i-1] + t*(eff[i]-eff[i-1]), true
		}
	}
	return eff[len(eff)-1], true
}

// PowerCurve 返回 [0, qMax] 上的流体功率曲线采样，供报告或图表。
func (f *Fan) PowerCurve(qMax float64, n int) ([]PowerSample, error) {
	if n < 1 {
		return nil, configErr("sample count", "must be >= 1")
	}
	if qMax <= 0 {
		return nil, configErr("max flow", "must be > 0")
	}
	out := make([]PowerSample, 0, n+1)
	for i := 0; i <= n; i++ {
		q := float64(i) / float64(n) * qMax
		p, err := f.PressureAt(q)
		if err != nil {
			return nil, err
		}
		out = append(out, PowerSample{Flow: q, Power: q * p})
	}
	return out, nil
}

// PowerSample 是风机功率曲线上的一个点。
type PowerSample struct {
	Flow  float64
	Power float64
}

// PeakPowerFlow 在 [0, qMax] 上扫描寻找流体功率最大的流量。
func (f *Fan) PeakPowerFlow(qMax float64, steps int) (float64, error) {
	pts, err := f.PowerCurve(qMax, steps)
	if err != nil {
		return 0, err
	}
	best := pts[0]
	for _, p := range pts[1:] {
		if p.Power > best.Power {
			best = p
		}
	}
	return best.Flow, nil
}

// SpecificPower 返回单位流量功率 P/Q = Δp，用于与风机压升一致性核对。
func (f *Fan) SpecificPower(q float64) (float64, error) {
	p, err := f.PressureAt(q)
	if err != nil {
		return 0, err
	}
	return p, nil
}

// String 返回效率模式的说明。
func (f *Fan) String() string {
	return fmt.Sprintf("fan curve: %d points, %s fit, %s extrapolation",
		len(f.cfg.Points), f.cfg.Fit.String(), f.cfg.Extrapolate.String())
}
