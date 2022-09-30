package accumulator

import (
	"fmt"
	"strconv"
	"time"
)

// ManualBench is used for manual benchmark. Because the running time can be long, the golang benchmark may not work
func ManualBench(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = AccAndProve(set, DIHashFromPoseidon, &setup)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProve with set size %v\nTakes [%.3f] Seconds \n",
		testSetSize, duration.Seconds())
}

// ManualBenchZKAcc is used for manual benchmark. Because the running time can be long,
// the golang benchmark may not work
func ManualBenchZKAcc(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = ZKAccumulate(set, DIHashFromPoseidon, &setup)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProve with set size %v\nTakes [%.3f] Seconds \n",
		testSetSize, duration.Seconds())
}

// ManualBenchIter is used for manual benchmark. Because the running time can be long, the golang benchmark may not work
func ManualBenchIter(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = AccAndProveIter(set, DIHashFromPoseidon, &setup)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProveIter with set size %v\nTakes [%.3f] Seconds \n",
		testSetSize, duration.Seconds())
}

// ManualBenchParallel is used for manual benchmark. Because the running time can be long,
// the golang benchmark may not work
func ManualBenchParallel(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = AccAndProveParallel(set, DIHashFromPoseidon, &setup)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProveParallel with set size %v\nTakes [%.3f] Seconds \n",
		testSetSize, duration.Seconds())
}

// ManualBenchIterParallel is used for manual benchmark. Because the running time can be long,
// the golang benchmark may not work
func ManualBenchIterParallel(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = AccAndProveIterParallel(set, DIHashFromPoseidon, &setup)
	endingTime := time.Now().UTC()
	var duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProveIterParallel with set size %v\nTakes [%.3f] Seconds \n",
		testSetSize, duration.Seconds())
}

// GenBenchSet generate one set where every element is different
func GenBenchSet(num int) []string {
	ret := make([]string, num)
	for i := 0; i < num; i++ {
		ret[i] = strconv.Itoa(i)
	}
	return ret
}
