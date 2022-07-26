package proof

import (
	"math"
	"math/big"
	"sync"

	comp "github.com/rsa_accumulator/complex"
)

var (
	// sync pool for big integers, lease GC and improve performance
	iPool = sync.Pool{
		New: func() interface{} { return new(big.Int) },
	}
	// sync pool for Gaussian integers
	giPool = sync.Pool{
		New: func() interface{} { return new(comp.GaussianInt) },
	}
	// sync pool for Hurwitz integers
	hiPool = sync.Pool{
		New: func() interface{} { return new(comp.HurwitzInt) },
	}
	// cache for precomputed prime numbers and prime number products
	pCache = newPrimeCache(32)
	// cache for computed Gaussian integers
	giCache = make(map[int]*comp.GaussianInt)
	// cache for precomputed square numbers
	sCache = newSquareCache(0)
)

// ResetPrimeCache resets the prime cache
func ResetPrimeCache() {
	pCache = newPrimeCache(0)
}

type primeCache struct {
	l   []int            // list of prime numbers
	m   map[int]*big.Int // map of prime numbers and the products
	max int              // the largest prime number included
}

// findPrimeProd finds the product of primes less than log n using binary search
func (p *primeCache) findPrimeProd(logN int) *big.Int {
	var (
		l int
		r = len(p.l) - 1
	)
	for l <= r {
		mid := (l-r)/2 + r
		pll := p.l[mid]
		if mid == len(p.l)-1 {
			return p.m[pll]
		}
		plr := p.l[mid+1]
		if pll < logN && plr >= logN {
			return p.m[pll]
		}
		if pll >= logN {
			r = mid - 1
		} else {
			l = mid + 1
		}
	}
	return big.NewInt(2)
}

func newPrimeCache(lmt int) *primeCache {
	ps := &primeCache{
		l:   []int{1, 2, 3, 5, 7},
		m:   make(map[int]*big.Int),
		max: 7,
	}
	ps.m[1] = big.NewInt(1)
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

func (p *primeCache) checkAddPrime(n int, prod, opt *big.Int) {
	isPrime := true
	sqrtN := int(math.Sqrt(float64(n)))
	for idx := 1; idx < len(p.l); idx++ {
		if sqrtN < p.l[idx] {
			break
		}
		if n%p.l[idx] == 0 && n != p.l[idx] {
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

// ResetGaussianIntCache resets the Gaussian integer cache
func ResetGaussianIntCache() {
	giCache = make(map[int]*comp.GaussianInt)
}

// CacheGaussianInt caches (1+i)^n, n <= e
func CacheGaussianInt(e int) {
	giCache = make(map[int]*comp.GaussianInt)
	gaussianProd := giPool.Get().(*comp.GaussianInt)
	defer giPool.Put(gaussianProd)
	gaussianProd.Update(big1, big0)
	for i := 0; i <= e; i++ {
		giCache[i] = gaussianProd.Copy()
		gaussianProd.Prod(gaussianProd, gaussianProd)
	}
}

// CacheSquareNums caches the square numbers of x, x <= sqrt(bit length)
func CacheSquareNums(bitLen int) {
	lmt := int(math.Sqrt(float64(bitLen)))
	sCache = newSquareCache(lmt)
}

// ResetSquareCache resets the cache for square numbers
func ResetSquareCache() {
	sCache = newSquareCache(0)
}

type squareCache struct {
	sm  map[string]*big.Int
	sl  []*big.Int
	max int
}

func newSquareCache(max int) *squareCache {
	ss := &squareCache{
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

func (s *squareCache) add(n, nsq *big.Int) {
	if _, ok := s.sm[nsq.String()]; !ok {
		s.sm[nsq.String()] = n
		s.sl = append(s.sl, nsq)
	}
}

func (s *squareCache) findXY(n *big.Int) (x, y *big.Int) {
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
