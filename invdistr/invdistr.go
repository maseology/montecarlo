package invdistr

import "math"

// Mapper : interface to MC distribution transforms
type Mapper interface {
	Inv(u float64) float64
}

// Map is a type used to contain sample mapping info
type Map struct {
	Low   float64
	High  float64
	Log   bool
	Distr Mapper
}

// P returns a paramater sample transformed from a given u[0,1] sample
func (m *Map) P(u float64) float64 {
	p := m.Distr.Inv(u)*(m.High-m.Low) + m.Low
	if m.Log {
		return math.Pow(10., p)
	}
	return p
}

// Uniform sampling distribution
type Uniform struct{}

// Inv : inverse function
func (t *Uniform) Inv(u float64) float64 { return u }
