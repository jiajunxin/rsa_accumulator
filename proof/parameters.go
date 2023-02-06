package proof

import (
	"crypto/rand"
	"math/big"
)

const (
	// security parameter for range proof and proof of exponentiation
	securityParam = 128
	int3Len       = 3
	int4Len       = 4
)

// PublicParameters holds public parameters initialized during the setup procedure
type PublicParameters struct {
	N *big.Int
	G *big.Int
	H *big.Int
}

// NewPublicParameters generates a new public parameter configuration
func NewPublicParameters(n, g, h *big.Int) *PublicParameters {
	return &PublicParameters{
		N: n,
		G: g,
		H: h,
	}
}

// Int4 is the 4-number big integer group
type Int4 [int4Len]*big.Int

// NewInt4 creates a new 4-number group, in descending order
func NewInt4(w1 *big.Int, w2 *big.Int, w3 *big.Int, w4 *big.Int) Int4 {
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
	return Int4{w1, w2, w3, w4}
}

// Mul multiplies all the 4 numbers by n
func (f *Int4) Mul(n *big.Int) {
	for i := 0; i < 4; i++ {
		f[i].Mul(f[i], n)
	}
}

// Div divides all the 4 numbers by n
func (f *Int4) Div(n *big.Int) {
	for i := 0; i < 4; i++ {
		f[i].Div(f[i], n)
	}
}

// String stringnifies the Int4 object
func (f *Int4) String() string {
	res := "{"
	for i := 0; i < 3; i++ {
		res += f[i].String()
		res += ", "
	}
	res += f[3].String()
	res += "}"
	return res
}

// newFourRandCoins creates a new random coins for range proof
func newFourRandCoins(n *big.Int) (coins Int4, err error) {
	for i := 0; i < 4; i++ {
		coins[i], err = freshRandCoin(n)
		if err != nil {
			return
		}
	}
	return
}

// newThreeRandCoins creates three new random coins for range proof
func newThreeRandCoins(n *big.Int) (coins Int3, err error) {
	for i := 0; i < 3; i++ {
		coins[i], err = freshRandCoin(n)
		if err != nil {
			return
		}
	}
	return
}

// freshRandCoin creates a new fresh random coin in [0, n]
func freshRandCoin(n *big.Int) (*big.Int, error) {
	lmt := iPool.Get().(*big.Int).Set(n)
	defer iPool.Put(lmt)
	lmt.Add(lmt, big1)
	res, err := rand.Int(rand.Reader, lmt)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Int3 is the 3-number big integer group
type Int3 [int3Len]*big.Int

// NewInt3 creates a new 3-number group, in descending order
func NewInt3(w1 *big.Int, w2 *big.Int, w3 *big.Int) Int3 {
	w1.Abs(w1)
	w2.Abs(w2)
	w3.Abs(w3)
	// sort the three big integers in descending order
	if w1.Cmp(w2) == -1 {
		w1, w2 = w2, w1
	}
	if w1.Cmp(w3) == -1 {
		w1, w3 = w3, w1
	}
	if w2.Cmp(w3) == -1 {
		w2, w3 = w3, w2
	}
	return Int3{w1, w2, w3}
}
