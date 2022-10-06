package precompute

import (
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

const (
	testSize           = 1000
	smallByteChunkSize = 1

	testByteChunkSize = 512
)

var (
	accSetup    *accumulator.Setup
	set         []string
	reps        []*big.Int
	repProd     *big.Int
	acc         *big.Int
	onceSetup   sync.Once
	onceSet     sync.Once
	onceReps    sync.Once
	onceRepProd sync.Once
	onceAcc     sync.Once
)

func getSetup() *accumulator.Setup {
	onceSetup.Do(func() {
		accSetup = accumulator.TrustedSetup()
	})
	return accSetup
}

func getSet() []string {
	onceSet.Do(func() {
		set = accumulator.GenBenchSet(testSize)
	})
	return set
}

func getRepresentations() []*big.Int {
	onceReps.Do(func() {
		reps = accumulator.HashEncode(getSet(), accumulator.EncodeTypePoseidonDIHash)
	})
	return reps
}

func getRepProd() *big.Int {
	onceRepProd.Do(func() {
		repProd = accumulator.SetProductRecursive(getRepresentations())
	})
	return repProd
}

func getAcc() *big.Int {
	onceAcc.Do(func() {
		setup := getSetup()
		set := getRepresentations()
		acc = new(big.Int).Set(setup.G)
		for _, v := range set {
			acc.Exp(acc, v, setup.N)
		}
	})
	return acc
}

func getSmallSetup() *accumulator.Setup {
	return &accumulator.Setup{
		G: big.NewInt(2),
		N: big.NewInt(1000003),
	}
}

func getSmallReps() []*big.Int {
	return []*big.Int{
		big.NewInt(21),
		big.NewInt(32),
		big.NewInt(15),
		big.NewInt(17),
	}
}

func getSmallRepProd() *big.Int {
	return big.NewInt(171360)
}

func getSmallAcc() *big.Int {
	setup := getSmallSetup()
	reps := getSmallReps()
	return accumulate(setup, reps)
}

func TestTable_Compute(t1 *testing.T) {
	type args struct {
		x          *big.Int
		numRoutine int
	}
	tests := []struct {
		name       string
		setupTable func() *Table
		args       args
		want       *big.Int
	}{
		{
			name: "TestTable_Compute_small",
			setupTable: func() *Table {
				setup := getSmallSetup()
				t := NewTable(setup.G, setup.N, 6, len(getSmallReps()), smallByteChunkSize)
				return t
			},
			args: args{
				x:          getSmallRepProd(),
				numRoutine: 4,
			},
			want: getSmallAcc(),
		},
		{
			name: "TestTable_Compute",
			setupTable: func() *Table {
				setup := getSetup()
				t := NewTable(setup.G, setup.N, 2048, testSize, testByteChunkSize)
				return t
			},
			args: args{
				x:          getRepProd(),
				numRoutine: 4,
			},
			want: getAcc(),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.setupTable()
			if got := t.Compute(tt.args.x, tt.args.numRoutine); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Compute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func accumulate(setup *accumulator.Setup, reps []*big.Int) *big.Int {
	acc := new(big.Int).Set(setup.G)
	for _, v := range reps {
		acc.Exp(acc, v, setup.N)
	}
	return acc
}

func BenchmarkAccumulate(b *testing.B) {
	setup := getSetup()
	reps := getRepresentations()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		accumulate(setup, reps)
	}
}

func BenchmarkPrecompute(b *testing.B) {
	setup := getSetup()
	reps := getRepresentations()
	t := NewTable(setup.G, setup.N, 2048, testSize, 1024)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		repProd := accumulator.SetProductRecursive(reps)
		t.Compute(repProd, 4)
	}
}

func TestComputeFromTable(t1 *testing.T) {
	setSize := 1000
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.HashEncode(set, accumulator.EncodeTypePoseidonDIHash)
	prod := accumulator.SetProductRecursive(rep)
	originalResult := accumulator.AccumulateNew(setup.G, prod, setup.N)

	table := GenPreTable(setup.G, setup.N, 10000, 100)
	result := ComputeFromTable(table, prod, setup.N)
	if result.Cmp(originalResult) != 0 {
		t1.Errorf("wrong result")

	}
}
