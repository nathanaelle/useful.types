package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
	"time"
)

func Test_Duration(t *testing.T) {
	l_inval := []string{
		"2h45j",
	}

	l_val := []valid_t{
		valid_t{"2h45m", Duration(165 * time.Minute)},
		valid_t{"127µs", Duration(127 * time.Microsecond)},
	}

	d := new(Duration)
	Has_All_Interfaces(t, d)

	for _, inv := range l_inval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range l_val {
		d := new(Duration)
		data := val.data.(Duration)
		err := d.Set(val.str)
		if err != nil {
			t.Errorf("[%v] parser invalid", val.str)
		}

		if data.String() != d.String() {
			t.Errorf("[%v] [%v] differs", data, d)
		}
	}
}
