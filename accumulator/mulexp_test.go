package accumulator

import (
	"crypto/rand"
	"fmt"
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
	fmt.Println("a = ", a.String())
	fmt.Println("b = ", b.String())

	var aCopy, bCopy big.Int
	aCopy.Set(a)
	bCopy.Set(b)
	fmt.Println("aCopy = ", aCopy.String())
	fmt.Println("bCopy = ", bCopy.String())

	result := GCB(&aCopy, &bCopy)

	var sum, sum2 big.Int
	sum.Add(&aCopy, &bCopy)
	fmt.Println("aCopy = ", aCopy.String())
	fmt.Println("bCopy = ", bCopy.String())
	fmt.Println("result = ", result.String())
	sum.Add(&sum, result)
	fmt.Println("test 3.1 ")
	sum.Add(&sum, result)

	sum2.Add(a, b)
	fmt.Println("test 4 ")
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
	fmt.Println("g = ", g.String())
	fmt.Println("x = ", x.String())
	fmt.Println("N = ", N.String())

	var result2 big.Int
	result2.Exp(g, x, N)
	result := SimpleExp(g, x, N)
	fmt.Println("result = ", result.String())
	fmt.Println("result2 = ", result2.String())

	fmt.Println("test 4 ")
	if result.Cmp(&result2) != 0 {
		t.Errorf("Wrong result for SimpleExp")
	}
}
