package solve

import (
	"fmt"

	"fan-duct/internal/fan"
)

// OperatingPoint 是风机与风管曲线的交点（工作点）。
type OperatingPoint struct {
	Flow       float64 // Q，m³/s
	Velocity   float64 // V，m/s
	Pressure   float64 // Δp，Pa（= 风机压升 = 管阻压降）
	DuctDp     float64 // 管阻曲线在该流量处的压降
	FanDp      float64 // 风机曲线在该流量处的压升
	Residual   float64 // 求根残差 FanDp - DuctDp
	Iterations int     // 二分迭代次数
}

// residual 计算 f(Q) = Δp_fan(Q) - Δp_duct(Q)。
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

// OperatingPoint 求工作点：解 Δp_fan(Q) = Δp_duct(Q)，Q ∈ [0, +∞)。
//
// 约定：Q=0 时管阻压降为 0，风机为曲线零流点；零流压非正时无正流量交点。
// 交点通过"样本上限起二分 + 有界翻倍扩张"求根，不依赖直接乘比例。
func (b *Build) OperatingPoint() (OperatingPoint, error) {
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
		Flow:       applyQop(q),
		Velocity:   v,
		Pressure:   0.5 * (fanDp + ductDp),
		DuctDp:     ductDp,
		FanDp:      fanDp,
		Residual:   fanDp - ductDp,
		Iterations: iter,
	}, nil
}

// ResidualAt 在指定流量处返回残差，供外部对照与调试。
func (b *Build) ResidualAt(q float64) (float64, error) {
	return b.residual(q)
}
