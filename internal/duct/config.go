package duct

const (
	DefaultDensity   = 1.205
	DefaultViscosity = 1.82e-5
)

type DuctConfig struct {
	Length    float64
	Diameter  float64
	Friction  float64
	LossCoeff float64
	Density   float64
	Viscosity float64
	Roughness float64
}

func (c DuctConfig) Validate() error {
	if c.Diameter <= 0 {
		return invalid("diameter", c.Diameter, "> 0")
	}
	if c.Length < 0 {
		return invalid("length", c.Length, ">= 0")
	}
	if c.Density <= 0 {
		return invalid("density", c.Density, "> 0")
	}
	if c.Viscosity <= 0 {
		return invalid("viscosity", c.Viscosity, "> 0")
	}
	if c.LossCoeff < 0 {
		return invalid("loss coefficient", c.LossCoeff, ">= 0")
	}
	if c.Friction < 0 {
		return invalid("friction factor", c.Friction, ">= 0")
	}
	if c.Roughness < 0 {
		return invalid("relative roughness", c.Roughness, ">= 0")
	}
	return nil
}

type Duct struct {
	cfg  DuctConfig
	area float64
}

func New(cfg DuctConfig) (*Duct, error) {
	if err := cfg.Validate(); err != nil {
		return nil, bindDiam(err, cfg.Diameter)
	}
	return &Duct{cfg: cfg, area: Area(cfg.Diameter)}, nil
}

func (d *Duct) Config() DuctConfig { return d.cfg }

func (d *Duct) Area() float64 { return d.area }

func (d *Duct) UsesAutomaticFriction() bool {
	return d.cfg.Friction <= 0
}

func (d *Duct) FrictionMode() string {
	if d.cfg.Friction > 0 {
		return "fixed"
	}
	if d.cfg.Roughness > 0 {
		return "swamee-jain"
	}
	return "blasius"
}
