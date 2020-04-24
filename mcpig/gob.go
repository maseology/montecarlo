package mcpig

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"

	"github.com/maseology/mmio"
	"github.com/maseology/montecarlo/sampler"
)

const nbins = 30

// ReadMCPIG reads a single sampling
func ReadMCPIG(gobfp string) {
	wcsv, _ := mmio.NewTXTwriter(mmio.RemoveExtension(gobfp) + ".csv")
	defer wcsv.Close()
	wcsv.WriteLine("station,bin,par,val")
	writeToCsv(wcsv, gobfp)
}

// ReadMCPIGs reads a set of samplings (i.e., from multiple stations)
func ReadMCPIGs(mcdir, prfx string) {
	wcsv, err := mmio.NewTXTwriter(mcdir + prfx + "summary.csv")
	if err != nil {
		log.Fatalln(err)
	}
	defer wcsv.Close()
	wcsv.WriteLine("station,bin,par,val")
	fmt.Println("\nBuilding MCPIG summary..")
	for _, fp := range mmio.FileListExt(mcdir, ".gob") {
		fmt.Println(" " + fp)
		writeToCsv(wcsv, fp)
	}
	fmt.Println("Results saved to " + mcdir + prfx + "summary.csv")
}

func readGob(fp string) (sampler.Set, [][]float64, error) {
	var d [][]float64
	var s sampler.Set
	f, err := os.Open(fp)
	defer f.Close()
	if err != nil {
		return sampler.Set{}, nil, err
	}
	enc := gob.NewDecoder(f)
	err = enc.Decode(&s)
	if err != nil {
		return sampler.Set{}, nil, err
	}
	err = enc.Decode(&d)
	if err != nil {
		return sampler.Set{}, nil, err
	}
	return s, d, nil
}

func writeToCsv(w *mmio.TXTwriter, gobfp string) {
	denom := 0.
	ss, coll, err := readGob(gobfp)
	if err != nil {
		log.Fatalln(err)
	}

	// intialize bins
	bins := make([][]float64, ss.Ndim)
	for i := 0; i < ss.Ndim; i++ {
		bins[i] = make([]float64, nbins)
		for j := 0; j < nbins; j++ {
			bins[i][j] = 0.
		}
	}

	for _, v := range coll {
		objfnc := v[0]
		for i := 0; i < ss.Ndim; i++ {
			ii := int(math.Floor(v[i+1] * float64(nbins)))
			bins[i][ii] += objfnc
		}
		denom++
	}

	denom /= float64(nbins)
	par, gnam := ss.ParameterNames(), mmio.FileName(gobfp, false)
	for j := 0; j < nbins; j++ {
		for i := 0; i < ss.Ndim; i++ {
			if ss.Samplers[i].Dist == sampler.Constant {
				continue
			}
			score := bins[i][j] / denom
			val := ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins))
			// w.WriteLine(fmt.Sprintf("%s,%d,%s,%f", gnam, j+1, par[i], score))
			w.WriteLine(fmt.Sprintf("%s,%f,%s,%f", gnam, val, par[i], score))
		}
	}
}
