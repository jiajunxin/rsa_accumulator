package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/rsa_accumulator/proof"
)

func main() {
	const bitLen = 500
	f, err := os.OpenFile("test_"+strconv.Itoa(bitLen)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	for i := 0; i < 1; i++ {
		_, err = f.WriteString(time.Now().String() + "\n")
		if err != nil {
			panic(err)
		}
		randLmt := new(big.Int).Exp(big.NewInt(2), big.NewInt(bitLen), nil)
		target, err := rand.Int(rand.Reader, randLmt)
		if err != nil {
			panic(err)
		}
		fmt.Println(target)
		_, err = f.WriteString(target.String() + "\n")
		if err != nil {
			panic(err)
		}
		start := time.Now()
		fs, err := proof.LagrangeFourSquares(target)
		if err != nil {
			panic(err)
		}
		currTime := time.Now()
		timeInterval := currTime.Sub(start)
		fmt.Println(timeInterval)
		secondsStr := fmt.Sprintf("%f", timeInterval.Seconds())
		_, err = f.WriteString(secondsStr + "\n")
		if err != nil {
			panic(err)
		}
		ok := proof.Verify(target, fs)
		if ok {
			fmt.Println("verification succeeded")
			_, err := f.WriteString("verification succeeded\n")
			if err != nil {
				panic(err)
			}
		} else {
			fmt.Println("verification failed")
			_, err := f.WriteString("verification failed\n")
			if err != nil {
				panic(err)
			}
		}
	}
}
