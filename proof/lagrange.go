package proof

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
)

var (
	bigInt0        = big.NewInt(0)
	bigInt1        = big.NewInt(1)
	bigInt2        = big.NewInt(2)
	uRetryExponent = big.NewInt(3)
)

// FourSquare is the LagrangeFourSquares representation of a positive integer
// w <- LagrangeFourSquares(mu), mu = w = W1^2 + W2^2 + W3^2 + W4^2
type FourSquare struct {
	W1 *big.Int
	W2 *big.Int
	W3 *big.Int
	W4 *big.Int
}

// NewFourSquare creates a new FourSquare
func NewFourSquare(w1 *big.Int, w2 *big.Int, w3 *big.Int, w4 *big.Int) FourSquare {
	if w1.Sign() == -1 {
		w1.Neg(w1)
	}
	if w2.Sign() == -1 {
		w2.Neg(w2)
	}
	if w3.Sign() == -1 {
		w3.Neg(w3)
	}
	if w4.Sign() == -1 {
		w4.Neg(w4)
	}
	return FourSquare{w1, w2, w3, w4}
}

// Mul multiplies all the square numbers by n
func (f *FourSquare) Mul(n *big.Int) {
	f.W1.Mul(f.W1, n)
	f.W2.Mul(f.W2, n)
	f.W3.Mul(f.W3, n)
	f.W4.Mul(f.W4, n)
}

// Print prints all the square numbers
func (f *FourSquare) Print() {
	fmt.Printf("Lagrange Four Square: {%d %d %d %d}\n",
		f.W1.Int64(),
		f.W2.Int64(),
		f.W3.Int64(),
		f.W4.Int64(),
	)
}

// LagrangeFourSquares calculates the LagrangeFourSquares representation of a positive integer
// Paper: On Diophantine Complexity and Statistical Zero-Knowledge Arguments
// Link: https://eprint.iacr.org/2003/105
func LagrangeFourSquares(mu *big.Int) (FourSquare, error) {
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
			return FourSquare{}, err
		}
		return FourSquare{w1, w2, w3, w4}, nil
	}

	// if t is odd but not 1
	if t%2 == 1 {
		muCopy.Mul(muCopy, bigInt2)
		w1, w2, w3, w4, err := calW1W2W3W4(muCopy)
		if err != nil {
			return FourSquare{}, err
		}
		s := new(big.Int).SetInt64(2)
		s.Exp(s, new(big.Int).SetInt64(int64((t-1)/2)), nil)
		fs := NewFourSquare(w1, w2, w3, w4)
		fs.Mul(s)
		return fs, nil
	}

	// if t is even
	fmt.Println(muCopy.Int64())
	muCopy.Sub(muCopy, bigInt1)
	muCopy.Div(muCopy, bigInt2)
	k := int(muCopy.Int64())
	muCopy.SetInt64(int64(2 * (2*k + 1)))
	fmt.Printf("mu: %d, t: %d, k: %d\n", mu.Int64(), t, k)
	w1, w2, w3, w4, err := calW1W2W3W4(muCopy)
	if err != nil {
		return FourSquare{}, err
	}
	w1Mod2 := new(big.Int).Mod(w1, bigInt2)
	if w1Mod2.Cmp(new(big.Int).Mod(w2, bigInt2)) != 0 {
		if w1Mod2.Cmp(new(big.Int).Mod(w3, bigInt2)) == 0 {
			w2, w3 = w3, w2
		} else {
			w2, w4 = w4, w2
		}
	}
	s := new(big.Int).SetInt64(2)
	s.Exp(s, new(big.Int).SetInt64(int64(t/2-1)), nil)
	fs := NewFourSquare(
		new(big.Int).Add(w1, w2),
		new(big.Int).Sub(w1, w2),
		new(big.Int).Add(w3, w4),
		new(big.Int).Sub(w3, w4),
	)
	fs.Mul(s)
	return fs, nil
}

func calPW1W2(mu *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	// if mu is 0
	if mu.Cmp(bigInt0) == 0 {
		p := new(big.Int).Set(bigInt0)
		w1 := new(big.Int).Set(bigInt0)
		w2 := new(big.Int).Set(bigInt0)
		return p, w1, w2, nil
	}
	// if mu is 1
	if mu.Cmp(bigInt1) == 0 {
		p := new(big.Int).Set(bigInt0)
		w1 := new(big.Int).Set(bigInt1)
		w2 := new(big.Int).Set(bigInt0)
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
	w2Lmt.Add(w2Lmt, bigInt1)
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
		if p.Cmp(bigInt0) == -1 {
			continue
		}
		// if p is 0
		if p.Cmp(bigInt0) == 0 {
			w3 = new(big.Int)
			w4 = new(big.Int)
			return w1, w2, w3, w4, nil
		}
		// if p is 1
		if p.Cmp(bigInt1) == 0 {
			w3 = new(big.Int).Set(bigInt1)
			w4 = new(big.Int).Set(bigInt0)
			return w1, w2, w3, w4, nil
		}

		// (b) hoping p is prime, try to express p = W3^2 + W4^2
		// find a solution u to the equation u^2 = -1 (mod p)
		mul := new(big.Int).Set(p)
		mul.Sub(mul, big.NewInt(1))
		targetMod := new(big.Int).Mod(big.NewInt(-1), p)
		u := new(big.Int).Set(targetMod)
		currMod := new(big.Int).Exp(u, bigInt2, p)
		doubleMU := big.NewInt(2)
		doubleMU.Mul(doubleMU, u)
		uLmt := new(big.Int).Exp(doubleMU, uRetryExponent, nil)
		var lmtFlg bool
		for currMod.Cmp(targetMod) != 0 {
			u.Add(u, bigInt1)
			currMod.Exp(u, bigInt2, p)
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
			sqrtP.Sub(sqrtP, bigInt1)
		}

		dividend := new(big.Int).Set(u)
		divisor := new(big.Int).Set(p)
		w3, err = euclideanStep(dividend, divisor)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		for w3.Cmp(sqrtP) != -1 {
			dividend = divisor
			divisor = w3
			w3, err = euclideanStep(dividend, divisor)
			if err != nil {
				return nil, nil, nil, nil, err
			}
		}
		dividend = divisor
		w4, err = euclideanStep(dividend, w3)
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

func isPerfectSquare(n *big.Int) (*big.Int, bool) {
	sqrt := new(big.Int).Sqrt(n)
	return sqrt, new(big.Int).Mul(sqrt, sqrt).Cmp(n) == 0
}

func euclideanStep(a, b *big.Int) (*big.Int, error) {
	if a.Cmp(b) == -1 {
		a, b = b, a
	}
	q := new(big.Int).Mod(b, a)
	return q, nil
}
