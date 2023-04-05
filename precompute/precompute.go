package precompute

import (
	"context"
	"math/big"
)

// const byteChunkSize = 125000

// Table is the precomputing table
type Table struct {
	g             *big.Int
	n             *big.Int
	byteChunkSize int
	table         []*big.Int
}

// NewTable creates a new precomputing table
func NewTable(g, n, elementUpperBound *big.Int, numElements uint64, byteChunkSize int) *Table {
	t := &Table{
		g:             g,
		n:             n,
		byteChunkSize: byteChunkSize,
	}
	maxBitLen := elementUpperBound.BitLen() * int(numElements)
	numByteChunks := maxBitLen / (t.byteChunkSize * 8)
	t.table = make([]*big.Int, numByteChunks)
	t.table[0] = new(big.Int).Set(g)
	opt := new(big.Int).Lsh(big1, uint(t.byteChunkSize*8))
	for i := 1; i < numByteChunks; i++ {
		t.table[i] = new(big.Int).Exp(t.table[i-1], opt, n)
	}
	return t
}

// Compute computes the result of base^x mod n with specified number of goroutines
func (t *Table) Compute(x *big.Int, numRoutine int) *big.Int {
	xBytes := x.Bytes()
	inputChan := make(chan input, numRoutine<<2)
	outputChan := make(chan *big.Int, numRoutine)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for i := 0; i < numRoutine; i++ {
		go t.routineCompute(ctx, t.n, xBytes, inputChan, outputChan)
	}
	resChan := make(chan *big.Int)
	go func() {
		res := big.NewInt(1)
		counter := len(xBytes) / t.byteChunkSize
		if len(xBytes)%t.byteChunkSize != 0 {
			counter++
		}
		for out := range outputChan {
			res.Mul(res, out)
			res.Mod(res, t.n)
			counter--
			if counter == 0 {
				resChan <- res
				return
			}
		}
	}()
	for i := len(xBytes); i > 0; i -= t.byteChunkSize {
		right := i
		left := right - t.byteChunkSize
		if left < 0 {
			left = 0
		}
		idx := (len(xBytes) - i) / t.byteChunkSize
		inputChan <- input{
			left:     left,
			right:    right,
			tableIdx: idx,
		}
	}
	return <-resChan
}

type input struct {
	left     int
	right    int
	tableIdx int
}

func (t *Table) routineCompute(ctx context.Context, n *big.Int, xBytes []byte,
	inputChan chan input, resChan chan *big.Int) {
	opt := new(big.Int)
	for {
		select {
		case <-ctx.Done():
			return
		case in := <-inputChan:
			opt.SetBytes(xBytes[in.left:in.right])
			res := new(big.Int).Exp(t.table[in.tableIdx], opt, n)
			resChan <- res
		default:
			continue
		}
	}
}
