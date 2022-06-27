package proof

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// LagrangeRepresentation is the Lagrange representation of a positive integer
// w <- Lagrange(mu), mu = w = w1^2 + w2^2 + w3^2 + w4^2
type LagrangeRepresentation struct {
	w1 *big.Int
	w2 *big.Int
	w3 *big.Int
	w4 *big.Int
}

// Lagrange calculates the Lagrange representation of a positive integer
// Paper: On Diophantine Complexity and Statistical Zero-Knowledge Arguments
// Link: https://eprint.iacr.org/2003/105
func Lagrange(mu *big.Int) (LagrangeRepresentation, error) {
	// write mu in the form mu = 2^t(2k + 1)
	var (
		t, k int
	)
	// copy mu for modification
	muCopy := new(big.Int).Set(mu)
	for muCopy.Bit(0) == 0 {
		t++
		// right shift
		muCopy.Rsh(muCopy, 1)
	}
	k = muCopy.BitLen() - 1
	fmt.Printf("mu: %d, t: %d, k: %d\n", mu, t, k)

	// if t = 1
	if t == 1 {
		w1, w2, w3, w4, err := calW1W2W3W4(mu)
		if err != nil {
			return LagrangeRepresentation{}, err
		}
		return LagrangeRepresentation{w1, w2, w3, w4}, nil
	}

	// if t is odd but not 1
	if t%2 == 1 {
		fmt.Println("t is odd but not 1")
		w1, w2, w3, w4, err := calW1W2W3W4(mu)
		if err != nil {
			return LagrangeRepresentation{}, err
		}
		s := new(big.Int).SetInt64(2)
		s.Exp(s, new(big.Int).SetInt64(int64((t-1)/2)), nil)
		w1.Mul(s, w1)
		w2.Mul(s, w2)
		w3.Mul(s, w3)
		w4.Mul(s, w4)
		return LagrangeRepresentation{w1, w2, w3, w4}, nil
	}
	// if t is even
	w1, w2, w3, w4, err := calW1W2W3W4(mu)
	if err != nil {
		return LagrangeRepresentation{}, err
	}
	bigInt2 := new(big.Int).SetInt64(2)
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
	w1.Mul(s, w1)
	w2.Mul(s, w2)
	w3.Mul(s, w3)
	w4.Mul(s, w4)
	return LagrangeRepresentation{w1, w2, w3, w4}, nil
}

func calPW1W2(mu *big.Int) (*big.Int, *big.Int, *big.Int, error) {
	fmt.Println("calPW1W2")
	// choose random w1, w2 such that exactly one of w1, w2 is even
	w1Lmt := new(big.Int)
	w1Lmt.Sqrt(mu)
	w1Lmt.Add(w1Lmt, big.NewInt(1))
	// randomly choose w1 within [0, sqrt(mu)]
	w1, err := rand.Int(rand.Reader, w1Lmt)
	if err != nil {
		return nil, nil, nil, err
	}
	w1Sq := new(big.Int).Mul(w1, w1)
	w2Lmt := new(big.Int).Set(w1Sq)
	w2Lmt.Sub(mu, w1Lmt)
	w2Lmt.Sqrt(w2Lmt)
	w2Lmt.Add(w2Lmt, big.NewInt(1))
	// randomly choose w2 within [0, sqrt(mu - w1^2)]
	w2, err := rand.Int(rand.Reader, w2Lmt)
	if err != nil {
		return nil, nil, nil, err
	}
	w2Sq := new(big.Int).Mul(w2, w2)
	// p <- mu - w1^2 - w2^2, now p = 1 (mod 4)
	p := new(big.Int).Sub(mu, w1Sq)
	p.Sub(p, w2Sq)

	fmt.Printf("p: %d\n", p.Int64())
	return p, w1, w2, nil
}

func calW1W2W3W4(mu *big.Int) (*big.Int, *big.Int, *big.Int, *big.Int, error) {
	fmt.Println("calW1W2W3W4")
	var (
		p, w1, w2, w3, w4 *big.Int
		err               error
	)
	for {
		// (a) choose random w1, w2, and calculate p <- mu - w1^2 - w2^2
		p, w1, w2, err = calPW1W2(mu)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		// if p is zero
		if p.Cmp(big.NewInt(0)) == 0 {
			w3 = new(big.Int)
			w4 = new(big.Int)
			return w1, w2, w3, w4, nil
		}

		// (b) hoping p is prime, try to express p = w3^2 + w4^2
		// find a solution u to the equation u^2 = -1 (mod p)
		mul := new(big.Int).Set(p)
		mul.Sub(mul, big.NewInt(1))
		// multiply mul by p until mul is the square of u
		//u, flg := isPerfectSquare(mul)
		//var cnt int
		//for !flg {
		//	mul.Mul(mul, p)
		//	fmt.Println(mul.Int64())
		//	u, flg = isPerfectSquare(mul)
		//	cnt++
		//	if cnt > 100 {
		//		return nil, nil, nil, nil, errors.New("cannot find u")
		//	}
		//}

		// brute force approach
		targetMod := new(big.Int).Mod(big.NewInt(-1), p)
		u := new(big.Int).Set(targetMod)
		//uSq := new(big.Int).Mul(u, u)
		bigInt2 := big.NewInt(2)
		currMod := new(big.Int).Exp(u, bigInt2, p)
		fmt.Printf("targetMod: %d\n", targetMod.Int64())
		//cnt := 0
		for currMod.Cmp(targetMod) != 0 {
			u.Add(u, big.NewInt(1))
			//uSq.Mul(u, u)
			currMod.Exp(u, bigInt2, p)
			fmt.Printf("u: %d, mod: %d\n", u.Int64(), currMod.Int64())
			//cnt++
			//if cnt > 100 {
			//	return nil, nil, nil, nil, errors.New("cannot find u")
			//}
		}
		fmt.Printf("u: %d\n", u.Int64())

		fmt.Printf("mul: %d\n", mul.Int64())
		// apply Euclidean algorithm to (u, p), and take the first two remainders that are less than sqrt(p)
		floatP := new(big.Float).SetInt(p)
		floatSqrtP := new(big.Float).Sqrt(floatP)
		sqrtP := new(big.Int).Sqrt(p)
		if floatSqrtP.IsInt() {
			sqrtP.Sub(sqrtP, big.NewInt(1))
		}

		dividend := new(big.Int).Set(p)
		divisor := new(big.Int).Set(u)
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
		// if p != w3^2 + w4^2, then p is not prime, so go back to the p calculation
		w3Sq := new(big.Int).Mul(w3, w3)
		w4Sq := new(big.Int).Mul(w4, w4)
		if p.Cmp(new(big.Int).Add(w3Sq, w4Sq)) == 0 {
			break
		}
	}
	return w1, w2, w3, w4, nil
}

func isPerfectSquare(n *big.Int) (*big.Int, bool) {
	sqrt := new(big.Int).Sqrt(n)
	return sqrt, new(big.Int).Mul(sqrt, sqrt).Cmp(n) == 0
}

func euclideanStep(a, b *big.Int) (*big.Int, error) {
	if a.Cmp(b) == -1 {
		return nil, fmt.Errorf("invalid input for Euclidean algorithm, a: %d < b %d", a.Int64(), b.Int64())
	}
	q := new(big.Int).Mod(b, a)
	return q, nil
}
