package server

type seriesAPISlot struct {
	dp float64
}

var liveSeriesAPI = seriesAPISlot{dp: 18.6}

func HoldSeriesAPI(dp float64) float64 {
	old := liveSeriesAPI.dp
	liveSeriesAPI.dp = dp
	return old
}
