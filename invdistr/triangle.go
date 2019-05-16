// triangle.go returns a Triangular probability distribution
// from u[0,1] with mode m.
// general trapezoid special case: m=n, a=2, b=2: Triangular transform

package invdistr

import "log"

// Triangle sampling distribution
type Triangle struct {
	m float64
}

// NewTriangle constructor (triangle is a special case of a trapezoid)
// Returns a Triangular probability distribution
// from u[0,1] with mode m. (special case: m=n, a=2, b=2: Triangular transform)
func NewTriangle(m float64) *Trapezoid {
	if m < 0. || m > 1. {
		log.Panicf("Inverse General Triangle: invalid arguments m = %v\n", m)
	}
	t := new(Trapezoid)
	t.m = m
	t.n = m
	t.a = 2.
	t.b = 2.
	return t
}

// Inv : inverse function
func (t *Triangle) Inv(u float64) float64 {
	trap := NewTrapezoid(t.m, t.m, 2., 2.)
	return trap.Inv(u)
}
