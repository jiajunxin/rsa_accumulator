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
	prodRecursive := *SetProductRecursive(rep)
	prodRecursiveParallel := *SetProductParallel(rep, 2)

	tmp := prodRecursive.Cmp(&prodRecursiveParallel)
	if tmp != 0 {
		t.Errorf("two ways have different result")
	}
}
