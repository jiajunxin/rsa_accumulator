package accumulator

import (
	"crypto/sha256"
	"fmt"
	"math/big"
	"math/bits"
	"math/rand"
	"strconv"
	"testing"

	"github.com/jiajunxin/multiexp"
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
	bitLen := Min1024.BitLen()
	if bitLen != 2048 {
		t.Errorf("Min2048 is not 2048 bits")
	}
}

func TestMultiDIHash(t *testing.T) {
	testSetSize := 2048
	set := GenTestSet(testSetSize)
	rep := GenRepresentatives(set, MultiDIHashFromPoseidon)
	if len(rep) != testSetSize*3 {
		t.Errorf("Representatives are not consistent for MultiDIHashFromPoseidon")
	}
	for i, v := range rep {
		if v.Cmp(big1) != 1 {
			t.Errorf("Representatives are not correct for position %d", i)
		}
	}
}

func TestAccAndProve(t *testing.T) {
	setup := TrustedSetup()

	testSetSize := 3072
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
	acc3 = AccumulateNew(proofs[1], rep[1], setup.N)
	if acc.Cmp(acc3) != 0 {
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

func TestAccAndProveWithMultiDI(t *testing.T) {
	setup := TrustedSetup()

	testSetSize := 256
	set := GenTestSet(testSetSize)
	//rep := GenRepresentatives(set, DIHashFromPoseidon)
	rep2 := GenRepresentatives(set, HashToPrimeFromSha256)
	rep3 := GenRepresentatives(set, HashToPrimeFromSha256)
	rep1 := make([]*big.Int, testSetSize)
	//rep2 := make([]*big.Int, testSetSize)
	//rep3 := make([]*big.Int, testSetSize)
	for i := range rep1 {
		rep1[i] = new(big.Int)
		rep1[i] = GenRandomizer()
		if rep1[i].Bit(0) == 0 {
			rep1[i].Add(rep1[i], big1)
		}
		//rep1[i] = getPrime256()
		//rep2[i] = GenRandomizer()
		//rep2[i] = getPrime256()
		//	fmt.Println("rep1[", i, "] = ", rep1[i].String())
	}

	proofs1 := ProveMembership(setup.G, setup.N, rep1)
	proofs2 := ProveMembership(setup.G, setup.N, rep2)
	proofs3 := ProveMembership(setup.G, setup.N, rep3)
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation

	if len(set) != len(proofs1) {
		t.Errorf("proofs have different size as the input set")
	}
	//fmt.Println("rep1[1] = ", rep1[1].String())
	acc1 := accumulateNew(setup.G, setup.N, rep1)
	//fmt.Println("rep1[1] = ", rep1[1].String())
	acc1Temp := AccumulateNew(proofs1[3], rep1[3], setup.N)
	// for i := range rep {
	// 	fmt.Println("rep1[", i, "] = ", rep1[i].String())
	// }
	if acc1.Cmp(acc1Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	acc2 := accumulateNew(setup.G, setup.N, rep2)
	acc2Temp := AccumulateNew(proofs2[1], rep2[1], setup.N)
	if acc2.Cmp(acc2Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	prod3 := SetProductRecursiveFast(rep3)
	acc3 := AccumulateNew(setup.G, prod3, setup.N)
	acc3Temp := AccumulateNew(proofs3[7], rep3[7], setup.N)
	if acc3.Cmp(acc3Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

}

func TestAccAndProveWithMultiDI2(t *testing.T) {
	setup := TrustedSetup()

	testSetSize := 80
	set := GenTestSet(testSetSize)
	rep := GenRepresentatives(set, MultiDIHashFromPoseidon)

	proofs1 := ProveMembership(setup.G, setup.N, rep[:testSetSize])
	proofs2 := ProveMembership(setup.G, setup.N, rep[testSetSize:2*testSetSize])
	proofs3 := ProveMembership(setup.G, setup.N, rep[2*testSetSize:])
	// we generate the accumulator by anyone of the membership proof raised to its power to save some calculation

	if len(set) != len(proofs1) {
		t.Errorf("proofs have different size as the input set")
	}
	acc1 := accumulateNew(setup.G, setup.N, rep[0:testSetSize])
	acc1Temp := AccumulateNew(proofs1[0], rep[0], setup.N)
	if acc1.Cmp(acc1Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	acc2 := accumulateNew(setup.G, setup.N, rep[testSetSize:2*testSetSize])
	acc2Temp := AccumulateNew(proofs2[5], rep[testSetSize+5], setup.N)
	if acc2.Cmp(acc2Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

	prod3 := SetProductRecursiveFast(rep[2*testSetSize:])
	acc3 := AccumulateNew(setup.G, prod3, setup.N)
	acc3Temp := AccumulateNew(proofs3[7], rep[2*testSetSize+7], setup.N)
	if acc3.Cmp(acc3Temp) != 0 {
		t.Errorf("proofs generated are not consistent")
	}

}

func TestProveMembershipParallelWithTableWithRandomizer(t *testing.T) {
	setSize := 256
	set := GenBenchSet(setSize)
	setup := *TrustedSetup()

	rep := GenRepresentatives(set, MultiDIHashFromPoseidon)
	fmt.Println("set size = ", len(rep))
	// generate a zero-knowledge RSA accumulator
	r1 := GenRandomizer()
	r2 := GenRandomizer()
	r3 := GenRandomizer()
	// fmt.Println("r1 = ", r1.String())
	// fmt.Println("r2 = ", r2.String())

	randomizedbase1 := AccumulateNew(setup.G, r1, setup.N)
	randomizedbase2 := AccumulateNew(setup.G, r2, setup.N)
	randomizedbase3 := AccumulateNew(setup.G, r3, setup.N)
	// calculate the exponentation
	exp1 := SetProductRecursiveFast(rep[:setSize])
	exp2 := SetProductRecursiveFast(rep[setSize : 2*setSize])
	exp3 := SetProductRecursiveFast(rep[2*setSize:])
	acc1 := AccumulateNew(randomizedbase1, exp1, setup.N)
	acc2 := AccumulateNew(randomizedbase2, exp2, setup.N)
	acc3 := AccumulateNew(randomizedbase3, exp3, setup.N)
	fmt.Println("acc1 = ", acc1.String())
	// fmt.Println("acc3 = ", acc3.String())
	// temp := append
	var tempProd1 big.Int
	tempProd1.Mul(exp1, r1)
	acc1temp := AccumulateNew(setup.G, &tempProd1, setup.N)
	fmt.Println("acc1temp = ", acc1temp.String())
	tempProd1.Mod(&tempProd1, setup.N)
	fmt.Println("tempProd1 = ", tempProd1.String())

	maxLen := setSize * 256 / bits.UintSize
	table := multiexp.NewPrecomputeTable(setup.G, setup.N, maxLen)
	proofs1 := ProveMembershipParallelWithTableWithRandomizer(setup.G, r1, setup.N, rep[:setSize], 0, table)
	proofs2 := ProveMembershipParallelWithTableWithRandomizer(setup.G, r2, setup.N, rep[setSize:2*setSize], 0, table)
	proofs3 := ProveMembershipParallelWithTableWithRandomizer(setup.G, r3, setup.N, rep[2*setSize:], 0, table)

	temp1 := AccumulateNew(proofs1[0], rep[0], setup.N)
	if temp1.Cmp(acc1) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	temp2 := AccumulateNew(proofs2[5], rep[setSize+5], setup.N)
	// fmt.Println("size of memberships2 = ", len(proofs2))
	// fmt.Println("temp2 = ", temp2.String())
	// fmt.Println("acc2 = ", acc2.String())
	if temp2.Cmp(acc2) != 0 {
		t.Errorf("proofs generated are not consistent")
	}
	temp3 := AccumulateNew(proofs3[9], rep[2*setSize+9], setup.N)
	if temp3.Cmp(acc3) != 0 {
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
