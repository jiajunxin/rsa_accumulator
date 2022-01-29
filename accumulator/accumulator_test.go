package accumulator

import (
	"crypto/sha256"
	"math/big"
	"testing"

	"github.com/rsa_accumulator/dihash"
)

const testString = "2021HKUST"

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

func TestSetup(t *testing.T) {
	setup := TrustedSetup()
	var gcd big.Int
	gcd.GCD(nil, nil, &setup.N, &setup.G)
	if gcd.Cmp(one) != 0 {
		// gcd != 1
		//this condition should never happen
		t.Errorf("g and N not co-prime! We win the RSA-2048 challenge!")
	}
}
