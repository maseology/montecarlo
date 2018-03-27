// knuthshuffle.go - an implementation of the permutation algorithm
// known as the Knuth shuffle; aka the Fisher and Yates' shuffle.
// produces a randomized set of [(i-1)/n,i/n), for integers i=1..n
// for orthogonality processing, see: Tang., B. (1993) Orthogonal
// Array-Based Latin Hypercubes. Journal of the American Statistical
// Association 88(424), pp. 1392-1397.

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

package smpln

import "math/rand"

// KnuthShuffle is the structure to hold the state of one
// instance of the LHC.  New instances can be allocated using
// the knuthShuffle.New() function.
type KnuthShuffle struct {
	Z    [][]int
	n, p int
}

// NewKS allocates a new instance of the Knuth Shuffle from n samples of p integers.
func NewKS(n, p int) *KnuthShuffle {
	ks := &KnuthShuffle{
		Z: make([][]int, p),
		n: n,
		p: p,
	}
	for i := range ks.Z {
		ks.Z[i] = make([]int, n)
	}
	return ks
}

// Make builds the sampling plan nxp matrix.
func (ks *KnuthShuffle) Make(rng *rand.Rand) {
	for i := 0; i < ks.n; i++ {
		ks.Z[0][i] = i
	}
	for j := 0; j < ks.p; j++ {
		if j > 0 {
			// copy over previous permutation
			for i := 0; i < ks.n; i++ {
				ks.Z[j][i] = ks.Z[j-1][i]
			}
		}
		for i := ks.n - 1; i > 0; i-- {
			swap := rng.Intn(i)
			zsv := ks.Z[j][swap]
			ks.Z[j][swap] = ks.Z[j][i]
			ks.Z[j][i] = zsv
		}
	}
}
