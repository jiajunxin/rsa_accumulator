package proof

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"github.com/rsa_accumulator/complex"
)

var (
	uRetryExponent = big.NewInt(3)
	// 0's precomputed Hurwitz GCRD: 0, 0, 0, 0
	hGCRD0 = complex.NewHurwitzInt(big.NewInt(0), big.NewInt(0), big.NewInt(0), big.NewInt(0), false)
	// 1's precomputed Hurwitz GCRD: 1, 0, 0, 0
	hGCRD1 = complex.NewHurwitzInt(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), false)
	// 2's precomputed Hurwitz GCRD: 1, 1, 0, 0
	hGCRD2 = complex.NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(0), false)
	// 3's precomputed Hurwitz GCRD: 1, 1, 1, 0
	hGCRD3 = complex.NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(0), false)
	// 4's precomputed Hurwitz GCRD: 2, 0, 0, 0
	hGCRD4 = complex.NewHurwitzInt(big.NewInt(2), big.NewInt(0), big.NewInt(0), big.NewInt(0), false)
	// 5's precomputed Hurwitz GCRD: 2, 1, 0, 0
	hGCRD5 = complex.NewHurwitzInt(big.NewInt(2), big.NewInt(1), big.NewInt(0), big.NewInt(0), false)
	// 6's precomputed Hurwitz GCRD: 2, 1, 1, 0
	hGCRD6 = complex.NewHurwitzInt(big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(0), false)
	// 7's precomputed Hurwitz GCRD: 2, 1, 1, 1
	hGCRD7 = complex.NewHurwitzInt(big.NewInt(2), big.NewInt(1), big.NewInt(1), big.NewInt(1), false)
	// 8's precomputed Hurwitz GCRD: 2, 2, 0, 0
	hGCRD8 = complex.NewHurwitzInt(big.NewInt(2), big.NewInt(2), big.NewInt(0), big.NewInt(0), false)
	// precomputed Hurwitz GCRDs for small integers
	precomputedHurwitzGCRDs = []*complex.HurwitzInt{hGCRD0, hGCRD1, hGCRD2, hGCRD3, hGCRD4, hGCRD5, hGCRD6, hGCRD7, hGCRD8}
	// LagrangeFourSquares is the function that computes the Lagrange four-squares
	LagrangeFourSquares = LagrangeFourSquaresPollack
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

// LagrangeFourSquareLipmaa calculates the Lagrange four square representation of a positive integer
// Paper: On Diophantine Complexity and Statistical Zero-Knowledge Arguments
// Link: https://eprint.iacr.org/2003/105
func LagrangeFourSquareLipmaa(mu *big.Int) (*FourSquare, error) {
	// write mu in the form mu = 2^t(2k + 1)
	var t int
	// copy mu for modification
	muCopy := new(big.Int).Set(mu)
	for muCopy.Bit(0) == 0 {
		t++
		// right shift
		muCopy.Rsh(muCopy, 1)
	}
	//muCopy.Lsh(muCopy, 1)
	fmt.Println(muCopy.Int64())

	// if t = 1
	if t == 1 {
		w1, w2, w3, w4, err := calW1W2W3W4(mu)
		if err != nil {
			return nil, err
		}
		fs := NewFourSquare(w1, w2, w3, w4)
		return fs, nil
	}

	// if t is odd but not 1
	if t%2 == 1 {
		muCopy.Mul(muCopy, big2)
		w1, w2, w3, w4, err := calW1W2W3W4(muCopy)
		if err != nil {
			return nil, err
		}
		s := new(big.Int).SetInt64(2)
		s.Exp(s, new(big.Int).SetInt64(int64((t-1)/2)), nil)
		fs := NewFourSquare(w1, w2, w3, w4)
		fs.Mul(s)
		return fs, nil
	}

	// if t is even
	fmt.Println(muCopy.Int64())
	muCopy.Sub(muCopy, big1)
	muCopy.Div(muCopy, big2)
	k := int(muCopy.Int64())
	muCopy.SetInt64(int64(2 * (2*k + 1)))
	fmt.Printf("mu: %d, t: %d, k: %d\n", mu.Int64(), t, k)
	w1, w2, w3, w4, err := calW1W2W3W4(muCopy)
	if err != nil {
		return nil, err
	}
	w1Mod2 := new(big.Int).Mod(w1, big2)
	if w1Mod2.Cmp(new(big.Int).Mod(w2, big2)) != 0 {
		if w1Mod2.Cmp(new(big.Int).Mod(w3, big2)) == 0 {
			w2, w3 = w3, w2
		} else {
			w2, w4 = w4, w2
		}
	}
	exp := int64(t/2 - 1)
	var isExpNegative bool
	if exp < 0 {
		exp = -exp
		isExpNegative = true
	}
	s := new(big.Int).SetInt64(2)
	s.Exp(s, big.NewInt(exp), nil)
	fmt.Println(s.Int64())
	fs := NewFourSquare(
		new(big.Int).Add(w1, w2),
		new(big.Int).Sub(w1, w2),
		new(big.Int).Add(w3, w4),
		new(big.Int).Sub(w3, w4),
	)
	if isExpNegative {
		fs.Div(s)
	} else {
		fs.Mul(s)
	}
	return fs, nil
}

func calPW1W2(mu *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	// if mu is 0
	if mu.Cmp(big0) == 0 {
		p := new(big.Int).Set(big0)
		w1 := new(big.Int).Set(big0)
		w2 := new(big.Int).Set(big0)
		return p, w1, w2, nil
	}
	// if mu is 1
	if mu.Cmp(big1) == 0 {
		p := new(big.Int).Set(big0)
		w1 := new(big.Int).Set(big1)
		w2 := new(big.Int).Set(big0)
		return p, w1, w2, nil
	}

	// choose random W1, W2 such that exactly one of W1, W2 is even
	w1Lmt := new(big.Int)
	w1Lmt.Sqrt(mu)
	w1Lmt.Add(w1Lmt, big.NewInt(1))
	// randomly choose W1 within [0, sqrt(mu)]
	w1, err := rand.Int(rand.Reader, w1Lmt)
	if err != nil {
		return nil, nil, nil, err
	}
	w1Sq := new(big.Int).Mul(w1, w1)
	w2Lmt := new(big.Int).Set(w1Sq)
	w2Lmt.Sub(mu, w1Lmt)
	w2Lmt.Sqrt(w2Lmt)
	w2Lmt.Add(w2Lmt, big1)
	// randomly choose W2 within [0, sqrt(mu - W1^2)]
	w2, err := rand.Int(rand.Reader, w2Lmt)
	if err != nil {
		return nil, nil, nil, err
	}
	w2Sq := new(big.Int).Mul(w2, w2)
	// p <- mu - W1^2 - W2^2, now p = 1 (mod 4)
	p := new(big.Int).Sub(mu, w1Sq)
	p.Sub(p, w2Sq)

	return p, w1, w2, nil
}

func calW1W2W3W4(mu *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	var (
		p, w1, w2, w3, w4 *big.Int
		err               error
	)
	for {
		// (a) choose random W1, W2, and calculate p <- mu - W1^2 - W2^2
		p, w1, w2, err = calPW1W2(mu)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		// if p is negative
		if p.Cmp(big0) == -1 {
			continue
		}
		// if p is 0
		if p.Cmp(big0) == 0 {
			w3 = new(big.Int)
			w4 = new(big.Int)
			return w1, w2, w3, w4, nil
		}
		// if p is 1
		if p.Cmp(big1) == 0 {
			w3 = new(big.Int).Set(big1)
			w4 = new(big.Int).Set(big0)
			return w1, w2, w3, w4, nil
		}

		// (b) hoping p is prime, try to express p = W3^2 + W4^2
		// find a solution u to the equation u^2 = -1 (mod p)
		mul := new(big.Int).Set(p)
		mul.Sub(mul, big.NewInt(1))
		targetMod := new(big.Int).Mod(big.NewInt(-1), p)
		u := new(big.Int).Set(targetMod)
		currMod := new(big.Int).Exp(u, big2, p)
		doubleMU := big.NewInt(2)
		doubleMU.Mul(doubleMU, u)
		uLmt := new(big.Int).Exp(doubleMU, uRetryExponent, nil)
		var lmtFlg bool
		for currMod.Cmp(targetMod) != 0 {
			u.Add(u, big1)
			currMod.Exp(u, big2, p)
			if u.Cmp(uLmt) == 1 {
				lmtFlg = true
				break
			}
		}
		if lmtFlg {
			log.Println("retrying finding q")
			continue
		}

		// apply Euclidean algorithm to (u, p), and take the first two remainders that are less than sqrt(p)
		floatP := new(big.Float).SetInt(p)
		floatSqrtP := new(big.Float).Sqrt(floatP)
		sqrtP := new(big.Int).Sqrt(p)
		if floatSqrtP.IsInt() {
			sqrtP.Sub(sqrtP, big1)
		}

		dividend := new(big.Int).Set(u)
		divisor := new(big.Int).Set(p)
		_, w3, err = euclideanDivision(dividend, divisor)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for w3.Cmp(sqrtP) != -1 {
			dividend = divisor
			divisor = w3
			_, w3, err = euclideanDivision(dividend, divisor)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		}
		dividend = divisor
		_, w4, err = euclideanDivision(dividend, w3)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		// if p != W3^2 + W4^2, then p is not prime, so go back to the p calculation
		w3Sq := new(big.Int).Mul(w3, w3)
		w4Sq := new(big.Int).Mul(w4, w4)
		if p.Cmp(new(big.Int).Add(w3Sq, w4Sq)) == 0 {
			continue
		}
	}
}

// LagrangeFourSquaresPollack calculates the Lagrange four squares representation of a positive integer
// Paper: Finding the Four Squares in Lagrange’s Theorem
// Link: http://pollack.uga.edu/finding4squares.pdf (page 6)
// The input should be an odd positive integer no less than 9
func LagrangeFourSquaresPollack(n *big.Int) (*FourSquare, error) {
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
	hurwitz1PlusI := complex.NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(0), big.NewInt(0), false)
	hurwitzProd := complex.NewHurwitzInt(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), false)
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
		primeProd = big.NewInt(1)
		idx       = big.NewInt(2)
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
	nPow5 := new(big.Int).Exp(n, big5, nil)
	preP := new(big.Int).Set(primeProd)
	preP.Mul(preP, n)
	for {
		var (
			k   = big.NewInt(0)
			u   = big.NewInt(0)
			err error
		)
		// choose an odd number k < n^5 at random, keep generating k until k is odd
		for new(big.Int).Mod(k, big2).Sign() == 0 {
			if k, err = rand.Int(rand.Reader, nPow5); err != nil {
				return nil, nil, err
			}
		}
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
func Verify(target, w1, w2, w3, w4 *big.Int) bool {
	sum := new(big.Int).Mul(w1, w1)
	sum.Add(sum, new(big.Int).Mul(w2, w2))
	sum.Add(sum, new(big.Int).Mul(w3, w3))
	sum.Add(sum, new(big.Int).Mul(w4, w4))
	return sum.Cmp(target) == 0
}
