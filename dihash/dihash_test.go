package dihash

import (
	"math/big"
	"math/rand"
	"testing"
)

func TestGcd(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var temp big.Int
	temp = *getDIHash(rng)
	flag := temp.Bit(0)
	if flag == 0 {
		t.Errorf("the last bit of DI hash is 0")
	}

	i := 0
	for i < 100 {
		temp = *getDIHash(rng)
		flag = temp.Bit(0)
		if flag == 0 {
			t.Errorf("the last bit of DI hash is 0")
		}
		i++
	}
}

func Test2048Rnd(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var temp big.Int
	temp = *get2048Rnd(rng)
	size := temp.BitLen()
	if size <= 2028 {
		t.Errorf("the size of 2048 rnd is to small then 2048, it is %d", size)
	}
}
