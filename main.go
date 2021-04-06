package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("test in main")
	newBases := accumulator.PreCompute()

	elements := make([]accumulator.Element, 1000)
	for i := 0; i < 1000; i++ {
		elements[i] = *accumulator.GetPseudoRandomElement(i)
	}
	result := accumulator.AccumulateSetWirhPreCompute(elements, newBases)
	fmt.Println("result = ", result.String())
}
