package proof

import (
	"math"
	"math/big"

	comp "github.com/rsa_accumulator/complex"
)

const (
	log10Base2    = 3.32192809489
	preComputeLmt = 20
)

var (
	big0 = big.NewInt(0)
	big1 = big.NewInt(1)
	big2 = big.NewInt(2)
	big3 = big.NewInt(3)
	big4 = big.NewInt(4)
	big8 = big.NewInt(8)

	// precomputed Hurwitz GCRDs for small integers
	precomputedHurwitzGCRDs = [preComputeLmt + 1]*comp.HurwitzInt{
		// 0's precomputed Hurwitz GCRD: 0, 0, 0, 0
		comp.NewHurwitzInt(big0, big0, big0, big0, false),
		// 1's precomputed Hurwitz GCRD: 1, 0, 0, 0
		comp.NewHurwitzInt(big1, big0, big0, big0, false),
		// 2's precomputed Hurwitz GCRD: 1, 1, 0, 0
		comp.NewHurwitzInt(big1, big1, big0, big0, false),
		// 3's precomputed Hurwitz GCRD: 1, 1, 1, 0
		comp.NewHurwitzInt(big1, big1, big1, big0, false),
		// 4's precomputed Hurwitz GCRD: 2, 0, 0, 0
		comp.NewHurwitzInt(big2, big0, big0, big0, false),
		// 5's precomputed Hurwitz GCRD: 2, 1, 0, 0
		comp.NewHurwitzInt(big2, big1, big0, big0, false),
		// 6's precomputed Hurwitz GCRD: 2, 1, 1, 0
		comp.NewHurwitzInt(big2, big1, big1, big0, false),
		// 7's precomputed Hurwitz GCRD: 2, 1, 1, 1
		comp.NewHurwitzInt(big2, big1, big1, big1, false),
		// 8's precomputed Hurwitz GCRD: 2, 2, 0, 0
		comp.NewHurwitzInt(big2, big2, big0, big0, false),
		// 9's precomputed Hurwitz GCRD: 2, 2, 1, 0
		comp.NewHurwitzInt(big2, big2, big1, big0, false),
		// 10's precomputed Hurwitz GCRD: 2, 2, 1, 1
		comp.NewHurwitzInt(big2, big2, big1, big1, false),
		// 11's precomputed Hurwitz GCRD: 3, 1, 1, 0
		comp.NewHurwitzInt(big3, big1, big1, big0, false),
		// 12's precomputed Hurwitz GCRD: 3, 1, 1, 1
		comp.NewHurwitzInt(big3, big1, big1, big1, false),
		// 13's precomputed Hurwitz GCRD: 3, 2, 0, 0
		comp.NewHurwitzInt(big3, big2, big0, big0, false),
		// 14's precomputed Hurwitz GCRD: 3, 2, 1, 0
		comp.NewHurwitzInt(big3, big2, big1, big0, false),
		// 15's precomputed Hurwitz GCRD: 3, 2, 1, 1
		comp.NewHurwitzInt(big3, big2, big1, big1, false),
		// 16's precomputed Hurwitz GCRD: 4, 0, 0, 0
		comp.NewHurwitzInt(big4, big0, big0, big0, false),
		// 17's precomputed Hurwitz GCRD: 4, 1, 0, 0
		comp.NewHurwitzInt(big4, big1, big0, big0, false),
		// 18's precomputed Hurwitz GCRD: 4, 1, 1, 0
		comp.NewHurwitzInt(big4, big1, big1, big0, false),
		// 19's precomputed Hurwitz GCRD: 4, 1, 1, 1
		comp.NewHurwitzInt(big4, big1, big1, big1, false),
		// 20's precomputed Hurwitz GCRD: 4, 2, 0, 0
		comp.NewHurwitzInt(big4, big2, big0, big0, false),
	}
	bigPreComputeLmt = big.NewInt(preComputeLmt)
)

func log2(n *big.Int) int {
	return n.BitLen() - 1
}

func log10(n *big.Int) int {
	log2 := log2(n)
	lnF := math.Round(float64(log2) / log10Base2)
	return int(lnF)
}
