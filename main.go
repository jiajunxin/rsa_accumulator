package main

import (
	"fmt"

	"github.com/rsa_accumulator/accumulator"
)

func main() {
	fmt.Println("test in main")
	accumulator.ManualBench(10000)
}
