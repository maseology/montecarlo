package invdistr

// Mapper : interface to MC distribution transforms
type Mapper interface {
	Inv(f float64) float64
}

// Map is a type used to contain sample mapping info
type Map struct {
	Low   float64
	High  float64
	Log   bool
	Distr Mapper
}

// Uniform sampling distribution
type Uniform struct{}

// Inv : inverse function
func (t *Uniform) Inv(f float64) float64 { return f }
