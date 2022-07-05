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
	DBLr *big.Int
	DBLi *big.Int
	DBLj *big.Int
	DBLk *big.Int
}

// NewHurwitzInt declares a new integral quaternion with the real, i, j, and k parts
// If isDouble is true, the arguments r, i, j, k are twice the original scalars
func NewHurwitzInt(r, i, j, k *big.Int, isDouble bool) *HurwitzInt {
	if isDouble {
		return &HurwitzInt{
			DBLr: new(big.Int).Set(r),
			DBLi: new(big.Int).Set(i),
			DBLj: new(big.Int).Set(j),
			DBLk: new(big.Int).Set(k),
		}
	}
	return &HurwitzInt{
		DBLr: new(big.Int).Mul(r, big2),
		DBLi: new(big.Int).Mul(i, big2),
		DBLj: new(big.Int).Mul(j, big2),
		DBLk: new(big.Int).Mul(k, big2),
	}
}

// Update updates the integral quaternion with the given real, i, j, and k parts
func (h *HurwitzInt) Update(r, i, j, k *big.Int, isDouble bool) *HurwitzInt {
	if isDouble {
		h.DBLr = r
		h.DBLi = i
		h.DBLj = j
		h.DBLk = k
	} else {
		h.DBLr = new(big.Int).Mul(r, big2)
		h.DBLi = new(big.Int).Mul(i, big2)
		h.DBLj = new(big.Int).Mul(j, big2)
		h.DBLk = new(big.Int).Mul(k, big2)
	}
	return h
}

// Add adds two integral quaternions
func (h *HurwitzInt) Add(a, b *HurwitzInt) *HurwitzInt {
	h.DBLr = new(big.Int).Add(a.DBLr, b.DBLr)
	h.DBLi = new(big.Int).Add(a.DBLi, b.DBLi)
	h.DBLj = new(big.Int).Add(a.DBLj, b.DBLj)
	h.DBLk = new(big.Int).Add(a.DBLk, b.DBLk)
	return h
}

// Sub subtracts two integral quaternions
func (h *HurwitzInt) Sub(a, b *HurwitzInt) *HurwitzInt {
	h.DBLr = new(big.Int).Sub(a.DBLr, b.DBLr)
	h.DBLi = new(big.Int).Sub(a.DBLi, b.DBLi)
	h.DBLj = new(big.Int).Sub(a.DBLj, b.DBLj)
	h.DBLk = new(big.Int).Sub(a.DBLk, b.DBLk)
	return h
}

// Conj obtains the conjugate of the original integral quaternion
func (h *HurwitzInt) Conj(origin *HurwitzInt) *HurwitzInt {
	h.DBLr = origin.DBLr
	h.DBLi = new(big.Int).Neg(origin.DBLi)
	h.DBLj = new(big.Int).Neg(origin.DBLj)
	h.DBLk = new(big.Int).Neg(origin.DBLk)
	return h
}

// Norm obtains the norm of the integral quaternion
func (h *HurwitzInt) Norm() *big.Int {
	norm := new(big.Int).Mul(h.DBLr, h.DBLr)
	norm.Add(norm, new(big.Int).Mul(h.DBLi, h.DBLi))
	norm.Add(norm, new(big.Int).Mul(h.DBLj, h.DBLj))
	norm.Add(norm, new(big.Int).Mul(h.DBLk, h.DBLk))
	norm.Div(norm, big4)
	return norm
}

// Copy copies the integral quaternion
func (h *HurwitzInt) Copy() *HurwitzInt {
	return NewHurwitzInt(h.DBLr, h.DBLi, h.DBLj, h.DBLk, true)
}

// Prod returns the Hamilton product of two integral quaternions
// the product (a1 + b1j + c1k + d1)(a2 + b2j + c2k + d2) is determined by the products of the
// basis elements and the distributive law
func (h *HurwitzInt) Prod(a, b *HurwitzInt) *HurwitzInt {
	opt := new(big.Int)
	// 1 part
	h.DBLr = new(big.Int).Mul(a.DBLr, b.DBLr)
	h.DBLr.Sub(h.DBLr, opt.Mul(a.DBLi, b.DBLi))
	h.DBLr.Sub(h.DBLr, opt.Mul(a.DBLj, b.DBLj))
	h.DBLr.Sub(h.DBLr, opt.Mul(a.DBLk, b.DBLk))
	h.DBLr.Div(h.DBLr, big2)

	// i part
	h.DBLi = new(big.Int).Mul(a.DBLr, b.DBLi)
	h.DBLi.Add(h.DBLi, opt.Mul(a.DBLi, b.DBLr))
	h.DBLi.Add(h.DBLi, opt.Mul(a.DBLj, b.DBLk))
	h.DBLi.Sub(h.DBLi, opt.Mul(a.DBLk, b.DBLj))
	h.DBLi.Div(h.DBLi, big2)

	// j part
	h.DBLj = new(big.Int).Mul(a.DBLr, b.DBLj)
	h.DBLj.Sub(h.DBLj, opt.Mul(a.DBLi, b.DBLk))
	h.DBLj.Add(h.DBLj, opt.Mul(a.DBLj, b.DBLr))
	h.DBLj.Add(h.DBLj, opt.Mul(a.DBLk, b.DBLi))
	h.DBLj.Div(h.DBLj, big2)

	// k part
	h.DBLk = new(big.Int).Mul(a.DBLr, b.DBLk)
	h.DBLk.Add(h.DBLk, opt.Mul(a.DBLi, b.DBLj))
	h.DBLk.Sub(h.DBLk, opt.Mul(a.DBLj, b.DBLi))
	h.DBLk.Add(h.DBLk, opt.Mul(a.DBLk, b.DBLr))
	h.DBLk.Div(h.DBLk, big2)

	return h
}

// Div performs Euclidean division of two Hurwitz integers, i.e. a/b
// the remainder is stored in the Hurwitz integer that calls the method
// the quotient is returned as a new Hurwitz integer
func (h *HurwitzInt) Div(a, b *HurwitzInt) *HurwitzInt {
	panic("not implemented")
}

// String returns the string representation of the integral quaternion
func (h *HurwitzInt) String() string {
	if h.DBLr.Sign() == 0 && h.DBLi.Sign() == 0 && h.DBLj.Sign() == 0 && h.DBLk.Sign() == 0 {
		return "0"
	}
	res := ""
	if h.DBLr.Sign() != 0 {
		res += new(big.Int).Div(h.DBLr, big2).String()
		if h.DBLr.Bit(0) == 1 {
			res += ".5"
		}
	}
	if h.DBLi.Sign() != 0 {
		if h.DBLr.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.DBLi, big2).String()
		if h.DBLi.Bit(0) == 1 {
			res += ".5"
		}
		res += "i"
	}
	if h.DBLj.Sign() != 0 {
		if h.DBLj.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.DBLj, big2).String()
		if h.DBLj.Bit(0) == 1 {
			res += ".5"
		}
		res += "j"
	}
	if h.DBLk.Sign() != 0 {
		if h.DBLk.Sign() == 1 {
			res += "+"
		}
		res += new(big.Int).Div(h.DBLk, big2).String()
		if h.DBLk.Bit(0) == 1 {
			res += ".5"
		}
		res += "k"
	}
	return res
}
