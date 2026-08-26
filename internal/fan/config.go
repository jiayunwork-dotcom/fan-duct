package fan

import (
	"fmt"
	"math"
)

type Point struct {
	Flow     float64
	Pressure float64
}

type FitMode int

const (
	FitPolyline FitMode = iota
	FitQuadratic
)

func ParseFit(s string) (FitMode, bool) {
	switch s {
	case "", "polyline", "linear":
		return FitPolyline, true
	case "quadratic":
		return FitQuadratic, true
	}
	return FitPolyline, false
}

func (m FitMode) String() string {
	switch m {
	case FitQuadratic:
		return "quadratic"
	default:
		return "polyline"
	}
}

type ExtrapolateMode int

const (
	ExtrapolateError ExtrapolateMode = iota
	ExtrapolateLinear
	ExtrapolateQuadratic
)

func ParseExtrapolate(s string) (ExtrapolateMode, bool) {
	switch s {
	case "", "error":
		return ExtrapolateError, true
	case "linear":
		return ExtrapolateLinear, true
	case "quadratic":
		return ExtrapolateQuadratic, true
	}
	return ExtrapolateError, false
}

func (m ExtrapolateMode) String() string {
	switch m {
	case ExtrapolateLinear:
		return "linear"
	case ExtrapolateQuadratic:
		return "quadratic"
	default:
		return "error"
	}
}

type FanConfig struct {
	Points      []Point
	Efficiency  []float64
	Fit         FitMode
	Extrapolate ExtrapolateMode
}

func (c FanConfig) Validate() error {
	if len(c.Points) < 2 {
		return configErr("points", "at least two sample points required")
	}
	first := c.Points[0]
	if first.Flow < -1e-12 || first.Flow > 1e-12 {
		return configErr("first point", "flow of the first sample point must be 0")
	}
	if first.Pressure < 0 || math.IsNaN(first.Pressure) || math.IsInf(first.Pressure, 0) {
		return configErr("shutoff pressure", "must be a finite non-negative number")
	}
	for i := 1; i < len(c.Points); i++ {
		p := c.Points[i]
		if math.IsNaN(p.Flow) || math.IsInf(p.Flow, 0) || math.IsNaN(p.Pressure) || math.IsInf(p.Pressure, 0) {
			return configErr("point", "flow and pressure must be finite")
		}
		if p.Flow <= c.Points[i-1].Flow {
			return configErr("point", "flow values must be strictly increasing")
		}
		if p.Pressure < 0 {
			return configErr("point", "pressure must be non-negative")
		}
	}
	if len(c.Efficiency) > 0 {
		if len(c.Efficiency) != len(c.Points) {
			return configErr("efficiency", "must have the same length as points")
		}
		for i, e := range c.Efficiency {
			if math.IsNaN(e) || math.IsInf(e, 0) || e <= 0 || e > 1 {
				return configErr("efficiency", fmt.Sprintf("point %d must be in (0, 1]", i))
			}
		}
	}
	return nil
}

type Fan struct {
	cfg FanConfig
	a   float64
	b   float64
	c   float64
}

func New(cfg FanConfig) (*Fan, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	a, b, c := FitQuadraticCoeffs(cfg.Points)
	return &Fan{cfg: cfg, a: a, b: b, c: c}, nil
}

func (f *Fan) Config() FanConfig { return f.cfg }

func (f *Fan) FlowRange() (float64, float64) {
	pts := f.cfg.Points
	return pts[0].Flow, pts[len(pts)-1].Flow
}

func (f *Fan) ZeroFlowPressure() float64 {
	return f.cfg.Points[0].Pressure
}

func (f *Fan) Coefficient() (float64, float64, float64) {
	return f.a, f.b, f.c
}
