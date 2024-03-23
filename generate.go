package montecarlo

import (
	"fmt"
	"math/rand"
	"sort"
	"time"

	mrg63k3a "github.com/maseology/goRNG/MRG63k3a"
	"github.com/maseology/mmaths"
	"github.com/maseology/mmaths/slice"
	"github.com/maseology/montecarlo/smpln"
)

// GenerateSamples returns the result from n evaluations of fun() sampling from p-hypercube
func GenerateSamples(fun func(u []float64, i int) float64, n, p, nthrd int) ([][]float64, []float64) { // ([][]float64, []float64, []int) {
	if n < nthrd {
		nthrd = n
	}
	fmt.Printf("generating %d samples of %d parameters, %d at a time..\n", n, p, nthrd)

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())
	sp := smpln.NewLHC(rng, n, p, false)
	// smpln.NewHalton(s, n)

	type sample struct {
		u []float64
		k int
	}
	type result struct {
		u []float64
		f float64
		k int
	}

	smpl := make(chan sample)
	res := make(chan result)
	done := make(chan any)

	// spin-up nthrd workers
	for i := 0; i < nthrd; i++ {
		go func() {
			for {
				select {
				case <-done:
					return
				case s := <-smpl:
					res <- result{s.u, fun(s.u, s.k), s.k}
				}
			}
		}()
	}

	go func() {
		for k := 0; k < n; k++ {
			s := make([]float64, p)
			for j := 0; j < p; j++ {
				s[j] = sp.U[j][k]
			}
			smpl <- sample{s, k}
		}
	}()

	f := make([]float64, n)   // function value
	u := make([][]float64, n) // sample points
	for k := 0; k < n; k++ {
		r := <-res
		u[r.k] = r.u
		f[r.k] = r.f
	}
	close(done)
	close(smpl)
	close(res)

	return u, f
}

// RankedUnBiased returns s n-dimensional samples of fun(), ranking samples accoring to evaluation value
func RankedUnBiased(fun func(u []float64, i int) float64, n, s, nthrd int) ([][]float64, []float64, []int) {
	fmt.Printf(" generating %d LHC samples from %d dimensions..\n", s, n)
	u, f := GenerateSamples(fun, n, s, nthrd)
	d := RankSamples(f, true)
	return u, f, d
}

// RankSamples ranks samples accoring to evaluation value
func RankSamples(f []float64, minimize bool) []int {
	f2 := make([]float64, len(f))
	copy(f2, f)
	d := slice.Sequential(len(f) - 1) // resetting d
	sort.Sort(mmaths.IndexedSlice{Indx: d, Val: f2})
	if !minimize {
		slice.Rev(d) // ordering from best (highest evaluated score) to worst
	}
	return d
}
