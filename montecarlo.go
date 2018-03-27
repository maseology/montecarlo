package montecarlo

import (
	"log"
	"math"

	id "github.com/maseology/montecarlo/invdistr"
)

// // SamplingPlan interface
// type SamplingPlan interface {
// 	SampleSize() int
// 	U() [][]float64
// }

// Parameter samples
type Parameter struct {
	id.Map
	P []float64
}

// NewParameter constructor
func NewParameter(n int, low, high float64, logscale bool, distro id.Mapper) *Parameter {
	p := new(Parameter)
	p.Low = low
	p.High = high
	p.Log = logscale
	p.Distr = distro
	p.P = make([]float64, n)
	return p
}

// BuildSampleSpace : from unit sample space, distribution mapping, build parameter space
func (p *Parameter) BuildSampleSpace(u []float64) {
	for i, f := range u {
		ut := p.Distr.Inv(f)
		p.P[i] = p.linearTransform(ut)
	}
}

func (p *Parameter) linearTransform(u float64) float64 {
	if u < 0.0 || u > 1.0 {
		log.Panicf("linear transform error, passing u = %v", u)
	}
	if p.Log {
		return math.Pow(10.0, math.Log10(p.Low*math.Pow(p.High/p.Low, u)))
	}
	return (p.High-p.Low)*u + p.Low
}
