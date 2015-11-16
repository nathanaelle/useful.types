package	types	// import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_BitSet(t *testing.T)  {
	t.Logf("Please don't use BitSet")

	b_val	:= []BitSet{
		BitSet{},
		BitSet{0,[]uint64{}},
		BitSet{1,[]uint64{0}},
		BitSet{1,[]uint64{1}},
		BitSet{1,[]uint64{10}},
		BitSet{1,[]uint64{64}},
		BitSet{1,[]uint64{0xffffffffffffffff}},
		BitSet{2,[]uint64{0,0}},
		BitSet{2,[]uint64{0xffffffffffffffff, 0x1}},
	}


	l_inval	:= []string{
		"1",
		"01",
		"ASZA=ERZ",
		"ASZA+ERZ",
	}

	l_val	:= []valid_t{
		valid_t{ "AAAAAAAAAAA=", BitSet{} },
		valid_t{ "AAAAAAAAAAEAAAAAAAAAAA==", BitSet{ 1, []uint64{ 0 }} },
		valid_t{ "AAAAAAAAAAEAAAAAAAAAAQ==", BitSet{ 1, []uint64{ 1 }} },
		valid_t{ "AAAAAAAAAAEAAAAAAAAACg==", BitSet{ 1, []uint64{ 10 }} },
		valid_t{ "AAAAAAAAAAEAAAAAAAAAQA==", BitSet{ 1, []uint64{ 64 }} },
		valid_t{ "AAAAAAAAAAH__________w==", BitSet{ 1, []uint64{ 0xffffffffffffffff }} },
		valid_t{ "AAAAAAAAAAIAAAAAAAAAAAAAAAAAAAAA", BitSet{ 2, []uint64{ 0,0 }} },
		valid_t{ "AAAAAAAAAAL__________wAAAAAAAAAB", BitSet{ 2, []uint64{ 0xffffffffffffffff, 0x1 }} },
	}

	d	:= new(BitSet)
	Has_All_Interfaces(t,d)

	for i,_ := range b_val {
		t.Logf("%d: [%s]\n", i, b_val[i].String())
	}



	for _,inv := range l_inval {
		err	:= d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid : %s", inv, d.String())
		}
		//t.Logf("[%s] %v", inv, err)
	}

	for _,val := range l_val {
		d	:= new(BitSet)
		data	:= val.data.(BitSet)
		err	:= d.Set(val.str)
		if err != nil {
			t.Logf("[%v]", err)
			t.Errorf("[%v] parser invalid", val.str)
		}

		if data.String() != d.String() {
			t.Errorf("[%v] [%v] differs", data, d)
		}
	}
}


func Test_BitSet_BitOperation(t *testing.T) {
	t.Logf("Please don't use BitSet")
}
