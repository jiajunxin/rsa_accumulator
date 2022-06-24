package proof

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// LagrangeRepresentation is the Lagrange representation of a positive integer
// w <- Lagrange(mu), mu = w = w1^2 + w2^2 + w3^2 + w4^2
type LagrangeRepresentation struct {
	w1 big.Int
	w2 big.Int
	w3 big.Int
	w4 big.Int
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
	k = int(muCopy.BitLen() - 1)
	fmt.Printf("mu: %d, t: %d, k: %d\n", mu, t, k)

	// if t = 1
	if t == 1 {
		p, err := calP(mu)
		if err != nil {
			return LagrangeRepresentation{}, err
		}
		// hoping p is prime, try to express p = w3^2 + w4^2
		// find a solution u to the equation u^2 = -1 (mod p)
		mul := new(big.Int).Set(p)
		mul.Sub(mul, big.NewInt(1))
		// multiply mul by p until it is the square of u
		u, flg := isPerfectSquare(mul)
		for !flg {
			mul.Mul(mul, p)
			u, flg = isPerfectSquare(mul)
		}
		// apply Euclidean algorithm to (u, p), and take the first two remainders that are less than sqrt(p)
		floatP := new(big.Float).SetInt(p)
		sqrtP := new(big.Float).Sqrt(floatP)

		fmt.Println("u: ", u)
		fmt.Println("sqrt of p: ", sqrtP)
	}

	// if t is odd but not 1

	// if t is even

	return LagrangeRepresentation{}, nil
}

func calP(mu *big.Int) (*big.Int, error) {
	// choose random w1, w2 such that exactly one of w1, w2 is even
	w1Lmt := new(big.Int)
	w1Lmt.Sqrt(mu)
	w1Lmt.Add(w1Lmt, big.NewInt(1))
	// randomly choose w1 within [0, sqrt(mu)]
	w1, err := rand.Int(rand.Reader, w1Lmt)
	if err != nil {
		return nil, err
	}
	w1Sq := new(big.Int).Mul(w1, w1)
	w2Lmt := new(big.Int).Set(w1Sq)
	w2Lmt.Sub(mu, w1Lmt)
	w2Lmt.Sqrt(w2Lmt)
	w2Lmt.Add(w2Lmt, big.NewInt(1))
	// randomly choose w2 within [0, sqrt(mu - w1^2)]
	w2, err := rand.Int(rand.Reader, w2Lmt)
	if err != nil {
		return nil, err
	}
	w2Sq := new(big.Int).Mul(w2, w2)
	// p <- mu - w1^2 - w2^2, now p = 1 (mod 4)
	p := new(big.Int).Sub(mu, w1Sq)
	p.Sub(p, w2Sq)

	return p, nil
}

func isPerfectSquare(n *big.Int) (*big.Int, bool) {
	sqrt := new(big.Int).Sqrt(n)
	return sqrt, sqrt.Mul(sqrt, sqrt).Cmp(n) == 0
}
