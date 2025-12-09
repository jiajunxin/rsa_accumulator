package main

import (
	"github.com/jiajunxin/rsa_accumulator/experiments"
	"github.com/jiajunxin/rsa_accumulator/zkmultiswap"
)

const (
	twoTo14 = 16384
	twoTo15 = 32768
	twoTo16 = 65536
	twoTo17 = 131072
	twoTo18 = 262144
	twoTo19 = 524288
)

func main() {
	println("Test for MultiSwap for 1024 swaps")
	zkmultiswap.TestMultiSwap(1024) //2^10

	println("The experiments below require a large memory for precomputation table")
	println("Test for set size ", twoTo14, " with 2^", 3, " cores")
	experiments.TestRSASubsetParallel(twoTo14, 1024, 3)
	println("Test for set size ", twoTo15, " with 2^", 3, " cores")
	experiments.TestRSASubsetParallel(twoTo15, 1024, 3)
	println("Test for set size ", twoTo16, " with 2^", 3, " cores")
	experiments.TestRSASubsetParallel(twoTo16, 1024, 3)
	println("Test for set size ", twoTo17, " with 2^", 3, " cores")
	experiments.TestRSASubsetParallel(twoTo17, 1024, 3)
	println("Test for set size ", twoTo18, " with 2^", 3, " cores")
	experiments.TestRSASubsetParallel(twoTo18, 1024, 3)
}
