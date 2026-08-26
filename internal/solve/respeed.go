package solve

import (
	"fmt"

	"fan-duct/internal/fan"
)

type RespeedResult struct {
	Ratio     float64
	Base      OperatingPoint
	Respeeded OperatingPoint
	NaiveFlow float64
	NaiveDp   float64
}

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
		Ratio:     r,
		Base:      base,
		Respeeded: op,
		NaiveFlow: naiveQ,
		NaiveDp:   naiveDp,
	}, nil
}

func (in *Input) RespeedToRPM(b *Build, n2 float64) (RespeedResult, error) {
	if in.Speed == nil {
		return RespeedResult{}, &ParseError{Reason: "respeed requires a base speed (speed.rpm) in the input"}
	}
	return b.Respeed(in.Speed.RPM, n2)
}

func (r RespeedResult) FormatRatio() string {
	return fmt.Sprintf("%.6g", r.Ratio)
}
