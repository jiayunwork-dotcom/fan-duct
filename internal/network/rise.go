package network

import (
	"fmt"

	"fan-duct/internal/fan"
)

type SeriesRise struct {
	Items []Riser
}

func (s SeriesRise) PressureAt(q float64) (float64, error) {
	if len(s.Items) == 0 {
		return 0, fmt.Errorf("network: series rise needs at least one fan")
	}
	sum := 0.0
	for _, it := range s.Items {
		p, err := it.PressureAt(q)
		if err != nil {
			return 0, err
		}
		sum += p
	}
	return sum, nil
}

func (s SeriesRise) ZeroFlowPressure() float64 {
	sum := 0.0
	for _, it := range s.Items {
		sum += it.ZeroFlowPressure()
	}
	return sum
}

func (s SeriesRise) MaxSampleFlow() float64 {
	if len(s.Items) == 0 {
		return 0
	}
	m := s.Items[0].MaxSampleFlow()
	for _, it := range s.Items[1:] {
		if x := it.MaxSampleFlow(); x < m {
			m = x
		}
	}
	return m
}

type ParallelRise struct {
	Items []Riser
}

func FanFlowAtPressure(f Riser, dp float64) (float64, error) {
	shut := f.ZeroFlowPressure()
	if dp < 0 {
		return 0, fmt.Errorf("network: fan pressure must be >= 0")
	}
	if dp >= shut {
		return 0, nil
	}
	lo, hi := 0.0, f.MaxSampleFlow()
	if hi <= 0 {
		return 0, fmt.Errorf("network: fan sample range is empty")
	}
	pHi, err := f.PressureAt(hi)
	if err != nil {
		return 0, err
	}
	if pHi > dp {
		expand := hi
		for i := 0; i < 16; i++ {
			expand *= 1.5
			p, err := f.PressureAt(expand)
			if err != nil {
				break
			}
			hi = expand
			pHi = p
			if pHi <= dp {
				break
			}
		}
		if pHi > dp {
			return hi, nil
		}
	}
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		p, err := f.PressureAt(mid)
		if err != nil {
			return 0, err
		}
		if p > dp {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

func (p ParallelRise) PressureAt(qTotal float64) (float64, error) {
	if len(p.Items) == 0 {
		return 0, fmt.Errorf("network: parallel rise needs at least one fan")
	}
	if qTotal < 0 {
		return 0, fmt.Errorf("network: flow must be >= 0")
	}
	if qTotal == 0 {
		return p.ZeroFlowPressure(), nil
	}
	lo := 0.0
	hi := p.ZeroFlowPressure()
	if hi <= 0 {
		return 0, fmt.Errorf("network: combined shutoff is not positive")
	}
	qLo, err := p.totalFlow(hi)
	if err != nil {
		return 0, err
	}
	if qLo > qTotal {
		return hi, nil
	}
	qHi, err := p.totalFlow(lo)
	if err != nil {
		return 0, err
	}
	if qHi < qTotal {
		return 0, fmt.Errorf("network: parallel fans cannot deliver Q=%v (free delivery %v)", qTotal, qHi)
	}
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		q, err := p.totalFlow(mid)
		if err != nil {
			return 0, err
		}
		if q > qTotal {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

func (p ParallelRise) totalFlow(dp float64) (float64, error) {
	sum := 0.0
	for _, it := range p.Items {
		q, err := FanFlowAtPressure(it, dp)
		if err != nil {
			return 0, err
		}
		sum += q
	}
	return sum, nil
}

func (p ParallelRise) ZeroFlowPressure() float64 {
	m := 0.0
	for _, it := range p.Items {
		if s := it.ZeroFlowPressure(); s > m {
			m = s
		}
	}
	return m
}

func (p ParallelRise) MaxSampleFlow() float64 {
	sum := 0.0
	for _, it := range p.Items {
		sum += it.MaxSampleFlow()
	}
	return sum
}

func SampleCombinedFan(r Riser, n int) (*fan.Fan, error) {
	if n < 4 {
		n = 16
	}
	qMax := r.MaxSampleFlow()
	if qMax <= 0 {
		return nil, fmt.Errorf("network: cannot sample an empty fan range")
	}
	pts := make([]fan.Point, 0, n+1)
	for i := 0; i <= n; i++ {
		q := float64(i) / float64(n) * qMax
		dp, err := r.PressureAt(q)
		if err != nil {
			return nil, err
		}
		if dp < 0 {
			dp = 0
		}
		if i > 0 && q <= pts[len(pts)-1].Flow {
			continue
		}
		pts = append(pts, fan.Point{Flow: q, Pressure: dp})
	}
	if pts[0].Flow != 0 {
		pts[0].Flow = 0
	}
	return fan.New(fan.FanConfig{Points: pts, Fit: fan.FitPolyline, Extrapolate: fan.ExtrapolateError})
}
