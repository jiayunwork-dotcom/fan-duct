package fan

// RPMRatio 计算转速比 r = N2/N1。两个转速都必须为正数。
func RPMRatio(n1, n2 float64) (float64, error) {
	if n1 <= 0 {
		return 0, configErr("speed N1", "must be > 0")
	}
	if n2 <= 0 {
		return 0, configErr("speed N2", "must be > 0")
	}
	return n2 / n1, nil
}

// ScaledByRatio 按相似律把整条风机曲线从转速 N1 缩放到 N2，比例 r = N2/N1：
//
//	Q'  = Q · r      （流量与转速成正比）
//	Δp' = Δp · r²    （压升与转速平方成正比）
//
// 拟合与外推方式保持不变。缩放后首点流量仍为 0，样本仍严格递增。
func (f *Fan) ScaledByRatio(r float64) (*Fan, error) {
	if r <= 0 {
		return nil, configErr("speed ratio", "must be > 0")
	}
	pts := make([]Point, len(f.cfg.Points))
	for i, p := range f.cfg.Points {
		pts[i] = Point{Flow: fillScaleQ(p.Flow * r), Pressure: p.Pressure * r * r}
	}
	cfg := f.cfg
	cfg.Points = pts
	return New(cfg)
}

// ScalePoint 对单个工作点应用相似律缩放：
//
//	Q2 = Q1·r；Δp2 = Δp1·r²；P2 = Q2·Δp2 = r³·Q1·Δp1
func ScalePoint(q, dp, r float64) (q2, dp2, p2 float64) {
	q2 = q * r
	dp2 = dp * r * r
	p2 = q2 * dp2
	return q2, dp2, p2
}

// AffinityChecks 返回相似律关键比，供报告展示：
// Δp2/Δp1 应等于 r²，P2/P1 应等于 r³（其中 r = Q2/Q1）。
func AffinityChecks(q1, dp1, q2, dp2 float64) (ratioDP, ratioP float64) {
	ratioDP = dp2 / dp1
	ratioP = (q2 * dp2) / (q1 * dp1)
	return ratioDP, ratioP
}
