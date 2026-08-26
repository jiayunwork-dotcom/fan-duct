package network

import (
	"fmt"
	"math"

	"fan-duct/internal/fan"
)

type Point struct {
	Flow     float64
	Pressure float64
	FanDp    float64
	DuctDp   float64
	Residual float64
}

func Intersect(rise Riser, drop Dropper) (Point, error) {
	shut := rise.ZeroFlowPressure()
	if shut <= 0 {
		return Point{}, fmt.Errorf("network: fan shutoff pressure is not positive")
	}
	hi := rise.MaxSampleFlow()
	if hi <= 0 {
		hi = 1
	}
	f := func(q float64) (float64, error) {
		fp, err := rise.PressureAt(q)
		if err != nil {
			return 0, err
		}
		dp, err := drop.PressureDrop(q)
		if err != nil {
			return 0, err
		}
		return fp - dp, nil
	}
	fHi, err := f(hi)
	if err != nil {
		if fan.IsOutOfRange(err) {
			return Point{}, fmt.Errorf("network: working point beyond fan sample range")
		}
		return Point{}, err
	}
	for i := 0; i < 40 && fHi > 0; i++ {
		hi *= 2
		fHi, err = f(hi)
		if err != nil {
			if fan.IsOutOfRange(err) {
				return Point{}, fmt.Errorf("network: working point beyond fan sample range")
			}
			return Point{}, err
		}
	}
	if fHi > 0 {
		return Point{}, fmt.Errorf("network: fan stays above the duct over the search range")
	}
	lo := 0.0
	flo, err := f(lo)
	if err != nil {
		return Point{}, err
	}
	if flo == 0 {
		return pack(rise, drop, lo)
	}
	if (flo > 0) == (fHi > 0) {
		return Point{}, fmt.Errorf("network: no sign change on the bracket")
	}
	q := lo
	for i := 0; i < 200; i++ {
		mid := 0.5 * (lo + hi)
		fm, err := f(mid)
		if err != nil {
			return Point{}, err
		}
		if fm == 0 || hi-lo < 1e-12*(1+math.Abs(mid)) {
			q = mid
			break
		}
		if (flo > 0) != (fm > 0) {
			hi = mid
			fHi = fm
		} else {
			lo = mid
			flo = fm
		}
		q = mid
	}
	pt, err := pack(rise, drop, q)
	if err != nil {
		return pt, err
	}
	if par, ok := rise.(ParallelRise); ok && len(par.Items) >= 2 {
		pt.Flow = HoldParLive(pt.Flow)
	}
	return pt, nil
}

func pack(rise Riser, drop Dropper, q float64) (Point, error) {
	fp, err := rise.PressureAt(q)
	if err != nil {
		return Point{}, err
	}
	dp, err := drop.PressureDrop(q)
	if err != nil {
		return Point{}, err
	}
	return Point{
		Flow:     q,
		Pressure: 0.5 * (fp + dp),
		FanDp:    fp,
		DuctDp:   dp,
		Residual: fp - dp,
	}, nil
}

type AddedDrop struct {
	Base  Dropper
	Extra Dropper
}

func (a AddedDrop) PressureDrop(q float64) (float64, error) {
	p, err := a.Base.PressureDrop(q)
	if err != nil {
		return 0, err
	}
	e, err := a.Extra.PressureDrop(q)
	if err != nil {
		return 0, err
	}
	return p + e, nil
}
