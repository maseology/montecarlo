package montecarlo

import (
	"fmt"
	"sort"
	"sync"
	"math/rand"
	"time"

	"github.com/maseology/mmaths"
	"github.com/maseology/montecarlo/smpln"
	mrg63k3a "github.com/maseology/pnrg/MRG63k3a"
)

// GenerateSamples returns the result from s evalutaions of fun() sampling from n-hypercube
func GenerateSamples(fun func(u []float64) float64, n, s int) ([][]float64, []float64) { // ([][]float64, []float64, []int) {
	var wg sync.WaitGroup
	smpls := make(chan []float64, s)
	results := make(chan []float64, s)
	wg.Add(s)
	for k := 0; k < s; k++ {
		go func() {
			defer wg.Done()
			s := <-smpls
			results <- append(s, fun(s))
		}()
	}

	rng := rand.New(mrg63k3a.New())
	rng.Seed(time.Now().UnixNano())

	sp := smpln.NewLHC(rng, s, n, false) // smpln.NewHalton(s, n)
	for k := 0; k < s; k++ {
		ut := make([]float64, n)
		for j := 0; j < n; j++ {
			ut[j] = sp.U[j][k]
		}
		smpls <- ut
	}
	wg.Wait()
	close(smpls)

	f := make([]float64, s) // function value
	// d := make([]int, s)       // function index, used for ranking
	u := make([][]float64, s) // sample points
	for k := 0; k < s; k++ {
		u[k] = make([]float64, n)
		r := <-results
		for j := 0; j < n; j++ {
			u[k][j] = r[j]
		}
		f[k] = r[n]
		// d[k] = k
	}
	return u, f //, d
}

// RankSamples ranks samples accoring to evaluation value
func RankSamples(f []float64, minimize bool) []int {
	f2 := make([]float64, len(f))
	copy(f2, f)
	d := mmaths.Sequential(len(f) - 1) // resetting d
	sort.Sort(mmaths.IndexedSlice{Indx: d, Val: f2})
	if !minimize {
		mmaths.Rev(d) // ordering from best (highest evaluated score) to worst
	}
	return d
}

// RankedUnBiased returns s n-dimensional samples of fun()
func RankedUnBiased(fun func(u []float64) float64, n, s int) ([][]float64, []float64, []int) {
	fmt.Printf(" generating %d LHC samples from %d dimensions..\n", s, n)
	u, f := GenerateSamples(fun, n, s)
	d := RankSamples(f, true)
	return u, f, d
}
