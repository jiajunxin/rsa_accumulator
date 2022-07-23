package proof

import (
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
)
