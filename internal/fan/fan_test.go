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

func TestFanRejectsEmptyCurve(t *testing.T) {
	for _, n := range []int{0, 1} {
		cfg := polylineConfig()
		cfg.Points = cfg.Points[:n]
		if _, err := New(cfg); err == nil {
			t.Errorf("New with %d sample points: expected error, got nil", n)
		}
	}
}

func TestFanRejectsNonIncreasingFlow(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[2].Flow = cfg.Points[1].Flow
	if _, err := New(cfg); err == nil {
		t.Error("New with non-increasing flow: expected error, got nil")
	}
}

func TestFanRejectsFirstPointNotZero(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[0].Flow = 0.02
	if _, err := New(cfg); err == nil {
		t.Error("New with first flow != 0: expected error, got nil")
	}
}

func TestFanRejectsNegativePressure(t *testing.T) {
	cfg := polylineConfig()
	cfg.Points[3].Pressure = -5
	if _, err := New(cfg); err == nil {
		t.Error("New with negative pressure: expected error, got nil")
	}
}

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
	if p, err := f.PressureAt(0.075); err != nil || math.Abs(p-1425) > 1e-9 {
		t.Errorf("PressureAt(0.075) = %v, %v; want 1425", p, err)
	}
}

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

func TestFanExtrapolationLinear(t *testing.T) {
	cfg := polylineConfig()
	cfg.Extrapolate = ExtrapolateLinear
	f, err := New(cfg)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	want := 850.0 + (850.0-1100.0)/(0.6-0.45)*(0.65-0.6)
	got, err := f.PressureAt(0.65)
	if err != nil {
		t.Fatalf("PressureAt(0.65): %v", err)
	}
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("linear extrapolation at 0.65 = %v, want %v", got, want)
	}
}

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

func TestFanShutoffPressure(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if got := f.ZeroFlowPressure(); got != 1450 {
		t.Errorf("ZeroFlowPressure = %v, want 1450", got)
	}
}

func TestFanZeroCrossing(t *testing.T) {
	f, err := New(polylineConfig())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	q, ok := f.ZeroCrossing(1425)
	if !ok {
		t.Fatal("ZeroCrossing(1425): expected a crossing")
	}
	want := 0.0 + (1425.0-1450.0)/(1400.0-1450.0)*(0.15-0.0)
	if math.Abs(q-want) > 1e-9 {
		t.Errorf("ZeroCrossing = %v, want %v", q, want)
	}
}

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

func TestFanEfficiencyLengthMismatch(t *testing.T) {
	cfg := polylineConfig()
	cfg.Efficiency = []float64{0.6, 0.7}
	if _, err := New(cfg); err == nil {
		t.Error("New with mismatched efficiency length: expected error, got nil")
	}
}
