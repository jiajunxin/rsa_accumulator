package proof

import (
	"context"
	"math/big"

	comp "github.com/txaty/go-bigcomplex"
	"lukechampine.com/frand"
)

const bitLenThreshold = 13

// ThreeSquares calculates the three square sum of a given integer
// i.e. target = t1^2 + t2^2 + t3^2
// Please note that we only consider the situation of target = 4N + 1,
// as in our range proof implementation, every integer passed to this function
// is in the form of 4N + 1.
func ThreeSquares(n *big.Int) (ThreeInt, error) {
	rt := new(big.Int).Sqrt(n)
	rt.Rsh(rt, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan ThreeInt)
	for i := 0; i < numRoutine; i++ {
		go routineFindTS(ctx, int64(i), int64(numRoutine), n, rt, resChan)
	}
	return <-resChan, nil
}

func routineFindTS(ctx context.Context, start, step int64, nn, rt *big.Int, resChan chan ThreeInt) {
	cnt := iPool.Get().(*big.Int).SetInt64(start)
	defer iPool.Put(cnt)
	stp := iPool.Get().(*big.Int).SetInt64(step)
	defer iPool.Put(stp)
	p := iPool.Get().(*big.Int)
	defer iPool.Put(p)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			x := new(big.Int).Sub(rt, cnt)
			if x.Sign() < 0 {
				return
			}
			x.Lsh(x, 1)
			p.Mul(x, x).Sub(nn, p)
			if p.Cmp(big2) < 0 {
				return
			}
			if p.BitLen() > bitLenThreshold && !p.ProbablyPrime(0) {
				cnt.Add(cnt, stp)
				break
			}
			gcd := findTwoSquares(p)
			for !isValidGaussianIntGCD(gcd) {
				gcd = findTwoSquares(p)
			}
			select {
			case resChan <- NewThreeInt(x, gcd.R, gcd.I):
				return
			default:
				return
			}
		}
	}
}

func findTwoSquares(n *big.Int) *comp.GaussianInt {
	nMin1 := iPool.Get().(*big.Int).Sub(n, big1)
	defer iPool.Put(nMin1)
	powU := iPool.Get().(*big.Int).Rsh(nMin1, 1)
	defer iPool.Put(powU)
	halfN := iPool.Get().(*big.Int).Rsh(n, 1)
	defer iPool.Put(halfN)
	u := iPool.Get().(*big.Int)
	defer iPool.Put(u)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	s := iPool.Get().(*big.Int)
	defer iPool.Put(s)
	for i := 0; i < maxUFindingIter; i++ {
		u = frand.BigIntn(halfN)
		u.Lsh(u, 1)

		// test if s^2 = -1 (mod p)
		// if so, continue, otherwise, repeat this step
		opt.Exp(u, powU, n)
		if opt.Cmp(nMin1) == 0 {
			// compute s = u^((n - 1) / 4) mod p
			powU.Rsh(powU, 1)
			s.Exp(u, powU, n)
			return gaussianIntGCD(s, n)
		}
	}
	return nil
}

// Verify checks if the three-square sum is equal to the original integer
// i.e. target = t1^2 + t2^2 + t3^2
func Verify(target *big.Int, ti ThreeInt) bool {
	sum := iPool.Get().(*big.Int).SetInt64(0)
	defer iPool.Put(sum)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for i := 0; i < 3; i++ {
		sum.Add(sum, opt.Mul(ti[i], ti[i]))
	}
	return sum.Cmp(target) == 0
}
