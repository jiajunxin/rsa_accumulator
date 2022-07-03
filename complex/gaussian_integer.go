package complex

import (
	"math/big"
)

// GaussianInt implements Gaussian Integer
// In number theory, a Gaussian Integer is a complex number whose real and imaginary parts are both integers
type GaussianInt struct {
	R *big.Int // real part
	I *big.Int // imaginary part
}

// NewGaussianInt declares a new Gaussian Integer with the real part and imaginary part
func NewGaussianInt(r *big.Int, i *big.Int) *GaussianInt {
	return &GaussianInt{
		R: r,
		I: i,
	}
}

// Add adds two Gaussian Integers
func (g *GaussianInt) Add(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Add(a.R, b.R)
	g.I = new(big.Int).Add(a.I, b.I)
	return g
}

// Sub subtracts two Gaussian Integers
func (g *GaussianInt) Sub(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Sub(a.R, b.R)
	g.I = new(big.Int).Sub(a.I, b.I)
	return g
}

// Prod returns the products of two Gaussian Integers
func (g *GaussianInt) Prod(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Mul(a.R, b.R)
	imgMul := new(big.Int).Mul(a.I, b.I)
	g.R.Sub(g.R, imgMul)
	g.I = new(big.Int).Mul(a.R, b.I)
	g.I.Add(g.I, new(big.Int).Mul(a.I, b.R))
	return g
}

// Conj obtains the conjugate of the original Gaussian Integer
func (g *GaussianInt) Conj(origin *GaussianInt) *GaussianInt {
	img := new(big.Int).Neg(origin.I)
	g.R = origin.R
	g.I = img
	return g
}

// Norm obtains the norm of the Gaussian Integer
func (g *GaussianInt) Norm() *big.Int {
	norm := new(big.Int).Mul(g.R, g.R)
	norm.Add(norm, new(big.Int).Mul(g.I, g.I))
	return norm
}

// Copy copies the Gaussian Integer
func (g *GaussianInt) Copy() *GaussianInt {
	return NewGaussianInt(
		new(big.Int).Set(g.R),
		new(big.Int).Set(g.I),
	)
}

// Div performs Euclidean division of two Gaussian Integers, i.e. a/b
// the remainder is stored in the Gaussian Integer that calls the method
// the quotient is returned as a new Gaussian Integer
func (g *GaussianInt) Div(a, b *GaussianInt) *GaussianInt {
	conjB := new(GaussianInt).Conj(b)
	numerator := new(GaussianInt).Prod(a, conjB)
	denominator := new(GaussianInt).Prod(b, conjB)
	deInt := denominator.R

	nuRealFloat := new(big.Float).SetInt(numerator.R)
	nuImagFloat := new(big.Float).SetInt(numerator.I)
	deFloat := new(big.Float).SetInt(deInt)
	realScalar := new(big.Float).Quo(nuRealFloat, deFloat)
	imagScalar := new(big.Float).Quo(nuImagFloat, deFloat)

	rsInt := roundFloat(realScalar)
	isInt := roundFloat(imagScalar)
	quotient := NewGaussianInt(rsInt, isInt)
	g.Sub(a, new(GaussianInt).Prod(b, quotient))
	return quotient
}

func (g *GaussianInt) String() string {
	sign := "+"
	if g.I.Sign() < 0 {
		sign = ""
	}
	return g.R.String() + sign + g.I.String() + "i"
}

func roundFloat(f *big.Float) *big.Int {
	delta := big.NewFloat(roundingDelta)
	if f.Sign() < 0 {
		delta.Neg(delta)
	}
	f.Add(f, delta)
	res := new(big.Int)
	f.Int(res)
	return res
}
