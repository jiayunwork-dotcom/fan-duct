package fan

// FitQuality 返回二次拟合的决定系数 R² 与最大绝对残差。
// R² 接近 1 表示拟合良好；maxRes 给出单个样本点上的最大偏差。
func FitQuality(pts []Point) (r2 float64, maxRes float64) {
	if len(pts) < 3 {
		return 0, 0
	}
	a, b, c := FitQuadraticCoeffs(pts)
	mean := 0.0
	for _, p := range pts {
		mean += p.Pressure
	}
	mean /= float64(len(pts))
	ssTot := 0.0
	ssRes := 0.0
	maxRes = 0
	for _, p := range pts {
		fit := a*p.Flow*p.Flow + b*p.Flow + c
		res := p.Pressure - fit
		ssTot += (p.Pressure - mean) * (p.Pressure - mean)
		ssRes += res * res
		if abs(res) > maxRes {
			maxRes = abs(res)
		}
	}
	if ssTot == 0 {
		return 1, maxRes
	}
	return 1 - ssRes/ssTot, maxRes
}

// MeanPressure 返回样本点压升的算术平均值。
func MeanPressure(pts []Point) float64 {
	if len(pts) == 0 {
		return 0
	}
	sum := 0.0
	for _, p := range pts {
		sum += p.Pressure
	}
	return sum / float64(len(pts))
}

// MaxPressure 返回样本点中的最大压升（对典型风机即零流压）。
func MaxPressure(pts []Point) float64 {
	if len(pts) == 0 {
		return 0
	}
	m := pts[0].Pressure
	for _, p := range pts[1:] {
		if p.Pressure > m {
			m = p.Pressure
		}
	}
	return m
}

// MinPressure 返回样本点中的最小压升。
func MinPressure(pts []Point) float64 {
	if len(pts) == 0 {
		return 0
	}
	m := pts[0].Pressure
	for _, p := range pts[1:] {
		if p.Pressure < m {
			m = p.Pressure
		}
	}
	return m
}

// Spread 返回压升相对零流压的下落比例 (p0 - pN)/p0，用于评估曲线陡峭程度。
func Spread(pts []Point) float64 {
	if len(pts) == 0 || pts[0].Pressure == 0 {
		return 0
	}
	return (pts[0].Pressure - pts[len(pts)-1].Pressure) / pts[0].Pressure
}

// RelativeResidual 返回第 i 个样本点的相对残差 (y - fit)/y。
func RelativeResidual(pts []Point, a, b, c float64) []float64 {
	out := make([]float64, len(pts))
	for i, p := range pts {
		fit := a*p.Flow*p.Flow + b*p.Flow + c
		if p.Pressure != 0 {
			out[i] = (p.Pressure - fit) / p.Pressure
		}
	}
	return out
}
