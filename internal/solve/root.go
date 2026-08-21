package solve

import "math"

// residualFunc 是求根目标函数：返回 f(Q)，并允许在越界等情形下报错。
type residualFunc func(q float64) (float64, error)

// findRoot 在 [lo, hi] 内二分求 f(Q) = 0。要求 f(lo) 与 f(hi) 异号
// （或某端恰为 0）。区间宽度收缩到相对 1e-12 或达到 maxIter 即停。
func findRoot(f residualFunc, lo, hi float64, maxIter int) (float64, int, error) {
	if lo >= hi {
		return 0, 0, noIntersection("empty bracket")
	}
	flo, err := f(lo)
	if err != nil {
		return 0, 0, err
	}
	fhi, err := f(hi)
	if err != nil {
		return 0, 0, err
	}
	if flo == 0 {
		return lo, 0, nil
	}
	if fhi == 0 {
		return hi, 0, nil
	}
	if (flo > 0) == (fhi > 0) {
		return 0, 0, noIntersection("no sign change on the bracket")
	}
	iter := 0
	for iter = 0; iter < maxIter; iter++ {
		mid := 0.5 * (lo + hi)
		fm, err := f(mid)
		if err != nil {
			return 0, 0, err
		}
		if fm == 0 || hi-lo < 1e-12*(1+math.Abs(mid)) {
			return mid, iter + 1, nil
		}
		if (flo > 0) != (fm > 0) {
			hi = mid
			fhi = fm
		} else {
			lo = mid
			flo = fm
		}
	}
	return 0.5 * (lo + hi), iter, nil
}

// expandUpperBound 从样本上限开始翻倍扩张上界，直到 f(q) <= 0。
// 翻倍过程中若外推被禁止而越过样本范围，返回该错误。
func (b *Build) expandUpperBound(start float64, f residualFunc, maxDoublings int) (float64, float64, error) {
	hi := start
	if hi <= 0 {
		hi = 1.0
	}
	fHi, err := f(hi)
	if err != nil {
		return 0, 0, err
	}
	for i := 0; i < maxDoublings && fHi > 0; i++ {
		hi *= 2
		fHi, err = f(hi)
		if err != nil {
			return 0, 0, err
		}
	}
	return hi, fHi, nil
}

// abs 返回绝对值。
func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
