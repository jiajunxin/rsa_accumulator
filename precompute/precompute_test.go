package precompute

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"reflect"
	"sync"
	"testing"

	"github.com/jiajunxin/rsa_accumulator/accumulator"
)

const testSize = 4

var (
	accSetup     *accumulator.Setup
	set          []string
	reps         []*big.Int
	repProd      *big.Int
	acc          *big.Int
	randReps     []*big.Int
	randAcc      *big.Int
	onceSetup    sync.Once
	onceSet      sync.Once
	onceReps     sync.Once
	onceRepProd  sync.Once
	onceAcc      sync.Once
	onceRandReps sync.Once
	onceRandAcc  sync.Once
	randLmt      = new(big.Int).Lsh(big.NewInt(1), 8)
)

func getSetup() *accumulator.Setup {
	onceSetup.Do(func() {
		//accSetup = accumulator.TrustedSetup()
		accSetup = &accumulator.Setup{
			G: big.NewInt(11),
			N: big.NewInt(23),
		}
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

func getRandomRepresentations() []*big.Int {
	onceRandReps.Do(func() {
		randReps = make([]*big.Int, testSize)
		var err error
		for i := 0; i < testSize; i++ {
			randReps[i], err = crand.Int(crand.Reader, randLmt)
			if err != nil {
				panic(err)
			}
			fmt.Println(randReps[i])
		}
	})
	return randReps
}

func getRandomAcc() *big.Int {
	onceRandAcc.Do(func() {
		randAcc = new(big.Int).Set(getSetup().G)
		for _, v := range getRandomRepresentations() {
			fmt.Println(v)
			randAcc.Exp(randAcc, v, getSetup().N)
		}
	})
	return randAcc
}

func getRandRepProd() *big.Int {
	randRepProd := big.NewInt(1)
	for _, v := range getRandomRepresentations() {
		randRepProd.Mul(randRepProd, v)
	}
	return randRepProd
}

func TestTable_Compute(t1 *testing.T) {
	type args struct {
		x          *big.Int
		numRoutine int
	}
	tests := []struct {
		name  string
		setup func() *Table
		args  args
		want  *big.Int
	}{
		{
			name: "TestTable_Compute",
			setup: func() *Table {
				setup := getSetup()
				t := NewTable(setup.G, setup.N, randLmt, testSize)
				return t
			},
			args: args{
				x:          getRandRepProd(),
				numRoutine: 4,
			},
			want: getRandomAcc(),
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := tt.setup()
			if got := t.Compute(tt.args.x, tt.args.numRoutine); !reflect.DeepEqual(got, tt.want) {
				t1.Errorf("Compute() = %v, want %v", got, tt.want)
			}
		})
	}
}
