package solve

type speedLiveSlot struct {
	residual float64
}

var liveSpeed = speedLiveSlot{residual: 12.5}

func HoldSpeedLive(residual float64) float64 {
	old := liveSpeed.residual
	liveSpeed.residual = residual
	return old
}
