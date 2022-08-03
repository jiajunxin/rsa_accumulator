package proof

import (
	"github.com/rsa_accumulator/dihash"
	"lukechampine.com/frand"
	"math/big"
)

const (
	randGenSize = 10000
)

func rand256BitList(size int) (res [][]byte, err error) {
	res = make([][]byte, size)
	for i := 0; i < size; i++ {
		bytes := make([]byte, 256/8) // 1 byte = 8 bits
		_, err = frand.Read(bytes)
		if err != nil {
			return nil, err
		}
		res[i] = bytes
	}
	return
}

func diHashList(byteList [][]byte) (res []*big.Int) {
	res = make([]*big.Int, len(byteList))
	for idx, bt := range byteList {
		h := dihash.DIHash(bt)
		res[idx] = h
	}
	return
}

func prodDIHashList(hashes []*big.Int) *big.Int {
	res := big.NewInt(1)
	for _, h := range hashes {
		res.Mul(res, h)
	}
	return res
}
