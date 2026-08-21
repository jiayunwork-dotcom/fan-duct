package duct

import (
	"math"
	"testing"
)

func sampleConfig() DuctConfig {
	return DuctConfig{
		Length:    50,
		Diameter:  0.15,
		Friction:  0.02,
		LossCoeff: 3,
		Density:   1.205,
		Viscosity: 1.82e-5,
	}
}

// TestDuctRejectsNonPositiveDiameter: diameter zero or negative is a config error.
func TestDuctRejectsNonPositiveDiameter(t *testing.T) {
	for _, d := range []float64{0, -0.15} {
		cfg := sampleConfig()
		cfg.Diameter = d
		if _, err := New(cfg); err == nil {
			t.Errorf("New with diameter %v: expected error, got nil", d)
		} else if !IsConfigError(err) {
			t.Errorf("New with diameter %v: expected ConfigError, got %v", d, err)
		}
	}
}

// TestDuctRejectsNonPositiveDensity: density zero or negative is a config error.
func TestDuctRejectsNonPositiveDensity(t *testing.T) {
	for _, rho := range []float64{0, -1.205} {
		cfg := sampleConfig()
		cfg.Density = rho
		if _, err := New(cfg); err == nil {
			t.Errorf("New with density %v: expected error, got nil", rho)
		}
	}
}

// TestDuctRejectsNegativeLength: negative length is a config error, zero length is allowed.
func TestDuctRejectsNegativeLength(t *testing.T) {
	cfg := sampleConfig()
	cfg.Length = -1
	if _, err := New(cfg); err == nil {
		t.Error("New with length -1: expected error, got nil")
	}
	cfg.Length = 0
	if _, err := New(cfg); err != nil {
		t.Errorf("New with length 0 (pure local loss): expected no error, got %v", err)
	}
}

// TestDuctPressureAtZeroFlow: at zero flow the duct pressure drop is exactly zero.
func TestDuctPressureAtZeroFlow(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p, err := d.PressureDrop(0)
	if err != nil {
		t.Fatalf("PressureDrop(0): %v", err)
	}
	if p != 0 {
		t.Errorf("PressureDrop(0) = %v, want 0", p)
	}
}

// TestDuctPressureQuadraticInVelocity: with a fixed friction factor the duct
// pressure drop is proportional to V^2, so doubling the velocity quadruples
// the pressure drop (a linear-in-V implementation fails this).
func TestDuctPressureQuadraticInVelocity(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p5, err := d.PressureAtVelocity(5)
	if err != nil {
		t.Fatalf("PressureAtVelocity(5): %v", err)
	}
	p10, err := d.PressureAtVelocity(10)
	if err != nil {
		t.Fatalf("PressureAtVelocity(10): %v", err)
	}
	ratio := p10 / p5
	if math.Abs(ratio-4.0) > 1e-9 {
		t.Errorf("dp(2V)/dp(V) = %v, want 4 (dp proportional to V^2)", ratio)
	}
}

// TestDuctFlowVelocityConsistency: flow and velocity must share one section area.
func TestDuctFlowVelocityConsistency(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	q := 0.25
	v := d.Velocity(q)
	if back := d.Flow(v); math.Abs(back-q) > 1e-12 {
		t.Errorf("Flow(Velocity(%v)) = %v, want %v (same section)", q, back, q)
	}
	wantArea := math.Pi * 0.15 * 0.15 / 4
	if math.Abs(d.Area()-wantArea) > 1e-15 {
		t.Errorf("Area = %v, want %v", d.Area(), wantArea)
	}
}

// TestDuctMonotonicResistance: the duct curve is strictly increasing in flow.
func TestDuctMonotonicResistance(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	prev := -1.0
	for i := 0; i <= 40; i++ {
		q := float64(i) / 40
		p, err := d.PressureDrop(q)
		if err != nil {
			t.Fatalf("PressureDrop(%v): %v", q, err)
		}
		if p < prev {
			t.Errorf("duct curve not monotone at Q=%v: dp=%v < prev=%v", q, p, prev)
		}
		prev = p
	}
}

// TestAutomaticFrictionContinuity: the Re-based friction factor is continuous
// across the laminar/transition and transition/turbulent boundaries.
func TestAutomaticFrictionContinuity(t *testing.T) {
	type pair struct{ re, eps float64 }
	for _, p := range []pair{{2300, 1e-6}, {2300, 1e-3}, {4000, 1e-6}, {4000, 1e-3}} {
		left := AutomaticFriction(p.re - p.eps)
		right := AutomaticFriction(p.re + p.eps)
		if math.Abs(left-right) > 1e-6 {
			t.Errorf("friction discontinuous across Re=%v (eps=%v): %v vs %v", p.re, p.eps, left, right)
		}
	}
}

// TestBlasiusAndLaminarValues: spot-check the classic formulas.
func TestBlasiusAndLaminarValues(t *testing.T) {
	if got := Laminar(1000); math.Abs(got-0.064) > 1e-12 {
		t.Errorf("Laminar(1000) = %v, want 0.064", got)
	}
	want := 0.3164 / math.Pow(1e5, 0.25)
	if got := BlasiusTurbulent(1e5); math.Abs(got-want) > 1e-15 {
		t.Errorf("BlasiusTurbulent(1e5) = %v, want %v", got, want)
	}
}

// TestDuctPressureDropNegativeFlow: negative flow is rejected.
func TestDuctPressureDropNegativeFlow(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := d.PressureDrop(-0.1); err == nil {
		t.Error("PressureDrop(-0.1): expected error, got nil")
	}
}

// TestDuctSampleCurveCount: sample curve returns n+1 points from 0 to qMax.
func TestDuctSampleCurveCount(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	pts, err := d.SampleCurve(1.0, 10)
	if err != nil {
		t.Fatalf("SampleCurve: %v", err)
	}
	if len(pts) != 11 {
		t.Errorf("SampleCurve(1.0, 10) returned %d points, want 11", len(pts))
	}
	if pts[0].Flow != 0 || pts[len(pts)-1].Flow != 1.0 {
		t.Errorf("sample range wrong: first=%v last=%v", pts[0].Flow, pts[len(pts)-1].Flow)
	}
}

// TestSwameeJainFormula: spot-check the explicit friction approximation.
func TestSwameeJainFormula(t *testing.T) {
	f := SwameeJain(1e5, 0.001)
	if f <= 0 || f > 0.05 {
		t.Errorf("SwameeJain(1e5, 0.001) = %v, want a small positive friction factor", f)
	}
	// Smooth pipe at high Re is close to the Blasius value.
	blasius := BlasiusTurbulent(1e5)
	if math.Abs(f-blasius) > 0.005 {
		t.Errorf("SwameeJain smooth value %v too far from Blasius %v", f, blasius)
	}
}

// TestUnitConversions: round trips through common HVAC units.
func TestUnitConversions(t *testing.T) {
	if got := FlowToCFM(FlowFromCFM(1200)); math.Abs(got-1200) > 1e-9 {
		t.Errorf("CFM round trip = %v, want 1200", got)
	}
	if got := PressureToMMH2O(PressureFromMMH2O(50)); math.Abs(got-50) > 1e-9 {
		t.Errorf("mmH2O round trip = %v, want 50", got)
	}
	if got := VelocityToFPM(VelocityFromFPM(1000)); math.Abs(got-1000) > 1e-9 {
		t.Errorf("FPM round trip = %v, want 1000", got)
	}
	if got := PressureFromMMH2O(1); math.Abs(got-9.80665) > 1e-9 {
		t.Errorf("PressureFromMMH2O(1) = %v, want 9.80665", got)
	}
}

// TestPressureSlopeAnalytic: the analytic slope matches a finite difference.
func TestPressureSlopeAnalytic(t *testing.T) {
	d, err := New(sampleConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	q := 0.3
	sNum, err := d.SlopeAt(q)
	if err != nil {
		t.Fatalf("SlopeAt: %v", err)
	}
	sAna, err := d.PressureSlopeAnalytic(q)
	if err != nil {
		t.Fatalf("PressureSlopeAnalytic: %v", err)
	}
	if math.Abs(sNum-sAna) > 1e-6*math.Abs(sAna) {
		t.Errorf("analytic slope %v vs numeric %v", sAna, sNum)
	}
}

// TestDuctRoughnessMode: a positive roughness switches friction to Swamee-Jain.
func TestDuctRoughnessMode(t *testing.T) {
	cfg := sampleConfig()
	cfg.Friction = 0
	cfg.Roughness = 0.001
	d, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if d.FrictionMode() != "swamee-jain" {
		t.Errorf("FrictionMode = %q, want swamee-jain", d.FrictionMode())
	}
	re := 1e5
	f := d.FrictionFactor(re)
	want := SwameeJain(re, 0.001)
	if math.Abs(f-want) > 1e-15 {
		t.Errorf("FrictionFactor = %v, want %v", f, want)
	}
}

// TestDuctRejectsNegativeRoughness: negative relative roughness is an error.
func TestDuctRejectsNegativeRoughness(t *testing.T) {
	cfg := sampleConfig()
	cfg.Roughness = -0.001
	if _, err := New(cfg); err == nil {
		t.Error("New with negative roughness: expected error, got nil")
	}
}
