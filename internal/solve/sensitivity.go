package solve

import (
	"fmt"
	"strings"

	"fan-duct/internal/duct"
)

type Sensitivity struct {
	Param     string
	BaseQ     float64
	NewQ      float64
	BaseDp    float64
	NewDp     float64
	DeltaQ    float64
	RelativeQ float64
}

type SensitivityResult struct {
	Items []Sensitivity
}

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
