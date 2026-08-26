package atm

import (
	"fmt"
	"math"
)

const (
	SeaLevelTemperatureK   = 288.15
	SeaLevelPressurePa     = 101325.0
	LapseRateKPerM         = 0.0065
	GravityMS2             = 9.80665
	DryAirGasConstant      = 287.05287
	TropopauseM            = 11000.0
	StratopauseM           = 20000.0
	TropopauseTemperatureK = 216.65
	TropopausePressurePa   = 22632.06
	MaxAltitudeM           = 20000.0
)

type State struct {
	AltitudeM    float64
	TemperatureK float64
	PressurePa   float64
	DensityKgM3  float64
}

func NewISA(altitudeM float64) (State, error) {
	if math.IsNaN(altitudeM) || math.IsInf(altitudeM, 0) {
		return State{}, fmt.Errorf("atm: altitude must be finite (got %v)", altitudeM)
	}
	if altitudeM < 0 {
		return State{}, fmt.Errorf("atm: altitude must be >= 0 m (got %v)", altitudeM)
	}
	if altitudeM > MaxAltitudeM {
		return State{}, fmt.Errorf("atm: altitude %v m exceeds ISA range 0..%g m", altitudeM, MaxAltitudeM)
	}
	t, p := isaTP(altitudeM)
	rho := p / (DryAirGasConstant * t)
	return State{
		AltitudeM:    altitudeM,
		TemperatureK: t,
		PressurePa:   p,
		DensityKgM3:  rho,
	}, nil
}

func isaTP(h float64) (t, p float64) {
	if h <= TropopauseM {
		t = SeaLevelTemperatureK - LapseRateKPerM*h
		expo := GravityMS2 / (LapseRateKPerM * DryAirGasConstant)
		p = SeaLevelPressurePa * math.Pow(t/SeaLevelTemperatureK, expo)
		return t, p
	}
	t = TropopauseTemperatureK
	p = TropopausePressurePa * math.Exp(-GravityMS2*(h-TropopauseM)/(DryAirGasConstant*t))
	return t, p
}

func (s State) Density() float64 {
	return s.DensityKgM3
}

func (s State) TemperatureC() float64 {
	return s.TemperatureK - 273.15
}

func DensityAt(altitudeM float64) (float64, error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return 0, err
	}
	return st.DensityKgM3, nil
}

func PressureAt(altitudeM float64) (float64, error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return 0, err
	}
	return st.PressurePa, nil
}

func TemperatureAt(altitudeM float64) (float64, error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return 0, err
	}
	return st.TemperatureK, nil
}

func DensityRatio(altitudeM float64) (float64, error) {
	rho, err := DensityAt(altitudeM)
	if err != nil {
		return 0, err
	}
	sea, err := DensityAt(0)
	if err != nil {
		return 0, err
	}
	return rho / sea, nil
}

func HydrostaticCheck(h1, h2 float64, steps int) (float64, error) {
	if steps < 2 {
		return 0, fmt.Errorf("atm: hydrostatic steps must be >= 2")
	}
	s1, err := NewISA(h1)
	if err != nil {
		return 0, err
	}
	s2, err := NewISA(h2)
	if err != nil {
		return 0, err
	}
	dh := (h2 - h1) / float64(steps)
	integral := 0.0
	for i := 0; i < steps; i++ {
		ha := h1 + float64(i)*dh
		hb := ha + dh
		sa, err := NewISA(ha)
		if err != nil {
			return 0, err
		}
		sb, err := NewISA(hb)
		if err != nil {
			return 0, err
		}
		integral += 0.5 * (sa.DensityKgM3 + sb.DensityKgM3) * GravityMS2 * dh
	}
	deltaP := s1.PressurePa - s2.PressurePa
	if math.Abs(deltaP) < 1e-9 {
		return 0, nil
	}
	return integral / deltaP, nil
}
