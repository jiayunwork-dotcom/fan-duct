package fan

import (
	"math"
	"testing"
)

func polylineConfig() FanConfig {
	return FanConfig{
		Points: []Point{
			{Flow: 0, Pressure: 1450},
			{Flow: 0.15, Pressure: 1400},
			{Flow: 0.3, Pressure: 1280},
			{Flow: 0.45, Pressure: 1100},
			{Flow: 0.6, Pressure: 850},
		},
		Fit:         FitPolyline,
		Extrapolate: ExtrapolateError,
	}
}

// TestFanRejectsEmptyCurve: a fan curve with fewer than two points is an error.
func TestFanRejectsEmptyCurve(t *testing.T) {
	for _, n := range []int{0, 1} {
		cfg := polylineConfig()
		cfg.Points = cfg.Points[:n]
		if _, err := New(cfg); err == nil {
			t.Errorf("New with %d sample points: expected error, got nil", n)
		}
	}
}

// TestFanRejectsNonIncreasingFlow: sample flows must be strictly increasing.
func TestFanRejectsNonIncreasingFlow(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[2].Flow = cfg.Points[1].Flow
	if _, err := New(cfg); err == nil {
		t.Error("New with non-increasing flow: expected error, got nil")
	}
}

// TestFanRejectsFirstPointNotZero: the first sample point must sit at zero flow.
func TestFanRejectsFirstPointNotZero(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[0].Flow = 0.02
	if _, err := New(cfg); err == nil {
		t.Error("New with first flow != 0: expected error, got nil")
	}
}

// TestFanRejectsNegativePressure: negative pressure samples are rejected.
func TestFanRejectsNegativePressure(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[3].Pressure = -5
	if _, err := New(cfg); err == nil {
		t.Error("New with negative pressure: expected error, got nil")
	}
}

// TestFanPolylineInterpolation: the curve passes through samples and interpolates linearly.
func TestFanPolylineInterpolation(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if p, err := f.PressureAt(0); err != nil || p != 1450 {
		t.Errorf("PressureAt(0) = %v, %v; want 1450", p, err)
	}
	if p, err := f.PressureAt(0.3); err != nil || p != 1280 {
		t.Errorf("PressureAt(0.3) = %v, %v; want 1280", p, err)
	}
	// 0.075 is midway between 0 and 0.15: pressure should be 1425.
	if p, err := f.PressureAt(0.075); err != nil || math.Abs(p-1425) > 1e-9 {
		t.Errorf("PressureAt(0.075) = %v, %v; want 1425", p, err)
	}
}

// TestFanQuadraticFit: quadratic fit reproduces a parabola within tolerance.
func TestFanQuadraticFit(t *testing.T) {
	pts := []Point{
		{Flow: 0, Pressure: 1000},
		{Flow: 0.2, Pressure: 900},
		{Flow: 0.4, Pressure: 700},
		{Flow: 0.6, Pressure: 400},
	}
	a, b, c := FitQuadraticCoeffs(pts)
	f, err := New(FanConfig{Points: pts, Fit: FitQuadratic, Extrapolate: ExtrapolateError})
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	for _, p := range pts {
		got, err := f.PressureAt(p.Flow)
		if err != nil {
			t.Fatalf("PressureAt(%v): %v", p.Flow, err)
		}
		if math.Abs(got-p.Pressure) > 1.0 {
			t.Errorf("quadratic fit at Q=%v = %v, want ~%v (a=%v b=%v c=%v)", p.Flow, got, p.Pressure, a, b, c)
		}
	}
}

// TestFanExtrapolationError: flow outside the sample range errors when disabled.
func TestFanExtrapolationError(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if _, err := f.PressureAt(0.7); err == nil {
		t.Error("PressureAt(0.7): expected out-of-range error, got nil")
	} else if !IsOutOfRange(err) {
		t.Errorf("PressureAt(0.7): expected OutOfRangeError, got %v", err)
	}
}

// TestFanExtrapolationLinear: with linear extrapolation the edge slope is continued.
func TestFanExtrapolationLinear(t *testing.T) {
	cfg := polylineConfig()
	cfg.Extrapolate = ExtrapolateLinear
	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Last segment slope: (850-1100)/(0.6-0.45) = -1666.67; at Q=0.65 -> 850-1666.67*0.05.
	want := 850.0 + (850.0-1100.0)/(0.6-0.45)*(0.65-0.6)
	got, err := f.PressureAt(0.65)
	if err != nil {
		t.Fatalf("PressureAt(0.65): %v", err)
	}
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("linear extrapolation at 0.65 = %v, want %v", got, want)
	}
}

// TestFanScalingLaws: affinity laws give Q2=Q1*r, dp2=dp1*r^2, P2=P1*r^3.
func TestFanScalingLaws(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	r := 1.1
	scaled, err := f.ScaledByRatio(r)
	if err != nil {
		t.Fatalf("ScaledByRatio: %v", err)
	}
	q := 0.3
	dp, err := f.PressureAt(q)
	if err != nil {
		t.Fatalf("PressureAt: %v", err)
	}
	p1 := q * dp
	q2 := q * r
	dp2, err := scaled.PressureAt(q2)
	if err != nil {
		t.Fatalf("scaled PressureAt: %v", err)
	}
	p2 := q2 * dp2
	if math.Abs(q2-q*r) > 1e-12 {
		t.Errorf("scaled flow = %v, want %v", q2, q*r)
	}
	if math.Abs(dp2-dp*r*r) > 1e-9 {
		t.Errorf("scaled pressure = %v, want %v", dp2, dp*r*r)
	}
	if math.Abs(p2-p1*r*r*r) > 1e-6 {
		t.Errorf("scaled power = %v, want %v", p2, p1*r*r*r)
	}
}

// TestFanShutoffPressure: the zero-flow point pressure is the curve shutoff.
func TestFanShutoffPressure(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := f.ZeroFlowPressure(); got != 1450 {
		t.Errorf("ZeroFlowPressure = %v, want 1450", got)
	}
}

// TestFanZeroCrossing: the first sample where pressure drops to target.
func TestFanZeroCrossing(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	// Pressure reaches 1425 midway on the first segment (0 to 0.15).
	q, ok := f.ZeroCrossing(1425)
	if !ok {
		t.Fatal("ZeroCrossing(1425): expected a crossing")
	}
	want := 0.0 + (1425.0-1450.0)/(1400.0-1450.0)*(0.15-0.0)
	if math.Abs(q-want) > 1e-9 {
		t.Errorf("ZeroCrossing = %v, want %v", q, want)
	}
}

// TestFanInRange: in-range checks respect the sample bounds.
func TestFanInRange(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if !f.InRange(0.3) {
		t.Error("InRange(0.3): expected true")
	}
	if f.InRange(0.7) {
		t.Error("InRange(0.7): expected false")
	}
}

// TestFanFitQuality: a perfect parabola gives R^2 of 1.
func TestFanFitQuality(t *testing.T) {
	pts := []Point{
		{Flow: 0, Pressure: 1000},
		{Flow: 0.2, Pressure: 960},
		{Flow: 0.4, Pressure: 840},
		{Flow: 0.6, Pressure: 640},
	}
	r2, maxRes := FitQuality(pts)
	if math.Abs(r2-1) > 1e-9 {
		t.Errorf("FitQuality R^2 = %v, want 1 for a perfect parabola", r2)
	}
	if maxRes > 1e-9 {
		t.Errorf("FitQuality max residual = %v, want near 0", maxRes)
	}
}

// TestFanEfficiencyInterpolation: efficiency samples interpolate linearly.
func TestFanEfficiencyInterpolation(t *testing.T) {
	cfg := polylineConfig()
	cfg.Efficiency = []float64{0.6, 0.7, 0.65, 0.6, 0.5}
	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	e, ok := f.EfficiencyAt(0.075)
	if !ok {
		t.Fatal("EfficiencyAt(0.075): expected an interpolated value")
	}
	if math.Abs(e-0.65) > 1e-12 {
		t.Errorf("EfficiencyAt(0.075) = %v, want 0.65", e)
	}
}

// TestFanEfficiencyLengthMismatch: efficiency samples must match points length.
func TestFanEfficiencyLengthMismatch(t *testing.T) {
	cfg := polylineConfig()
	cfg.Efficiency = []float64{0.6, 0.7}
	if _, err := New(cfg); err == nil {
		t.Error("New with mismatched efficiency length: expected error, got nil")
	}
}
