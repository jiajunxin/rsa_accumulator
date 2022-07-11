package proof

import (
	"fmt"
	"math/big"
)

var (
	bigNeg1 = big.NewInt(-1)
	big0    = big.NewInt(0)
	big1    = big.NewInt(1)
	big2    = big.NewInt(2)
	big3    = big.NewInt(3)
	big4    = big.NewInt(4)
	big5    = big.NewInt(5)
	big7    = big.NewInt(7)
	big8    = big.NewInt(8)
	big10   = big.NewInt(10)
)

func isPerfectSquare(n *big.Int) (*big.Int, bool) {
	sqrt := new(big.Int).Sqrt(n)
	return sqrt, new(big.Int).Mul(sqrt, sqrt).Cmp(n) == 0
}

func euclideanDivision(a, b *big.Int) (*big.Int, *big.Int, error) {
	if a.Cmp(b) == -1 {
		a, b = b, a
	}
	if b.Cmp(big0) == 0 {
		return nil, nil, fmt.Errorf("euclideanDivision: valid input should be no less than 0")
	}
	if b.Cmp(big1) == 0 {
		return a, big.NewInt(0), nil
	}
	quotient := new(big.Int)
	remainder := new(big.Int)
	quotient.DivMod(a, b, remainder)
	return quotient, remainder, nil
}

func log2(n *big.Int) *big.Int {
	return big.NewInt(int64(n.BitLen() - 1))
}

func log10(n *big.Int) (*big.Int, error) {
	nCopy := new(big.Int).Set(n)
	if nCopy.Cmp(big1) == -1 {
		return nil, fmt.Errorf("log10: valid input should be no less than 1")
	}
	res := big.NewInt(0)
	for nCopy.Cmp(big10) == 1 {
		nCopy.Div(nCopy, big10)
		res.Add(res, big1)
	}
	if nCopy.Cmp(big10) == 0 {
		res.Add(res, big1)
	}
	return res, nil
}

func rangeDiv(n *big.Int, numDIv int) (res [][2]*big.Int) {
	nd := big.NewInt(int64(numDIv))
	if n.Cmp(nd) < 0 {
		return [][2]*big.Int{
			{big0, new(big.Int).Add(n, big1)},
		}
	}
	batchSize := new(big.Int).Div(n, nd)
	idx := big.NewInt(0)
	for idx.Cmp(nd) == -1 {
		start := new(big.Int).Mul(idx, batchSize)
		end := new(big.Int).Add(start, batchSize)
		if new(big.Int).Add(idx, big1).Cmp(nd) == 0 {
			end = n
		}
		res = append(res, [2]*big.Int{start, end})
		idx.Add(idx, big1)
	}
	return
}