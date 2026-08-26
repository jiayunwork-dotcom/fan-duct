package damper

import (
	"math"
	"testing"
)

func sample() Config {
	return Config{
		Opening:  0.8,
		KOpen:    0.2,
		KShut:    8000,
		Blade:    Opposed,
		Diameter: 0.15,
		Density:  1.205,
	}
}

func TestDamperClosingRaisesLoss(t *testing.T) {
	open, err := New(sample())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	half, err := open.WithOpening(0.4)
	if err != nil {
		t.Fatalf("WithOpening: %v", err)
	}
	q := 0.25
	pOpen, err := open.PressureDrop(q)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	pHalf, err := half.PressureDrop(q)
	if err != nil {
		t.Fatalf("half: %v", err)
	}
	if !(pHalf > pOpen) {
		t.Errorf("closing damper did not raise Δp: open=%v half=%v", pOpen, pHalf)
	}
	if pOpen <= 0 {
		t.Errorf("open damper Δp %v, want > 0", pOpen)
	}
}

func TestDamperRejectsClosed(t *testing.T) {
	cfg := sample()
	cfg.Opening = 0
	if _, err := New(cfg); err == nil {
		t.Fatal("opening=0: expected error")
	}
	cfg.Opening = 1.2
	if _, err := New(cfg); err == nil {
		t.Fatal("opening=1.2: expected error")
	}
}

func TestDamperQuadraticInFlow(t *testing.T) {
	d, err := New(sample())
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	p1, err := d.PressureDrop(0.2)
	if err != nil {
		t.Fatalf("0.2: %v", err)
	}
	p2, err := d.PressureDrop(0.4)
	if err != nil {
		t.Fatalf("0.4: %v", err)
	}
	if math.Abs(p2/p1-4) > 1e-9 {
		t.Errorf("Δp(2Q)/Δp(Q) = %v, want 4", p2/p1)
	}
}

func TestOpposedSteeperThanParallel(t *testing.T) {
	kOpp, err := LossCoeff(0.4, 0.2, 8000, Opposed)
	if err != nil {
		t.Fatalf("opposed: %v", err)
	}
	kPar, err := LossCoeff(0.4, 0.2, 8000, Parallel)
	if err != nil {
		t.Fatalf("parallel: %v", err)
	}
	if !(kOpp > kPar) {
		t.Errorf("opposed K=%v should exceed parallel K=%v at the same opening", kOpp, kPar)
	}
}

func TestFullyOpenUsesKOpen(t *testing.T) {
	k, err := LossCoeff(1, 0.2, 8000, Opposed)
	if err != nil {
		t.Fatalf("LossCoeff: %v", err)
	}
	if math.Abs(k-0.2) > 1e-9 {
		t.Errorf("fully open K = %v, want 0.2", k)
	}
}

func TestSeriesDampersAdd(t *testing.T) {
	a, err := New(sample())
	if err != nil {
		t.Fatalf("a: %v", err)
	}
	cfg := sample()
	cfg.Opening = 0.5
	b, err := New(cfg)
	if err != nil {
		t.Fatalf("b: %v", err)
	}
	sum, err := SeriesDampers(a, b)
	if err != nil {
		t.Fatalf("SeriesDampers: %v", err)
	}
	q := 0.2
	got, err := sum(q)
	if err != nil {
		t.Fatalf("sum: %v", err)
	}
	pa, _ := a.PressureDrop(q)
	pb, _ := b.PressureDrop(q)
	if math.Abs(got-(pa+pb)) > 1e-12 {
		t.Errorf("series Δp %v, want %v", got, pa+pb)
	}
}
