package proof

import (
	"context"
	crand "crypto/rand"
	"math"
	"math/big"

	bc "github.com/tommytim0515/go-bigcomplex"
)

// UnconditionalLagrangeFourSquares calculates the Lagrange four squares for a given non-positive integer
// the method doesn't rely on the Extended Riemann Hypothesis (ERH)
func UnconditionalLagrangeFourSquares(n *big.Int) (FourInt, error) {
	if n.Sign() == 0 {
		res := NewFourInt(precomputedHurwitzGCRDs[0].ValInt())
		return res, nil
	}
	nc, e := divideN(n)
	var hurwitzGCRD *bc.HurwitzInt

	if nc.Cmp(big8) <= 0 {
		hurwitzGCRD = precomputedHurwitzGCRDs[nc.Int64()]
	} else {
		x, y, p, r1, s, primes, err := initTrail(nc)
		if err != nil {
			return FourInt{}, err
		}
		// compute u, v
		u, v, err := computeUV(r1, nc, primes)
		if err != nil {
			return FourInt{}, err
		}
		var up, vp *big.Int
		// compute U -> up, V -> vp
		if s == nil {
			up = big.NewInt(1)
			vp = big.NewInt(0)
		} else {
			gcd := giPool.Get().(*bc.GaussianInt)
			defer giPool.Put(gcd)
			gOpt1 := giPool.Get().(*bc.GaussianInt)
			defer giPool.Put(gOpt1)
			gOpt2 := giPool.Get().(*bc.GaussianInt)
			defer giPool.Put(gOpt2)
			gOpt1.Update(p, big0)
			gOpt2.Update(s, big1)
			gcd.GCD(gOpt1, gOpt2)
			up = gcd.R
			vp = gcd.I
		}
		uvi := bc.NewGaussianInt(u, v)
		uPvPI := bc.NewGaussianInt(up, vp)
		zwi := new(bc.GaussianInt).Prod(uvi, uPvPI)
		hOpt1 := hiPool.Get().(*bc.HurwitzInt)
		defer hiPool.Put(hOpt1)
		hOpt1.Update(n, big0, big0, big0, false)
		hOpt2 := hiPool.Get().(*bc.HurwitzInt)
		defer hiPool.Put(hOpt2)
		hOpt2.Update(x, y, zwi.R, zwi.I, false)
		hurwitzGCRD = new(bc.HurwitzInt).GCRD(hOpt1, hOpt2)
	}

	// if x'^2 + Y'^2 + Z'^2 + W'^2 = n'
	// then x^2 + Y^2 + Z^2 + W^2 = n for x, Y, Z, W defined by
	// (1 + i)^e * (x' + Y'i + Z'j + W'k) = (x + Yi + Zj + Wk)
	// Gaussian integer: 1 + i
	gi := gaussian1PlusIPow(e)
	hurwitzProd := bc.NewHurwitzInt(gi.R, gi.I, big0, big0, false)
	hurwitzProd.Prod(hurwitzProd, hurwitzGCRD)
	w1, w2, w3, w4 := hurwitzProd.ValInt()
	fs := NewFourInt(w1, w2, w3, w4)
	return fs, nil
}

func initTrail(n *big.Int) (x, y, p, r1, s *big.Int, primes []*big.Int, err error) {
	logN := log2(n)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	i := 1
	for i < logN {
		for _, prime := range primes {
			bigI := big.NewInt(int64(i))
			if opt.Mod(bigI, prime).Cmp(big0) == 0 {
				break
			}
			primes = append(primes, bigI)
		}
		i += 4
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan initTrailResult)
	for i := 0; i < numRoutine; i++ {
		go initTrailRoutine(ctx, n, primes, resChan)
	}
	res := <-resChan
	return res.x, res.y, res.p, res.r1, res.s, primes, res.err
}

func initTrailRoutine(ctx context.Context, n *big.Int, primes []*big.Int, resChan chan<- initTrailResult) {
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	mod := iPool.Get().(*big.Int)
	defer iPool.Put(mod)
	r := iPool.Get().(*big.Int)
	defer iPool.Put(r)
	pMinus1 := iPool.Get().(*big.Int)
	defer iPool.Put(pMinus1)
	u := iPool.Get().(*big.Int)
	defer iPool.Put(u)
	powU := iPool.Get().(*big.Int)
	defer iPool.Put(powU)
	var (
		x, y, p, r1, s *big.Int
		err            error
	)
	for {
		select {
		case <-ctx.Done():
			return
		default:
			r, x, y, err = randomChoiceXY(n)
			r1 = big.NewInt(1)
			for _, prime := range primes {
				opt = big.NewInt(1)
				for mod.Mod(r, opt).Cmp(big0) == 0 {
					r1.Mul(r1, opt)
					opt.Mul(opt, prime)
				}
			}
			p = new(big.Int).Div(r, r1)
			if p.Cmp(big1) == 0 {
				select {
				case resChan <- initTrailResult{x, y, p, r1, s, err}:
					ctx.Done()
					return
				default:
					return
				}
			}
			pMinus1.Sub(p, big1)
			u, err = crand.Int(crand.Reader, pMinus1)
			if err != nil {
				select {
				case resChan <- initTrailResult{x, y, p, r1, s, err}:
					return
				default:
					return
				}
			}
			u.Add(u, big1)
			powU.Rsh(pMinus1, 2)
			s = new(big.Int).Exp(u, powU, p)
			if s.Mul(s, s).Cmp(pMinus1) == 0 {
				select {
				case resChan <- initTrailResult{x, y, p, r1, s, err}:
					ctx.Done()
					return
				default:
					return
				}
			}
		}
	}
}

type initTrailResult struct {
	x, y, p, r1, s *big.Int
	err            error
}

func randomChoiceXY(n *big.Int) (r, x, y *big.Int, err error) {
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	gcd := iPool.Get().(*big.Int)
	defer iPool.Put(gcd)
	r = new(big.Int)

	for {
		// randomly pick x, y from [1, n]
		if x, err = crand.Int(crand.Reader, n); err != nil {
			return nil, nil, nil, err
		}
		x.Add(x, big1)
		if y, err = crand.Int(crand.Reader, n); err != nil {
			return nil, nil, nil, err
		}
		y.Add(y, big1)
		// compute r := (-(x^2 + y^2)) mod n
		r.Mul(x, x)
		opt.Mul(y, y)
		r.Add(r, opt)
		r.Neg(r)
		r.Mod(r, n)
		// check if r = 1 (mod n)
		if r.Cmp(big1) != 0 {
			continue
		}
		// check if gcd(r, n) = 1
		if gcd.GCD(nil, nil, r, n).Cmp(big1) != 0 {
			continue
		}
		break
	}
	return
}

func computeUV(r1, n *big.Int, primes []*big.Int) (u, v *big.Int, err error) {
	uvi := bc.NewGaussianInt(big1, big0)
	ss, err := computeSquares(n)
	if err != nil {
		return nil, nil, err
	}
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	mod := iPool.Get().(*big.Int)
	defer iPool.Put(mod)
	gOpt := giPool.Get().(*bc.GaussianInt)
	defer giPool.Put(gOpt)
	for _, prime := range primes {
		x, y := ss.findXY(prime)
		gOpt.Update(x, y)
		// determine vl
		opt.Set(prime)
		for mod.Mod(r1, opt).Cmp(big0) == 0 {
			uvi.Prod(uvi, gOpt)
			opt.Mul(opt, prime)
		}
	}
	u = uvi.R
	v = uvi.I
	return
}

func computeSquares(n *big.Int) (*squareCache, error) {
	lmt := log2(n)
	lmt = int(math.Sqrt(float64(lmt)))
	for i := sCache.max; i <= lmt; i++ {
		bigI := big.NewInt(int64(i))
		sq := new(big.Int).Mul(bigI, bigI)
		sCache.add(bigI, sq)
		sCache.max = i
	}
	return sCache, nil
}
