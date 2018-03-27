// lhc.go - an implementation of the Latin Hyper-cube
// (LHC) sampling plan for use with multi-dimensional pseudo-random
// sampling, based on: Lemieux, C. (2009) Monte Carlo and Quasi-
// Monte Carlo Sampling. Springer Science. 2009. 373pp.

// Copyright (C) 2018 Mason Marchildon <mason@riffle.ca>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// LHC has a good 'space-filling design', has low-discrepancy (not
// as low as with digital nets and sequences), and is projection
// regular.

// The output array is a nxp matrix (where n=number of iterations
// (model runs) and p is the number of parameters) that produces
// a randomized parameter value between the uniform interval U[0,1).
// U[0,1) can then be applied to each parameter user-defined limits,
// i.e., [ll,hl), where ll=lower limit, and hl=upper limit: Apply
// output (u) for this function's matrix in the form: p=ll+(hl-ll)*u,
// where u=LHS(n,p)

package smpln

import (
	"log"
	"math/rand"
)

// LatinHyperCube is the structure to hold the state of one
// instance of the LHC.  New instances can be allocated using
// the latinHyperCube.New() function.
type LatinHyperCube struct {
	U    [][]float64
	n, p int
}

// NewLHC allocates a new instance of the LHC from n samples of p dimensions.
func NewLHC(n, p int) *LatinHyperCube {
	lhc := &LatinHyperCube{
		U: make([][]float64, p),
		n: n,
		p: p,
	}
	for i := range lhc.U {
		lhc.U[i] = make([]float64, n)
	}
	return lhc
}

// Make builds the sampling plan nxp matrix.
// Setting midpoint to False adds an additional random jitter
// to the position of the sample within sample space.
func (lhc *LatinHyperCube) Make(rng *rand.Rand, midpoint bool) {
	nf := float64(lhc.n)
	w := 1.0 / (2.0 * nf)
	ks := NewKS(lhc.n, lhc.p)
	ks.Make(rng)

	for j := 0; j < lhc.p; j++ {
		for i := 0; i < lhc.n; i++ {
			if !midpoint {
				w = rng.Float64() / nf
			}
			lhc.U[j][i] = float64(ks.Z[j][i])/nf + w
			if lhc.U[j][i] > 1.0 || lhc.U[j][i] < 0.0 {
				log.Panicf("LHC error: value out of range U[0,1): %v", lhc.U[j][i])
			}
		}
	}
}

// SampleSize simply returns the number of samples
func (lhc *LatinHyperCube) SampleSize() int { return lhc.n }
