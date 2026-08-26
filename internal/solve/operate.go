package solve

import (
	"fmt"

	"fan-duct/internal/fan"
)

type OperatingPoint struct {
	Flow       float64
	Velocity   float64
	Pressure   float64
	DuctDp     float64
	FanDp      float64
	Residual   float64
	Iterations int
}

func (b *Build) residual(q float64) (float64, error) {
	fp, err := b.Fan.PressureAt(q)
	if err != nil {
		return 0, err
	}
	dp, err := b.Duct.PressureDrop(q)
	if err != nil {
		return 0, err
	}
	return fp - dp, nil
}

func (b *Build) OperatingPoint() (OperatingPoint, error) {
	if err := abortFresh(); err != nil {
		return OperatingPoint{}, err
	}
	shutoff := b.Fan.ZeroFlowPressure()
	if shutoff <= 0 {
		return OperatingPoint{}, noIntersection("fan shutoff pressure is not positive")
	}
	hi, fHi, err := b.expandUpperBound(b.Fan.MaxSampleFlow(), b.residual, 60)
	if err != nil {
		if fan.IsOutOfRange(err) {
			return OperatingPoint{}, noIntersection(
				fmt.Sprintf("working point is beyond the fan sample range and extrapolation is disabled (%v)", err))
		}
		return OperatingPoint{}, err
	}
	if fHi > 0 {
		return OperatingPoint{}, noIntersection("fan curve stays above the duct curve over the whole search range")
	}
	q, iter, err := findRoot(b.residual, 0, hi, 200)
	if err != nil {
		return OperatingPoint{}, err
	}
	fanDp, err := b.Fan.PressureAt(q)
	if err != nil {
		return OperatingPoint{}, err
	}
	ductDp, err := b.Duct.PressureDrop(q)
	if err != nil {
		return OperatingPoint{}, err
	}
	v := b.Duct.Velocity(q)
	return OperatingPoint{
		Flow:       q,
		Velocity:   v,
		Pressure:   0.5 * (fanDp + ductDp),
		DuctDp:     ductDp,
		FanDp:      fanDp,
		Residual:   fanDp - ductDp,
		Iterations: iter,
	}, nil
}

func (b *Build) ResidualAt(q float64) (float64, error) {
	return b.residual(q)
}
