package accumulator

import (
	"math/big"
	"testing"

	"github.com/rsa_accumulator/dihash"
)

func BenchmarkHashToPrime(b *testing.B) {
	testBytes := []byte(testString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		HashToPrime(testBytes)
	}
}

func BenchmarkAccAndProve(b *testing.B) {
	testSetSize := 10000
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = AccAndProve(set, HashToPrimeFromSha256, &setup)
	}
}

func BenchmarkDIHash(b *testing.B) {
	testBytes := []byte(testString)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = dihash.DIHash(testBytes)
	}
}

func BenchmarkAccumulate256bits(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)
	prime256bits := HashToPrime(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Accumulate(&testObject.G, prime256bits, &testObject.N)
	}
}

func BenchmarkAccumulateDIHash(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)
	dihashResult := dihash.DIHash(testBytes)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Accumulate(&testObject.G, dihashResult, &testObject.N)
	}
}

func BenchmarkAccumulateDIHashWithPreCompute(b *testing.B) {
	testObject := *TrustedSetup()
	testBytes := []byte(testString)

	B := Accumulate(&testObject.G, dihash.Delta, &testObject.N)
	tempInt := *SHA256ToInt(testBytes)
	var BCSum big.Int

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		C := Accumulate(&testObject.G, &tempInt, &testObject.N)
		BCSum.Mul(B, C)
		BCSum.Mod(&BCSum, &testObject.N)
	}
}
