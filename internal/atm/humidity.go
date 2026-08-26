package atm

import (
	"fmt"
	"math"
)

const (
	WaterGasConstant = 461.5
	MagnusA          = 17.625
	MagnusB          = 243.04
)

func SaturationVaporPa(tempC float64) (float64, error) {
	if tempC <= -40 || tempC >= 50 {
		return 0, fmt.Errorf("atm: Magnus saturation range is -40..50 C (got %v)", tempC)
	}
	es := 610.94 * math.Exp(MagnusA*tempC/(tempC+MagnusB))
	return es, nil
}

func MoistDensity(pressurePa, tempK, relHumidity float64) (float64, error) {
	if pressurePa <= 0 {
		return 0, fmt.Errorf("atm: pressure must be > 0 (got %v)", pressurePa)
	}
	if tempK <= 0 {
		return 0, fmt.Errorf("atm: temperature must be > 0 K (got %v)", tempK)
	}
	if relHumidity < 0 || relHumidity > 1 {
		return 0, fmt.Errorf("atm: relative humidity must be in [0, 1] (got %v)", relHumidity)
	}
	es, err := SaturationVaporPa(tempK - 273.15)
	if err != nil {
		return 0, err
	}
	pv := relHumidity * es
	if pv >= pressurePa {
		return 0, fmt.Errorf("atm: vapor pressure %v exceeds total pressure %v", pv, pressurePa)
	}
	pd := pressurePa - pv
	return pd/(DryAirGasConstant*tempK) + pv/(WaterGasConstant*tempK), nil
}

func MoistCorrection(altitudeM, relHumidity float64) (dry, moist, ratio float64, err error) {
	st, err := NewISA(altitudeM)
	if err != nil {
		return 0, 0, 0, err
	}
	dry = st.DensityKgM3
	moist, err = MoistDensity(st.PressurePa, st.TemperatureK, relHumidity)
	if err != nil {
		return 0, 0, 0, err
	}
	ratio = moist / dry
	return dry, moist, ratio, nil
}
