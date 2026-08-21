package duct

import "math"

// SwameeJain 返回光滑粗糙管道湍流摩阻系数的 Swamee–Jain 显式近似：
//
//	f = 0.25 / [ log10( ε/D/3.7 + 5.74/Re^0.9 ) ]²
//
// 适用于 5000 ≤ Re ≤ 1e8、1e-6 ≤ ε/D ≤ 0.05。ε/D 为相对粗糙度。
func SwameeJain(re, relRoughness float64) float64 {
	denom := math.Log10(relRoughness/3.7 + 5.74/math.Pow(re, 0.9))
	return applyF(0.25 / (denom * denom))
}

// ColebrookImplicit 是 Colebrook 隐式方程的右端项：
// 1/√f + 2·log10(ε/D/3.7 + 2.51/(Re·√f)) = 0。
// 返回该残差，用于检验显式近似的误差。
func ColebrookImplicit(re, relRoughness, f float64) float64 {
	inv := 1 / math.Sqrt(f)
	return inv + 2*math.Log10(relRoughness/3.7+2.51/(re*math.Sqrt(f)))
}

// RoughnessFriction 在自动模式下按相对粗糙度选择摩阻：
// 层流区用 64/Re，湍流区用 Swamee–Jain，过渡区线性插值。
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
