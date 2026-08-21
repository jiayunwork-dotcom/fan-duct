package solve

import "fmt"

// BuildModelReport 返回模型描述文本：风管摘要与风机摘要分开渲染，
// 动压与管阻压降分开展示，避免与风机曲线的全压混淆。
func BuildModelReport(b *Build, q float64) (string, error) {
	dr, err := b.Duct.Report(q)
	if err != nil {
		return "", err
	}
	fr := b.Fan.Report()
	return fmt.Sprintf("--- duct model ---\n%s--- fan model ---\n%s", dr.String(), fr.String()), nil
}

// ModelString 便捷入口：默认用样本中点流量渲染模型描述。
func ModelString(b *Build) (string, error) {
	lo, hi := b.Fan.FlowRange()
	return BuildModelReport(b, 0.5*(lo+hi))
}

// ModelBlock 渲染模型描述并追加到输出。
func ModelBlock(b *Build) string {
	s, err := ModelString(b)
	if err != nil {
		return ""
	}
	return s + "\n"
}
