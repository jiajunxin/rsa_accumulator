package main

import (
	"fmt"
	"math/big"
	"time"

	"github.com/jiajunxin/rsa_accumulator/experiments"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
	"github.com/jiajunxin/rsa_accumulator/precompute"
)

func testPreCompute() {
	setSize := 65536 // 2 ^ 16
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)

	elementUpperBound := new(big.Int).Lsh(big.NewInt(1), 2047)
	startingTime := time.Now().UTC()
	t := precompute.NewTable(setup.G, setup.N, elementUpperBound, uint64(setSize), 102400)
	duration := time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running NewTable Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	experiments.ProveMembershipParallel2(t, setup.G, setup.N, rep, 4, 16)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel2 Takes [%.3f] Seconds \n", duration.Seconds())

	startingTime = time.Now().UTC()
	experiments.ProveMembershipParallel3(t, setup.G, setup.N, rep, 4, 16)
	duration = time.Now().UTC().Sub(startingTime)
	fmt.Printf("Running ProveMembershipParallel3 Takes [%.3f] Seconds \n", duration.Seconds())
}

func testBigInt() {
	var temp, temp2 big.Int
	temp.SetInt64(1024)
	bytes := temp.Bytes()
	fmt.Println("byte[0] = ", bytes[0])
	fmt.Println("byte[1] = ", bytes[1])
	//fmt.Println("byte[2] = ", bytes[2])
	//fmt.Println("byte[4] = ", bytes[3])
	temp2.SetBytes(bytes)
	fmt.Println("temp2 = ", temp2.String())

	bits := temp.Bits()
	fmt.Println("bit[0] = ", bits[0])
	bits[0]++
	//fmt.Println("bit[1] = ", bits[1])
	temp2.SetBits(bits)
	fmt.Println("temp = ", temp.String())
	fmt.Println("temp2 = ", temp2.String())
}

func testExp() {
	setup := *accumulator.TrustedSetup()
	var ret1, ret2 big.Int
	ret1.Exp(setup.G, setup.G, setup.N)
	ret2 = *big.MultiExp(setup.G, setup.G, setup.G, setup.N)[0]
	fmt.Println("ret1 in main = ", ret1.String())
	fmt.Println("ret2 in main = ", ret2.String())
}

func main() {

	//testPreCompute()
	//testBigInt()
	testExp()
}
