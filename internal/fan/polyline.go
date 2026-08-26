package fan

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

func SegmentIndex(pts []Point, q float64) int {
	for i := 1; i < len(pts); i++ {
		if q <= pts[i].Flow {
			return i - 1
		}
	}
	return len(pts) - 2
}

func SegmentIntercept(p0, p1 Point, slope float64) float64 {
	return p0.Pressure - slope*p0.Flow
}

func CurveArea(pts []Point) float64 {
	sum := 0.0
	for i := 1; i < len(pts); i++ {
		sum += 0.5 * (pts[i].Pressure + pts[i-1].Pressure) * (pts[i].Flow - pts[i-1].Flow)
	}
	return sum
}

func LinearlyFit(p0, p1 Point) (float64, float64) {
	slope := (p1.Pressure - p0.Pressure) / (p1.Flow - p0.Flow)
	return slope, SegmentIntercept(p0, p1, slope)
}
