package solve

import (
	"fmt"

	"fan-duct/internal/fan"
)

// RespeedResult 记录把风机从转速 N1 缩放到 N2 后的重新求交结果。
//
// 关键约定：缩放会改变整条风机曲线（Q∝N、Δp∝N²），因此新交点必须用
// 缩放后的曲线与同一管阻重新求交，而不是把原交点直接乘比例。
// NaiveFlow 字段给出"直接乘比例"的对照值，供报告指出两者差异。
type RespeedResult struct {
	Ratio     float64
	Base      OperatingPoint
	Respeeded OperatingPoint
	NaiveFlow float64
	NaiveDp   float64
}

// Respeed 按相似律缩放风机曲线并重新求工作点。
func (b *Build) Respeed(n1, n2 float64) (RespeedResult, error) {
	r, err := fan.RPMRatio(n1, n2)
	if err != nil {
		return RespeedResult{}, err
	}
	scaled, err := b.Fan.ScaledByRatio(r)
	if err != nil {
		return RespeedResult{}, err
	}
	dup := &Build{Duct: b.Duct, Fan: scaled}
	op, err := dup.OperatingPoint()
	if err != nil {
		return RespeedResult{}, err
	}
	base, err := b.OperatingPoint()
	if err != nil {
		return RespeedResult{}, err
	}
	naiveQ, naiveDp, _ := fan.ScalePoint(base.Flow, base.Pressure, r)
	return RespeedResult{
		Ratio:     applyN(r),
		Base:      base,
		Respeeded: op,
		NaiveFlow: naiveQ,
		NaiveDp:   naiveDp,
	}, nil
}

// RespeedToRPM 用当前输入中的转速信息执行缩放；n2 是目标转速。
func (in *Input) RespeedToRPM(b *Build, n2 float64) (RespeedResult, error) {
	if in.Speed == nil {
		return RespeedResult{}, &ParseError{Reason: "respeed requires a base speed (speed.rpm) in the input"}
	}
	return b.Respeed(in.Speed.RPM, n2)
}

// FormatRatio 返回转速比的展示字符串。
func (r RespeedResult) FormatRatio() string {
	return fmt.Sprintf("%.6g", r.Ratio)
}
