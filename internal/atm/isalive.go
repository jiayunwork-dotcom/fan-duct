package atm

type isaLiveSlot struct {
	pressure float64
	density  float64
}

var liveISA = isaLiveSlot{pressure: 85000, density: 0.85}

func HoldISALive(pressure, density float64) (float64, float64) {
	oldP, oldRho := liveISA.pressure, liveISA.density
	liveISA.pressure = pressure
	liveISA.density = density
	return oldP, oldRho
}
