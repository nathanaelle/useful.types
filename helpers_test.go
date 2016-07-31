package types // import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"testing"
)

type hsb_test struct {
	value  float64
	base   float64
	unit   []byte
	result []byte
}

func Test_HumanScale(t *testing.T) {
	tests := []hsb_test{
		{0, 1000, []byte("m"), []byte("0m")},
		{1000, 1000, []byte("m"), []byte("1km")},
		{1000000, 1000, []byte("m"), []byte("1Mm")},
		{0.001, 1000, []byte("m"), []byte("1mm")},
	}

	for i, cas := range tests {
		res_b := HumanScaleBytes(cas.value, cas.base, cas.unit)
		if !bytes.Equal(cas.result, res_b) {
			t.Errorf("t_%02d : [%v] expected, got [%v]", i, cas.result, res_b)
		}
		res_s := HumanScaleString(cas.value, cas.base, string(cas.unit))

		if res_s != string(cas.result) {
			t.Errorf("t_%02d : [%v] expected, got [%v]", i, string(cas.result), res_s)
		}

	}

}

func Benchmark_HumanScaleBytes(b *testing.B) {
	unit := []byte("iB")
	for i := 0; i < b.N; i++ {
		HumanScaleBytes(13573471044894720, 1024, unit)
	}
}

func Benchmark_HumanScaleString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		HumanScaleString(13573471044894720, 1024, "iB")
	}
}
