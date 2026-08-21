package solve

import (
	"fmt"
	"strings"
)

// Output 是一次 operate 调用要打印的完整结果。
type Output struct {
	Input     *Input
	Build     *Build
	Base      OperatingPoint
	Respeeded *RespeedResult
}

// String 渲染 operate 结果的多行文本：先描述两模型，再给工作点，
// 有目标转速时再给重新求交的新工作点与纯缩放对照。
func (o Output) String() string {
	var sb strings.Builder
	sb.WriteString("fan-duct operating point\n")
	sb.WriteString(ModelBlock(o.Build))
	sb.WriteString(o.BaseBlock())

	if o.Respeeded != nil {
		sb.WriteString(o.RespeededBlock())
	}
	return sb.String()
}

// BaseBlock 渲染基础工作点段落。
func (o Output) BaseBlock() string {
	op := o.Base
	return fmt.Sprintf(
		"operating point:\n"+
			"  flow        Q = %.6g m3/s\n"+
			"  velocity    V = %.6g m/s\n"+
			"  pressure   Δp = %.6g Pa\n"+
			"  fan  Δp        = %.6g Pa\n"+
			"  duct Δp        = %.6g Pa\n"+
			"  residual       = %.3g Pa\n"+
			"  iterations     = %d\n",
		op.Flow, op.Velocity, op.Pressure, op.FanDp, op.DuctDp, op.Residual, op.Iterations)
}

// RespeededBlock 渲染目标转速下的重新求交段落。
func (o Output) RespeededBlock() string {
	r := o.Respeeded
	op := r.Respeeded
	return fmt.Sprintf(
		"respeed N2/N1 = %s (N1=%.6g rpm -> N2=%.6g rpm):\n"+
			"  flow        Q2 = %.6g m3/s\n"+
			"  velocity    V2 = %.6g m/s\n"+
			"  pressure   Δp2 = %.6g Pa\n"+
			"  residual       = %.3g Pa\n"+
			"  naive Q1*ratio = %.6g m3/s (for reference, not the solution)\n"+
			"  naive Δp1*r^2  = %.6g Pa\n",
		r.FormatRatio(), o.Input.Speed.RPM, o.Input.NewSpeed.RPM,
		op.Flow, op.Velocity, op.Pressure, op.Residual,
		r.NaiveFlow, r.NaiveDp)
}

// Compact 返回单行摘要，便于脚本解析。
func (o Output) Compact() string {
	op := o.Base
	line := fmt.Sprintf("Q=%.6g V=%.6g dP=%.6g", op.Flow, op.Velocity, op.Pressure)
	if o.Respeeded != nil {
		line += fmt.Sprintf(" Q2=%.6g dP2=%.6g ratio=%s", o.Respeeded.Respeeded.Flow, o.Respeeded.Respeeded.Pressure, o.Respeeded.FormatRatio())
	}
	return line
}
