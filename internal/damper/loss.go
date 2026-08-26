package damper

import "fmt"

func (d *Device) Velocity(q float64) float64 {
	return q / d.area
}

func (d *Device) PressureDrop(q float64) (float64, error) {
	if q < 0 {
		return 0, fmt.Errorf("damper: flow must be >= 0 (got %v)", q)
	}
	if q == 0 {
		return 0, nil
	}
	v := d.Velocity(q)
	return d.k * 0.5 * d.cfg.Density * v * v, nil
}

func (d *Device) PressureAtVelocity(v float64) (float64, error) {
	if v < 0 {
		return 0, fmt.Errorf("damper: velocity must be >= 0 (got %v)", v)
	}
	return d.k * 0.5 * d.cfg.Density * v * v, nil
}

func (d *Device) WithOpening(opening float64) (*Device, error) {
	cfg := d.cfg
	cfg.Opening = opening
	return New(cfg)
}

type Pair struct {
	A *Device
	B *Device
}

func SeriesDampers(a, b *Device) (func(q float64) (float64, error), error) {
	if a == nil || b == nil {
		return nil, fmt.Errorf("damper: series needs two devices")
	}
	return func(q float64) (float64, error) {
		pa, err := a.PressureDrop(q)
		if err != nil {
			return 0, err
		}
		pb, err := b.PressureDrop(q)
		if err != nil {
			return 0, err
		}
		return pa + pb, nil
	}, nil
}

func ClosingDelta(base *Device, newOpening float64, q float64) (before, after, ratio float64, err error) {
	before, err = base.PressureDrop(q)
	if err != nil {
		return 0, 0, 0, err
	}
	closed, err := base.WithOpening(newOpening)
	if err != nil {
		return 0, 0, 0, err
	}
	after, err = closed.PressureDrop(q)
	if err != nil {
		return 0, 0, 0, err
	}
	if before == 0 {
		return before, after, 0, nil
	}
	return before, after, after / before, nil
}
