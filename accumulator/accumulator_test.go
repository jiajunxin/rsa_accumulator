package accumulator

import (
	"crypto/sha256"
	"math/big"
	"math/rand"
	"strconv"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/dihash"
)

const testString = "2021HKUST"

func TestDIHash(t *testing.T) {
	// test two different ways of generating DI hash
	// we need to check if A ?= B + C
	testObject := TrustedSetup()

	diHashValue := dihash.DIHash([]byte(testString))
	A := AccumulateNew(testObject.G, diHashValue, testObject.N)

	B := AccumulateNew(testObject.G, dihash.Delta, testObject.N)
	var tempInt big.Int
	h := sha256.New()
	h.Write([]byte(testString))
	hashTemp := h.Sum(nil)
	tempInt.SetBytes(hashTemp)
	C := AccumulateNew(testObject.G, &tempInt, testObject.N)

	var BCSum big.Int
	BCSum.Mul(B, C)
	BCSum.Mod(&BCSum, testObject.N)

	tmp := A.Cmp(&BCSum)
	if tmp != 0 {
		t.Errorf("two ways have different result")
	}
}

func TestSetup(t *testing.T) {
	setup := TrustedSetup()
	var gcd big.Int
	gcd.GCD(nil, nil, setup.N, setup.G)
	if gcd.Cmp(big1) != 0 {
		// gcd != 1
		//this condition should never happen
		t.Errorf("g and N not co-prime! We win the RSA-2048 challenge!")
	}
	bitLen := Min2048.BitLen()
	if bitLen != 2048 {
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
	rep := GenRepresentatives(set, HashToPrimeFromSha256)
	acc2 := accumulateNew(setup.G, setup.N, rep)
	acc3 := AccumulateNew(proofs[5], rep[5], setup.N)
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
	rep = GenRepresentatives(set, HashToPrimeFromSha256)
	acc2 = accumulateNew(setup.G, setup.N, rep)
	acc3 = AccumulateNew(proofs[7], rep[7], setup.N)
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
	rep = GenRepresentatives(set, HashToPrimeFromSha256)
	acc2 = accumulateNew(setup.G, setup.N, rep)
	acc3 = AccumulateNew(proofs[252], rep[252], setup.N)
	if acc.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	if acc2.Cmp(acc3) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
}

func genAccts(set []string, setup *Setup, proofs []*big.Int, idx int) (acc1, acc2 *big.Int) {
	rep := GenRepresentatives(set, HashToPrimeFromSha256)
	acc1 = accumulateNew(setup.G, setup.N, rep)
	acc2 = AccumulateNew(proofs[idx], rep[idx], setup.N)
	return
}

func TestAccAndProveIter(t *testing.T) {

	type args struct {
		set        []string
		encodeType EncodeType
		setup      *Setup
	}
	tests := []struct {
		name         string
		args         args
		idx          int
		wantProofLen int
	}{
		{
			name: "set_size_16",
			args: args{
				set:        GenTestSet(16),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          5,
			wantProofLen: 16,
		},
		{
			name: "set_size_17",
			args: args{
				set:        GenTestSet(17),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          7,
			wantProofLen: 17,
		},
		{
			name: "set_size_254",
			args: args{
				set:        GenTestSet(254),
				encodeType: HashToPrimeFromSha256,
				setup:      TrustedSetup(),
			},
			idx:          253,
			wantProofLen: 254,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			acc, proofs := AccAndProveIter(tt.args.set, tt.args.encodeType, tt.args.setup)
			acc1, acc2 := genAccts(tt.args.set, tt.args.setup, proofs, tt.idx)
			if len(proofs) != tt.wantProofLen {
				t.Errorf("AccAndProveIter() got proof len = %v, want %v", len(proofs), tt.wantProofLen)
				return
			}
			if acc.Cmp(acc2) != 0 {
				t.Errorf("proofs generated are not consistent acc = %v, acc2 %v", acc, acc2)
				return
			}
			if acc1.Cmp(acc2) != 0 {
				t.Errorf("proofs generated are not consistent acc1 = %v, acc2 %v", acc1, acc2)
				return
			}
		})
	}
}

func GenTestSet(num int) []string {
	ret := make([]string, num)
	for i := 0; i < num; i++ {
		temp := rand.Intn(100000000)
		ret[i] = strconv.Itoa(temp)
	}
	return ret
}
