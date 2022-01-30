package accumulator

import (
	"crypto/sha256"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/rsa_accumulator/dihash"
)

const testString = "2021HKUST"

func TestDIHash(t *testing.T) {
	// test two different ways of generating DI hash
	// we need to check if A ?= B + C
	testObject := *TrustedSetup()

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
	len := Min2048.BitLen()
	if len != 2048 {
		t.Errorf("Min2048 is not 2048 bits")
	}
}

func TestAccAndProve(t *testing.T) {
	setup := TrustedSetup()

	testSetSize := 16
	set := GenTestSet(testSetSize)
	acc, proofs := AccAndProve(set, HashToPrimeFromSha256, setup)
	if len(set) != len(proofs) {
		t.Errorf("proofs have different size as the input set")
	}
	rep := GenRepersentatives(set, HashToPrimeFromSha256)
	acc2 := accumulate(rep, &setup.G, &setup.N)
	acc3 := Accumulate(&proofs[5], &rep[5], &setup.N)
	if acc.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	if acc2.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	// test another set size not a power of 2
	testSetSize = 17
	set = GenTestSet(testSetSize)
	acc, proofs = AccAndProve(set, HashToPrimeFromSha256, setup)
	if len(set) != len(proofs) {
		t.Errorf("proofs have different size as the input set")
	}
	rep = GenRepersentatives(set, HashToPrimeFromSha256)
	acc2 = accumulate(rep, &setup.G, &setup.N)
	acc3 = Accumulate(&proofs[7], &rep[7], &setup.N)
	if acc.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	if acc2.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	// test another set size not a power of 2
	testSetSize = 254
	set = GenTestSet(testSetSize)
	acc, proofs = AccAndProve(set, HashToPrimeFromSha256, setup)
	if len(set) != len(proofs) {
		t.Errorf("proofs have different size as the input set")
	}
	rep = GenRepersentatives(set, HashToPrimeFromSha256)
	acc2 = accumulate(rep, &setup.G, &setup.N)
	acc3 = Accumulate(&proofs[253], &rep[253], &setup.N)
	if acc.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	if acc2.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

}

func GenTestSet(num int) []string {
	ret := make([]string, num)
	for i := 0; i < num; i++ {
		temp := rand.Intn(100000)
		ret[i] = strconv.Itoa(temp)
	}
	return ret
}
