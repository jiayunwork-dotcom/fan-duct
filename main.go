package main

import (
	"flag"
	"fmt"
	"os"

	"fan-duct/internal/solve"
)

const usage = `fan-duct: fan-duct operating point calculator.

Given a circular duct (length, diameter, friction factor, local loss sum,
air density) and a fan curve (Q-dp sample points with polyline or quadratic
fit), the tool finds the operating point where the duct resistance pressure
drop equals the fan pressure rise, and optionally rescales the fan to a new
speed using affinity laws (Q proportional to N, dp proportional to N^2).

usage:
  fan-duct operate <input.json> [--new-speed rpm]
  fan-duct help

The input JSON has the shape:

  {
    "duct":   { "length": 50, "diameter": 0.15, "friction": 0.02,
                 "lossCoeff": 3, "density": 1.205, "viscosity": 1.82e-5 },
    "fan":    { "points": [ {"q": 0, "dp": 1450}, {"q": 0.3, "dp": 1280} ],
                 "fit": "polyline", "extrapolate": "error" },
    "speed":  { "rpm": 1450 },
    "newSpeed": { "rpm": 1595 }
  }

The --new-speed flag overrides "newSpeed" from the JSON. Illegal inputs
(diameter <= 0, length < 0, density <= 0, empty fan curve, flow outside the
sample range with extrapolation disabled, no intersection) are reported on
stderr and exit non-zero.
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	var err error
	switch os.Args[1] {
	case "operate":
		err = runOperate(os.Args[2:])
	case "help", "-h", "--help":
		fmt.Print(usage)
		return
	default:
		fmt.Fprintf(os.Stderr, "fan-duct: unknown command %q\n\n%s", os.Args[1], usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "fan-duct: %v\n", err)
		os.Exit(1)
	}
}

// runOperate 执行 operate 子命令：解析输入 → 求工作点 →（可选）转速缩放重求交 → 打印。
func runOperate(args []string) error {
	fs := flag.NewFlagSet("operate", flag.ContinueOnError)
	newSpeedFlag := fs.Float64("new-speed", 0, "target RPM; overrides newSpeed from the JSON")
	compact := fs.Bool("compact", false, "print a single-line summary")
	sensitivity := fs.Bool("sensitivity", false, "also print +10% duct-parameter sensitivity")
	fs.Usage = func() {
		fmt.Fprint(fs.Output(), "usage: fan-duct operate <input.json> [--new-speed rpm] [--compact] [--sensitivity]\n")
	}
	if err := fs.Parse(args); err != nil {
		return err
	}
	if fs.NArg() != 1 {
		return fmt.Errorf("operate needs exactly one input JSON file")
	}
	in, err := solve.ParseInputFile(fs.Arg(0))
	if err != nil {
		return err
	}
	if *newSpeedFlag > 0 {
		in.NewSpeed = &solve.SpeedSpec{RPM: *newSpeedFlag}
	}
	if in.NewSpeed != nil && in.Speed == nil {
		return fmt.Errorf("operate: newSpeed requires a base speed (speed.rpm) in the input")
	}
	b, err := in.Build()
	if err != nil {
		return err
	}
	base, err := b.OperatingPoint()
	if err != nil {
		return err
	}
	out := solve.Output{Input: in, Build: b, Base: base}
	if in.NewSpeed != nil {
		rr, err := in.RespeedToRPM(b, in.NewSpeed.RPM)
		if err != nil {
			return err
		}
		out.Respeeded = &rr
	}
	if *compact {
		fmt.Println(out.Compact())
		return nil
	}
	fmt.Print(out.String())
	if *sensitivity {
		sr, err := solve.ComputeSensitivities(b)
		if err != nil {
			return err
		}
		fmt.Print(sr.String())
	}
	return nil
}
