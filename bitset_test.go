package	types	// import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_BitSet(t *testing.T)  {
	t.Skip()
	l_inval	:= []string{
		"1",
		"01",
		"ASZA=ERZ",
		"ASZA+ERZ",
	}

	l_val	:= []valid_t{
		valid_t{ "", BitSet{} },
		valid_t{ "0==", BitSet{} },
		valid_t{ "10=", BitSet{ 1, []uint64{ 0 }} },
		valid_t{ "101", BitSet{ 1, []uint64{ 1 }} },
		valid_t{ "111", BitSet{ 2, []uint64{ 1 }} },
		valid_t{ "1111", BitSet{ 9, []uint64{ 0xffffff, 0xf }} },
		valid_t{ "00", BitSet{} },
		valid_t{ "10", BitSet{} },
		valid_t{ "01", BitSet{} },
		valid_t{ "10", BitSet{} },
	}

	d	:= new(BitSet)
	Has_All_Interfaces(t,d)

	for _,inv := range l_inval {
		err	:= d.Set(inv)
		if err == nil {
			t.Logf("[%+v]", d)
			t.Errorf("[%v] parser invalid", inv)
		}
		t.Logf("[%s] %v", inv, err)
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
