// halton.go this fillows the 'generalized' Halton
// digital sequence that improves 'space-filling' performance in
// higher dimensions.  See pg. 153 and 164 in Lemieux, C. (2009)
// Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 373pp.

package smpln

import (
	"log"

	"github.com/maseology/mmaths"
)

// from Faure, H., C. Lemieux (2008) Generalized Halton Sequences
// in 2008: A Comparative Study
// the Faure-Lemieux sequence <http://www.math.uwaterloo.ca/~clemieux/FLFactors.html>:
var fls = [360]int{
	1, 1, 3, 3, 4, 9, 7, 5, 9, 18, 18, 8, 13, 31, 9, 19, 36, 33, 21, 44, 43, 61, 60, 56, 26, 71, 32, 77, 26, 95, 92, 47, 29, 61, 57, 69, 115, 63,
	92, 31, 104, 126, 50, 80, 55, 152, 114, 80, 83, 97, 95, 150, 148, 55, 80, 192, 71, 76, 82, 109, 105, 173, 58, 143, 56, 177, 203, 239, 196,
	143, 278, 227, 87, 274, 264, 84, 226, 163, 231, 177, 95, 116, 165, 131, 156, 105, 188, 142, 105, 125, 269, 292, 215, 182, 294, 152, 148, 144,
	382, 194, 346, 323, 220, 174, 133, 324, 215, 246, 159, 337, 254, 423, 484, 239, 440, 362, 464, 376, 398, 174, 149, 418, 306, 282, 434,
	196, 458, 313, 512, 450, 161, 315, 441, 549, 555, 431, 295, 557, 172, 343, 472, 604, 297, 524, 251, 514, 385, 531, 663, 674, 255, 519,
	324, 391, 394, 533, 253, 717, 651, 399, 596, 676, 425, 261, 404, 691, 604, 274, 627, 777, 269, 217, 599, 447, 581, 640, 666, 595, 669,
	686, 305, 460, 599, 335, 258, 649, 771, 619, 666, 669, 707, 737, 854, 925, 818, 424, 493, 463, 535, 782, 476, 451, 520, 886, 340, 793,
	390, 381, 274, 500, 581, 345, 363, 1024, 514, 773, 932, 556, 954, 793, 294, 863, 393, 827, 527, 1007, 622, 549, 613, 799, 408, 856, 601,
	1072, 938, 322, 1142, 873, 629, 1071, 1063, 1205, 596, 973, 984, 875, 918, 1133, 1223, 933, 1110, 1228, 1017, 701, 480, 678, 1172, 689,
	1138, 1022, 682, 613, 635, 984, 526, 1311, 459, 1348, 477, 716, 1075, 682, 1245, 401, 774, 1026, 499, 1314, 743, 693, 1282, 1003, 1181,
	1079, 765, 815, 1350, 1144, 1449, 718, 805, 1203, 1173, 737, 562, 579, 701, 1104, 1105, 1379, 827, 1256, 759, 540, 1284, 1188, 776,
	853, 1140, 445, 1265, 802, 932, 632, 1504, 856, 1229, 1619, 774, 1229, 1300, 1563, 1551, 1265, 905, 1333, 493, 913, 1397, 1250, 612,
	1251, 1765, 1303, 595, 981, 671, 1403, 820, 1404, 1661, 973, 1340, 1015, 1649, 855, 1834, 1621, 1704, 893, 1033, 721, 1737, 1507,
	1851, 1006, 994, 923, 872, 1860}

// HaltonDigitalSequence is the structure to hold the state of one
// instance of the Halton digital sequence.  New instances can be
// allocated using the HaltonDigitalSequence.New() function.
type HaltonDigitalSequence struct {
	U    [][]float64
	n, p int
}

// NewHalton allocates a new instance of the LHC
func NewHalton(n, p int) *HaltonDigitalSequence {
	hds := &HaltonDigitalSequence{
		U: make([][]float64, p),
		n: n,
		p: p,
	}
	for i := range hds.U {
		hds.U[i] = make([]float64, n)
	}

	// original Halton sequence; multiplicative factor added to the van der Corput forms the FL generalized Halton
	b := mmaths.Primes(p)
	for i := 0; i < n; i++ {
		for j := 0; j < p; j++ {
			// hds.U[j][i] = vanderCorput(i+1, b[j], 1) // original Halton sequence
			hds.U[j][i] = vanderCorput(i+1, b[j], fls[j]) // Faure-Lemieux generalized Halton sequence (using i+1 to avoid 0,0 returned for small dimensions)
			if hds.U[j][i] > 1.0 || hds.U[j][i] < 0.0 {
				log.Panicf("Halton digital sequence error: value out of range U[0,1): %v", hds.U[j][i])
			}
		}
	}
	return hds
}

// SampleSize simply returns the number of samples
func (hds *HaltonDigitalSequence) SampleSize() int { return hds.n }

func vanderCorputSeq(k, b, m int) []float64 {
	// van der Corput sequence (i.e., radical inverse function of a given base)
	// see pg. 145 in Lemieux, C., Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 2009. 373pp.
	if b <= 1 || m < 1 {
		panic("vanderCorput_seq input error")
	}
	p := make([]float64, 0, k)
	for i := 0; i < k; i++ {
		p = append(p, vanderCorput(i, b, m))
	}
	return p
}

func vanderCorput(i, b, m int) float64 {
	// van der Corput sequence (i.e., radical inverse function of base b)
	// see pg. 145 in Lemieux, C., Monte Carlo and Quasi-Monte Carlo Sampling. Springer Science. 2009. 373pp.
	if b <= 1 || m < 1 {
		panic("vanderCorput input error")
	}
	inv := 0.
	ins := 1
loop:
	ins *= b
	ina := i % b // a_l pg. 153
	if m > 1 {
		ina = (m * ina) % b // pi_j,r pg. 165
	}
	inv += float64(ina) / float64(ins)
	i /= b
	if i > 0 {
		goto loop
	}
	return inv
}
