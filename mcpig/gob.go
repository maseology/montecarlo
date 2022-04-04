package mcpig

import (
	"encoding/gob"
	"fmt"
	"log"
	"math"
	"os"

	mmplt "github.com/maseology/mmPlot"
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

	// // save to PNG
	// writePNG(mmio.GetFileDir(gobfp)+"/", mmio.FileName(gobfp, false)+"_")
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

	// save to PNG
	writePNG(mcdir, prfx)
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

func collectBins(gobfp string) (sampler.Set, [][]float64) {
	ss, coll, err := readGob(gobfp)
	if err != nil {
		log.Fatalln(err)
	}

	// intialize bins
	bins, denom := make([][]float64, ss.Ndim), 0.
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

	// normalize
	denom /= float64(nbins)
	for i := 0; i < ss.Ndim; i++ {
		for j := 0; j < nbins; j++ {
			bins[i][j] /= denom
		}
	}
	return ss, bins
}

func writeToCsv(w *mmio.TXTwriter, gobfp string) {
	ss, bins := collectBins(gobfp)
	par, gnam := ss.ParameterNames(), mmio.FileName(gobfp, false)
	for j := 0; j < nbins; j++ {
		for i := 0; i < ss.Ndim; i++ {
			pnam := par[i]
			var val float64
			switch ss.Samplers[i].Dist {
			case sampler.Constant:
				continue
			case sampler.LogLinear:
				pnam += " (log)"
				val = math.Log10(ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins)))
			default:
				val = ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins))
			}
			score := bins[i][j]
			// w.WriteLine(fmt.Sprintf("%s,%d,%s,%f", gnam, j+1, pnam, score))
			w.WriteLine(fmt.Sprintf("%s,%f,%s,%f", gnam, val, pnam, score))
		}
	}
}

func writePNG(mcdir, prfx string) {
	var pars []string
	var scr [][]float64
	var xlab [][]string
	for _, fp := range mmio.FileListExt(mcdir, ".gob") {
		ss, bins := collectBins(fp)
		if len(pars) == 0 {
			pars = ss.ParameterNames()
			scr = make([][]float64, ss.Ndim)
			xlab = make([][]string, ss.Ndim)
			for i := 0; i < ss.Ndim; i++ {
				if ss.Samplers[i].Dist == sampler.Constant {
					continue
				}
				mfmt := "%.5f"
				scr[i] = make([]float64, nbins)
				xlab[i] = make([]string, nbins)
				for j := 0; j < nbins; j++ {
					var v float64
					// switch ss.Samplers[i].Dist {
					// case sampler.LogLinear:
					// 	v = math.Log10(ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins)))
					// default:
					// 	v = ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins))
					// }
					v = ss.Samplers[i].Sample((float64(j) + .5) / float64(nbins))
					xlab[i][j] = fmt.Sprintf(mfmt, v)
				}
			}
		}
		for i := 0; i < ss.Ndim; i++ {
			if ss.Samplers[i].Dist == sampler.Constant {
				continue
			}
			for j := 0; j < nbins; j++ {
				scr[i][j] += bins[i][j]
			}
		}
	}

	for i, n := range pars {
		if xlab[i] == nil {
			continue
		}
		mmplt.Bar(mcdir+prfx+n+".png", scr[i], xlab[i])
	}
}
