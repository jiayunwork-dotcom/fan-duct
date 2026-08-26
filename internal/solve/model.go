package solve

import "fmt"

func BuildModelReport(b *Build, q float64) (string, error) {
	dr, err := b.Duct.Report(q)
	if err != nil {
		return "", err
	}
	fr := b.Fan.Report()
	return fmt.Sprintf("--- duct model ---\n%s--- fan model ---\n%s", dr.String(), fr.String()), nil
}

func ModelString(b *Build) (string, error) {
	lo, hi := b.Fan.FlowRange()
	return BuildModelReport(b, 0.5*(lo+hi))
}

func ModelBlock(b *Build) string {
	s, err := ModelString(b)
	if err != nil {
		return ""
	}
	return s + "\n"
}
