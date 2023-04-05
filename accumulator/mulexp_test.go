package accumulator

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestCommonBits(t *testing.T) {
	var a, b uint
	a = 1024
	b = 1025

	result := CommonBits(big.Word(a), big.Word(b))
	if result != 1024 {
		t.Errorf("Wrong result for CommonBits, result for (1024, 1025) = %d", result)
	}

	a = 7 // binary 111
	b = 1025

	result = CommonBits(big.Word(a), big.Word(b))
	if result != 1 {
		t.Errorf("Wrong result for CommonBits, result for (1024, 1025) = %d", result)
	}
}

func TestGCB(t *testing.T) {
	var max big.Int
	max.SetInt64(1000000)

	a, err := rand.Int(rand.Reader, &max)
	if err != nil {
		t.Errorf(err.Error())
	}
	b, err := rand.Int(rand.Reader, &max)
	if err != nil {
		t.Errorf(err.Error())
	}

	var aCopy, bCopy big.Int
	aCopy.Set(a)
	bCopy.Set(b)

	result := GCB(&aCopy, &bCopy)

	var sum, sum2 big.Int
	sum.Add(&aCopy, &bCopy)
	sum.Add(&sum, result)
	sum.Add(&sum, result)

	sum2.Add(a, b)
	if sum.Cmp(&sum2) != 0 {
		t.Errorf("Wrong result for GCB")
	}
}

func TestSimpleExp(t *testing.T) {
	var max big.Int
	max.SetInt64(1000000)

	g, err := rand.Int(rand.Reader, &max)
	if err != nil {
		t.Errorf(err.Error())
	}
	x, err := rand.Int(rand.Reader, &max)
	if err != nil {
		t.Errorf(err.Error())
	}
	N, err := rand.Int(rand.Reader, &max)
	if err != nil {
		t.Errorf(err.Error())
	}

	var result2 big.Int
	result2.Exp(g, x, N)
	result := SimpleExp(g, x, N)

	if result.Cmp(&result2) != 0 {
		t.Errorf("Wrong result for SimpleExp")
	}
}
