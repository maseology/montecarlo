// johnsonb.go returns a uni-modal Johnson-bounded
// probability distribution from u[0,1] with mode m.
// see pg.297 & 457 of: Law, A.M., 2007. Simulation Modeling and Analysis. McGraw-Hill, fourth ed. New York. 768pp.

package invdistr

import (
	"log"
	"math"
)

// JohnsonB (bounded) sampling distribution
type JohnsonB struct {
	m float64
}

// NewJohnsonB constructor
func NewJohnsonB(m float64) *JohnsonB {
	if m < 0.0 || m > 1.0 {
		log.Panicf("Invalid JohnsonB parameter m: %v\n", m)
	} else if m == 0. {
		m = 0.01
	} else if m == 1. {
		m = 0.99
	}
	j := new(JohnsonB)
	j.m = m
	return j
}

// Inv : inverse function
// setting parameter alpha2 to 2.0 (increase number for smaller variance about the mode)
func (t *JohnsonB) Inv(f float64) float64 {
	var a1, y float64
	a2 := 0.69 // a2 < 0.7 results in bimodal distributions
loop:
	a2 += 0.01
	a1 = (2.*t.m-1.)/a2 - a2*math.Log10(t.m/(1.-t.m))
	y = a2 / t.m / (1. - t.m) / math.Sqrt(2.*math.Pi) * math.Exp(-0.5*math.Pow(a1+a2*math.Log10(t.m/(1.-t.m)), 2.0))
	if y <= 4. { // search for alpha2 such that f(x)=4
		goto loop
	}

	z := math.Sqrt(2.) * math.Erfinv(2.*f-1.)
	y = math.Exp((z - a1) / a2)
	return y / (y + 1.)
}
