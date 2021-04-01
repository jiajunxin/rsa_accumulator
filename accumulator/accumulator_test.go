package accumulator

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/rsa_accumulator/dihash"
)

const testString = "2021HKUST"

func TestGcd(t *testing.T) {
	var testObject AccumulatorSetup
	testObject = *TrustedSetup()

	var N big.Int
	N.Mul(&testObject.P, &testObject.Q)

	tmp := N.Cmp(&testObject.N)
	if tmp != 0 {
		t.Errorf("the N is not calculated correctly")
	}
	flag := isQR(&testObject.G, &testObject.P, &testObject.Q)
	if flag == false {
		t.Errorf("G is QR")
	}
}

func TestDIHash(t *testing.T) {
	// test two different ways of generating DI hash
	// we need to check if A ?= B + C
	var testObject AccumulatorSetup
	testObject = *TrustedSetup()

	dihashValue := dihash.DIHash([]byte(testString))
	A := Accumulate(&testObject.G, dihashValue, &testObject.N)

	B := Accumulate(&testObject.G, dihash.Delta, &testObject.N)
	var tempInt big.Int
	h := sha256.New()
	h.Write([]byte(testString))
	hashTemp := h.Sum(nil)
	tempInt.SetBytes(hashTemp)
	C := Accumulate(&testObject.G, &tempInt, &testObject.N)

	var BCSum big.Int
	BCSum.Mul(B, C)
	BCSum.Mod(&BCSum, &testObject.N)

	tmp := A.Cmp(&BCSum)
	if tmp != 0 {
		t.Errorf("two ways have different result")
	}
}
