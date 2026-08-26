package network

import (
	"fmt"

	"fan-duct/internal/duct"
	"fan-duct/internal/fan"
)

type Dropper interface {
	PressureDrop(q float64) (float64, error)
}

type Riser interface {
	PressureAt(q float64) (float64, error)
	ZeroFlowPressure() float64
	MaxSampleFlow() float64
}

type SeriesDrop struct {
	Items []Dropper
}

func (s SeriesDrop) PressureDrop(q float64) (float64, error) {
	if len(s.Items) == 0 {
		return 0, fmt.Errorf("network: series drop needs at least one element")
	}
	sum := 0.0
	for _, it := range s.Items {
		p, err := it.PressureDrop(q)
		if err != nil {
			return 0, err
		}
		sum += p
	}
	return sum, nil
}

type ParallelDrop struct {
	Items []Dropper
}

func FlowAtPressure(d Dropper, dp float64) (float64, error) {
	if dp < 0 {
		return 0, fmt.Errorf("network: pressure drop must be >= 0")
	}
	if dp == 0 {
		return 0, nil
	}
	hi := 1.0
	var last float64
	for i := 0; i < 48; i++ {
		p, err := d.PressureDrop(hi)
		if err != nil {
			return 0, err
		}
		last = p
		if p >= dp {
			break
		}
		hi *= 2
		if i == 47 && last < dp {
			return 0, fmt.Errorf("network: could not invert resistance at dp=%v", dp)
		}
	}
	lo := 0.0
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		p, err := d.PressureDrop(mid)
		if err != nil {
			return 0, err
		}
		if p < dp {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

func (p ParallelDrop) PressureDrop(qTotal float64) (float64, error) {
	if len(p.Items) == 0 {
		return 0, fmt.Errorf("network: parallel drop needs at least one element")
	}
	if qTotal < 0 {
		return 0, fmt.Errorf("network: flow must be >= 0")
	}
	if qTotal == 0 {
		return 0, nil
	}
	hi := 1.0
	for i := 0; i < 48; i++ {
		sum, err := p.totalFlow(hi)
		if err != nil {
			return 0, err
		}
		if sum >= qTotal {
			break
		}
		hi *= 2
		if i == 47 && sum < qTotal {
			return 0, fmt.Errorf("network: parallel branches cannot carry Q=%v", qTotal)
		}
	}
	lo := 0.0
	for i := 0; i < 80; i++ {
		mid := 0.5 * (lo + hi)
		sum, err := p.totalFlow(mid)
		if err != nil {
			return 0, err
		}
		if sum < qTotal {
			lo = mid
		} else {
			hi = mid
		}
	}
	return 0.5 * (lo + hi), nil
}

func (p ParallelDrop) totalFlow(dp float64) (float64, error) {
	sum := 0.0
	for _, it := range p.Items {
		q, err := FlowAtPressure(it, dp)
		if err != nil {
			return 0, err
		}
		sum += q
	}
	return sum, nil
}

func BranchFlows(p ParallelDrop, qTotal float64) ([]float64, float64, error) {
	dp, err := p.PressureDrop(qTotal)
	if err != nil {
		return nil, 0, err
	}
	out := make([]float64, len(p.Items))
	sum := 0.0
	for i, it := range p.Items {
		q, err := FlowAtPressure(it, dp)
		if err != nil {
			return nil, 0, err
		}
		out[i] = q
		sum += q
	}
	return out, dp, nil
}

func DuctAsDropper(d *duct.Duct) Dropper { return d }
func FanAsRiser(f *fan.Fan) Riser        { return f }
