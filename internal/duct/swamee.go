package duct

import "math"

func SwameeJain(re, relRoughness float64) float64 {
	denom := math.Log10(relRoughness/3.7 + 5.74/math.Pow(re, 0.9))
	return 0.25 / (denom * denom)
}

func ColebrookImplicit(re, relRoughness, f float64) float64 {
	inv := 1 / math.Sqrt(f)
	return inv + 2*math.Log10(relRoughness/3.7+2.51/(re*math.Sqrt(f)))
}

func RoughnessFriction(re, relRoughness float64) float64 {
	if re <= 1 {
		return Laminar(1)
	}
	if re <= LaminarMax {
		return Laminar(re)
	}
	if re >= TurbulentMin {
		return SwameeJain(re, relRoughness)
	}
	fLo := Laminar(LaminarMax)
	fHi := SwameeJain(TurbulentMin, relRoughness)
	t := (re - LaminarMax) / (TurbulentMin - LaminarMax)
	return fLo + t*(fHi-fLo)
}
