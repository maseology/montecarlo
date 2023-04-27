package montecarlo

import (
	"fmt"
	"math/rand"
	"sort"
	"sync"
	"time"

	mrg63k3a "github.com/maseology/goRNG/MRG63k3a"
	"github.com/maseology/mmaths"
	"github.com/maseology/mmaths/slice"
	"github.com/maseology/montecarlo/smpln"
)

// GenerateSamples returns the result from n evaluations of fun() sampling from p-hypercube
func GenerateSamples(fun func(u []float64, i int) float64, n, p, nthrd int) ([][]float64, []float64) { // ([][]float64, []float64, []int) {
	fmt.Printf("generating %d samples of %d parameters, %d at a time..\n", n, p, nthrd)
	var wg sync.WaitGroup
	smpls := make(chan []float64, nthrd)
	results := make(chan []float64, n)

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())
	sp := smpln.NewLHC(rng, n, p, false) // smpln.NewHalton(s, n)

	for k := 0; k < n; k++ {
		wg.Add(1)
		go func(k int) {
			s := <-smpls
			results <- append(s, fun(s, k))
			wg.Done()
		}(k)
	}

	for k := 0; k < n; k++ {
		ut := make([]float64, p)
		for j := 0; j < p; j++ {
			ut[j] = sp.U[j][k]
		}
		smpls <- ut
	}
	wg.Wait()
	close(smpls)

	f := make([]float64, n) // function value
	// d := make([]int, n)       // function index, used for ranking
	u := make([][]float64, n) // sample points
	for k := 0; k < n; k++ {
		u[k] = make([]float64, p)
		r := <-results
		for j := 0; j < p; j++ {
			u[k][j] = r[j]
		}
		f[k] = r[p]
		// d[k] = k
	}
	close(results)
	return u, f //, d
}

// RankedUnBiased returns s n-dimensional samples of fun()
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
