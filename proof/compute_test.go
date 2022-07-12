package proof

import (
	"math/big"
	"reflect"
	"testing"
)

func Test_log10(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name    string
		args    args
		want    *big.Int
		wantErr bool
	}{
		{
			name: "log10_0",
			args: args{
				n: big.NewInt(0),
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "log10_1",
			args: args{
				n: big.NewInt(1),
			},
			want:    big.NewInt(0),
			wantErr: false,
		},
		{
			name: "log10_10",
			args: args{
				n: big.NewInt(10),
			},
			want:    big.NewInt(1),
			wantErr: false,
		},
		{
			name: "log10_20",
			args: args{
				n: big.NewInt(20),
			},
			want:    big.NewInt(1),
			wantErr: false,
		},
		{
			name: "log10_10000",
			args: args{
				n: big.NewInt(10000),
			},
			want:    big.NewInt(4),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := log10(tt.args.n)
			if (err != nil) != tt.wantErr {
				t.Errorf("log10() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("log10() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_log2(t *testing.T) {
	type args struct {
		n *big.Int
	}
	tests := []struct {
		name string
		args args
		want *big.Int
	}{
		{
			name: "test_0",
			args: args{
				n: big.NewInt(0),
			},
			want: big.NewInt(-1),
		},
		{
			name: "test_1",
			args: args{
				n: big.NewInt(1),
			},
			want: big.NewInt(0),
		},
		{
			name: "test_32",
			args: args{
				n: big.NewInt(32),
			},
			want: big.NewInt(5),
		},
		{
			name: "test_35",
			args: args{
				n: big.NewInt(35),
			},
			want: big.NewInt(5),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := log2(tt.args.n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("log2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_euclideanDivision(t *testing.T) {
	type args struct {
		a *big.Int
		b *big.Int
	}
	tests := []struct {
		name          string
		args          args
		wantQuotient  *big.Int
		wantRemainder *big.Int
		wantErr       bool
	}{
		{
			name: "test_5_3",
			args: args{
				a: big.NewInt(5),
				b: big.NewInt(3),
			},
			wantQuotient:  big.NewInt(1),
			wantRemainder: big.NewInt(2),
			wantErr:       false,
		},
		{
			name: "test_3_5",
			args: args{
				a: big.NewInt(3),
				b: big.NewInt(5),
			},
			wantQuotient:  big.NewInt(1),
			wantRemainder: big.NewInt(2),
		},
		{
			name: "test_5_0",
			args: args{
				a: big.NewInt(5),
				b: big.NewInt(0),
			},
			wantQuotient:  nil,
			wantRemainder: nil,
			wantErr:       true,
		},
		{
			name: "test_5_1",
			args: args{
				a: big.NewInt(5),
				b: big.NewInt(1),
			},
			wantQuotient:  big.NewInt(5),
			wantRemainder: big.NewInt(0),
			wantErr:       false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quotient, remainder, err := euclideanDivision(tt.args.a, tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("euclideanDivision() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(quotient, tt.wantQuotient) {
				t.Errorf("euclideanDivision() got quotient = %v, want %v", quotient, tt.wantQuotient)
			}
			if !reflect.DeepEqual(remainder, tt.wantRemainder) {
				t.Errorf("euclideanDivision() got remainder = %v, want %v", remainder, tt.wantRemainder)
			}
		})
	}
}

func Test_rangeDiv(t *testing.T) {
	type args struct {
		n      *big.Int
		numDIv int
	}
	tests := []struct {
		name    string
		args    args
		wantRes [][2]*big.Int
	}{
		{
			name: "test_1_1",
			args: args{
				n:      big.NewInt(1),
				numDIv: 1,
			},
			wantRes: [][2]*big.Int{
				{big.NewInt(0), big.NewInt(1)},
			},
		},
		{
			name: "test_10_3",
			args: args{
				n:      big.NewInt(10),
				numDIv: 3,
			},
			wantRes: [][2]*big.Int{
				{big.NewInt(0), big.NewInt(3)},
				{big.NewInt(3), big.NewInt(6)},
				{big.NewInt(6), big.NewInt(10)},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := rangeDiv(tt.args.n, tt.args.numDIv); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("rangeDiv() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
