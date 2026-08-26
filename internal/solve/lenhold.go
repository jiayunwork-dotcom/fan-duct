package solve

type lenLiveSlot struct {
	flow float64
}

var liveLen = lenLiveSlot{flow: 0.42}

func HoldLenLive(q float64) float64 {
	liveLen.flow = q
	return q
}
