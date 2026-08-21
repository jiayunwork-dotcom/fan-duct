package fan

import (
	"fmt"
	"math"
)

// Point 是风机样本曲线上的一个点：流量与风机提供的压升。
type Point struct {
	Flow     float64 // m³/s
	Pressure float64 // Pa
}

// FitMode 选择样本点之间的拟合方式。
type FitMode int

const (
	// FitPolyline 逐点线性插值。
	FitPolyline FitMode = iota
	// FitQuadratic 最小二乘二次多项式。
	FitQuadratic
)

// ParseFit 把字符串转换为拟合方式，"polyline"/"linear" → FitPolyline，
// "quadratic" → FitQuadratic；其他值返回 false。
func ParseFit(s string) (FitMode, bool) {
	switch s {
	case "", "polyline", "linear":
		return FitPolyline, true
	case "quadratic":
		return FitQuadratic, true
	}
	return FitPolyline, false
}

// String 返回拟合方式的字符串表示。
func (m FitMode) String() string {
	switch m {
	case FitQuadratic:
		return "quadratic"
	default:
		return "polyline"
	}
}

// ExtrapolateMode 选择流量超出样本范围时的处理策略。
type ExtrapolateMode int

const (
	// ExtrapolateError 禁止外推，越界直接报错。
	ExtrapolateError ExtrapolateMode = iota
	// ExtrapolateLinear 按首尾两点斜率线性外推。
	ExtrapolateLinear
	// ExtrapolateQuadratic 按二次多项式外推。
	ExtrapolateQuadratic
)

// ParseExtrapolate 把字符串转换为外推策略，默认 error。
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

// String 返回外推策略的字符串表示。
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

// FanConfig 描述风机曲线的样本点与拟合/外推方式。
//
//	Points      样本点，按流量严格递增；首点流量必须为 0（零流点）
//	Efficiency  可选的总效率样本（0~1），与 Points 等长；空表示未知
//	Fit         点内插值方式
//	Extrapolate 点外策略
type FanConfig struct {
	Points      []Point
	Efficiency  []float64
	Fit         FitMode
	Extrapolate ExtrapolateMode
}

// Validate 检查风机曲线配置：至少两个点、首点流量为零、流量严格递增、
// 压力非负且有限、效率样本（若提供）与点等长且在 (0,1]。
// 任一违反返回 *ConfigError。
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

// Fan 是一条可求值的风机 Δp(Q) 曲线。用 New 构造后只读。
type Fan struct {
	cfg FanConfig
	a   float64 // 二次系数
	b   float64
	c   float64
}

// New 校验配置并构造风机曲线。
func New(cfg FanConfig) (*Fan, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	a, b, c := FitQuadraticCoeffs(cfg.Points)
	return &Fan{cfg: cfg, a: a, b: b, c: c}, nil
}

// Config 返回配置副本。
func (f *Fan) Config() FanConfig { return f.cfg }

// FlowRange 返回样本点流量范围 [lo, hi]。
func (f *Fan) FlowRange() (float64, float64) {
	pts := f.cfg.Points
	return pts[0].Flow, pts[len(pts)-1].Flow
}

// ZeroFlowPressure 返回零流量时的静压（曲线零流点）。
func (f *Fan) ZeroFlowPressure() float64 {
	return f.cfg.Points[0].Pressure
}

// Coefficient 返回二次多项式 Δp = aQ² + bQ + c 的系数。
func (f *Fan) Coefficient() (float64, float64, float64) {
	return f.a, f.b, f.c
}
