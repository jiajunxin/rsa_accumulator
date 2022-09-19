package proof

import (
	"context"
	"math/big"

	comp "github.com/rsa_accumulator/complex"
	"lukechampine.com/frand"
)

var big10000 = big.NewInt(10000)

// ThreeSquareNew calculates the three square sum for 4N + 1
func ThreeSquareNew(n *big.Int) (ThreeInt, error) {
	nn := iPool.Get().(*big.Int).Lsh(n, 2)
	defer iPool.Put(nn)
	nn.Add(nn, big1)
	rt := iPool.Get().(*big.Int).Sqrt(nn)
	defer iPool.Put(rt)
	rt.Sub(rt, big1)
	rt.Rsh(rt, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan ThreeInt)
	for i := 0; i < numRoutine; i++ {
		go routineFindTS(ctx, int64(i), int64(numRoutine), nn, rt, resChan)
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
	x := iPool.Get().(*big.Int)
	defer iPool.Put(x)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			x.Sub(rt, cnt)
			x.Lsh(x, 1)
			p.Mul(x, x)
			p.Sub(nn, p)
			if p.Cmp(big10000) >= 0 && !p.ProbablyPrime(0) {
				cnt.Add(cnt, stp)
				continue
			}
			gGCD, ok := findTwoSquares(p)
			if !ok {
				cnt.Add(cnt, stp)
				continue
			}
			ts := NewThreeInt(x, gGCD.R, gGCD.I)
			if !verifyTS(nn, ts) {
				cnt.Add(cnt, stp)
				continue
			}
			select {
			case resChan <- ts:
				return
			default:
				return
			}
		}
	}
}

func findTwoSquares(n *big.Int) (*comp.GaussianInt, bool) {
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
		if halfN.Cmp(big0) <= 0 {
			return nil, false
		}
		u = frand.BigIntn(halfN)
		u.Lsh(u, 1)

		// test if s^2 = -1 (mod p)
		// if so, continue, otherwise, repeat this step
		opt.Exp(u, powU, n)
		if opt.Cmp(nMin1) == 0 {
			// compute s = u^((n - 1) / 4) mod p
			powU.Rsh(powU, 1)
			s.Exp(u, powU, n)
			return gaussianIntGCD(s, n), true
		}
	}
	return nil, false
}
