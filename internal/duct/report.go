package duct

import "fmt"

// Report 汇总给定流量下风管的几何、流动与阻力信息。
// 动压与管阻压降分开列出，避免与风机曲线的全压混淆。
type Report struct {
	Config          DuctConfig
	Area            float64
	Flow            float64
	Velocity        float64
	Reynolds        float64
	Friction        float64
	Resistance      float64
	DynamicPressure float64
	PressureDrop    float64
}

// Report 计算给定流量下的风管摘要。
func (d *Duct) Report(q float64) (Report, error) {
	dp, err := d.PressureDrop(q)
	if err != nil {
		return Report{}, err
	}
	v := d.Velocity(q)
	re := d.Reynolds(v)
	return Report{
		Config:          d.cfg,
		Area:            d.area,
		Flow:            q,
		Velocity:        v,
		Reynolds:        re,
		Friction:        d.FrictionFactor(re),
		Resistance:      d.ResistanceFactor(re),
		DynamicPressure: d.VelocityPressure(v),
		PressureDrop:    dp,
	}, nil
}

// String 渲染风管摘要为多行文本。
func (r Report) String() string {
	mode := "auto (Re-based)"
	if r.Config.Friction > 0 {
		mode = "fixed"
	}
	return fmt.Sprintf(
		"Duct model:\n"+
			"  length L      = %.6g m\n"+
			"  diameter D    = %.6g m\n"+
			"  area A        = %.6g m2\n"+
			"  friction f    = %.6g (%s)\n"+
			"  local sum K   = %.6g\n"+
			"  density rho   = %.6g kg/m3\n"+
			"  viscosity mu  = %.6g Pa.s\n"+
			"At Q = %.6g m3/s:\n"+
			"  velocity V    = %.6g m/s\n"+
			"  Re            = %.6g\n"+
			"  total R       = %.6g\n"+
			"  dynamic q     = %.6g Pa\n"+
			"  dP duct       = %.6g Pa\n",
		r.Config.Length, r.Config.Diameter, r.Area, r.Friction, mode,
		r.Config.LossCoeff, r.Config.Density, r.Config.Viscosity,
		r.Flow, r.Velocity, r.Reynolds, r.Resistance, r.DynamicPressure, r.PressureDrop)
}
