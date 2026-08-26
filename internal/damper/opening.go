package damper

import (
	"fmt"
	"math"
)

type Blade int

const (
	Opposed Blade = iota
	Parallel
)

func ParseBlade(s string) (Blade, bool) {
	switch s {
	case "", "opposed":
		return Opposed, true
	case "parallel":
		return Parallel, true
	}
	return Opposed, false
}

func (b Blade) String() string {
	if b == Parallel {
		return "parallel"
	}
	return "opposed"
}

func (b Blade) exponent() float64 {
	if b == Parallel {
		return 1.6
	}
	return 2.5
}

type Config struct {
	Opening  float64
	KOpen    float64
	KShut    float64
	Blade    Blade
	Diameter float64
	Density  float64
}

func (c Config) Validate() error {
	if c.Opening <= 0 || c.Opening > 1 {
		return fmt.Errorf("damper: opening must be in (0, 1] (got %v)", c.Opening)
	}
	if c.KOpen < 0 {
		return fmt.Errorf("damper: open-loss coefficient must be >= 0 (got %v)", c.KOpen)
	}
	if c.KShut <= c.KOpen {
		return fmt.Errorf("damper: shut-loss coefficient must exceed open value")
	}
	if c.Diameter <= 0 {
		return fmt.Errorf("damper: diameter must be > 0 (got %v)", c.Diameter)
	}
	if c.Density <= 0 {
		return fmt.Errorf("damper: density must be > 0 (got %v)", c.Density)
	}
	return nil
}

func LossCoeff(opening float64, kOpen, kShut float64, blade Blade) (float64, error) {
	if opening <= 0 || opening > 1 {
		return 0, fmt.Errorf("damper: opening must be in (0, 1] (got %v)", opening)
	}
	if kOpen < 0 || kShut <= kOpen {
		return 0, fmt.Errorf("damper: Kshut must exceed Kopen >= 0")
	}
	closed := 1 - opening
	n := blade.exponent()
	k := kOpen / math.Pow(opening, n)
	cap := kOpen + (kShut-kOpen)*math.Pow(closed, 0.4)
	if k > cap {
		k = cap
	}
	if k > kShut {
		k = kShut
	}
	return k, nil
}

type Device struct {
	cfg  Config
	k    float64
	area float64
}

func New(cfg Config) (*Device, error) {
	if cfg.KOpen == 0 && cfg.KShut == 0 {
		cfg.KOpen = 0.2
		cfg.KShut = 8000
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	k, err := LossCoeff(cfg.Opening, cfg.KOpen, cfg.KShut, cfg.Blade)
	if err != nil {
		return nil, err
	}
	return &Device{cfg: cfg, k: k, area: math.Pi * cfg.Diameter * cfg.Diameter / 4}, nil
}

func (d *Device) Coefficient() float64 { return d.k }

func (d *Device) Opening() float64 { return d.cfg.Opening }

func (d *Device) Area() float64 { return d.area }

func (d *Device) Config() Config { return d.cfg }
