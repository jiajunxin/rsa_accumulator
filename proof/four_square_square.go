package proof

import (
	"context"
	"math/big"

	"lukechampine.com/frand"

	comp "github.com/rsa_accumulator/complex"
)

var (
	sqLagThreshold = new(big.Int).Lsh(big1, 500)
)

// SqLagFourSquares finds the Lagrange four square solution for a very large integer
func SqLagFourSquares(n *big.Int) (FourInt, error) {
	if n.Cmp(sqLagThreshold) < 0 {
		return LagFourSquares(n)
	}

	nc, e := divideN(n)
	var hurwitzGCRD *comp.HurwitzInt
	gcd, l := sqRandTrails(nc)
	var err error
	hurwitzGCRD, err = sqDenouement(nc, l, gcd)
	if err != nil {
		return FourInt{}, err
	}

	// if x'^2 + Y'^2 + Z'^2 + W'^2 = n'
	// then x^2 + Y^2 + Z^2 + W^2 = n for x, Y, Z, W defined by
	// (1 + i)^e * (x' + Y'i + Z'j + W'k) = (x + Yi + Zj + Wk)
	// Gaussian integer: 1 + i
	gi := gaussian1PlusIPow(e)
	hurwitzProd := comp.NewHurwitzInt(gi.R, gi.I, big0, big0, false)
	hurwitzProd.Prod(hurwitzProd, hurwitzGCRD)
	w1, w2, w3, w4 := hurwitzProd.ValInt()
	fi := NewFourInt(w1, w2, w3, w4)
	return fi, nil
}

func sqRandTrails(nc *big.Int) (*comp.GaussianInt, *big.Int) {
	preP := iPool.Get().(*big.Int).Lsh(nc, 1)
	defer iPool.Put(preP)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan sqFindResult)
	randLmt := iPool.Get().(*big.Int).Lsh(big1, sqSetRandBitLen(preP))
	defer iPool.Put(randLmt)
	for i := 0; i < numRoutine; i++ {
		go sqFindSRoutine(ctx, randLmt, preP, resChan)
	}
	res := <-resChan
	return res.gcd, res.l
}

func sqSetRandBitLen(n *big.Int) uint {
	bitLen := n.BitLen()
	ret := uint(float32(bitLen) / 2)
	if ret < 10 {
		ret = 10
	}
	return ret
}

type sqFindResult struct {
	gcd *comp.GaussianInt
	l   *big.Int
}

func sqFindSRoutine(ctx context.Context, randLmt, preP *big.Int, resChan chan<- sqFindResult) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s, p, l, ok, err := sqPickS(randLmt, preP)
			if err != nil {
				panic(err)
			}
			if !ok {
				continue
			}
			gcd := gaussianIntGCD(s, p)
			if !isValidGaussianIntGCD(gcd) {
				continue
			}
			ctx.Done()
			select {
			case resChan <- sqFindResult{gcd: gcd, l: l}:
				return
			default:
				return
			}
		}
	}
}

func sqPickS(randLmt, preP *big.Int) (s, p, l *big.Int, found bool, err error) {
	l = frand.BigIntn(randLmt)
	l.Lsh(l, 1)
	l.Add(l, big1)
	lSq := iPool.Get().(*big.Int).Mul(l, l)
	defer iPool.Put(lSq)
	p = new(big.Int).Sub(preP, lSq)
	if p.Sign() <= 0 {
		return nil, nil, nil, false, nil
	}
	if !p.ProbablyPrime(0) {
		return nil, nil, nil, false, nil
	}

	pMinus1 := iPool.Get().(*big.Int).Sub(p, big1)
	defer iPool.Put(pMinus1)
	powU := iPool.Get().(*big.Int).Rsh(pMinus1, 1)
	defer iPool.Put(powU)
	halfP := iPool.Get().(*big.Int).Rsh(p, 1)
	defer iPool.Put(halfP)
	u := iPool.Get().(*big.Int)
	defer iPool.Put(u)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for i := 0; i < maxUFindingIter; i++ {
		u = frand.BigIntn(halfP)
		u.Lsh(u, 1)

		// test if s^2 = -1 (mod p)
		// if so, continue to the next step, otherwise, repeat this step
		opt.Exp(u, powU, p)
		if opt.Cmp(pMinus1) == 0 {
			found = true
			break
		}
	}
	if !found {
		return nil, nil, nil, false, nil
	}

	// compute s = u^((p - 1) / 4) mod p
	powU.Rsh(powU, 1)
	s = new(big.Int).Exp(u, powU, p)
	return
}

func sqDenouement(n, l *big.Int, gcd *comp.GaussianInt) (*comp.HurwitzInt, error) {
	// compute gcrd(A + Bi + Lj, n), normalized to have integer component
	// Hurwitz integer: A + Bi + Lj
	hurwitzInt := hiPool.Get().(*comp.HurwitzInt).Update(gcd.R, gcd.I, l, big0, false)
	defer hiPool.Put(hurwitzInt)
	// Hurwitz integer: n
	hurwitzN := hiPool.Get().(*comp.HurwitzInt).Update(n, big0, big0, big0, false)
	defer hiPool.Put(hurwitzN)
	gcrd := new(comp.HurwitzInt).GCRD(hurwitzInt, hurwitzN)

	return gcrd, nil
}
