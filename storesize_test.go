package	types	// import "github.com/nathanaelle/useful.types"

import (
	"testing"
)



func Test_StoreSize(t *testing.T)  {
	l_inval	:= []string{
		"0.12",
		"M0o",
		"12Mb",
		"12yo",
		"1.2Mo",
		"1ki",
	}

	l_val	:= []valid_t{
		valid_t{ "123M", StoreSize(123000000) },
		valid_t{ "1kiB", StoreSize(1024) },
		valid_t{ "1ko", StoreSize(1000) },
		valid_t{ "100", StoreSize(100) },
		valid_t{ "0", StoreSize(0) },
		valid_t{ "0o", StoreSize(0) },
		valid_t{ "0B", StoreSize(0) },
		valid_t{ "0kB", StoreSize(0) },
		valid_t{ "12345To", StoreSize(12345000000000000) },
		valid_t{ "12345Tio", StoreSize(13573471044894720) },
	}

	d	:= new(StoreSize)
	Has_All_Interfaces(t,d)

	for _,inv := range l_inval {
		err	:= d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _,val := range l_val {
		d	:= new(StoreSize)
		data	:= val.data.(StoreSize)
		err	:= d.Set(val.str)
		if err != nil {
			t.Errorf("[%v] parser invalid : [%v]", val.str, err)
		}

		if data.String() != d.String() {
			t.Errorf("source=[%v] parsed=[%v] differs", data, d)
		}
	}
}


func Benchmark_StoreSize_Set(b *testing.B) {
	d	:= new(StoreSize)
	for i := 0; i < b.N; i++ {
		d.Set("12345Tio")
	}
}
