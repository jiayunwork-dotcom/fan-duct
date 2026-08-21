package solve

import (
	"fmt"
	"strings"

	"fan-duct/internal/duct"
)

// Sensitivity 记录对风管某个参数做小扰动后工作点的变化。
type Sensitivity struct {
	Param     string  // 参数名
	BaseQ     float64 // 原工作点流量
	NewQ      float64 // 扰动后工作点流量
	BaseDp    float64 // 原工作点压降
	NewDp     float64 // 扰动后工作点压降
	DeltaQ    float64 // NewQ - BaseQ
	RelativeQ float64 // DeltaQ / BaseQ
}

// SensitivityResult 是一组参数的灵敏度汇总。
type SensitivityResult struct {
	Items []Sensitivity
}

// ComputeSensitivity 把风管参数 param 乘以 (1+frac) 后重新求交，
// 返回工作点流量的变化。param 可取 "length"、"diameter"、"lossCoeff"、
// "density"、"friction"。frac 是相对扰动量（如 0.1 表示 +10%）。
func ComputeSensitivity(b *Build, param string, frac float64) (Sensitivity, error) {
	dc := b.Duct.Config()
	base, err := b.OperatingPoint()
	if err != nil {
		return Sensitivity{}, err
	}
	switch param {
	case "length":
		dc.Length *= (1 + frac)
	case "diameter":
		dc.Diameter *= (1 + frac)
	case "lossCoeff":
		dc.LossCoeff *= (1 + frac)
	case "density":
		dc.Density *= (1 + frac)
	case "friction":
		dc.Friction *= (1 + frac)
	default:
		return Sensitivity{}, fmt.Errorf("solve: unknown sensitivity parameter %q", param)
	}
	nd, err := duct.New(dc)
	if err != nil {
		return Sensitivity{}, err
	}
	nb := &Build{Duct: nd, Fan: b.Fan}
	np, err := nb.OperatingPoint()
	if err != nil {
		return Sensitivity{}, err
	}
	rel := 0.0
	if base.Flow != 0 {
		rel = (np.Flow - base.Flow) / base.Flow
	}
	return Sensitivity{
		Param:     param,
		BaseQ:     base.Flow,
		NewQ:      np.Flow,
		BaseDp:    base.Pressure,
		NewDp:     np.Pressure,
		DeltaQ:    np.Flow - base.Flow,
		RelativeQ: rel,
	}, nil
}

// ComputeSensitivities 对常用参数批量做 +10% 扰动灵敏度。
func ComputeSensitivities(b *Build) (SensitivityResult, error) {
	var out SensitivityResult
	for _, p := range []string{"length", "diameter", "lossCoeff", "density"} {
		s, err := ComputeSensitivity(b, p, 0.1)
		if err != nil {
			return out, err
		}
		out.Items = append(out.Items, s)
	}
	return out, nil
}

// String 渲染灵敏度汇总为多行文本。
func (r SensitivityResult) String() string {
	var sb strings.Builder
	sb.WriteString("sensitivity (+10%% on duct parameter):\n")
	for _, s := range r.Items {
		sb.WriteString(fmt.Sprintf(
			"  %-10s Q %.6g -> %.6g m3/s (%.2f%%), dP %.6g -> %.6g Pa\n",
			s.Param, s.BaseQ, s.NewQ, s.RelativeQ*100, s.BaseDp, s.NewDp))
	}
	return sb.String()
}
