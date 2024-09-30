package accumulator

import (
	crand "crypto/rand"
	"math/big"
	"testing"

	"github.com/jiajunxin/multiexp"
)

func BenchmarkHashToPrime(b *testing.B) {
	testBytes := []byte(testString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashToPrime(testBytes)
	}
}

func BenchmarkAccAndProve(b *testing.B) {
	testSetSize := 1000
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AccAndProve(set, HashToPrimeFromSha256, &setup)
	}
}

func BenchmarkProveMembership(b *testing.B) {
	setSize := 1000
	set := GenBenchSet(setSize)
	rep := GenRepresentatives(set, DIHashFromPoseidon)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProveMembership(setup.G, setup.N, rep)
	}
}

func BenchmarkProveMembershipIter(b *testing.B) {
	setSize := 1000
	set := GenBenchSet(setSize)
	rep := GenRepresentatives(set, DIHashFromPoseidon)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProveMembershipIter(*setup.G, setup.N, rep)
	}
}

func BenchmarkProveMembershipParallel(b *testing.B) {
	setSize := 1000
	set := GenBenchSet(setSize)
	rep := GenRepresentatives(set, HashToPrimeFromSha256)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProveMembershipParallel(setup.G, setup.N, rep, 8)
	}
}

func BenchmarkAccumulateNew256bits(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)
	prime256bits := HashToPrime(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AccumulateNew(testObject.G, prime256bits, testObject.N)
	}
}

func BenchmarkGroupElementMul(b *testing.B) {
	setup := *TrustedSetup()

	setSize := 10000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, setup.N)
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		groupElementMul(setup.N, set)
	}
}

func BenchmarkGroupElementSquare(b *testing.B) {
	setup := *TrustedSetup()

	setSize := 10000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, setup.N)
	}
	if err != nil {
		b.Fatal(err)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		groupElementSquare(setup.N, set)
	}
}

func groupElementMul(N *big.Int, set []*big.Int) {
	var temp big.Int
	for i := range set {
		if i > 0 {
			temp.Mul(set[i], set[i-1])
			temp.Mod(&temp, N)
		}
	}
}

func groupElementSquare(N *big.Int, set []*big.Int) {
	var temp big.Int
	for i := range set {
		if i > 0 {
			temp.Exp(set[i], big2, N)
		}
	}
}

func BenchmarkExp(b *testing.B) {
	setup := *TrustedSetup()

	setSize := 1000000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, setup.N)
		if err != nil {
			b.Fatal(err)
		}
	}

	b.ResetTimer()
	var temp big.Int
	for j := 0; j < setSize; j++ {
		temp.Exp(setup.G, set[j], setup.N)
	}
}

func BenchmarkDoubleExp(b *testing.B) {
	setup := *TrustedSetup()

	var largeTestNum big.Int
	largeTestNum.Mul(setup.N, setup.N)
	setSize := 10000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, &largeTestNum)
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 1; i < b.N; i++ {
		multiexp.DoubleExp(setup.G, [2]*big.Int{set[i-1], set[i]}, setup.N)
	}
}

func BenchmarkFourfoldExp(b *testing.B) {
	setup := *TrustedSetup()

	var largeTestNum big.Int
	largeTestNum.Mul(setup.N, setup.N)
	setSize := 10000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, &largeTestNum)
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 4; i < b.N; i++ {
		multiexp.FourfoldExp(setup.G, setup.N, [4]*big.Int{set[0], set[1], set[2], set[3]})
	}
}

func BenchmarkSimpleExp(b *testing.B) {
	setup := *TrustedSetup()

	setSize := 100
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, setup.N)
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 0; i < setSize; i++ {
		_ = SimpleExp(setup.G, set[i], setup.N)
	}
}

func BenchmarkGCB(b *testing.B) {
	setup := *TrustedSetup()

	setSize := 10000
	set := make([]*big.Int, setSize)
	var err error
	for i := range set {
		set[i], err = crand.Int(crand.Reader, setup.N)
	}
	if err != nil {
		b.Fatal(err)
	}

	b.ResetTimer()
	for i := 1; i < setSize; i++ {
		_ = GCB(set[i-1], set[i])
	}
}
