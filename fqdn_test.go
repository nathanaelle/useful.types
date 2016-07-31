package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_FQDN(t *testing.T) {
	l_inval := []string{
		"-",
		"-example.invalid.",
		"example-.invalid",
		"example..invalid",
		"www.-example.invalid",
	}

	l_val := []valid_t{
		valid_t{".", FQDN(".")},
		valid_t{"invalid", FQDN("invalid")},
		valid_t{"invalid.", FQDN("invalid.")},
		valid_t{".invalid", FQDN(".invalid")},
		valid_t{".invalid.", FQDN(".invalid.")},
		valid_t{"www.example.", FQDN("www.example.")},
		valid_t{"www.sub-zone.example.", FQDN("www.sub-zone.example.")},
	}

	d := new(FQDN)
	Has_All_Interfaces(t, d)

	for _, inv := range l_inval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range l_val {
		d := new(FQDN)
		data := val.data.(FQDN)
		err := d.Set(val.str)
		if err != nil {
			t.Errorf("[%v] parser invalid", val.str)
		}

		if data.String() != d.String() {
			t.Errorf("[%v] [%v] differs", &data, d)
		}
	}
}

func Benchmark_FQDN_Set(b *testing.B) {
	d := new(FQDN)
	for i := 0; i < b.N; i++ {
		d.Set("www.sub-zone.example.")
	}
}
