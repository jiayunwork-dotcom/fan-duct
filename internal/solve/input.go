package solve

import "fmt"

// Input 是 `operate` 子命令的完整 JSON 输入。
//
//	duct     风管参数
//	fan      风机样本曲线
//	speed    当前转速（可省略，仅作展示）
//	newSpeed 目标转速；给出时按相似律缩放风机并重新求交
type Input struct {
	Duct     DuctSpec   `json:"duct"`
	Fan      FanSpec    `json:"fan"`
	Speed    *SpeedSpec `json:"speed,omitempty"`
	NewSpeed *SpeedSpec `json:"newSpeed,omitempty"`
}

// DuctSpec 是风管部分的 JSON 字段。带指针的字段缺省时使用库默认值。
type DuctSpec struct {
	Length    float64  `json:"length"`              // 管长 m
	Diameter  float64  `json:"diameter"`            // 内径 m
	Friction  *float64 `json:"friction,omitempty"`  // 固定 f；缺省按 Re 自动
	LossCoeff *float64 `json:"lossCoeff,omitempty"` // 局部阻力系数和 ΣK
	Density   *float64 `json:"density,omitempty"`   // 缺省 1.205 kg/m³
	Viscosity *float64 `json:"viscosity,omitempty"` // 缺省 1.82e-5 Pa·s
	Roughness *float64 `json:"roughness,omitempty"` // 相对粗糙度 ε/D；>0 时自动模式用 Swamee–Jain
}

// FanSpec 是风机部分的 JSON 字段。
type FanSpec struct {
	Points      []PointSpec `json:"points"`                // 样本点，至少 2 个
	Efficiency  []float64   `json:"efficiency,omitempty"`  // 可选，与 points 等长（0~1）
	Fit         string      `json:"fit,omitempty"`         // polyline | quadratic
	Extrapolate string      `json:"extrapolate,omitempty"` // error | linear | quadratic
}

// PointSpec 是单个样本点 {q, dp}。
type PointSpec struct {
	Q  float64 `json:"q"`  // m³/s
	Dp float64 `json:"dp"` // Pa
}

// SpeedSpec 是转速 {rpm}。
type SpeedSpec struct {
	RPM float64 `json:"rpm"`
}

// Validate 检查输入中除风机曲线以外的参数约束。
func (in *Input) Validate() error {
	if in.Speed != nil && in.Speed.RPM <= 0 {
		return &ParseError{Reason: fmt.Sprintf("speed rpm must be > 0 (got %v)", in.Speed.RPM)}
	}
	if in.NewSpeed != nil && in.NewSpeed.RPM <= 0 {
		return &ParseError{Reason: fmt.Sprintf("newSpeed rpm must be > 0 (got %v)", in.NewSpeed.RPM)}
	}
	return nil
}

// HasRespeed 报告输入是否包含目标转速。
func (in *Input) HasRespeed() bool {
	return in.NewSpeed != nil
}
