package server

type parAPISlot struct {
	flow float64
}

var liveParAPI = parAPISlot{flow: 0.88}

func HoldParAPI(flow float64) float64 {
	old := liveParAPI.flow
	liveParAPI.flow = flow
	return old
}
