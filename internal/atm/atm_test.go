package atm

import (
	"math"
	"testing"
)

func TestISADensityDecreasesWithAltitude(t *testing.T) {
	sea, err := NewISA(0)
	if err != nil {
		t.Fatalf("NewISA(0): %v", err)
	}
	high, err := NewISA(2000)
	if err != nil {
		t.Fatalf("NewISA(2000): %v", err)
	}
	if !(high.DensityKgM3 < sea.DensityKgM3) {
		t.Errorf("density at 2000 m %v not lower than sea level %v", high.DensityKgM3, sea.DensityKgM3)
	}
	if !(high.PressurePa < sea.PressurePa) {
		t.Errorf("pressure at 2000 m %v not lower than sea level %v", high.PressurePa, sea.PressurePa)
	}
	if math.Abs(sea.TemperatureK-288.15) > 1e-9 {
		t.Errorf("sea-level T = %v, want 288.15", sea.TemperatureK)
	}
	if math.Abs(sea.PressurePa-101325) > 1e-6 {
		t.Errorf("sea-level P = %v, want 101325", sea.PressurePa)
	}
}

func TestISARejectsNegativeAltitude(t *testing.T) {
	if _, err := NewISA(-1); err == nil {
		t.Fatal("NewISA(-1): expected error")
	}
	if _, err := NewISA(25000); err == nil {
		t.Fatal("NewISA(25000): expected error")
	}
}

func TestTropopauseContinuity(t *testing.T) {
	below, err := NewISA(TropopauseM - 1)
	if err != nil {
		t.Fatalf("below: %v", err)
	}
	above, err := NewISA(TropopauseM + 1)
	if err != nil {
		t.Fatalf("above: %v", err)
	}
	at, err := NewISA(TropopauseM)
	if err != nil {
		t.Fatalf("at: %v", err)
	}
	if math.Abs(at.TemperatureK-TropopauseTemperatureK) > 1e-6 {
		t.Errorf("tropopause T = %v, want %v", at.TemperatureK, TropopauseTemperatureK)
	}
	if math.Abs(below.PressurePa-above.PressurePa) > 40 {
		t.Errorf("pressure jump across tropopause: %v vs %v", below.PressurePa, above.PressurePa)
	}
}

func TestSutherlandIncreasesWithTemperature(t *testing.T) {
	mu273, err := SutherlandViscosity(273.15)
	if err != nil {
		t.Fatalf("273.15: %v", err)
	}
	mu288, err := SutherlandViscosity(288.15)
	if err != nil {
		t.Fatalf("288.15: %v", err)
	}
	if math.Abs(mu273-SutherlandRefMu) > 1e-12 {
		t.Errorf("μ(273.15) = %v, want %v", mu273, SutherlandRefMu)
	}
	if !(mu288 > mu273) {
		t.Errorf("μ(288.15)=%v not greater than μ(273.15)=%v", mu288, mu273)
	}
}

func TestMoistAirDensityBelowDry(t *testing.T) {
	dry, moist, ratio, err := MoistCorrection(0, 0.8)
	if err != nil {
		t.Fatalf("MoistCorrection: %v", err)
	}
	if !(moist < dry) {
		t.Errorf("moist density %v not below dry %v", moist, dry)
	}
	if ratio >= 1 || ratio <= 0.95 {
		t.Errorf("moist/dry ratio %v, want in (0.95, 1)", ratio)
	}
	same, err := MoistDensity(101325, 288.15, 0)
	if err != nil {
		t.Fatalf("RH=0: %v", err)
	}
	if math.Abs(same-dry)/dry > 1e-9 {
		t.Errorf("RH=0 density %v, want dry %v", same, dry)
	}
}

func TestHydrostaticIntegralMatchesPressureDrop(t *testing.T) {
	ratio, err := HydrostaticCheck(0, 1000, 200)
	if err != nil {
		t.Fatalf("HydrostaticCheck: %v", err)
	}
	if math.Abs(ratio-1) > 0.02 {
		t.Errorf("hydrostatic integral / ΔP = %v, want ~1", ratio)
	}
}

func TestViscosityAtAltitudeFollowsTemperature(t *testing.T) {
	mu0, err := ViscosityAtAltitude(0)
	if err != nil {
		t.Fatalf("0 m: %v", err)
	}
	mu8, err := ViscosityAtAltitude(8000)
	if err != nil {
		t.Fatalf("8000 m: %v", err)
	}
	if !(mu8 < mu0) {
		t.Errorf("colder air at 8000 m should be less viscous: %v vs %v", mu8, mu0)
	}
}
