package fan

// SegmentSlopes 返回相邻样本点之间的线性段斜率数组，长度为 n-1。
func SegmentSlopes(pts []Point) []float64 {
	if len(pts) < 2 {
		return nil
	}
	out := make([]float64, 0, len(pts)-1)
	for i := 1; i < len(pts); i++ {
		slope := (pts[i].Pressure - pts[i-1].Pressure) / (pts[i].Flow - pts[i-1].Flow)
		out = append(out, slope)
	}
	return out
}

// SegmentIndex 返回流量 q 落在的样本段下标 i（表示区间 [i, i+1]）。
// q 超过最大样本流量时返回最后一段的 i；q 小于最小流量时返回 0。
func SegmentIndex(pts []Point, q float64) int {
	for i := 1; i < len(pts); i++ {
		if q <= pts[i].Flow {
			return i - 1
		}
	}
	return len(pts) - 2
}

// SegmentIntercept 返回第 i 段直线的截距 b = y_i - slope·x_i，即 Δp = slope·Q + b。
func SegmentIntercept(p0, p1 Point, slope float64) float64 {
	return p0.Pressure - slope*p0.Flow
}

// CurveArea 用梯形法则求样本范围 [lo, hi] 上风机曲线下的面积，单位 Pa·m³/s。
// 面积正比于风机在该流量区间可做的功，用于曲线形态检验。
func CurveArea(pts []Point) float64 {
	sum := 0.0
	for i := 1; i < len(pts); i++ {
		sum += 0.5 * (pts[i].Pressure + pts[i-1].Pressure) * (pts[i].Flow - pts[i-1].Flow)
	}
	return sum
}

// LinearlyFit 对两个点做直线插值并返回 (slope, intercept)。
func LinearlyFit(p0, p1 Point) (float64, float64) {
	slope := (p1.Pressure - p0.Pressure) / (p1.Flow - p0.Flow)
	return slope, SegmentIntercept(p0, p1, slope)
}
