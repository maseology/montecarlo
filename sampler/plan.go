package sampler

// Set holds a set of Samplers (i.e., a sampling plan)
type Set struct {
	Samplers []*Sampler
	Ndim     int
}

// NewSet constructs a new sampling set
func NewSet(s []*Sampler) *Set {
	ss := make([]*Sampler, len(s))
	for i := range s {
		ss[i] = s[i]
	}
	return &Set{Samplers: ss, Ndim: len(ss)}
}

// Sample returns sample from U^n
func (s *Set) Sample(u []float64) []float64 {
	v := make([]float64, s.Ndim)
	for i, uu := range u {
		v[i] = s.Samplers[i].Sample(uu)
	}
	return v
}

// ParameterNames returns the names of Sampler parameters
func (s *Set) ParameterNames() []string {
	p := make([]string, s.Ndim)
	for i, m := range s.Samplers {
		p[i] = m.Name
	}
	return p
}
