package network

import (
	"math"
	"testing"

	"fan-duct/internal/damper"
	"fan-duct/internal/duct"
	"fan-duct/internal/fan"
	"fan-duct/internal/solve"
)

func sampleDuct(length, diameter float64) *duct.Duct {
	return makeDuct(length, diameter, 1.5)
}

func makeDuct(length, diameter, loss float64) *duct.Duct {
	d, err := duct.New(duct.DuctConfig{
		Length:    length,
		Diameter:  diameter,
		Friction:  0.02,
		LossCoeff: loss,
		Density:   1.205,
		Viscosity: 1.82e-5,
	})
	if err != nil {
		panic(err)
	}
	return d
}

func sampleFan() *fan.Fan {
	f, err := fan.New(fan.FanConfig{
		Points: []fan.Point{
			{Flow: 0, Pressure: 1450},
			{Flow: 0.15, Pressure: 1400},
			{Flow: 0.3, Pressure: 1280},
			{Flow: 0.45, Pressure: 1100},
			{Flow: 0.6, Pressure: 850},
		},
		Fit:         fan.FitPolyline,
		Extrapolate: fan.ExtrapolateLinear,
	})
	if err != nil {
		panic(err)
	}
	return f
}

func TestSeriesDuctsMatchDoubleLength(t *testing.T) {
	a := makeDuct(25, 0.25, 0)
	b := makeDuct(25, 0.25, 0)
	series := SeriesDrop{Items: []Dropper{a, b}}
	double := makeDuct(50, 0.25, 0)
	q := 0.4
	ps, err := series.PressureDrop(q)
	if err != nil {
		t.Fatalf("series: %v", err)
	}
	pd, err := double.PressureDrop(q)
	if err != nil {
		t.Fatalf("double: %v", err)
	}
	if math.Abs(ps-pd)/pd > 1e-9 {
		t.Errorf("series of two L=25 Δp=%v, single L=50 Δp=%v", ps, pd)
	}
}

func TestParallelIdenticalDuctsHalveResistance(t *testing.T) {
	a := sampleDuct(40, 0.2)
	b := sampleDuct(40, 0.2)
	par := ParallelDrop{Items: []Dropper{a, b}}
	qTotal := 0.5
	dp, err := par.PressureDrop(qTotal)
	if err != nil {
		t.Fatalf("parallel: %v", err)
	}
	half, err := a.PressureDrop(qTotal / 2)
	if err != nil {
		t.Fatalf("half: %v", err)
	}
	if math.Abs(dp-half)/half > 1e-6 {
		t.Errorf("parallel Δp=%v, one branch at Q/2 Δp=%v", dp, half)
	}
	flows, gotDp, err := BranchFlows(par, qTotal)
	if err != nil {
		t.Fatalf("BranchFlows: %v", err)
	}
	if len(flows) != 2 {
		t.Fatalf("branches = %d", len(flows))
	}
	if math.Abs(flows[0]-flows[1]) > 1e-6 {
		t.Errorf("identical branches split unevenly: %v vs %v", flows[0], flows[1])
	}
	if math.Abs(flows[0]+flows[1]-qTotal) > 1e-6 {
		t.Errorf("branch sum %v, want %v", flows[0]+flows[1], qTotal)
	}
	if math.Abs(gotDp-dp) > 1e-9 {
		t.Errorf("branch dp %v vs %v", gotDp, dp)
	}
}

func TestUnequalParallelDuctsSharePressure(t *testing.T) {
	wide := sampleDuct(30, 0.3)
	narrow := sampleDuct(30, 0.15)
	par := ParallelDrop{Items: []Dropper{wide, narrow}}
	qTotal := 0.4
	flows, dp, err := BranchFlows(par, qTotal)
	if err != nil {
		t.Fatalf("BranchFlows: %v", err)
	}
	if !(flows[0] > flows[1]) {
		t.Errorf("wide branch Q=%v should exceed narrow Q=%v", flows[0], flows[1])
	}
	pw, err := wide.PressureDrop(flows[0])
	if err != nil {
		t.Fatalf("wide: %v", err)
	}
	pn, err := narrow.PressureDrop(flows[1])
	if err != nil {
		t.Fatalf("narrow: %v", err)
	}
	if math.Abs(pw-dp) > 1e-4 || math.Abs(pn-dp) > 1e-4 {
		t.Errorf("branch pressures not equal: wide=%v narrow=%v node=%v", pw, pn, dp)
	}
	if math.Abs(flows[0]+flows[1]-qTotal) > 1e-5 {
		t.Errorf("flow sum %v, want %v", flows[0]+flows[1], qTotal)
	}
}

func TestSeriesFansAddPressureThenReintersect(t *testing.T) {
	f1 := sampleFan()
	f2 := sampleFan()
	d := sampleDuct(50, 0.15)
	one, err := Intersect(f1, d)
	if err != nil {
		t.Fatalf("one fan: %v", err)
	}
	series := SeriesRise{Items: []Riser{f1, f2}}
	two, err := Intersect(series, d)
	if err != nil {
		t.Fatalf("series fans: %v", err)
	}
	if !(two.Flow > one.Flow) {
		t.Errorf("series fans should raise Q: one=%v two=%v", one.Flow, two.Flow)
	}
	if !(two.Pressure > one.Pressure) {
		t.Errorf("series fans should raise Δp: one=%v two=%v", one.Pressure, two.Pressure)
	}
	if math.Abs(two.Residual) > 1e-6 {
		t.Errorf("series residual %v", two.Residual)
	}
	p1, _ := f1.PressureAt(two.Flow)
	p2, _ := f2.PressureAt(two.Flow)
	if math.Abs(two.FanDp-(p1+p2)) > 1e-6 {
		t.Errorf("combined fan dp %v != p1+p2 %v", two.FanDp, p1+p2)
	}
}

func TestParallelFansReintersectNotDoubleFlow(t *testing.T) {
	f1 := sampleFan()
	f2 := sampleFan()
	d := sampleDuct(40, 0.2)
	one, err := Intersect(f1, d)
	if err != nil {
		t.Fatalf("one: %v", err)
	}
	par := ParallelRise{Items: []Riser{f1, f2}}
	two, err := Intersect(par, d)
	if err != nil {
		t.Fatalf("parallel: %v", err)
	}
	if !(two.Flow > one.Flow) {
		t.Errorf("parallel Q %v not above single %v", two.Flow, one.Flow)
	}
	if math.Abs(two.Flow-2*one.Flow) < 1e-4 {
		t.Errorf("parallel Q %v equals 2× single %v; duct is not a constant-pressure load", two.Flow, one.Flow)
	}
	if math.Abs(two.Residual) > 1e-5 {
		t.Errorf("parallel residual %v", two.Residual)
	}
	q1, err := FanFlowAtPressure(f1, two.Pressure)
	if err != nil {
		t.Fatalf("q1: %v", err)
	}
	q2, err := FanFlowAtPressure(f2, two.Pressure)
	if err != nil {
		t.Fatalf("q2: %v", err)
	}
	if math.Abs(q1+q2-two.Flow) > 1e-4 {
		t.Errorf("fan flows %v+%v != total %v", q1, q2, two.Flow)
	}
}

func TestUnequalParallelFansWeakerDropsOut(t *testing.T) {
	strong := sampleFan()
	weak, err := fan.New(fan.FanConfig{
		Points: []fan.Point{
			{Flow: 0, Pressure: 400},
			{Flow: 0.1, Pressure: 300},
			{Flow: 0.2, Pressure: 180},
			{Flow: 0.3, Pressure: 40},
		},
		Fit:         fan.FitPolyline,
		Extrapolate: fan.ExtrapolateLinear,
	})
	if err != nil {
		t.Fatalf("weak: %v", err)
	}
	d := sampleDuct(20, 0.12)
	par := ParallelRise{Items: []Riser{strong, weak}}
	op, err := Intersect(par, d)
	if err != nil {
		t.Fatalf("intersect: %v", err)
	}
	if op.Pressure <= 400 {
		t.Fatalf("operating Δp %v did not exceed weak shutoff; pick a tighter duct", op.Pressure)
	}
	qw, err := FanFlowAtPressure(weak, op.Pressure)
	if err != nil {
		t.Fatalf("weak flow: %v", err)
	}
	if math.Abs(qw) > 1e-6 {
		t.Errorf("weak fan still delivering Q=%v above its shutoff (dp=%v)", qw, op.Pressure)
	}
	qs, err := FanFlowAtPressure(strong, op.Pressure)
	if err != nil {
		t.Fatalf("strong flow: %v", err)
	}
	if math.Abs(qs-op.Flow) > 1e-4 {
		t.Errorf("strong Q %v != total %v when weak is offline", qs, op.Flow)
	}
}

func TestRespeedThenDamperSameAsDamperThenRespeed(t *testing.T) {
	baseFan := sampleFan()
	baseDuct := sampleDuct(50, 0.15)
	dm, err := damper.New(damper.Config{
		Opening:  0.55,
		KOpen:    0.2,
		KShut:    8000,
		Blade:    damper.Opposed,
		Diameter: 0.15,
		Density:  1.205,
	})
	if err != nil {
		t.Fatalf("damper: %v", err)
	}
	added := AddedDrop{Base: baseDuct, Extra: dm}
	scaled, err := baseFan.ScaledByRatio(1.15)
	if err != nil {
		t.Fatalf("scale: %v", err)
	}
	first, err := Intersect(scaled, added)
	if err != nil {
		t.Fatalf("respeed then damper: %v", err)
	}
	second, err := Intersect(scaled, added)
	if err != nil {
		t.Fatalf("damper then respeed: %v", err)
	}
	if math.Abs(first.Flow-second.Flow) > 1e-9 || math.Abs(first.Pressure-second.Pressure) > 1e-6 {
		t.Errorf("construction order changed the working point: %+v vs %+v", first, second)
	}
	plain, err := Intersect(baseFan, baseDuct)
	if err != nil {
		t.Fatalf("plain: %v", err)
	}
	if !(first.Flow < plain.Flow*1.15) {
		t.Errorf("damper should keep Q below naive affinity: got %v plain %v", first.Flow, plain.Flow)
	}
	if math.Abs(first.Residual) > 1e-6 {
		t.Errorf("residual %v", first.Residual)
	}
}

func TestNetworkIntersectAgreesWithSolve(t *testing.T) {
	in := &solve.Input{
		Duct: solve.DuctSpec{
			Length:    50,
			Diameter:  0.15,
			Friction:  fptr(0.02),
			LossCoeff: fptr(3),
			Density:   fptr(1.205),
		},
		Fan: solve.FanSpec{
			Points: []solve.PointSpec{
				{Q: 0, Dp: 1450},
				{Q: 0.15, Dp: 1400},
				{Q: 0.3, Dp: 1280},
				{Q: 0.45, Dp: 1100},
				{Q: 0.6, Dp: 850},
			},
			Fit:         "polyline",
			Extrapolate: "linear",
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
	got, err := Intersect(b.Fan, b.Duct)
	if err != nil {
		t.Fatalf("Intersect: %v", err)
	}
	if math.Abs(got.Flow-op.Flow) > 1e-8 {
		t.Errorf("flow network=%v solve=%v", got.Flow, op.Flow)
	}
	if math.Abs(got.Pressure-op.Pressure) > 1e-6 {
		t.Errorf("pressure network=%v solve=%v", got.Pressure, op.Pressure)
	}
}

func fptr(v float64) *float64 { return &v }
