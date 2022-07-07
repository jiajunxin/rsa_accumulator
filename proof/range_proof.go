package proof

import (
	"crypto/rand"
	"math/big"
)

// RangeProofProver refers to the Prover in zero-knowledge integer range proof
type RangeProofProver struct {
	G *big.Int
	H *big.Int
	X *big.Int
	N *big.Int
}

func (r *RangeProofProver) commitX() (c1, c2, c3, c4 *big.Int, err error) {
	// calculate lagrange four squares for x
	fs, err := LagrangeFourSquares(r.X)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	// calculate commitment for x
	rc, err := NewPseudoRangeProofRandomCoins(r.N)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	c1, c2, c3, c4 = fs.RangeProofCommit(r.G, r.H, *rc)
	return
}

// RangeProofVerifier refers to the Verifier in zero-knowledge integer range proof
type RangeProofVerifier struct {
}

// RangeProofRandomCoins is the random coins used in range proof
type RangeProofRandomCoins struct {
	r1 *big.Int
	r2 *big.Int
	r3 *big.Int
	r4 *big.Int
}

// NewPseudoRangeProofRandomCoins creates a new random coins for range proof
func NewPseudoRangeProofRandomCoins(n *big.Int) (*RangeProofRandomCoins, error) {
	r1, err := pseudoFreshRandomCoins(n)
	if err != nil {
		return nil, err
	}
	r2, err := pseudoFreshRandomCoins(n)
	if err != nil {
		return nil, err
	}
	r3, err := pseudoFreshRandomCoins(n)
	if err != nil {
		return nil, err
	}
	r4, err := pseudoFreshRandomCoins(n)
	if err != nil {
		return nil, err
	}
	return &RangeProofRandomCoins{r1, r2, r3, r4}, nil
}

func pseudoFreshRandomCoins(n *big.Int) (*big.Int, error) {
	n = new(big.Int).Set(n).Add(n, big1)
	res, err := rand.Int(rand.Reader, n)
	if err != nil {
		return nil, err
	}
	return res, nil
}
