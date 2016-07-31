package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_Path(t *testing.T) {
	l_inval := []string{
		"_test_data/dontexist",
		"_test_data/dontexist/",
		"_test_data/exist/no-file",
	}

	l_val := []valid_t{
		valid_t{"_test_data/exist/", Path("_test_data/exist/")},
		valid_t{"_test_data/exist", Path("_test_data/exist")},
		valid_t{"_test_data/exist/some-file", Path("_test_data/exist/some-file")},
	}

	d := new(Path)
	Has_All_Interfaces(t, d)

	for _, inv := range l_inval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range l_val {
		d := new(Path)
		data := val.data.(Path)
		err := d.Set(val.str)
		if err != nil {
			t.Errorf("[%v] parser invalid [%v]", val.str, err)
		}

		if data.String() != d.String() {
			t.Errorf("[%v] [%v] differs", data, d)
		}
	}
}

func Benchmark_Path_Set(b *testing.B) {
	d := new(Path)
	for i := 0; i < b.N; i++ {
		d.Set("_test_data/exist/some-file")
	}
}

func Benchmark_Path_String(b *testing.B) {
	d := new(Path)
	d.Set("_test_data/exist/some-file")
	for i := 0; i < b.N; i++ {
		d.String()
	}
}

func Benchmark_Path_MarshalText(b *testing.B) {
	d := new(Path)
	d.Set("_test_data/exist/some-file")
	for i := 0; i < b.N; i++ {
		d.MarshalText()
	}
}
