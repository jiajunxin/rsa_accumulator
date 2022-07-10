package complex

import (
	"math/big"
)

// HurwitzInt implements Hurwitz quaternion (or Hurwitz integer) a + bi + cj + dk
// The set of all Hurwitz quaternion is H = {a + bi + cj + dk | a, b, c, d are all integers or all half-integers}
// A mixture of integers and half-integers is excluded
// In the struct each scalar is twice the original scalar so that all the scalars can be stored using
// big integers for computation efficiency
type HurwitzInt struct {
	dblR *big.Int // r part doubled
	dblI *big.Int // i part doubled
	dblJ *big.Int // j part doubled
	dblK *big.Int // k part doubled
}

// Init initialize a Hurwitz integer
func (h *HurwitzInt) Init() *HurwitzInt {
	h.dblR = new(big.Int)
	h.dblI = new(big.Int)
	h.dblJ = new(big.Int)
	h.dblK = new(big.Int)
	return h
}

// String returns the string representation of the integral quaternion
func (h *HurwitzInt) String() string {
	if h.dblR.Sign() == 0 && h.dblI.Sign() == 0 && h.dblJ.Sign() == 0 && h.dblK.Sign() == 0 {
		return "0"
	}
	res := ""
	if h.dblR.Cmp(big1) == 0 {
		res += "0.5"
	} else if h.dblR.Cmp(bigNeg1) == 0 {
		res += "-0.5"
	} else if h.dblR.Sign() != 0 {
		res += new(big.Int).Rsh(h.dblR, 1).String()
		if h.dblR.Bit(0) == 1 {
			res += ".5"
		}
	}
	if h.dblI.Cmp(big1) == 0 {
		res += "0.5i"
	} else if h.dblI.Cmp(bigNeg1) == 0 {
		res += "-0.5i"
	} else if h.dblI.Sign() != 0 {
		if h.dblR.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Rsh(h.dblI, 1).String()
		if h.dblI.Bit(0) == 1 {
			res += ".5"
		}
		res += "i"
	}
	if h.dblJ.Cmp(big1) == 0 {
		res += "0.5j"
	} else if h.dblJ.Cmp(bigNeg1) == 0 {
		res += "-0.5j"
	} else if h.dblJ.Sign() != 0 {
		if h.dblJ.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Rsh(h.dblJ, 1).String()
		if h.dblJ.Bit(0) == 1 {
			res += ".5"
		}
		res += "j"
	}
	if h.dblK.Cmp(big1) == 0 {
		res += "0.5j"
	} else if h.dblK.Cmp(bigNeg1) == 0 {
		res += "-0.5j"
	} else if h.dblK.Sign() != 0 {
		if h.dblK.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Rsh(h.dblK, 1).String()
		if h.dblK.Bit(0) == 1 {
			res += ".5"
		}
		res += "k"
	}
	return res
}

// NewHurwitzInt declares a new integral quaternion with the real, i, j, and k parts
// If isDouble is true, the arguments r, i, j, k are twice the original scalars
func NewHurwitzInt(r, i, j, k *big.Int, isDouble bool) *HurwitzInt {
	if isDouble {
		return &HurwitzInt{
			dblR: new(big.Int).Set(r),
			dblI: new(big.Int).Set(i),
			dblJ: new(big.Int).Set(j),
			dblK: new(big.Int).Set(k),
		}
	}
	return &HurwitzInt{
		dblR: new(big.Int).Lsh(r, 1),
		dblI: new(big.Int).Lsh(i, 1),
		dblJ: new(big.Int).Lsh(j, 1),
		dblK: new(big.Int).Lsh(k, 1),
	}
}

// Set sets the Hurwitz integer to the given Hurwitz integer
func (h *HurwitzInt) Set(a *HurwitzInt) *HurwitzInt {
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Set(a.dblR)
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Set(a.dblI)
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Set(a.dblJ)
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Set(a.dblK)
	return h
}

// SetFloat set scalars of a Hurwitz integer by big float variables
//func (h *HurwitzInt) SetFloat(r, i, j, k *big.Float) *HurwitzInt {
//	panic("not implemented")
//}

// Val reveals value of a Hurwitz integer
func (h *HurwitzInt) Val() (r, i, j, k *big.Float) {
	r = new(big.Float).SetInt(h.dblR)
	r.Quo(r, big2f)
	i = new(big.Float).SetInt(h.dblI)
	i.Quo(i, big2f)
	j = new(big.Float).SetInt(h.dblJ)
	j.Quo(j, big2f)
	k = new(big.Float).SetInt(h.dblK)
	k.Quo(k, big2f)
	return
}

// ValInt reveals value of a Hurwitz integer in integer
func (h *HurwitzInt) ValInt() (r, i, j, k *big.Int) {
	rF, iF, jF, kF := h.Val()
	r = roundFloat(rF)
	i = roundFloat(iF)
	j = roundFloat(jF)
	k = roundFloat(kF)
	return
}

// Update updates the integral quaternion with the given real, i, j, and k parts
func (h *HurwitzInt) Update(r, i, j, k *big.Int, isDouble bool) *HurwitzInt {
	if isDouble {
		h.dblR = r
		h.dblI = i
		h.dblJ = j
		h.dblK = k
	} else {
		if h.dblR == nil {
			h.dblR = new(big.Int)
		}
		h.dblR.Lsh(r, 1)
		if h.dblI == nil {
			h.dblI = new(big.Int)
		}
		h.dblI.Lsh(i, 1)
		if h.dblJ == nil {
			h.dblJ = new(big.Int)
		}
		h.dblJ.Lsh(j, 1)
		if h.dblK == nil {
			h.dblK = new(big.Int)
		}
		h.dblK.Lsh(k, 1)
	}
	return h
}

// Zero sets the Hurwitz integer to zero
func (h *HurwitzInt) Zero() *HurwitzInt {
	h.dblR = big.NewInt(0)
	h.dblI = big.NewInt(0)
	h.dblJ = big.NewInt(0)
	h.dblK = big.NewInt(0)
	return h
}

// Add adds two integral quaternions
func (h *HurwitzInt) Add(a, b *HurwitzInt) *HurwitzInt {
	//h.dblR = new(big.Int).Add(a.dblR, b.dblR)
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Add(a.dblR, b.dblR)
	h.dblI = new(big.Int).Add(a.dblI, b.dblI)
	h.dblJ = new(big.Int).Add(a.dblJ, b.dblJ)
	h.dblK = new(big.Int).Add(a.dblK, b.dblK)
	return h
}

// Sub subtracts two integral quaternions
func (h *HurwitzInt) Sub(a, b *HurwitzInt) *HurwitzInt {
	h.dblR = new(big.Int).Sub(a.dblR, b.dblR)
	h.dblI = new(big.Int).Sub(a.dblI, b.dblI)
	h.dblJ = new(big.Int).Sub(a.dblJ, b.dblJ)
	h.dblK = new(big.Int).Sub(a.dblK, b.dblK)
	return h
}

// Conj obtains the conjugate of the original integral quaternion
func (h *HurwitzInt) Conj(origin *HurwitzInt) *HurwitzInt {
	h.dblR = origin.dblR
	h.dblI = new(big.Int).Neg(origin.dblI)
	h.dblJ = new(big.Int).Neg(origin.dblJ)
	h.dblK = new(big.Int).Neg(origin.dblK)
	return h
}

// Norm obtains the norm of the integral quaternion
func (h *HurwitzInt) Norm() *big.Int {
	norm := new(big.Int).Mul(h.dblR, h.dblR)
	norm.Add(norm, new(big.Int).Mul(h.dblI, h.dblI))
	norm.Add(norm, new(big.Int).Mul(h.dblJ, h.dblJ))
	norm.Add(norm, new(big.Int).Mul(h.dblK, h.dblK))
	norm.Div(norm, big4)
	return norm
}

// Copy copies the integral quaternion
func (h *HurwitzInt) Copy() *HurwitzInt {
	return NewHurwitzInt(h.dblR, h.dblI, h.dblJ, h.dblK, true)
}

// Prod returns the Hamilton product of two integral quaternions
// the product (a1 + b1j + c1k + d1)(a2 + b2j + c2k + d2) is determined by the products of the
// basis elements and the distributive law
func (h *HurwitzInt) Prod(a, b *HurwitzInt) *HurwitzInt {
	a = a.Copy()
	b = b.Copy()
	opt := new(big.Int)
	// 1 part
	if h.dblR == nil {
		h.dblR = new(big.Int)
	}
	h.dblR.Mul(a.dblR, b.dblR)
	h.dblR.Sub(h.dblR, opt.Mul(a.dblI, b.dblI))
	h.dblR.Sub(h.dblR, opt.Mul(a.dblJ, b.dblJ))
	h.dblR.Sub(h.dblR, opt.Mul(a.dblK, b.dblK))
	h.dblR.Div(h.dblR, big2)

	// i part
	if h.dblI == nil {
		h.dblI = new(big.Int)
	}
	h.dblI.Mul(a.dblR, b.dblI)
	h.dblI.Add(h.dblI, opt.Mul(a.dblI, b.dblR))
	h.dblI.Add(h.dblI, opt.Mul(a.dblJ, b.dblK))
	h.dblI.Sub(h.dblI, opt.Mul(a.dblK, b.dblJ))
	h.dblI.Div(h.dblI, big2)

	// j part
	if h.dblJ == nil {
		h.dblJ = new(big.Int)
	}
	h.dblJ.Mul(a.dblR, b.dblJ)
	h.dblJ.Sub(h.dblJ, opt.Mul(a.dblI, b.dblK))
	h.dblJ.Add(h.dblJ, opt.Mul(a.dblJ, b.dblR))
	h.dblJ.Add(h.dblJ, opt.Mul(a.dblK, b.dblI))
	h.dblJ.Div(h.dblJ, big2)

	// k part
	if h.dblK == nil {
		h.dblK = new(big.Int)
	}
	h.dblK.Mul(a.dblR, b.dblK)
	h.dblK.Add(h.dblK, opt.Mul(a.dblI, b.dblJ))
	h.dblK.Sub(h.dblK, opt.Mul(a.dblJ, b.dblI))
	h.dblK.Add(h.dblK, opt.Mul(a.dblK, b.dblR))
	h.dblK.Div(h.dblK, big2)

	return h
}

// Div performs Euclidean division of two Hurwitz integers, i.e. a/b
// the remainder is stored in the Hurwitz integer that calls the method
// the quotient is returned as a new Hurwitz integer
func (h *HurwitzInt) Div(a, b *HurwitzInt) *HurwitzInt {
	a = a.Copy()
	b = b.Copy()
	bConj := new(HurwitzInt).Conj(b)
	numerator := new(HurwitzInt).Prod(a, bConj)
	denominator := new(HurwitzInt).Prod(b, bConj)
	deInt := denominator.dblR
	deFloat := new(big.Float).SetInt(deInt)

	nuRFloat := new(big.Float).SetInt(numerator.dblR)
	nuIFloat := new(big.Float).SetInt(numerator.dblI)
	nuJFloat := new(big.Float).SetInt(numerator.dblJ)
	nuKFloat := new(big.Float).SetInt(numerator.dblK)
	rScalar := new(big.Float).Quo(nuRFloat, deFloat)
	iScalar := new(big.Float).Quo(nuIFloat, deFloat)
	jScalar := new(big.Float).Quo(nuJFloat, deFloat)
	kScalar := new(big.Float).Quo(nuKFloat, deFloat)

	rsInt := roundFloat(rScalar)
	isInt := roundFloat(iScalar)
	jsInt := roundFloat(jScalar)
	ksInt := roundFloat(kScalar)
	quotient := NewHurwitzInt(rsInt, isInt, jsInt, ksInt, false)
	h.Sub(a, new(HurwitzInt).Prod(quotient, b))
	return quotient
}

// GCRD calculates the greatest common right-divisor of two Hurwitz integers using Euclidean algorithm
// The GCD is unique only up to multiplication by a unit (multiplication on the left in the case
// of a GCRD, and on the right in the case of a GCLD)
// the result is stored in the Hurwitz integer that calls the method and returned
func (h *HurwitzInt) GCRD(a, b *HurwitzInt) *HurwitzInt {
	a = a.Copy()
	b = b.Copy()
	if a.CmpNorm(b) < 0 {
		a, b = b, a
	}
	remainder := new(HurwitzInt)
	for {
		remainder.Div(a, b)
		if remainder.IsZero() {
			h.Set(b)
			return b
		}
		a.Set(b)
		b.Set(remainder)
	}
}

// IsZero returns true if the Hurwitz integer is zero
func (h *HurwitzInt) IsZero() bool {
	return h.dblR.Sign() == 0 &&
		h.dblI.Sign() == 0 &&
		h.dblJ.Sign() == 0 &&
		h.dblK.Sign() == 0
}

// CmpNorm compares the norm of two Hurwitz integers
func (h *HurwitzInt) CmpNorm(a *HurwitzInt) int {
	return h.Norm().Cmp(a.Norm())
}
