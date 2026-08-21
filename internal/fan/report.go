package fan

import "fmt"

// Report 汇总风机曲线的配置与特征量。
type Report struct {
	Config           FanConfig
	Fit              string
	Extrapolate      string
	ZeroFlowPressure float64
	Lo               float64
	Hi               float64
	A                float64
	B                float64
	C                float64
	FirstSlope       float64
	LastSlope        float64
}

// Report 构造风机摘要。
func (f *Fan) Report() Report {
	lo, hi := f.FlowRange()
	slopes := SegmentSlopes(f.cfg.Points)
	r := Report{
		Config:           f.cfg,
		Fit:              f.cfg.Fit.String(),
		Extrapolate:      f.cfg.Extrapolate.String(),
		ZeroFlowPressure: f.ZeroFlowPressure(),
		Lo:               lo,
		Hi:               hi,
		A:                f.a,
		B:                f.b,
		C:                f.c,
	}
	if len(slopes) > 0 {
		r.FirstSlope = slopes[0]
		r.LastSlope = slopes[len(slopes)-1]
	}
	return r
}

// String 渲染风机摘要为多行文本。
func (r Report) String() string {
	return fmt.Sprintf(
		"Fan model:\n"+
			"  sample points  = %d\n"+
			"  flow range     = [%.6g, %.6g] m3/s\n"+
			"  fit            = %s\n"+
			"  extrapolate    = %s\n"+
			"  shutoff dp     = %.6g Pa\n"+
			"  quadratic      = %.6g Q^2 + %.6g Q + %.6g\n"+
			"  first slope    = %.6g Pa/(m3/s)\n"+
			"  last slope     = %.6g Pa/(m3/s)\n",
		len(r.Config.Points), r.Lo, r.Hi, r.Fit, r.Extrapolate,
		r.ZeroFlowPressure, r.A, r.B, r.C, r.FirstSlope, r.LastSlope)
}

// OperatingHint 返回估计的工作点提示：把样本中压升首次低于零流压 90% 的位置
// 作为"曲线快速下降区"的入口，帮助用户理解交点可能出现的位置。
func (f *Fan) OperatingHint() (float64, bool) {
	target := 0.9 * f.ZeroFlowPressure()
	return f.ZeroCrossing(target)
}
