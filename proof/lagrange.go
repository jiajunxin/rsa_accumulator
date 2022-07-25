package proof

import (
	"context"
	crand "crypto/rand"
	"errors"
	"math"
	"math/big"
	"math/rand"
	"runtime"
	"time"

	comp "github.com/rsa_accumulator/complex"
)

const (
	randLmtThreshold0 = 16
	randLmtThreshold1 = 32
	randLmtThreshold2 = 64
	preComputeLmt     = 20
)

var (
	// precomputed Hurwitz GCRDs for small integers
	precomputedHurwitzGCRDs = [preComputeLmt + 1]*comp.HurwitzInt{
		// 0's precomputed Hurwitz GCRD: 0, 0, 0, 0
		comp.NewHurwitzInt(big0, big0, big0, big0, false),
		// 1's precomputed Hurwitz GCRD: 1, 0, 0, 0
		comp.NewHurwitzInt(big1, big0, big0, big0, false),
		// 2's precomputed Hurwitz GCRD: 1, 1, 0, 0
		comp.NewHurwitzInt(big1, big1, big0, big0, false),
		// 3's precomputed Hurwitz GCRD: 1, 1, 1, 0
		comp.NewHurwitzInt(big1, big1, big1, big0, false),
		// 4's precomputed Hurwitz GCRD: 2, 0, 0, 0
		comp.NewHurwitzInt(big2, big0, big0, big0, false),
		// 5's precomputed Hurwitz GCRD: 2, 1, 0, 0
		comp.NewHurwitzInt(big2, big1, big0, big0, false),
		// 6's precomputed Hurwitz GCRD: 2, 1, 1, 0
		comp.NewHurwitzInt(big2, big1, big1, big0, false),
		// 7's precomputed Hurwitz GCRD: 2, 1, 1, 1
		comp.NewHurwitzInt(big2, big1, big1, big1, false),
		// 8's precomputed Hurwitz GCRD: 2, 2, 0, 0
		comp.NewHurwitzInt(big2, big2, big0, big0, false),
		// 9's precomputed Hurwitz GCRD: 2, 2, 1, 0
		comp.NewHurwitzInt(big2, big2, big1, big0, false),
		// 10's precomputed Hurwitz GCRD: 2, 2, 1, 1
		comp.NewHurwitzInt(big2, big2, big1, big1, false),
		// 11's precomputed Hurwitz GCRD: 3, 1, 1, 0
		comp.NewHurwitzInt(big3, big1, big1, big0, false),
		// 12's precomputed Hurwitz GCRD: 3, 1, 1, 1
		comp.NewHurwitzInt(big3, big1, big1, big1, false),
		// 13's precomputed Hurwitz GCRD: 3, 2, 0, 0
		comp.NewHurwitzInt(big3, big2, big0, big0, false),
		// 14's precomputed Hurwitz GCRD: 3, 2, 1, 0
		comp.NewHurwitzInt(big3, big2, big1, big0, false),
		// 15's precomputed Hurwitz GCRD: 3, 2, 1, 1
		comp.NewHurwitzInt(big3, big2, big1, big1, false),
		// 16's precomputed Hurwitz GCRD: 4, 0, 0, 0
		comp.NewHurwitzInt(big4, big0, big0, big0, false),
		// 17's precomputed Hurwitz GCRD: 4, 1, 0, 0
		comp.NewHurwitzInt(big4, big1, big0, big0, false),
		// 18's precomputed Hurwitz GCRD: 4, 1, 1, 0
		comp.NewHurwitzInt(big4, big1, big1, big0, false),
		// 19's precomputed Hurwitz GCRD: 4, 1, 1, 1
		comp.NewHurwitzInt(big4, big1, big1, big1, false),
		// 20's precomputed Hurwitz GCRD: 4, 2, 0, 0
		comp.NewHurwitzInt(big4, big2, big0, big0, false),
	}
	bigPreComputeLmt = big.NewInt(preComputeLmt)
	numRoutine       = runtime.NumCPU()
)

// LagrangeFourSquares calculates the Lagrange four squares representation of a positive integer
// Paper: Finding the Four Squares in Lagrange’s Theorem
// Link: http://pollack.uga.edu/finding4squares.pdf (page 6)
// The input should be an odd positive integer no less than 9
func LagrangeFourSquares(n *big.Int) (FourNum, error) {
	if n.Sign() == 0 {
		res := NewFourNum(precomputedHurwitzGCRDs[0].ValInt())
		return res, nil
	}
	nc, e := divideN(n)
	var hurwitzGCRD *comp.HurwitzInt

	if nc.Cmp(bigPreComputeLmt) <= 0 {
		hurwitzGCRD = precomputedHurwitzGCRDs[nc.Int64()]
	} else {
		primeProd, err := preCompute(nc)
		if err != nil {
			return FourNum{}, err
		}
		var gcd *comp.GaussianInt
		for {
			s, p, err := randTrails(nc, primeProd)
			if err != nil {
				return FourNum{}, err
			}
			gcd = gaussianIntGCD(s, p)
			// continue if the GCD is valid
			if isValidGaussianIntGCD(gcd) {
				break
			}
		}
		//fmt.Println(gcd)
		hurwitzGCRD, err = denouement(nc, gcd)
		if err != nil {
			return FourNum{}, err
		}
	}

	// if x'^2 + Y'^2 + Z'^2 + W'^2 = n'
	// then x^2 + Y^2 + Z^2 + W^2 = n for x, Y, Z, W defined by
	// (1 + i)^e * (x' + Y'i + Z'j + W'k) = (x + Yi + Zj + Wk)
	// Gaussian integer: 1 + i
	gi := gaussian1PlusIPow(e)
	hurwitzProd := comp.NewHurwitzInt(gi.R, gi.I, big0, big0, false)
	hurwitzProd.Prod(hurwitzProd, hurwitzGCRD)
	w1, w2, w3, w4 := hurwitzProd.ValInt()
	fs := NewFourNum(w1, w2, w3, w4)
	return fs, nil
}

func isValidGaussianIntGCD(gcd *comp.GaussianInt) bool {
	absR := iPool.Get().(*big.Int)
	defer iPool.Put(absR)
	absR.Abs(gcd.R)
	absI := iPool.Get().(*big.Int)
	defer iPool.Put(absI)
	absI.Abs(gcd.I)
	rCmp1 := absR.Cmp(big1)
	rSign := absR.Sign()
	iCmp1 := absI.Cmp(big1)
	iSign := absI.Sign()
	if rCmp1 == 0 && iSign == 0 {
		return false
	}
	if rSign == 0 && iCmp1 == 0 {
		return false
	}
	if rCmp1 == 0 && iCmp1 == 0 {
		return false
	}
	return true
}

func divideN(n *big.Int) (*big.Int, int) {
	// n = 2^e * n', n' is odd
	nc := new(big.Int).Set(n)
	var e int
	for nc.Bit(0) == 0 {
		nc.Rsh(nc, 1)
		e++
	}
	return nc, e
}

// gaussian1PlusIPow calculates Gaussian integer (1 + i)^e
func gaussian1PlusIPow(e int) *comp.GaussianInt {
	if e == 0 {
		return comp.NewGaussianInt(big1, big0)
	}
	if gi, ok := giCache[e]; ok {
		return gi
	}
	gaussian1PlusI := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gaussian1PlusI)
	gaussian1PlusI.Update(big1, big1)

	gaussianProd := comp.NewGaussianInt(big1, big0)
	idx := e
	for idx > 0 {
		gaussianProd.Prod(gaussianProd, gaussian1PlusI)
		idx--
	}
	gi := new(comp.GaussianInt)
	gi.Update(gaussianProd.R, gaussianProd.I)
	giCache[e] = gi
	return gaussianProd
}

// preCompute determine the primes not exceeding log n and compute their product
// the function only handles positive integers larger than 8
func preCompute(n *big.Int) (*big.Int, error) {
	if n.Cmp(big8) <= 0 {
		return nil, errors.New("n should be larger than 8")
	}
	logN := log(n)
	if logN <= pCache.max {
		prod, err := pCache.findPrimeProd(logN)
		if err != nil {
			return nil, err
		}
		return prod, nil
	}
	prod := iPool.Get().(*big.Int)
	defer iPool.Put(prod)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	prod.Set(pCache.m[pCache.max])
	for idx := pCache.max + 2; idx < logN; idx += 2 {
		pCache.checkAddPrime(idx, prod, opt)
	}
	return new(big.Int).Set(prod), nil
}

func randTrails(n, primeProd *big.Int) (*big.Int, *big.Int, error) {
	// use goroutines to choose a random number between [0, n^5 / 2 / numRoutine]
	// then construct k based on the random number
	// and check the validity of the trails
	// p = M * n * k - 1, pre-p = M * n
	preP := new(big.Int).Mul(primeProd, n)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	resChan := make(chan findSResult)
	randLmtBitLen := n.BitLen()
	if randLmtBitLen < randLmtThreshold2 {
		randLmt := setInitRandLmt(n)
		//randLmt := new(big.Int).Sqrt(n)
		randLmt.Rsh(randLmt, 1)
		randLmt.Div(randLmt, big.NewInt(int64(numRoutine)))
		randLmt.Add(randLmt, big1)

		var (
			mul  = big.NewInt(int64(2 * numRoutine)) // 2 * numRoutine
			adds []*big.Int
		)
		for i := 0; i <= numRoutine; i++ {
			adds = append(adds, big.NewInt(int64(2*i+1))) // 2i+1
		}
		for _, add := range adds {
			go findSRoutine(ctx, add, mul, randLmt, preP, resChan)
		}
	} else {
		bl := setInitRandBitLen(randLmtBitLen)
		for i := 0; i < numRoutine; i++ {
			go findLargeSRoutine(ctx, bl, preP, resChan)
		}
	}
	res := <-resChan
	return res.s, res.p, res.err
}

func setInitRandLmt(n *big.Int) *big.Int {
	bitLen := n.BitLen()
	if bitLen < randLmtThreshold0 {
		return new(big.Int).Exp(n, big4, nil)
	}
	if bitLen < randLmtThreshold1 {
		return new(big.Int).Exp(n, big3, nil)
	}
	if bitLen < randLmtThreshold2 {
		return new(big.Int).Exp(n, big2, nil)
	}
	return new(big.Int).Set(n)
}

func setInitRandBitLen(bitLen int) int {
	lenF := 20 + 2*math.Log(float64(bitLen))
	return int(math.Round(lenF))
}

func findSRoutine(ctx context.Context, mul, add, randLmt, preP *big.Int, resChan chan<- findSResult) {
	rg := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s, p, ok, err := pickS(rg, mul, add, randLmt, preP)
			if err != nil {
				select {
				case resChan <- findSResult{err: err}:
					return
				default:
					return
				}
			}
			if !ok {
				continue
			}
			ctx.Done()
			select {
			case resChan <- findSResult{s: s, p: p}:
				return
			default:
				return
			}
		}
	}
}

type findSResult struct {
	s, p *big.Int
	err  error
}

func pickS(rg *rand.Rand, mul, add, randLmt, preP *big.Int) (*big.Int, *big.Int, bool, error) {
	// choose k' in [0, randLmt)
	k, err := crand.Int(crand.Reader, randLmt)
	if err != nil {
		return nil, nil, false, err
	}
	// construct k, k = k' * mul + add
	k.Mul(k, mul)
	k.Add(k, add)
	return determineSAndP(rg, k, preP)
}

func findLargeSRoutine(ctx context.Context, randBitLen int, preP *big.Int, resChan chan<- findSResult) {
	rg := rand.New(rand.NewSource(time.Now().UnixNano()))
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s, p, ok, err := pickLargeS(rg, randBitLen, preP)
			if err != nil {
				select {
				case resChan <- findSResult{err: err}:
					return
				default:
					return
				}
			}
			if !ok {
				continue
			}
			ctx.Done()
			select {
			case resChan <- findSResult{s: s, p: p}:
				return
			default:
				return
			}
		}
	}
}

func pickLargeS(rg *rand.Rand, randBitLen int, preP *big.Int) (*big.Int, *big.Int, bool, error) {
	k, err := crand.Prime(crand.Reader, randBitLen)
	if err != nil {
		return nil, nil, false, err
	}
	return determineSAndP(rg, k, preP)
}

func determineSAndP(rg *rand.Rand, k, preP *big.Int) (*big.Int, *big.Int, bool, error) {
	// p = {Product of primes} * n * k - 1 = preP * k - 1
	p := iPool.Get().(*big.Int)
	defer iPool.Put(p)
	p.Mul(preP, k)
	p.Sub(p, big1)

	// choose u from [1, p - 1]
	// here we can pick u in [0, p)
	// if u is 0, then the accepting condition will not pass
	u := iPool.Get().(*big.Int)
	defer iPool.Put(u)
	u.Rand(rg, p)

	// test if s^2 = -1 (mod p)
	// if so, continue to the next step, otherwise, repeat this step
	pMinus1 := iPool.Get().(*big.Int)
	defer iPool.Put(pMinus1)
	pMinus1.Sub(p, big1)
	powU := iPool.Get().(*big.Int)
	defer iPool.Put(powU)
	powU.Rsh(pMinus1, 1)

	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	if opt.Exp(u, powU, p).Cmp(pMinus1) != 0 {
		return nil, nil, false, nil
	}

	// compute s = u^((p - 1) / 4) mod p
	powU.Rsh(powU, 1)
	s := new(big.Int).Exp(u, powU, p)
	return s, new(big.Int).Set(p), true, nil
}

func gaussianIntGCD(s, p *big.Int) *comp.GaussianInt {
	// compute A + Bi := gcd(s + i, p)
	// Gaussian integer: s + i
	gaussianInt := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gaussianInt)
	gaussianInt.Update(s, big1)
	// Gaussian integer: p
	gaussianP := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gaussianP)
	gaussianP.Update(p, big0)
	// compute gcd(s + i, p)
	gcd := new(comp.GaussianInt)
	gcd.GCD(gaussianInt, gaussianP)
	return gcd
}

func denouement(n *big.Int, gcd *comp.GaussianInt) (*comp.HurwitzInt, error) {
	// compute gcrd(A + Bi + j, n), normalized to have integer component
	// Hurwitz integer: A + Bi + j
	hurwitzInt := hiPool.Get().(*comp.HurwitzInt)
	defer hiPool.Put(hurwitzInt)
	hurwitzInt.Update(gcd.R, gcd.I, big1, big0, false)
	// Hurwitz integer: n
	hurwitzN := hiPool.Get().(*comp.HurwitzInt)
	defer hiPool.Put(hurwitzN)
	hurwitzN.Update(n, big0, big0, big0, false)
	gcrd := new(comp.HurwitzInt).GCRD(hurwitzInt, hurwitzN)

	return gcrd, nil
}

// Verify checks if the four-square sum is equal to the original integer
// i.e. target = w1^2 + w2^2 + w3^2 + w4^2
func Verify(target *big.Int, fs FourNum) bool {
	sum := iPool.Get().(*big.Int)
	defer iPool.Put(sum)
	sum.Set(big0)
	for i := 0; i < 4; i++ {
		sum.Add(sum, new(big.Int).Mul(fs[i], fs[i]))
	}
	return sum.Cmp(target) == 0
}
