package dihash

import (
	"math/big"
	"math/rand"
	"testing"
)

func Test2048Rnd(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	var temp big.Int
	temp = *get2048Rnd(rng)
	size := temp.BitLen()
	if size <= 2028 {
		t.Errorf("the size of 2048 rnd is smaller then 2048, it is %d", size)
	}
}
