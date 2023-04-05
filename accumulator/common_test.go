package accumulator

import (
	"testing"
)

func TestProduct(t *testing.T) {
	// test two different ways of generating DI hash
	// we need to check if A ?= B + C
	setSize := 1000
	set := GenBenchSet(setSize)
	rep := GenRepresentatives(set, DIHashFromPoseidon)
	prodRecursive := *SetProductRecursiveFast(rep)
	prodRecursiveParallel := *SetProductParallel(rep, 2)
	prodNative := SetProduct2(rep)

	tmp := prodRecursive.Cmp(&prodRecursiveParallel)
	if tmp != 0 {
		t.Errorf("two ways have different result")
	}
	tmp = prodRecursive.Cmp(prodNative)
	if tmp != 0 {
		t.Errorf("two ways have different result")
	}
}
