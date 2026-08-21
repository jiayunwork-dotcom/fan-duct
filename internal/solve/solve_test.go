package solve

import (
	"math"
	"strings"
	"testing"
)

func exampleInput() *Input {
	return &Input{
		Duct: DuctSpec{
			Length:    50,
			Diameter:  0.15,
			Friction:  fptr(0.02),
			LossCoeff: fptr(3),
			Density:   fptr(1.205),
			Viscosity: fptr(1.82e-5),
		},
		Fan: FanSpec{
			Points: []PointSpec{
				{Q: 0, Dp: 1450},
				{Q: 0.15, Dp: 1400},
				{Q: 0.3, Dp: 1280},
				{Q: 0.45, Dp: 1100},
				{Q: 0.6, Dp: 850},
			},
			Fit:         "polyline",
			Extrapolate: "error",
		},
		Speed:    &SpeedSpec{RPM: 1450},
		NewSpeed: &SpeedSpec{RPM: 1595},
	}
}

func fptr(v float64) *float64 { return &v }

// TestOperatingPointOnBothCurves: the working point satisfies both the fan
// curve and the duct curve simultaneously within the root tolerance.
func TestOperatingPointOnBothCurves(t *testing.T) {
	b, err := exampleInput().Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	if math.Abs(op.FanDp-op.DuctDp) > 1e-6 {
		t.Errorf("fan dp %v vs duct dp %v differ by %v, want < 1e-6", op.FanDp, op.DuctDp, math.Abs(op.FanDp-op.DuctDp))
	}
	if math.Abs(op.Pressure-op.FanDp) > 1e-6 {
		t.Errorf("reported pressure %v does not match fan dp %v", op.Pressure, op.FanDp)
	}
	if math.Abs(op.Residual) > 1e-6 {
		t.Errorf("residual %v, want close to 0", op.Residual)
	}
}

// TestSpeedChangeReintersects: after raising the speed to 1.1N the new
// operating point must come from re-intersecting the scaled fan curve with
// the duct curve. With a Reynolds-dependent friction factor the duct curve
// is not a pure parabola, so the re-solved flow differs from the naive
// Q1*1.1 value and the residual stays near zero.
func TestSpeedChangeReintersects(t *testing.T) {
	in := &Input{
		Duct: DuctSpec{
			Length:    40,
			Diameter:  0.5,
			LossCoeff: fptr(2.5),
			Density:   fptr(1.205),
			Viscosity: fptr(1.82e-5),
			// Friction omitted -> Reynolds-based automatic value.
		},
		Fan: FanSpec{
			Points: []PointSpec{
				{Q: 0, Dp: 6},
				{Q: 0.15, Dp: 5.5},
				{Q: 0.3, Dp: 5},
				{Q: 0.45, Dp: 4},
				{Q: 0.6, Dp: 3},
			},
			Fit:         "polyline",
			Extrapolate: "error",
		},
		Speed: &SpeedSpec{RPM: 1000},
	}
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	rr, err := b.Respeed(1000, 1100)
	if err != nil {
		t.Fatalf("Respeed: %v", err)
	}
	if math.Abs(rr.Ratio-1.1) > 1e-12 {
		t.Errorf("ratio = %v, want 1.1", rr.Ratio)
	}
	op := rr.Respeeded
	if math.Abs(op.Residual) > 1e-6 {
		t.Errorf("respeeded residual %v, want close to 0 (must re-intersect)", op.Residual)
	}
	naive := rr.Base.Flow * 1.1
	if math.Abs(op.Flow-naive) < 1e-4 {
		t.Errorf("respeeded flow %v equals the naive Q1*1.1 = %v; the duct curve is not a parabola here", op.Flow, naive)
	}
}

// TestLongerDuctLowersFlow: lengthening the duct lowers the working flow and
// raises the pressure drop, consistent with a monotone duct curve.
func TestLongerDuctLowersFlow(t *testing.T) {
	fanSpec := FanSpec{
		Points: []PointSpec{
			{Q: 0, Dp: 1200},
			{Q: 0.2, Dp: 1150},
			{Q: 0.4, Dp: 1050},
			{Q: 0.6, Dp: 880},
			{Q: 0.8, Dp: 650},
			{Q: 1.0, Dp: 360},
		},
		Fit:         "polyline",
		Extrapolate: "error",
	}
	base := &Input{
		Duct: DuctSpec{Length: 25, Diameter: 0.25, Friction: fptr(0.02), LossCoeff: fptr(1.5), Density: fptr(1.205)},
		Fan:  fanSpec,
	}
	long := &Input{
		Duct: DuctSpec{Length: 50, Diameter: 0.25, Friction: fptr(0.02), LossCoeff: fptr(1.5), Density: fptr(1.205)},
		Fan:  fanSpec,
	}
	op1b, err := base.Build()
	if err != nil {
		t.Fatalf("Build(base): %v", err)
	}
	op2b, err := long.Build()
	if err != nil {
		t.Fatalf("Build(long): %v", err)
	}
	op1, err := op1b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint(base): %v", err)
	}
	op2, err := op2b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint(long): %v", err)
	}
	if !(op2.Flow < op1.Flow) {
		t.Errorf("longer duct flow %v not lower than %v", op2.Flow, op1.Flow)
	}
	if !(op2.Pressure > op1.Pressure) {
		t.Errorf("longer duct pressure %v not higher than %v", op2.Pressure, op1.Pressure)
	}
	if math.Abs(op1.Residual) > 1e-6 || math.Abs(op2.Residual) > 1e-6 {
		t.Errorf("residuals not near zero: base=%v long=%v", op1.Residual, op2.Residual)
	}
}

// TestHigherDensityLowersFlow: raising the air density increases the duct
// resistance at the same flow, so the working flow drops.
func TestHigherDensityLowersFlow(t *testing.T) {
	fanSpec := FanSpec{
		Points: []PointSpec{
			{Q: 0, Dp: 1200},
			{Q: 0.2, Dp: 1150},
			{Q: 0.4, Dp: 1050},
			{Q: 0.6, Dp: 880},
			{Q: 0.8, Dp: 650},
			{Q: 1.0, Dp: 360},
		},
		Fit:         "polyline",
		Extrapolate: "error",
	}
	air := &Input{
		Duct: DuctSpec{Length: 25, Diameter: 0.25, Friction: fptr(0.02), LossCoeff: fptr(1.5), Density: fptr(1.205)},
		Fan:  fanSpec,
	}
	dense := &Input{
		Duct: DuctSpec{Length: 25, Diameter: 0.25, Friction: fptr(0.02), LossCoeff: fptr(1.5), Density: fptr(1.5)},
		Fan:  fanSpec,
	}
	op1b, err := air.Build()
	if err != nil {
		t.Fatalf("Build(air): %v", err)
	}
	op2b, err := dense.Build()
	if err != nil {
		t.Fatalf("Build(dense): %v", err)
	}
	op1, err := op1b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint(air): %v", err)
	}
	op2, err := op2b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint(dense): %v", err)
	}
	if !(op2.Flow < op1.Flow) {
		t.Errorf("higher-density flow %v not lower than %v", op2.Flow, op1.Flow)
	}
}

// TestNoIntersectionAtZeroShutoff: a fan with non-positive shutoff pressure
// has no positive-flow operating point.
func TestNoIntersectionAtZeroShutoff(t *testing.T) {
	in := &Input{
		Duct: DuctSpec{Length: 50, Diameter: 0.15, Friction: fptr(0.02), LossCoeff: fptr(3)},
		Fan: FanSpec{
			Points:      []PointSpec{{Q: 0, Dp: 0}, {Q: 0.3, Dp: 100}},
			Fit:         "polyline",
			Extrapolate: "error",
		},
	}
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if _, err := b.OperatingPoint(); err == nil {
		t.Error("OperatingPoint: expected no-intersection error, got nil")
	} else if !IsSolveError(err) {
		t.Errorf("OperatingPoint: expected SolveError, got %v", err)
	}
}

// TestParseInputValid: a well-formed JSON parses without error.
func TestParseInputValid(t *testing.T) {
	json := `{
		"duct": {"length": 50, "diameter": 0.15, "friction": 0.02, "lossCoeff": 3},
		"fan": {"points": [{"q": 0, "dp": 1450}, {"q": 0.3, "dp": 1280}]},
		"speed": {"rpm": 1450}
	}`
	in, err := ParseInput([]byte(json))
	if err != nil {
		t.Fatalf("ParseInput: %v", err)
	}
	if in.Duct.Diameter != 0.15 {
		t.Errorf("parsed diameter = %v, want 0.15", in.Duct.Diameter)
	}
	if in.Speed == nil || in.Speed.RPM != 1450 {
		t.Errorf("parsed speed = %+v, want rpm 1450", in.Speed)
	}
}

// TestParseRejectsUnknownField: unknown JSON fields are rejected.
func TestParseRejectsUnknownField(t *testing.T) {
	json := `{
		"duct": {"length": 50, "diameter": 0.15, "banana": 1},
		"fan": {"points": [{"q": 0, "dp": 1}, {"q": 0.3, "dp": 0.5}]}
	}`
	if _, err := ParseInput([]byte(json)); err == nil {
		t.Error("ParseInput with unknown field: expected error, got nil")
	}
}

// TestParseRejectsMissingDuct: a missing duct section is rejected at build time.
func TestParseRejectsMissingDuct(t *testing.T) {
	json := `{"fan": {"points": [{"q": 0, "dp": 1}, {"q": 0.3, "dp": 0.5}]}}`
	in, err := ParseInput([]byte(json))
	if err != nil {
		t.Fatalf("ParseInput: %v", err)
	}
	if _, err := in.Build(); err == nil {
		t.Error("Build with missing duct: expected error, got nil")
	}
}

// TestRespeedRequiresBaseSpeed: respeed without a base speed in the input is an error.
func TestRespeedRequiresBaseSpeed(t *testing.T) {
	in := &Input{
		Duct:     DuctSpec{Length: 50, Diameter: 0.15, Friction: fptr(0.02)},
		Fan:      FanSpec{Points: []PointSpec{{Q: 0, Dp: 1450}, {Q: 0.3, Dp: 1280}}},
		NewSpeed: &SpeedSpec{RPM: 1595},
	}
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if _, err := in.RespeedToRPM(b, in.NewSpeed.RPM); err == nil {
		t.Error("RespeedToRPM without speed.rpm: expected error, got nil")
	}
}

// TestExampleInlineFanInSampleRange: the example operating flow sits strictly
// between two fan sample points (0.15 and 0.30 m3/s).
func TestExampleInlineFanInSampleRange(t *testing.T) {
	b, err := exampleInput().Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	if !(op.Flow > 0.15 && op.Flow < 0.30) {
		t.Errorf("example working flow %v not between sample points 0.15 and 0.30", op.Flow)
	}
}

// TestExampleOperateValues: spot-check the printed numbers of the example.
func TestExampleOperateValues(t *testing.T) {
	in := exampleInput()
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	if rel := math.Abs(op.Flow-0.264838) / 0.264838; rel > 1e-3 {
		t.Errorf("example flow %v, want ~0.264838", op.Flow)
	}
	if rel := math.Abs(op.Pressure-1308.13) / 1308.13; rel > 1e-3 {
		t.Errorf("example pressure %v, want ~1308.13", op.Pressure)
	}
	rr, err := b.Respeed(1450, 1595)
	if err != nil {
		t.Fatalf("Respeed: %v", err)
	}
	if rel := math.Abs(rr.Respeeded.Flow-0.291322) / 0.291322; rel > 1e-3 {
		t.Errorf("respeeded flow %v, want ~0.291322", rr.Respeeded.Flow)
	}
}

// TestOutputContainsKeyNumbers: the rendered output carries Q, V and dp.
func TestOutputContainsKeyNumbers(t *testing.T) {
	b, err := exampleInput().Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	out := Output{Input: exampleInput(), Build: b, Base: op}
	rr, err := b.Respeed(1450, 1595)
	if err != nil {
		t.Fatalf("Respeed: %v", err)
	}
	out.Respeeded = &rr
	s := out.String()
	for _, want := range []string{"operating point", "flow", "velocity", "respeed"} {
		if !strings.Contains(s, want) {
			t.Errorf("output missing %q", want)
		}
	}
}

// TestSensitivityLength: lengthening the duct lowers the working flow.
func TestSensitivityLength(t *testing.T) {
	b, err := exampleInput().Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	s, err := ComputeSensitivity(b, "length", 0.1)
	if err != nil {
		t.Fatalf("ComputeSensitivity: %v", err)
	}
	if !(s.NewQ < s.BaseQ) {
		t.Errorf("longer duct sensitivity flow %v not lower than %v", s.NewQ, s.BaseQ)
	}
	if s.RelativeQ >= 0 {
		t.Errorf("sensitivity relative flow %v, want negative for length+10%%", s.RelativeQ)
	}
}

// TestModelReport: the model description mentions both duct and fan.
func TestModelReport(t *testing.T) {
	b, err := exampleInput().Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	s, err := ModelString(b)
	if err != nil {
		t.Fatalf("ModelString: %v", err)
	}
	if !strings.Contains(s, "Duct model") || !strings.Contains(s, "Fan model") {
		t.Errorf("model report missing sections:\n%s", s)
	}
}

// TestQuadraticFitOperatingPoint: a quadratic fan curve still gives a valid intersection.
func TestQuadraticFitOperatingPoint(t *testing.T) {
	in := &Input{
		Duct: DuctSpec{Length: 50, Diameter: 0.15, Friction: fptr(0.02), LossCoeff: fptr(3)},
		Fan: FanSpec{
			Points: []PointSpec{
				{Q: 0, Dp: 1450},
				{Q: 0.15, Dp: 1400},
				{Q: 0.3, Dp: 1280},
				{Q: 0.45, Dp: 1100},
				{Q: 0.6, Dp: 850},
			},
			Fit:         "quadratic",
			Extrapolate: "error",
		},
	}
	b, err := in.Build()
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	op, err := b.OperatingPoint()
	if err != nil {
		t.Fatalf("OperatingPoint: %v", err)
	}
	if math.Abs(op.Residual) > 1e-6 {
		t.Errorf("quadratic-fit residual %v, want close to 0", op.Residual)
	}
}
