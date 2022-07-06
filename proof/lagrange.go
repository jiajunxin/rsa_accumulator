package proof

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/rsa_accumulator/complex"
)

var (
	// 0's precomputed Hurwitz GCRD: 0, 0, 0, 0
	hGCRD0 = complex.NewHurwitzInt(big0, big0, big0, big0, false)
	// 1's precomputed Hurwitz GCRD: 1, 0, 0, 0
	hGCRD1 = complex.NewHurwitzInt(big1, big0, big0, big0, false)
	// 2's precomputed Hurwitz GCRD: 1, 1, 0, 0
	hGCRD2 = complex.NewHurwitzInt(big1, big1, big0, big0, false)
	// 3's precomputed Hurwitz GCRD: 1, 1, 1, 0
	hGCRD3 = complex.NewHurwitzInt(big1, big1, big1, big0, false)
	// 4's precomputed Hurwitz GCRD: 2, 0, 0, 0
	hGCRD4 = complex.NewHurwitzInt(big2, big0, big0, big0, false)
	// 5's precomputed Hurwitz GCRD: 2, 1, 0, 0
	hGCRD5 = complex.NewHurwitzInt(big2, big1, big0, big0, false)
	// 6's precomputed Hurwitz GCRD: 2, 1, 1, 0
	hGCRD6 = complex.NewHurwitzInt(big2, big1, big1, big0, false)
	// 7's precomputed Hurwitz GCRD: 2, 1, 1, 1
	hGCRD7 = complex.NewHurwitzInt(big2, big1, big1, big1, false)
	// 8's precomputed Hurwitz GCRD: 2, 2, 0, 0
	hGCRD8 = complex.NewHurwitzInt(big2, big2, big0, big0, false)
	// precomputed Hurwitz GCRDs for small integers
	precomputedHurwitzGCRDs = []*complex.HurwitzInt{hGCRD0, hGCRD1, hGCRD2, hGCRD3, hGCRD4, hGCRD5, hGCRD6, hGCRD7, hGCRD8}
)

// FourSquare is the LagrangeFourSquareLipmaa representation of a positive integer
// w <- LagrangeFourSquareLipmaa(mu), mu = w = W1^2 + W2^2 + W3^2 + W4^2
type FourSquare struct {
	W1 *big.Int
	W2 *big.Int
	W3 *big.Int
	W4 *big.Int
}

// NewFourSquare creates a new FourSquare
func NewFourSquare(w1 *big.Int, w2 *big.Int, w3 *big.Int, w4 *big.Int) *FourSquare {
	w1.Abs(w1)
	w2.Abs(w2)
	w3.Abs(w3)
	w4.Abs(w4)
	// sort the four big integers in descending order
	if w1.Cmp(w2) == -1 {
		w1, w2 = w2, w1
	}
	if w1.Cmp(w3) == -1 {
		w1, w3 = w3, w1
	}
	if w1.Cmp(w4) == -1 {
		w1, w4 = w4, w1
	}
	if w2.Cmp(w3) == -1 {
		w2, w3 = w3, w2
	}
	if w2.Cmp(w4) == -1 {
		w2, w4 = w4, w2
	}
	if w3.Cmp(w4) == -1 {
		w3, w4 = w4, w3
	}
	return &FourSquare{w1, w2, w3, w4}
}

// Mul multiplies all the square numbers by n
func (f *FourSquare) Mul(n *big.Int) {
	f.W1.Mul(f.W1, n)
	f.W2.Mul(f.W2, n)
	f.W3.Mul(f.W3, n)
	f.W4.Mul(f.W4, n)
}

// Div divides all the square numbers by n
func (f *FourSquare) Div(n *big.Int) {
	f.W1.Div(f.W1, n)
	f.W2.Div(f.W2, n)
	f.W3.Div(f.W3, n)
	f.W4.Div(f.W4, n)
}

// String stringnifies the FourSquare object
func (f *FourSquare) String() string {
	return fmt.Sprintf("{%s %s %s %s}",
		f.W1.String(),
		f.W2.String(),
		f.W3.String(),
		f.W4.String(),
	)
}

// LagrangeFourSquares calculates the Lagrange four squares representation of a positive integer
// Paper: Finding the Four Squares in Lagrange’s Theorem
// Link: http://pollack.uga.edu/finding4squares.pdf (page 6)
// The input should be an odd positive integer no less than 9
func LagrangeFourSquares(n *big.Int) (*FourSquare, error) {
	n = new(big.Int).Set(n)
	// n = 2^e * n', n' is odd
	e := big.NewInt(0)
	for n.Bit(0) == 0 {
		n.Rsh(n, 1)
		e.Add(e, big1)
	}
	var hurwitzGCRD *complex.HurwitzInt

	if n.Cmp(big8) <= 0 {
		hurwitzGCRD = precomputedHurwitzGCRDs[n.Int64()]
	} else {
		primeProd, err := preCompute(n)
		if err != nil {
			return nil, err
		}
		for {
			s, p, err := randomTrails(n, primeProd)
			if err != nil {
				return nil, err
			}
			hurwitzGCRD, err = denouement(n, s, p)
			if err != nil {
				return nil, err
			}
			w1, w2, w3, w4 := hurwitzGCRD.ValInt()
			if Verify(n, w1, w2, w3, w4) {
				break
			}
		}
	}

	// if X'^2 + Y'^2 + Z'^2 + W'^2 = n'
	// then X^2 + Y^2 + Z^2 + W^2 = n for X, Y, Z, W defined by
	// (1 + i)^e * (X' + Y'i + Z'j + W'k) = (X + Yi + Zj + Wk)
	// Hurwitz integer: i + i
	hurwitz1PlusI := complex.NewHurwitzInt(big1, big1, big0, big0, false)
	hurwitzProd := complex.NewHurwitzInt(big1, big0, big0, big0, false)
	for e.Sign() > 0 {
		hurwitzProd.Prod(hurwitzProd, hurwitz1PlusI)
		e.Sub(e, big1)
	}
	hurwitzProd.Prod(hurwitzProd, hurwitzGCRD)
	w1, w2, w3, w4 := hurwitzProd.ValInt()
	fs := NewFourSquare(w1, w2, w3, w4)
	return fs, nil
}

// preCompute determine the primes not exceeding log n and compute their product
func preCompute(n *big.Int) (*big.Int, error) {
	logN := log2(n)
	var (
		primes    []*big.Int
		primeProd = big.NewInt(1) // 1
		idx       = big.NewInt(2) // 2
	)
	for idx.Cmp(logN) < 1 {
		isPrime := true
		for _, prime := range primes {
			mod := new(big.Int).Mod(idx, prime)
			if mod.Sign() == 0 {
				isPrime = false
				break
			}
		}
		if isPrime {
			newPrime := new(big.Int).Set(idx)
			primes = append(primes, newPrime)
			primeProd.Mul(primeProd, newPrime)
		}
		idx.Add(idx, big1)
	}
	return primeProd, nil
}

func randomTrails(n, primeProd *big.Int) (*big.Int, *big.Int, error) {
	nPow5Div2 := new(big.Int).Exp(n, big5, nil)
	nPow5Div2.Div(nPow5Div2, big2)
	preP := new(big.Int).Set(primeProd)
	preP.Mul(preP, n)
	for {
		var (
			k   = big.NewInt(0)
			u   *big.Int
			err error
		)
		// choose an odd number k < n^5 at random
		// let k = 2k' + 1, then 2k' + 1 < n^5
		// 2k' <= n^5 - 2
		// k' < ceiling{n^5 / 2}
		// start finding k' in [0, n^5 / 2)
		for new(big.Int).Mod(k, big2).Sign() == 0 {
			if k, err = rand.Int(rand.Reader, nPow5Div2); err != nil {
				return nil, nil, err
			}
		}
		// construct k
		k.Mul(k, big2)
		k.Add(k, big1)
		// p = {Product of primes} * n * k - 1
		p := new(big.Int).Set(preP)
		p.Mul(p, k)
		p.Sub(p, big1)
		pMinus1 := new(big.Int).Set(p)
		pMinus1.Sub(pMinus1, big1)
		// choose u from [1, p - 1]
		if u, err = rand.Int(rand.Reader, pMinus1); err != nil {
			return nil, nil, err
		}
		u.Add(u, big1)
		// compute s = u^((p - 1) / 4) mod p
		powU := new(big.Int).Set(pMinus1)
		powU.Div(powU, big4)
		s := new(big.Int).Exp(u, powU, p)
		targetMod := new(big.Int).Mod(bigNeg1, p)
		// test if s^2 = -1 (mod p)
		// if so, continue to the next step, otherwise, repeat this step
		if new(big.Int).Exp(s, big2, p).Cmp(targetMod) == 0 {
			return s, p, nil
		}
	}
}

func denouement(n, s, p *big.Int) (*complex.HurwitzInt, error) {
	// compute A + Bi := gcd(s + i, p)
	// Gaussian integer: s + i
	gaussianInt := complex.NewGaussianInt(s, big1)
	// Gaussian integer: p
	gaussianP := complex.NewGaussianInt(p, big0)
	gcd := new(complex.GaussianInt)
	gcd.GCD(gaussianInt, gaussianP)
	// compute gcrd(A + Bi + j, n), normalized to have integer component
	// Hurwitz integer: A + Bi + j
	hurwitzInt := complex.NewHurwitzInt(gcd.R, gcd.I, big1, big0, false)
	// Hurwitz integer: n
	hurwitzN := complex.NewHurwitzInt(n, big0, big0, big0, false)
	gcrd := new(complex.HurwitzInt)
	gcrd.GCRD(hurwitzInt, hurwitzN)

	return gcrd, nil
}

// Verify checks if the four-square sum is equal to the original integer
// i.e. target = w1^2 + w2^2 + w3^2 + w4^2
func Verify(target, w1, w2, w3, w4 *big.Int) bool {
	sum := new(big.Int).Mul(w1, w1)
	sum.Add(sum, new(big.Int).Mul(w2, w2))
	sum.Add(sum, new(big.Int).Mul(w3, w3))
	sum.Add(sum, new(big.Int).Mul(w4, w4))
	return sum.Cmp(target) == 0
}
