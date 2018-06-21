// this package contains 3 forms of symmetric and invertable copulae used to create joint sampling distributions
// see Kurowicka, D. and R. Cooke. Uncertainty Analysis with High Dimensional Dependence Modelling. John Wiley & Sons, Ltd. 2006. 284pp.

package jointdist

import (
	"math"
	"math/rand"
)

// Elliptical copula
// see pg 44 in Kurowicka, D. and R. Cooke. Uncertainty Analysis with High Dimensional Dependence Modelling. John Wiley & Sons, Ltd. 2006. 284pp.
func Elliptical(u1, u2, spearmanRho float64) (float64, float64) {
	v1 := u1 - 0.5
	v2 := math.Sqrt(1.-math.Pow(spearmanRho, 2.))*math.Sqrt(0.25-math.Pow(v1, 2.))*math.Sin(math.Pi*u2) + spearmanRho*v1 + 0.5
	return u1, v2
}

// DiagonalBand copula
// see pg 39 in Kurowicka, D. and R. Cooke. Uncertainty Analysis with High Dimensional Dependence Modelling. John Wiley & Sons, Ltd. 2006. 284pp.
func DiagonalBand(u1, u2, correlation float64) (float64, float64) {
	acor := math.Abs(correlation)
	neg := correlation < 0.0
	if acor > 1. {
		panic("DiagonalBand copula input error")
	}

	v1, v2 := u1, u2
	if neg {
		v1 = 1. - v1
	}
	if v1 < 1.-acor && v2 < 1.-v1/(1.-acor) {
		u2 = (1. - acor) * v2
	} else if v1 > acor && v2 > (1.-v1)/(1.-acor) {
		u2 = (1.-acor)*v2 + acor
	} else {
		u2 = 2.*(1.-acor)*v2 + v1 + acor - 1.
	}
	return u1, u2
}

// Franks copula
// a good choice for theta = 10.
// causion, u2 is not informed by the previous sampling plan, and thus low sampling dicsrepancy (i.e., LHC) will not be preserved
// see pg 47 in Kurowicka, D. and R. Cooke. Uncertainty Analysis with High Dimensional Dependence Modelling. John Wiley & Sons, Ltd. 2006. 284pp.
func Franks(u1, theta float64, rng *rand.Rand) (float64, float64) {
	if theta == 0. {
		panic("Franks copula input error")
	}

	u2 := -math.Log(1.-(1.-math.Exp(-theta))/((1./rng.Float64()-1.)*math.Exp(-theta*u1)+1.)) / theta
	return u1, u2
}
