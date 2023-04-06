package jointdist

import "math"

// Nested transforms n variables such that 0.0 <= u1 <= ... <= un <= 1.0
func Nested(u ...float64) []float64 {
	f, d := 1., len(u)
	o := make([]float64, d)
	for i, v := range u {
		x := f * math.Pow(v, 1./float64(d-i))
		o[i] = x
		f *= x
	}
	return o
}

// Nested2 transforms variables such that 0.0 <= u1 <= u2 <= 1.0
// page 53 in Lemieux, C., Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 2009. 373pp.
func Nested2(u1, u2 float64) (float64, float64) {
	u2 = math.Sqrt(u2)
	u1 *= u2
	return u1, u2
}

// SumToOne creates n variable such that sum(u) = 1
func SumToOne(u ...float64) []float64 {
	s, o := 0., make([]float64, len(u))
	for _, v := range u {
		s += v
	}
	for i := range u {
		// if i == 0 {
		o[i] = u[i] / s
		// } else {
		// 	o[i] = o[i-1] + u[i]/s
		// }
	}
	return o
}
