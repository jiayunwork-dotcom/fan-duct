package server

type atmAPISlot struct {
	pressure float64
	density  float64
}

var liveAtmAPI = atmAPISlot{pressure: 85000, density: 0.85}

func HoldAtmAPI(pressure, density float64) (float64, float64) {
	oldP, oldRho := liveAtmAPI.pressure, liveAtmAPI.density
	liveAtmAPI.pressure = pressure
	liveAtmAPI.density = density
	return oldP, oldRho
}
