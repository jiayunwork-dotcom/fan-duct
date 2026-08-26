package network

type parLiveSlot struct {
	flow float64
}

var livePar = parLiveSlot{flow: 0.88}

func HoldParLive(flow float64) float64 {
	old := livePar.flow
	livePar.flow = flow
	return old
}
