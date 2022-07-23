package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/rsa_accumulator/proof"
)

func main() {
	bitLen := flag.Int("bit", 896, "bit length of the modulus")
	tries := flag.Int("try", 10, "number of tries")
	flag.Parse()
	f, err := os.OpenFile("test_"+strconv.Itoa(*bitLen)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	handleError(err)
	defer func(f *os.File) {
		err := f.Close()
		handleError(err)
	}(f)

	var totalTime float64
	for i := 0; i < *tries; i++ {
		fmt.Println("No. ", i)
		_, err = f.WriteString(time.Now().String() + "\n")
		handleError(err)
		target := randOddGen(*bitLen)
		handleError(err)
		_, err = f.WriteString(fmt.Sprintf("%d\n", target.BitLen()))
		handleError(err)
		_, err = f.WriteString(target.String() + "\n")
		handleError(err)
		start := time.Now()
		//fs, err := proof.UnconditionalLagrangeFourSquares(target)
		fs, err := proof.LagrangeFourSquares(target)
		handleError(err)
		currTime := time.Now()
		timeInterval := currTime.Sub(start)
		fmt.Println(timeInterval)
		totalTime += timeInterval.Seconds()
		secondsStr := fmt.Sprintf("%f", timeInterval.Seconds())
		_, err = f.WriteString(secondsStr + "\n")
		handleError(err)
		if ok := proof.Verify(target, fs); !ok {
			panic("verification failed")
		}
	}
	fmt.Printf("average: %f\n", totalTime/float64(*tries))
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func randOddGen(bitLen int) *big.Int {
	randLmt := new(big.Int).Lsh(big.NewInt(1), uint(bitLen-2))
	target, err := rand.Int(rand.Reader, randLmt)
	target.Lsh(target, 1)
	handleError(err)
	target.Add(target, big.NewInt(1))
	target.Add(target, new(big.Int).Lsh(big.NewInt(1), uint(bitLen-1)))
	return target
}
