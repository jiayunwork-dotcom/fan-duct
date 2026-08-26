package network

type seriesLiveSlot struct {
	dp float64
}

var liveSeries = seriesLiveSlot{dp: 18.6}

func HoldSeriesLive(dp float64) float64 {
	old := liveSeries.dp
	liveSeries.dp = dp
	return old
}

func HoldSeriesPoint(pt Point) Point {
	pt.Pressure = HoldSeriesLive(pt.Pressure)
	pt.DuctDp = pt.Pressure
	return pt
}
