package accumulator

import (
	"fmt"
	"strconv"
	"time"
)

// ManualBench is used for manul benchmark. Because the running time can be long, the golang benchmark may not work
func ManualBench(testSetSize int) {
	set := GenBenchSet(testSetSize)
	setup := *TrustedSetup()
	startingTime := time.Now().UTC()
	_, _ = AccAndProve(set, HashToPrimeFromSha256, &setup)
	endingTime := time.Now().UTC()
	var duration time.Duration = endingTime.Sub(startingTime)
	fmt.Printf("Running AccAndProve with set size %v\n Takes [%.3f] Seconds \n", testSetSize, duration.Seconds())
}

// GenBenchSet generate one set where every element is identical
func GenBenchSet(num int) []string {
	ret := make([]string, num)
	for i := 0; i < num; i++ {
		ret[i] = strconv.Itoa(i)
	}
	return ret
}
