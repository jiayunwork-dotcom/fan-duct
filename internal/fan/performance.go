package fan

import "fmt"

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

type PowerSample struct {
	Flow  float64
	Power float64
}

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

func (f *Fan) SpecificPower(q float64) (float64, error) {
	p, err := f.PressureAt(q)
	if err != nil {
		return 0, err
	}
	return p, nil
}

func (f *Fan) String() string {
	return fmt.Sprintf("fan curve: %d points, %s fit, %s extrapolation",
		len(f.cfg.Points), f.cfg.Fit.String(), f.cfg.Extrapolate.String())
}
