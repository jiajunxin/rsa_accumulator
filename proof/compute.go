package proof

import (
	"errors"
	"io"
	"math/big"

	comp "github.com/rsa_accumulator/complex"
)

const preComputeLmt = 20

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
	tinyPrimeProd    = big.NewInt(2310) // 2 * 3 * 5 * 7 * 11
)

func log2(n *big.Int) int {
	return n.BitLen() - 1
}

//var smallPrimes = []uint8{
//	3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53,
//}
var smallPrimes = []uint8{
	3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37,
}

// smallPrimesProduct is the product of the values in smallPrimes and allows us
// to reduce a candidate prime by this number and then determine whether it's
// co-prime to all the elements of smallPrimes without further big.Int
// operations.
//var smallPrimesProduct = new(big.Int).SetUint64(16294579238595022365)
var smallPrimesProduct = new(big.Int).SetUint64(3710369067405)

// probPrime is modified crypto/rand Prime function for our use cases
// it returns error for any error returned by rand.Read or if bits < 2.
func probPrime(rand io.Reader, bits int) (p *big.Int, err error) {
	if bits < 2 {
		err = errors.New("bit size must be at least 2")
		return
	}

	b := uint(bits % 8)
	if b == 0 {
		b = 8
	}

	bytes := make([]byte, (bits+7)/8)
	p = new(big.Int)

	bigMod := new(big.Int)

	_, err = io.ReadFull(rand, bytes)
	if err != nil {
		return nil, err
	}

	// Clear bits in the first byte to make sure the candidate has a size <= bits.
	bytes[0] &= uint8(int(1<<b) - 1)
	// Don't let the value be too small, i.e, set the most significant two bits.
	// Setting the top two bits, rather than just the top bit,
	// means that when two of these values are multiplied together,
	// the result isn't ever one bit short.
	if b >= 2 {
		bytes[0] |= 3 << (b - 2)
	} else {
		// Here b==1, because b cannot be zero.
		bytes[0] |= 1
		if len(bytes) > 1 {
			bytes[1] |= 0x80
		}
	}
	// Make the value odd since an even number this large certainly isn't prime.
	bytes[len(bytes)-1] |= 1

	p.SetBytes(bytes)

	bigMod.Mod(p, smallPrimesProduct)
	mod := bigMod.Uint64()

NextDelta:
	for delta := uint64(0); delta < 1<<20; delta += 2 {
		m := mod + delta
		for _, prime := range smallPrimes {
			if m%uint64(prime) == 0 && (bits > 6 || m != uint64(prime)) {
				continue NextDelta
			}
		}

		if delta > 0 {
			bigMod.SetUint64(delta)
			p.Add(p, bigMod)
		}
		break
	}
	return
}
