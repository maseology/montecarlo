package smpln

import (
	"encoding/binary"
	"math"

	"github.com/maseology/mmaths"
	"github.com/maseology/mmio"
)

func Choose(nsmpl, nchoose int) [][]bool {
	n := mmaths.Binomial(nsmpl, nchoose)
	o := make([][]bool, n)
	cnt := func(b []bool) int {
		v := 0
		for _, bb := range b {
			if bb {
				v++
			}
		}
		return v
	}
	j := 0
	for i := 0; i < int(math.Pow(2, float64(nsmpl))); i++ {
		bs := make([]byte, 2)
		binary.LittleEndian.PutUint16(bs, uint16(i))
		ba := mmio.BitArray(bs)
		c := cnt(ba)
		if c != nchoose {
			continue
		}
		o[j] = ba[:nsmpl]
		// fmt.Println(j, i, c, o[j])
		j++
	}
	return o
}
