package sampler

import (
	"log"

	mm "github.com/maseology/mmaths"
)

// Distribution enum type
type Distribution int

// Distribution enums
const (
	Uniform Distribution = iota
	Linear
	LogLinear
)

// String needed to return a Distribution type as string
func (d Distribution) String() string {
	return [...]string{"uniform", "linear", "log-linear"}[d]
}

// Sampler is a general struct used
type Sampler struct {
	Dist       Distribution
	Rmin, Rmax float64
	Name       string
}

// New Sampler constructor
func New(name string, d Distribution, rangeMin, rangeMax float64) *Sampler {
	if rangeMin > rangeMax {
		log.Fatalf("Sampler.New error: invalid input range for %s: min > max\n", name)
	}
	switch d {
	case LogLinear:
		if rangeMin <= 0. || rangeMax <= 0. {
			log.Fatalf("Sampler.New error: invalid input range for %s (%s distribution) (min = %f; max = %f\n", name, d, rangeMin, rangeMax)
		}
	default:
	}
	return &Sampler{Dist: d, Rmin: rangeMin, Rmax: rangeMax, Name: name}
}

// Sample returns a value from the distribution based on a U[0,1] sample
func (s *Sampler) Sample(u float64) float64 {
	switch s.Dist {
	case Uniform:
		return u
	case Linear:
		return mm.LinearTransform(s.Rmin, s.Rmax, u)
	case LogLinear:
		return mm.LogLinearTransform(s.Rmin, s.Rmax, u)
	default:
		log.Fatalln("Sampler.Sample error: unknown distribution used")
		return -9999.
	}
}
