// trapezoid.go returns a general Trapezoidal probability distribution
// from u[0,1] with modes m and n and shape factors a and b.
// special case: m=n, a=2, b=2: Triangular transform

package invdistr

import (
	"log"
	"math"
)

// Trapezoid sampling distribution
type Trapezoid struct {
	m, n float64 // modes [0,1]
	a, b float64 // shapeparameters 2 for straight
}

// NewTrapezoid constructor
func NewTrapezoid(m, n, a, b float64) *Trapezoid {
	if m < 0. || m > n || n > 1. || a < 0 || b < 0 {
		log.Panicf("General Trapezoid distribution: invalid arguments m, n, a, b = %v, %v, %v, %v\n", m, n, a, b)
	}
	t := new(Trapezoid)
	t.m = m
	t.n = n
	t.a = a
	t.b = b
	return t
}

// Inv : inverse function
func (t *Trapezoid) Inv(u float64) float64 {
	m, n, a, b := t.properties()
	if m < 0. || m > n || n > 1. || a > m || b < n {
		panic("Inverse General Trapezoid: invalid arguments")
	}
	pd := b*m + a*b*(n-m) + a*(1.-n)
	p1 := b * m / pd
	p2 := a * b * (n - m) / pd
	p3 := a * (1. - n) / pd
	if u < p1 {
		return m * math.Pow(u/p1, 1./a)
	} else if u <= 1.-p3 {
		return u*(n-m)/p2 + m*(1.-1./a)
	} else if u > 1.-p3 {
		return 1. - (1.-n)*math.Pow((1.-u)/p3, 1./b)
	}
	panic("error in Inverse General Trapezoid")
}

func (t *Trapezoid) properties() (float64, float64, float64, float64) {
	return t.m, t.n, t.a, t.b
}
