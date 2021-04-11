package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("test in main")
	newBases := accumulator.PreCompute()
	fmt.Println("precompute finished")
	setSize := 3
	elements := make([]accumulator.Element, setSize)
	for i := 0; i < setSize; i++ {
		elements[i] = *accumulator.GetPseudoRandomElement(i + 1)
		//fmt.Println("elements = ", elements[i])
	}
	fmt.Println("GetPseudoRandomElement finished")
	result := accumulator.AccumulateSetWirhPreCompute(elements, newBases)
	//fmt.Println("result = ", result.String())

	result2 := accumulator.AccumulateSetWithoutPreCompute(elements)
	if result.Cmp(result2) != 0 {
		fmt.Println("result2 != result1")
	} else {
		fmt.Println("result2 == result1")
	}
}
