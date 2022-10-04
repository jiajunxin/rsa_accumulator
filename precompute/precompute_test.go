package precompute

import (
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

const testSize = 128

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
		reps = accumulator.GenRepresentatives(getSet(), accumulator.DIHashFromPoseidon)
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
			name: "TestTable_Compute",
			setupTable: func() *Table {
				setup := getSetup()
				elemUpperBound := new(big.Int).Lsh(big.NewInt(1), 2048)
				elemUpperBound.Sub(elemUpperBound, big.NewInt(1))
				t := NewTable(setup.G, setup.N, elemUpperBound, testSize)
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

func TestComputeFromTable(t1 *testing.T) {
	setSize := 100
	set := accumulator.GenBenchSet(setSize)
	setup := *accumulator.TrustedSetup()
	rep := accumulator.GenRepresentatives(set, accumulator.DIHashFromPoseidon)
	prod := accumulator.SetProductRecursive(rep)
	originalResult := accumulator.AccumulateNew(setup.G, prod, setup.N)

	table := GenPreTable(setup.G, setup.N, 100, 10)
	result := ComputeFromTable(table, prod, setup.N)
	if result.Cmp(originalResult) != 0 {
		t1.Errorf("wrong result")
	}
}
