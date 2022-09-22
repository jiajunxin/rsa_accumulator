package accumulator

import (
	"math/big"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/dihash"
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
	rep := GenRepresentatives(set, HashToPrimeFromSha256)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProveMembership(setup.G, setup.N, rep)
	}
}

func BenchmarkProveMembershipIter(b *testing.B) {
	setSize := 1000
	set := GenBenchSet(setSize)
	rep := GenRepresentatives(set, HashToPrimeFromSha256)
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

func BenchmarkDIHash(b *testing.B) {
	testBytes := []byte(testString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dihash.DIHash(testBytes)
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

func BenchmarkAccumulateNewDIHash(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)
	diHashResult := dihash.DIHash(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AccumulateNew(testObject.G, diHashResult, testObject.N)
	}
}

func BenchmarkAccumulateDIHashWithPreCompute(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)

	B := AccumulateNew(testObject.G, dihash.Delta, testObject.N)
	tempInt := *SHA256ToInt(testBytes)
	var BCSum big.Int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C := AccumulateNew(testObject.G, &tempInt, testObject.N)
		BCSum.Mul(B, C)
		BCSum.Mod(&BCSum, testObject.N)
	}
}
