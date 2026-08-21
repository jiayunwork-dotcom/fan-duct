package solve

import (
	"fmt"
	"os"

	"fan-duct/internal/duct"
	"fan-duct/internal/fan"
)

// osReadFile 是 os.ReadFile 的可替换引用，便于测试注入。
var osReadFile = os.ReadFile

// Build 是把已解析的输入组装成可计算的 duct 与 fan 模型。
type Build struct {
	Duct *duct.Duct
	Fan  *fan.Fan
}

// Build 从 Input 构造 Build。风管缺省物性在这里填充，
// 所有非法参数（直径 ≤ 0、长度 < 0、密度 ≤ 0、空风机曲线等）在此报错。
func (in *Input) Build() (*Build, error) {
	dc := duct.DuctConfig{
		Length:    in.Duct.Length,
		Diameter:  in.Duct.Diameter,
		Density:   duct.DefaultDensity,
		Viscosity: duct.DefaultViscosity,
	}
	if in.Duct.Friction != nil {
		dc.Friction = *in.Duct.Friction
	}
	if in.Duct.LossCoeff != nil {
		dc.LossCoeff = *in.Duct.LossCoeff
	}
	if in.Duct.Density != nil {
		dc.Density = *in.Duct.Density
	}
	if in.Duct.Viscosity != nil {
		dc.Viscosity = *in.Duct.Viscosity
	}
	if in.Duct.Roughness != nil {
		dc.Roughness = *in.Duct.Roughness
	}
	d, err := duct.New(dc)
	if err != nil {
		return nil, err
	}
	fm, ok := fan.ParseFit(in.Fan.Fit)
	if !ok {
		return nil, &ParseError{Reason: fmt.Sprintf("unknown fit mode %q (want polyline or quadratic)", in.Fan.Fit)}
	}
	em, ok := fan.ParseExtrapolate(in.Fan.Extrapolate)
	if !ok {
		return nil, &ParseError{Reason: fmt.Sprintf("unknown extrapolate mode %q (want error, linear or quadratic)", in.Fan.Extrapolate)}
	}
	pts := make([]fan.Point, len(in.Fan.Points))
	for i, p := range in.Fan.Points {
		pts[i] = fan.Point{Flow: p.Q, Pressure: p.Dp}
	}
	f, err := fan.New(fan.FanConfig{
		Points:      pts,
		Efficiency:  append([]float64(nil), in.Fan.Efficiency...),
		Fit:         fm,
		Extrapolate: em,
	})
	if err != nil {
		return nil, err
	}
	return &Build{Duct: d, Fan: f}, nil
}

// MustBuild 从文件路径直接解析并构建，供 CLI 接线使用。
func MustBuild(path string) (*Build, error) {
	in, err := ParseInputFile(path)
	if err != nil {
		return nil, err
	}
	return in.Build()
}
