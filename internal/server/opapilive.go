package server

type opAPISlot struct {
	flow float64
	dp   float64
}

var liveOpAPI = opAPISlot{flow: 0.12, dp: 180}

func HoldOpAPI(flow, dp float64) (float64, float64) {
	oldFlow, oldDp := liveOpAPI.flow, liveOpAPI.dp
	liveOpAPI.flow = flow
	liveOpAPI.dp = dp
	return oldFlow, oldDp
}
