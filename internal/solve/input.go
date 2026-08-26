package solve

import "fmt"

type Input struct {
	Duct     DuctSpec   `json:"duct"`
	Fan      FanSpec    `json:"fan"`
	Speed    *SpeedSpec `json:"speed,omitempty"`
	NewSpeed *SpeedSpec `json:"newSpeed,omitempty"`
}

type DuctSpec struct {
	Length    float64  `json:"length"`
	Diameter  float64  `json:"diameter"`
	Friction  *float64 `json:"friction,omitempty"`
	LossCoeff *float64 `json:"lossCoeff,omitempty"`
	Density   *float64 `json:"density,omitempty"`
	Viscosity *float64 `json:"viscosity,omitempty"`
	Roughness *float64 `json:"roughness,omitempty"`
}

type FanSpec struct {
	Points      []PointSpec `json:"points"`
	Efficiency  []float64   `json:"efficiency,omitempty"`
	Fit         string      `json:"fit,omitempty"`
	Extrapolate string      `json:"extrapolate,omitempty"`
}

type PointSpec struct {
	Q  float64 `json:"q"`
	Dp float64 `json:"dp"`
}

type SpeedSpec struct {
	RPM float64 `json:"rpm"`
}

func (in *Input) Validate() error {
	if in.Speed != nil && in.Speed.RPM <= 0 {
		return &ParseError{Reason: fmt.Sprintf("speed rpm must be > 0 (got %v)", in.Speed.RPM)}
	}
	if in.NewSpeed != nil && in.NewSpeed.RPM <= 0 {
		return &ParseError{Reason: fmt.Sprintf("newSpeed rpm must be > 0 (got %v)", in.NewSpeed.RPM)}
	}
	return nil
}

func (in *Input) HasRespeed() bool {
	return in.NewSpeed != nil
}
