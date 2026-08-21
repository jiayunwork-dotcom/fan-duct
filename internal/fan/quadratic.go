package fan

// FitQuadraticCoeffs 用最小二乘法对样本点拟合二次多项式 y = a·q² + b·q + c，
// 返回系数 (a, b, c)。正规方程为 3×3 线性系统，用带部分主元的高斯消元求解。
func FitQuadraticCoeffs(pts []Point) (a, b, c float64) {
	var s0, s1, s2, s3, s4, sy, s1y, s2y float64
	for _, p := range pts {
		q := p.Flow
		y := p.Pressure
		q2 := q * q
		q3 := q2 * q
		s0++
		s1 += q
		s2 += q2
		s3 += q3
		s4 += q2 * q2
		sy += y
		s1y += q * y
		s2y += q2 * y
	}
	// 矩阵 [s4 s3 s2 | s2y; s3 s2 s1 | s1y; s2 s1 s0 | sy]
	mat := [3][4]float64{
		{s4, s3, s2, s2y},
		{s3, s2, s1, s1y},
		{s2, s1, s0, sy},
	}
	sol := solve3(mat)
	return sol[0], sol[1], sol[2]
}

// solve3 用带部分主元的高斯消元解 3 元一次方程组，返回解向量。
func solve3(m [3][4]float64) [3]float64 {
	for col := 0; col < 3; col++ {
		pivot := col
		for r := col + 1; r < 3; r++ {
			if abs(m[r][col]) > abs(m[pivot][col]) {
				pivot = r
			}
		}
		if pivot != col {
			m[col], m[pivot] = m[pivot], m[col]
		}
		pv := m[col][col]
		for k := col; k <= 3; k++ {
			m[col][k] /= pv
		}
		for r := 0; r < 3; r++ {
			if r == col {
				continue
			}
			factor := m[r][col]
			for k := col; k <= 3; k++ {
				m[r][k] -= factor * m[col][k]
			}
		}
	}
	return [3]float64{m[0][3], m[1][3], m[2][3]}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

// FitResiduals 返回样本点与二次拟合值的压升差，供拟合质量校验。
func FitResiduals(pts []Point) []float64 {
	a, b, c := FitQuadraticCoeffs(pts)
	out := make([]float64, len(pts))
	for i, p := range pts {
		out[i] = p.Pressure - (a*p.Flow*p.Flow + b*p.Flow + c)
	}
	return out
}
