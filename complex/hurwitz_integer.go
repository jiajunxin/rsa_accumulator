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

// String returns the string representation of the integral quaternion
func (h *HurwitzInt) String() string {
	if h.dblR.Sign() == 0 && h.dblI.Sign() == 0 && h.dblJ.Sign() == 0 && h.dblK.Sign() == 0 {
		return "0"
	}
	res := ""
	if h.dblR.Sign() != 0 {
		res += new(big.Int).Div(h.dblR, big2).String()
		if h.dblR.Bit(0) == 1 {
			res += ".5"
		}
	}
	if h.dblI.Sign() != 0 {
		if h.dblR.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.dblI, big2).String()
		if h.dblI.Bit(0) == 1 {
			res += ".5"
		}
		res += "i"
	}
	if h.dblJ.Sign() != 0 {
		if h.dblJ.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.dblJ, big2).String()
		if h.dblJ.Bit(0) == 1 {
			res += ".5"
		}
		res += "j"
	}
	if h.dblK.Sign() != 0 {
		if h.dblK.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.dblK, big2).String()
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
		dblR: new(big.Int).Mul(r, big2),
		dblI: new(big.Int).Mul(i, big2),
		dblJ: new(big.Int).Mul(j, big2),
		dblK: new(big.Int).Mul(k, big2),
	}
}

// Set sets the Hurwitz integer to the given Hurwitz integer
func (h *HurwitzInt) Set(a *HurwitzInt) *HurwitzInt {
	h.dblR = a.dblR
	h.dblI = a.dblI
	h.dblJ = a.dblJ
	h.dblK = a.dblK
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
		h.dblR = new(big.Int).Mul(r, big2)
		h.dblI = new(big.Int).Mul(i, big2)
		h.dblJ = new(big.Int).Mul(j, big2)
		h.dblK = new(big.Int).Mul(k, big2)
	}
	return h
}

// Add adds two integral quaternions
func (h *HurwitzInt) Add(a, b *HurwitzInt) *HurwitzInt {
	h.dblR = new(big.Int).Add(a.dblR, b.dblR)
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
	h.dblR = new(big.Int).Mul(a.dblR, b.dblR)
	h.dblR.Sub(h.dblR, opt.Mul(a.dblI, b.dblI))
	h.dblR.Sub(h.dblR, opt.Mul(a.dblJ, b.dblJ))
	h.dblR.Sub(h.dblR, opt.Mul(a.dblK, b.dblK))
	h.dblR.Div(h.dblR, big2)

	// i part
	h.dblI = new(big.Int).Mul(a.dblR, b.dblI)
	h.dblI.Add(h.dblI, opt.Mul(a.dblI, b.dblR))
	h.dblI.Add(h.dblI, opt.Mul(a.dblJ, b.dblK))
	h.dblI.Sub(h.dblI, opt.Mul(a.dblK, b.dblJ))
	h.dblI.Div(h.dblI, big2)

	// j part
	h.dblJ = new(big.Int).Mul(a.dblR, b.dblJ)
	h.dblJ.Sub(h.dblJ, opt.Mul(a.dblI, b.dblK))
	h.dblJ.Add(h.dblJ, opt.Mul(a.dblJ, b.dblR))
	h.dblJ.Add(h.dblJ, opt.Mul(a.dblK, b.dblI))
	h.dblJ.Div(h.dblJ, big2)

	// k part
	h.dblK = new(big.Int).Mul(a.dblR, b.dblK)
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
	if a.Norm().Cmp(b.Norm()) < 0 {
		a, b = b, a
	}
	remainder := new(HurwitzInt)
	for {
		remainder.Div(a, b)
		if remainder.IsZero() {
			h.Update(b.dblR, b.dblI, b.dblJ, b.dblK, true)
			return b
		}
		a.Update(b.dblR, b.dblI, b.dblJ, b.dblK, true)
		b.Update(remainder.dblR, remainder.dblI, remainder.dblJ, remainder.dblK, true)
	}
}

// IsZero returns true if the Hurwitz integer is zero
func (h *HurwitzInt) IsZero() bool {
	return h.dblR.Sign() == 0 && h.dblI.Sign() == 0 && h.dblJ.Sign() == 0 && h.dblK.Sign() == 0
}

// CmpNorm compares the norm of two Hurwitz integers
func (h *HurwitzInt) CmpNorm(a *HurwitzInt) int {
	return h.Norm().Cmp(a.Norm())
}
