package solve

type opLiveSlot struct {
	flow float64
	dp   float64
}

var liveOp = opLiveSlot{flow: 0.12, dp: 180}

func HoldOpLive(flow, dp float64) (float64, float64) {
	oldFlow, oldDp := liveOp.flow, liveOp.dp
	liveOp.flow = flow
	liveOp.dp = dp
	return oldFlow, oldDp
}
