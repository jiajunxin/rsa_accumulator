package proof

import (
	"context"
	"math/big"

	"lukechampine.com/frand"
)

// ThreeSquares calculates the three square sum for 4N + 1
func ThreeSquares(n *big.Int) (ThreeInt, error) {
	nc := iPool.Get().(*big.Int).Lsh(n, 2)
	defer iPool.Put(nc)
	nc.Add(nc, big1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan ThreeInt)
	ncBitLen := nc.BitLen()
	randLmt := iPool.Get().(*big.Int).Lsh(big1, uint(ncBitLen/2))
	//randLmt := iPool.Get().(*big.Int).Sqrt(nc)
	defer iPool.Put(randLmt)
	for i := 0; i < numRoutine; i++ {
		go findRoutineTS(ctx, randLmt, nc, resChan)
	}
	return <-resChan, nil
}

func findRoutineTS(ctx context.Context, randLmt, preP *big.Int, resChan chan<- ThreeInt) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			a, s, p, ok, err := pickASP(randLmt, preP)
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
			case resChan <- NewThreeInt(a, gcd.R, gcd.I):
				return
			default:
				return
			}
		}
	}
}

func pickASP(randLmt, preP *big.Int) (a, s, p *big.Int, found bool, err error) {
	a = frand.BigIntn(randLmt)
	a.Lsh(a, 1)
	aSq := iPool.Get().(*big.Int).Mul(a, a)
	defer iPool.Put(aSq)
	p = new(big.Int).Sub(preP, aSq)
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

// VerifyTS checks if the three-square sum is equal to the original integer
// i.e. target = t1^2 + t2^2 + t3^2
func VerifyTS(target *big.Int, ti ThreeInt) bool {
	check := iPool.Get().(*big.Int).Lsh(target, 2)
	defer iPool.Put(check)
	check.Add(check, big1)
	return verifyTS(check, ti)
}

func verifyTS(target *big.Int, ti ThreeInt) bool {
	sum := iPool.Get().(*big.Int).SetInt64(0)
	defer iPool.Put(sum)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for i := 0; i < 3; i++ {
		sum.Add(sum, opt.Mul(ti[i], ti[i]))
	}
	return sum.Cmp(target) == 0
}
