package montecarlo

import (
	"encoding/gob"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/maseology/montecarlo/sampler"
)

const maxtrials = 10

// GenerateTop returns nsamples of function evaluations that exceed the minOF
func GenerateTop(fp string, eval func(u []float64, i int) float64, s sampler.Set, nsamples int, minOF float64) {
	tim := time.Now()
	cnt, iter := 0, 0
	coll := make([][]float64, 0, nsamples*maxtrials)
	for {
		uFinal, rFinal := GenerateSamples(eval, s.Ndim, nsamples, runtime.GOMAXPROCS(0))
		for i, f := range rFinal {
			if f > minOF {
				lst := make([]float64, s.Ndim+1)
				lst[0] = f
				for ii, v := range uFinal[i] {
					lst[ii+1] = v
				}
				coll = append(coll, lst)
				cnt++
			}
		}
		if cnt >= nsamples {
			fmt.Printf("\n  %d samples in %d iterations -- %v\n", cnt, iter+1, time.Now().Sub(tim))
			saveGob(fp, s, coll)
			break
		}
		iter++
		if iter == maxtrials && cnt == 0 {
			fmt.Printf("\n  no sample found to meet minimum objective -- %v\n", time.Now().Sub(tim))
			break
		}
	}
}

// saveGob can be read using mcpig
func saveGob(fp string, s sampler.Set, coll [][]float64) {
	f, err := os.Create(fp)
	defer f.Close()
	if err != nil {
		fmt.Println(err)
	}
	enc := gob.NewEncoder(f)
	err = enc.Encode(s)
	if err != nil {
		fmt.Println(err)
	}
	err = enc.Encode(coll)
	if err != nil {
		fmt.Println(err)
	}
}
