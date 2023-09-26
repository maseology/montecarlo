// permute.go - used to create a complete sample set of
// every possible permutation of p dimensions and w discrete
// values.

package smpln

import (
	"fmt"

	"github.com/maseology/mmaths"
)

// Permutations returns a sampling plan [0,1]
// p: dimension; w: number of discrete samples
func Permutations(p, w int) [][]float64 {
	ip, n, d := intPermute(p, w), mmaths.IntPow(w, p), float64(w-1)
	fmt.Println(ip)
	u := make([][]float64, n)
	for i := 0; i < n; i++ {
		u[i] = make([]float64, p)
		for j := 0; j < p; j++ {
			u[i][j] = float64(ip[i][j]) / d
		}
	}
	return u
}

// w (width): number of d-nary digits; d (depth)
func intPermute(w, d int) [][]int {
	s := make([][]int, 0)
	var recurs func([]int, int)
	recurs = func(c []int, i int) {
		if i == 0 {
			s = append(s, c)
		} else {
			for j := 0; j < d; j++ {
				a := append(c, j)
				recurs(a, i-1)
			}
		}
	}

	c0 := make([]int, 0)
	recurs(c0, w)
	return s
}
