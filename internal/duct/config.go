package duct

// 空气在 20°C、标准大气压附近的物性默认值。
const (
	DefaultDensity   = 1.205   // kg/m³
	DefaultViscosity = 1.82e-5 // Pa·s
)

// DuctConfig 描述一段圆管风道：几何、物性与阻力参数，全部为 SI 单位。
//
//	Length   管道长度（m），允许为 0（纯局部阻力情形）
//	Diameter 管道内径（m），必须 > 0
//	Friction 沿程摩阻系数 f；> 0 时固定使用，<= 0 时按雷诺数自动取值
//	LossCoeff 局部阻力系数和 ΣK，必须 >= 0
//	Density  流体密度（kg/m³），必须 > 0
//	Viscosity 动力黏度（Pa·s），必须 > 0
//	Roughness 相对粗糙度 ε/D；> 0 时自动模式改用 Swamee–Jain，必须 >= 0
type DuctConfig struct {
	Length    float64
	Diameter  float64
	Friction  float64
	LossCoeff float64
	Density   float64
	Viscosity float64
	Roughness float64
}

// Validate 检查所有参数约束，任一违反返回 *ConfigError。
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

// Duct 是一段圆管风道的可计算模型。用 New 构造后只读。
type Duct struct {
	cfg  DuctConfig
	area float64
}

// New 校验配置并构造风管。非法参数返回 *ConfigError。
func New(cfg DuctConfig) (*Duct, error) {
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return &Duct{cfg: cfg, area: Area(cfg.Diameter)}, nil
}

// Config 返回不可变的配置副本。
func (d *Duct) Config() DuctConfig { return d.cfg }

// Area 返回截面积（m²），在构造时预计算。
func (d *Duct) Area() float64 { return d.area }

// UsesAutomaticFriction 报告沿程摩阻是否按雷诺数自动取值。
func (d *Duct) UsesAutomaticFriction() bool {
	return d.cfg.Friction <= 0
}

// FrictionMode 返回摩阻取值方式："fixed"、"swamee-jain" 或 "blasius"。
func (d *Duct) FrictionMode() string {
	if d.cfg.Friction > 0 {
		return "fixed"
	}
	if d.cfg.Roughness > 0 {
		return "swamee-jain"
	}
	return "blasius"
}
