package complex

import (
	"math/big"
)

// GaussianInt implements Gaussian integer
// In number theory, a Gaussian integer is a complex number whose real and imaginary parts are both integers
type GaussianInt struct {
	R *big.Int // real part
	I *big.Int // imaginary part
}

// String returns the string representation of the Gaussian integer
func (g *GaussianInt) String() string {
	res := ""
	if g.R.Sign() != 0 {
		res += g.R.String()
	}
	gISign := g.I.Sign()
	if gISign == 0 {
		if res == "" {
			return "0"
		}
		return res
	}
	if gISign == 1 {
		res += "+"
	}
	if g.I.Cmp(bigNeg1) == 0 {
		res += "-"
	} else if g.I.Cmp(big1) != 0 {
		res += g.I.String()
	}
	res += "i"
	return res
}

// NewGaussianInt declares a new Gaussian integer with the real part and imaginary part
func NewGaussianInt(r *big.Int, i *big.Int) *GaussianInt {
	return &GaussianInt{
		R: new(big.Int).Set(r),
		I: new(big.Int).Set(i),
	}
}

// Set sets the Gaussian integer to the given Gaussian integer
func (g *GaussianInt) Set(a *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Set(a.R)
	g.I = new(big.Int).Set(a.I)
	return g
}

// Update updates the Gaussian integer with the given real and imaginary parts
func (g *GaussianInt) Update(r, i *big.Int) {
	g.R = r
	g.I = i
}

// Add adds two Gaussian integers
func (g *GaussianInt) Add(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Add(a.R, b.R)
	g.I = new(big.Int).Add(a.I, b.I)
	return g
}

// Sub subtracts two Gaussian integers
func (g *GaussianInt) Sub(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Sub(a.R, b.R)
	g.I = new(big.Int).Sub(a.I, b.I)
	return g
}

// Prod returns the products of two Gaussian integers
func (g *GaussianInt) Prod(a, b *GaussianInt) *GaussianInt {
	g.R = new(big.Int).Mul(a.R, b.R)
	imgMul := new(big.Int).Mul(a.I, b.I)
	g.R.Sub(g.R, imgMul)
	g.I = new(big.Int).Mul(a.R, b.I)
	g.I.Add(g.I, new(big.Int).Mul(a.I, b.R))
	return g
}

// Conj obtains the conjugate of the original Gaussian integer
func (g *GaussianInt) Conj(origin *GaussianInt) *GaussianInt {
	img := new(big.Int).Neg(origin.I)
	g.Update(origin.R, img)
	return g
}

// Norm obtains the norm of the Gaussian integer
func (g *GaussianInt) Norm() *big.Int {
	norm := new(big.Int).Mul(g.R, g.R)
	norm.Add(norm, new(big.Int).Mul(g.I, g.I))
	return norm
}

// Copy copies the Gaussian integer
func (g *GaussianInt) Copy() *GaussianInt {
	return NewGaussianInt(
		new(big.Int).Set(g.R),
		new(big.Int).Set(g.I),
	)
}

// Div performs Euclidean division of two Gaussian integers, i.e. a/b
// the remainder is stored in the Gaussian integer that calls the method
// the quotient is returned as a new Gaussian integer
func (g *GaussianInt) Div(a, b *GaussianInt) *GaussianInt {
	bConj := new(GaussianInt).Conj(b)
	numerator := new(GaussianInt).Prod(a, bConj)
	denominator := new(GaussianInt).Prod(b, bConj)
	deInt := denominator.R
	deFloat := new(big.Float).SetInt(deInt)

	nuRealFloat := new(big.Float).SetInt(numerator.R)
	nuImagFloat := new(big.Float).SetInt(numerator.I)
	realScalar := new(big.Float).Quo(nuRealFloat, deFloat)
	imagScalar := new(big.Float).Quo(nuImagFloat, deFloat)

	rsInt := roundFloat(realScalar)
	isInt := roundFloat(imagScalar)
	quotient := NewGaussianInt(rsInt, isInt)
	g.Sub(a, new(GaussianInt).Prod(quotient, b))
	return quotient
}

// IsZero returns true if the Gaussian integer is zero
func (g *GaussianInt) IsZero() bool {
	return g.R.Sign() == 0 && g.I.Sign() == 0
}

// CmpNorm compares the norm of two Gaussian integers
func (g *GaussianInt) CmpNorm(a *GaussianInt) int {
	return g.Norm().Cmp(a.Norm())
}

// GCD calculates the greatest common divisor of two Gaussian integers using Euclidean algorithm
// the result is stored in the Gaussian integer that calls the method and returned
func (g *GaussianInt) GCD(a, b *GaussianInt) *GaussianInt {
	a = a.Copy()
	b = b.Copy()
	remainder := new(GaussianInt)
	for {
		remainder.Div(a, b)
		if remainder.IsZero() {
			g.Update(b.R, b.I)
			return b
		}
		a.Update(b.R, b.I)
		b.Update(remainder.R, remainder.I)
	}
}
