package dihash

import (
	"math/big"
	"math/rand"
	"reflect"
	"testing"

	"lukechampine.com/frand"
)

const (
	randGenSize = 30000
)

func Test2048Rnd(t *testing.T) {
	rng := rand.New(rand.NewSource(99))
	temp := *get2048Rnd(rng)
	size := temp.BitLen()
	if size <= 2028 {
		t.Errorf("the size of 2048 rnd is smaller then 2048, it is %d", size)
	}
}

func TestProd(t *testing.T) {
	hashes := []*big.Int{big.NewInt(1), big.NewInt(2), big.NewInt(3)}
	res := Prod(hashes)
	if res.Cmp(big.NewInt(6)) != 0 {
		t.Errorf("the product of [1, 2, 3] is 6, but it is %d", res)
	}
}

func TestProdParallel(t *testing.T) {
	type args struct {
		hashes []*big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "test_[1, 2, 3]",
			args: args{
				hashes: []*big.Int{
					big.NewInt(1), big.NewInt(2), big.NewInt(3),
				},
			},
			want: big.NewInt(6),
		},
		{
			name: "test_[1, 2, 3, 4, 5, 6, 7, 8, 9, 10]",
			args: args{
				hashes: []*big.Int{
					big.NewInt(1), big.NewInt(2), big.NewInt(3),
					big.NewInt(4), big.NewInt(5), big.NewInt(6),
					big.NewInt(7), big.NewInt(8), big.NewInt(9),
					big.NewInt(10),
				},
			},
			want: big.NewInt(3628800),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ProdParallel(tt.args.hashes); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ProdParallel() = %v, want %v", got, tt.want)
			}
		})
	}
}

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
		h := DIHash(bt)
		res[idx] = h
	}
	return
}

func BenchmarkProdParallel(b *testing.B) {
	randBitList, err := rand256BitList(randGenSize)
	if err != nil {
		b.Error(err)
	}
	hashes := diHashList(randBitList)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ProdParallel(hashes)
	}
}
