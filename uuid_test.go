package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_UUID(t *testing.T) {
	l_inval := []string{
		"de305d54-75b4f-431b-adb2-eb6b9e546014",
	}

	l_val := []valid_t{
		valid_t{"de305d54-75b4-431b-adb2-eb6b9e546014", UUID([16]byte{0xde, 0x30, 0x5d, 0x54, 0x75, 0xb4, 0x43, 0x1b, 0xad, 0xb2, 0xeb, 0x6b, 0x9e, 0x54, 0x60, 0x14})},
	}

	d := new(UUID)
	Has_All_Interfaces(t, d)

	for _, inv := range l_inval {
		err := d.Set(inv)
		if err == nil {
			t.Errorf("[%v] parser invalid", inv)
		}
	}

	for _, val := range l_val {
		d := new(UUID)
		data := val.data.(UUID)
		err := d.Set(val.str)
		if err != nil {
			t.Errorf("[%v] parser invalid [%v]", val.str, err)
		}
		if !data.IsValid() {
			t.Errorf("[%v] invalid source", data)
		}

		if !d.IsValid() {
			t.Errorf("[%v] invalid parse", d)
		}

		if val.str != d.String() {
			t.Errorf("[%v] [%v] differs", val.str, d)
		}

		if data.String() != d.String() {
			t.Errorf("[%v] [%v] differs", data, d)
		}
		if !data.IsEqual(*d) {
			t.Errorf("[%v] [%v] differs", data, d)
		}
	}

	uu, err := NewUUID(UUIDv4)
	if err != nil {
		t.Errorf("[%v] generator invalid", uu)
	}

	if !uu.IsValid() {
		t.Errorf("[%v] invalid parse", uu)
	}

	uu, err = NewUUID(UUIDv1_timestamp)
	if err != nil {
		t.Errorf("[%v] generator invalid", uu)
	}

	if !uu.IsValid() {
		t.Errorf("[%v] invalid parse", uu)
	}

}

func Benchmark_UUID_Set(b *testing.B) {
	d := new(UUID)
	for i := 0; i < b.N; i++ {
		d.Set("de305d54-75b4-431b-adb2-eb6b9e546014")
	}
}

func Benchmark_UUID_String(b *testing.B) {
	d := new(UUID)
	d.Set("de305d54-75b4-431b-adb2-eb6b9e546014")
	for i := 0; i < b.N; i++ {
		d.String()
	}
}

func Benchmark_UUIDv1_TS(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID(UUIDv1_timestamp)
	}
}

func Benchmark_UUIDv1MacRand(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID(UUIDv1MacRand)
	}
}

func Benchmark_UUIDv4(b *testing.B) {
	for i := 0; i < b.N; i++ {
		NewUUID(UUIDv4)
	}
}
