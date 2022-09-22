package main

import (
	"crypto/rand"
	"flag"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"time"

	"github.com/jiajunxin/rsa_accumulator/proof"
)

func main() {
	bitLen := flag.Int("bit", 1792, "bit length of the modulus")
	tries := flag.Int("try", 1000, "number of tries")
	flag.Parse()
	f, err := os.OpenFile("test_"+strconv.Itoa(*bitLen)+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	handleError(err)
	defer func(f *os.File) {
		err := f.Close()
		handleError(err)
	}(f)

	//randLmt := new(big.Int).Lsh(big.NewInt(1), uint(*bitLen))
	var totalTime float64
	for i := 0; i < *tries; i++ {
		_, err = f.WriteString(time.Now().String() + "\n")
		handleError(err)
		randLmt := new(big.Int).Lsh(big.NewInt(1), uint(*bitLen-2))
		target := randGen(randLmt)
		_, err = f.WriteString(fmt.Sprintf("%d\n", target.BitLen()))
		handleError(err)
		_, err = f.WriteString(target.String() + "\n")
		handleError(err)
		start := time.Now()
		ts, err := proof.ThreeSquares(target)
		handleError(err)
		currTime := time.Now()
		timeInterval := currTime.Sub(start)
		fmt.Println(i, timeInterval)
		totalTime += timeInterval.Seconds()
		secondsStr := fmt.Sprintf("%f", timeInterval.Seconds())
		_, err = f.WriteString(secondsStr + "\n")
		handleError(err)
		if ok := proof.Verify(target, ts); !ok {
			fmt.Println(target)
			fmt.Println(ts)
			panic("verification failed")
		}
	}
	fmt.Printf("average: %f\n", totalTime/float64(*tries))
	//n := new(big.Int)
	//n.SetString(accumulator.N2048String, 10)
	//g := new(big.Int)
	//g.SetString(accumulator.G2048String, 10)
	//h := new(big.Int)
	//h.SetString(accumulator.H2048String, 10)
	//pp := proof.NewPublicParameters(n, g, h)
	//u := big.NewInt(123)
	//x := big.NewInt(3)
	//w := new(big.Int)
	//w.Exp(u, x, nil)
	//prover := proof.NewZKPoKEProver(pp)
	//verifier := proof.NewZKPoKEVerifier(pp)
	//pf, err := prover.Prove(u, x)
	//handleError(err)
	//ok, err := verifier.Verify(pf, u, w)
	//handleError(err)
	//if !ok {
	//	panic("verification failed")
	//}
}

func handleError(err error) {
	if err != nil {
		panic(err)
	}
}

func randGen(randLmt *big.Int) *big.Int {
	x, err := rand.Int(rand.Reader, randLmt)
	handleError(err)
	x.Lsh(x, 2)
	x.Add(x, big.NewInt(1))
	return x
}
