package atm

import (
	"fmt"
	"math"
)

const (
	SutherlandRefMu    = 1.716e-5
	SutherlandRefT     = 273.15
	SutherlandConstant = 110.4
)

func SutherlandViscosity(tempK float64) (float64, error) {
	if math.IsNaN(tempK) || math.IsInf(tempK, 0) || tempK <= 0 {
		return 0, fmt.Errorf("atm: temperature must be > 0 K (got %v)", tempK)
	}
	ratio := tempK / SutherlandRefT
	return SutherlandRefMu * math.Pow(ratio, 1.5) * (SutherlandRefT + SutherlandConstant) / (tempK + SutherlandConstant), nil
}

func ViscosityAtAltitude(altitudeM float64) (float64, error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return 0, err
	}
	return SutherlandViscosity(st.TemperatureK)
}

func KinematicViscosity(density, dynamicViscosity float64) (float64, error) {
	if density <= 0 {
		return 0, fmt.Errorf("atm: density must be > 0 (got %v)", density)
	}
	if dynamicViscosity <= 0 {
		return 0, fmt.Errorf("atm: dynamic viscosity must be > 0 (got %v)", dynamicViscosity)
	}
	return dynamicViscosity / density, nil
}

func PropertiesAt(altitudeM float64) (State, float64, float64, error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return State{}, 0, 0, err
	}
	mu, err := SutherlandViscosity(st.TemperatureK)
	if err != nil {
		return State{}, 0, 0, err
	}
	nu, err := KinematicViscosity(st.DensityKgM3, mu)
	if err != nil {
		return State{}, 0, 0, err
	}
	return st, mu, nu, nil
}
