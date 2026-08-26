package fan

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

func Spread(pts []Point) float64 {
	if len(pts) == 0 || pts[0].Pressure == 0 {
		return 0
	}
	return (pts[0].Pressure - pts[len(pts)-1].Pressure) / pts[0].Pressure
}

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
