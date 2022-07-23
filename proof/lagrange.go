package proof

import (
	"context"
	crand "crypto/rand"
	"errors"
	comp "github.com/rsa_accumulator/complex"
	"math"
	"math/big"
	"runtime"
	"sync"
)

const squareNum = 4

var (
	// precomputed Hurwitz GCRDs for small integers
	precomputedHurwitzGCRDs = [9]*comp.HurwitzInt{
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
	}
	numCPU     = runtime.NumCPU()
	ps         = newPrimeStore(1792)
	ss         = newSquareStore(0)
	pickSCache = sync.Map{}
)

// FourSquare is the LagrangeFourSquareLipmaa representation of a positive integer
type FourSquare [squareNum]*big.Int

// NewFourSquare creates a new FourSquare
func NewFourSquare(w1 *big.Int, w2 *big.Int, w3 *big.Int, w4 *big.Int) FourSquare {
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
	return FourSquare{w1, w2, w3, w4}
}

// Mul multiplies all the square numbers by n
func (f *FourSquare) Mul(n *big.Int) {
	for i := 0; i < squareNum; i++ {
		f[i].Mul(f[i], n)
	}
}

// Div divides all the square numbers by n
func (f *FourSquare) Div(n *big.Int) {
	for i := 0; i < squareNum; i++ {
		f[i].Div(f[i], n)
	}
}

// String stringnifies the FourSquare object
func (f *FourSquare) String() string {
	res := "{"
	for i := 0; i < squareNum-1; i++ {
		res += f[i].String()
		res += ", "
	}
	res += f[squareNum-1].String()
	res += "}"
	return res
}

// RPCommit generates a range proof commitment for a given integer
func (f *FourSquare) RPCommit(pp *PublicParameters, coins rpRandCoins) (cList [squareNum]*big.Int) {
	for i := 0; i < squareNum; i++ {
		cList[i] = new(big.Int).Exp(pp.G, f[i], pp.N)
		cList[i].Mul(cList[i], new(big.Int).Exp(pp.H, coins[i], pp.N))
	}
	return
}

// LagrangeFourSquares calculates the Lagrange four squares representation of a positive integer
// Paper: Finding the Four Squares in Lagrange’s Theorem
// Link: http://pollack.uga.edu/finding4squares.pdf (page 6)
// The input should be an odd positive integer no less than 9
func LagrangeFourSquares(n *big.Int) (FourSquare, error) {
	if n.Sign() == 0 {
		res := NewFourSquare(precomputedHurwitzGCRDs[0].ValInt())
		return res, nil
	}
	nc, e := divideN(n)
	var hurwitzGCRD *comp.HurwitzInt

	if nc.Cmp(big8) <= 0 {
		hurwitzGCRD = precomputedHurwitzGCRDs[nc.Int64()]
	} else {
		primeProd, err := preCompute(nc)
		if err != nil {
			return FourSquare{}, err
		}
		s, p, err := randTrails(nc, primeProd)
		if err != nil {
			return FourSquare{}, err
		}
		hurwitzGCRD, err = denouement(nc, s, p)
		if err != nil {
			return FourSquare{}, err
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
	fs := NewFourSquare(w1, w2, w3, w4)
	return fs, nil
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
	gaussian1PlusI := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gaussian1PlusI)
	gaussian1PlusI.Update(big1, big1)

	gaussianProd := comp.NewGaussianInt(big1, big0)
	for e > 0 {
		gaussianProd.Prod(gaussianProd, gaussian1PlusI)
		e--
	}
	return gaussianProd
}

// preCompute determine the primes not exceeding log n and compute their product
// the function only handles positive integers larger than 8
func preCompute(n *big.Int) (*big.Int, error) {
	if n.Cmp(big8) <= 0 {
		return nil, errors.New("n should be larger than 8")
	}
	logN := log2(n)
	if logN <= ps.max {
		//for idx := len(ps.l) - 1; idx >= 0; idx-- {
		//	psl := ps.l[idx]
		//	if psl < logN {
		//		return ps.m[psl], nil
		//	}
		//}
		//return nil, errors.New("precomputed primes not found")
		prod, err := ps.findPrimeProd(logN)
		if err != nil {
			return nil, err
		}
		return prod, nil
	}
	prod := iPool.Get().(*big.Int)
	defer iPool.Put(prod)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	prod.Set(ps.m[ps.max])
	for idx := ps.max + 2; idx < logN; idx += 2 {
		ps.checkAddPrime(idx, prod, opt)
	}
	return new(big.Int).Set(prod), nil
}

type primeStore struct {
	l   []int
	m   map[int]*big.Int
	max int
}

func newPrimeStore(lmt int) *primeStore {
	ps := &primeStore{
		l:   []int{2, 3, 5, 7},
		m:   make(map[int]*big.Int),
		max: 7,
	}
	ps.m[2] = big.NewInt(2)
	ps.m[3] = big.NewInt(6)
	ps.m[5] = big.NewInt(30)
	ps.m[7] = big.NewInt(210)

	prod := iPool.Get().(*big.Int)
	defer iPool.Put(prod)
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	prod.SetInt64(210)
	for idx := 9; idx <= lmt; idx += 2 {
		ps.checkAddPrime(idx, prod, opt)
	}
	return ps
}

func (p *primeStore) checkAddPrime(n int, prod, opt *big.Int) {
	isPrime := true
	for _, prime := range p.l {
		if n%prime == 0 && n != prime {
			isPrime = false
			break
		}
	}
	if !isPrime {
		return
	}
	p.l = append(p.l, n)
	opt.SetInt64(int64(n))
	prod.Mul(prod, opt)
	p.m[n] = new(big.Int).Set(prod)
	p.max = n
}

// findPrimeProd finds the product of primes less than log n using binary search
func (p *primeStore) findPrimeProd(logN int) (*big.Int, error) {
	var (
		l int
		r = len(p.l) - 1
	)
	for l <= r {
		mid := (l-r)/2 + r
		pll := p.l[mid]
		if mid == len(p.l)-1 {
			return p.m[pll], nil
		}
		plr := p.l[mid+1]
		if pll < logN && plr >= logN {
			return p.m[pll], nil
		}
		if pll >= logN {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return nil, errors.New("precomputed primes not found")
}

func ResetPrimeStore() {
	ps = newPrimeStore(0)
}

func randTrails(n, primeProd *big.Int) (*big.Int, *big.Int, error) {
	// use goroutines to choose a random number between [0, n^5 / 2 / numCPU]
	// then construct k based on the random number
	// and check the validity of the trails
	randLmt := new(big.Int).Exp(n, big5, nil)
	randLmt.Rsh(randLmt, 1)
	randLmt.Div(randLmt, big.NewInt(int64(numCPU)))
	randLmt.Add(randLmt, big1)
	// p = M * n * k - 1, pre-p = M * n
	preP := new(big.Int).Mul(primeProd, n)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var (
		mul  = big.NewInt(int64(2 * numCPU)) // 2 * numCPU
		adds []*big.Int
	)
	for i := 0; i <= numCPU; i++ {
		adds = append(adds, big.NewInt(int64(2*i+1))) // 2i+1
	}
	resChan := make(chan findSResult)
	for _, add := range adds {
		go findSRoutine(ctx, add, mul, randLmt, preP, resChan)
	}
	res := <-resChan
	return res.s, res.p, res.err
}

func findSRoutine(ctx context.Context, mul, add, randLmt, preP *big.Int, resChan chan<- findSResult) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			s, p, ok, err := pickS(mul, add, randLmt, preP)
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
			case resChan <- findSResult{
				s: s, p: p,
			}:
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

func pickS(mul, add, randLmt, preP *big.Int) (*big.Int, *big.Int, bool, error) {
	var err error
	// choose k' in [0, randLmt)
	k := iPool.Get().(*big.Int)
	defer iPool.Put(k)
	if k, err = crand.Int(crand.Reader, randLmt); err != nil {
		return nil, nil, false, err
	}
	// construct k, k = k' * mul + add
	k.Mul(k, mul)
	k.Add(k, add)

	// p = {Product of primes} * n * k - 1 = preP * k - 1
	p := iPool.Get().(*big.Int)
	defer iPool.Put(p)
	p.Mul(preP, k)
	p.Sub(p, big1)
	pMinus1 := iPool.Get().(*big.Int)
	defer iPool.Put(pMinus1)
	pMinus1.Sub(p, big1)

	// choose u from [1, p - 1]
	u := iPool.Get().(*big.Int)
	defer iPool.Put(u)
	if u, err = crand.Int(crand.Reader, pMinus1); err != nil {
		return nil, nil, false, err
	}
	u.Add(u, big1)

	// compute s = u^((p - 1) / 4) mod p
	powU := iPool.Get().(*big.Int)
	defer iPool.Put(powU)
	powU.Rsh(pMinus1, 2)
	s := iPool.Get().(*big.Int)
	defer iPool.Put(s)
	s.Exp(u, powU, p)

	// test if s^2 = -1 (mod p)
	// if so, continue to the next step, otherwise, repeat this step
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	if opt.Exp(s, big2, p).Cmp(pMinus1) != 0 {
		return nil, nil, false, nil
	}
	return new(big.Int).Set(s), new(big.Int).Set(p), true, nil
}

func denouement(n, s, p *big.Int) (*comp.HurwitzInt, error) {
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
	gcd := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gcd)
	gcd.GCD(gaussianInt, gaussianP)
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

// UnconditionalLagrangeFourSquares calculates the Lagrange four squares for a given non-positive integer
// the method doesn't rely on the Extended Riemann Hypothesis (ERH
func UnconditionalLagrangeFourSquares(n *big.Int) (FourSquare, error) {
	if n.Sign() == 0 {
		res := NewFourSquare(precomputedHurwitzGCRDs[0].ValInt())
		return res, nil
	}
	nc, e := divideN(n)
	var hurwitzGCRD *comp.HurwitzInt

	if nc.Cmp(big8) <= 0 {
		hurwitzGCRD = precomputedHurwitzGCRDs[nc.Int64()]
	} else {
		x, y, p, r1, s, primes, err := initTrail(nc)
		if err != nil {
			return FourSquare{}, err
		}
		// compute u, v
		u, v, err := computeUV(r1, nc, primes)
		if err != nil {
			return FourSquare{}, err
		}
		var up, vp *big.Int
		// compute U -> up, V -> vp
		if s == nil {
			up = big.NewInt(1)
			vp = big.NewInt(0)
		} else {
			gcd := giPool.Get().(*comp.GaussianInt)
			defer giPool.Put(gcd)
			gopt1 := giPool.Get().(*comp.GaussianInt)
			defer giPool.Put(gopt1)
			gopt2 := giPool.Get().(*comp.GaussianInt)
			defer giPool.Put(gopt2)
			gopt1.Update(p, big0)
			gopt2.Update(s, big1)
			gcd.GCD(gopt1, gopt2)
			up = gcd.R
			vp = gcd.I
		}
		uvi := comp.NewGaussianInt(u, v)
		uPvPI := comp.NewGaussianInt(up, vp)
		zwi := new(comp.GaussianInt).Prod(uvi, uPvPI)
		hopt1 := hiPool.Get().(*comp.HurwitzInt)
		defer hiPool.Put(hopt1)
		hopt1.Update(n, big0, big0, big0, false)
		hopt2 := hiPool.Get().(*comp.HurwitzInt)
		defer hiPool.Put(hopt2)
		hopt2.Update(x, y, zwi.R, zwi.I, false)
		hurwitzGCRD = new(comp.HurwitzInt).GCRD(hopt1, hopt2)
	}

	// if x'^2 + Y'^2 + Z'^2 + W'^2 = n'
	// then x^2 + Y^2 + Z^2 + W^2 = n for x, Y, Z, W defined by
	// (1 + i)^e * (x' + Y'i + Z'j + W'k) = (x + Yi + Zj + Wk)
	// Gaussian integer: 1 + i
	gi := gaussian1PlusIPow(e)
	hurwitzProd := comp.NewHurwitzInt(gi.R, gi.I, big0, big0, false)
	hurwitzProd.Prod(hurwitzProd, hurwitzGCRD)
	w1, w2, w3, w4 := hurwitzProd.ValInt()
	fs := NewFourSquare(w1, w2, w3, w4)
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
	for i := 0; i < numCPU; i++ {
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
	uvi := comp.NewGaussianInt(big1, big0)
	ss, err := computeSquares(n)
	if err != nil {
		return nil, nil, err
	}
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	mod := iPool.Get().(*big.Int)
	defer iPool.Put(mod)
	gopt := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gopt)
	for _, prime := range primes {
		x, y := ss.findXY(prime)
		gopt.Update(x, y)
		// determine vl
		opt.Set(prime)
		for mod.Mod(r1, opt).Cmp(big0) == 0 {
			uvi.Prod(uvi, gopt)
			opt.Mul(opt, prime)
		}
	}
	u = uvi.R
	v = uvi.I
	return
}

func computeSquares(n *big.Int) (*squareStore, error) {
	lmt := log2(n)
	lmt = int(math.Sqrt(float64(lmt)))
	for i := ss.max; i <= lmt; i++ {
		bigI := big.NewInt(int64(i))
		sq := new(big.Int).Mul(bigI, bigI)
		ss.add(bigI, sq)
		ss.max = i
	}
	return ss, nil
}

func CacheSquareNums(bitLen int) {
	lmt := int(math.Sqrt(float64(bitLen)))
	ss = newSquareStore(lmt)
}

type squareStore struct {
	sm  map[string]*big.Int
	sl  []*big.Int
	max int
}

func newSquareStore(max int) *squareStore {
	ss := &squareStore{
		sm: make(map[string]*big.Int),
	}
	if max > 0 {
		ss.sl = make([]*big.Int, max)
		for i := 1; i <= max; i++ {
			bigI := big.NewInt(int64(i))
			sq := new(big.Int).Mul(bigI, bigI)
			ss.add(bigI, sq)
		}
		ss.max = max
	}
	return ss
}

func (s *squareStore) add(n, nsq *big.Int) {
	if _, ok := s.sm[nsq.String()]; !ok {
		s.sm[nsq.String()] = n
		s.sl = append(s.sl, nsq)
	}
}

func (s *squareStore) findXY(n *big.Int) (x, y *big.Int) {
	opt := iPool.Get().(*big.Int)
	defer iPool.Put(opt)
	for _, sq := range s.sl {
		if sq.Cmp(n) == 1 {
			break
		}
		opt.Sub(n, sq)
		if resY, ok := s.sm[opt.String()]; ok {
			x = new(big.Int).Set(s.sm[sq.String()])
			y = new(big.Int).Set(resY)
			return
		}
	}
	return
}

// Verify checks if the four-square sum is equal to the original integer
// i.e. target = w1^2 + w2^2 + w3^2 + w4^2
func Verify(target *big.Int, fs [squareNum]*big.Int) bool {
	sum := iPool.Get().(*big.Int)
	defer iPool.Put(sum)
	sum.Set(big0)
	for i := 0; i < squareNum; i++ {
		sum.Add(sum, new(big.Int).Mul(fs[i], fs[i]))
	}
	return sum.Cmp(target) == 0
}
