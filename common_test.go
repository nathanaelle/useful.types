package types // import "github.com/nathanaelle/useful.types"

import (
	"encoding"
	"flag"
	"fmt"
	"testing"
)

type valid_t struct {
	str  string
	data interface{}
}

func Has_All_Interfaces(t *testing.T, usefultype interface{}) {
	if !Is_fmt_stringer(usefultype) {
		t.Errorf("not %25s  : %#v", "fmt.Stringer", usefultype)
	}

	if !Is_flag_value(usefultype) {
		t.Errorf("not %25s  : %#v", "flag.Value", usefultype)
	}
	if !Is_flag_getter(usefultype) {
		t.Errorf("not %25s  : %#v", "flag.Getter", usefultype)
	}

	if !Is_encoding_textmarshaler(usefultype) {
		t.Errorf("not %25s  : %#v", "encoding.TextMarshaler", usefultype)
	}
	if !Is_encoding_textunmarshaler(usefultype) {
		t.Errorf("not %25s  : %#v", "encoding.TextUnmarshaler", usefultype)
	}
}

func Is_encoding_textmarshaler(i interface{}) bool {
	_, ok := i.(encoding.TextMarshaler)
	return ok
}

func Is_encoding_textunmarshaler(i interface{}) bool {
	_, ok := i.(encoding.TextUnmarshaler)
	return ok
}

func Is_fmt_stringer(i interface{}) bool {
	_, ok := i.(fmt.Stringer)
	return ok
}

func Is_flag_value(i interface{}) bool {
	_, ok := i.(flag.Value)
	return ok
}

func Is_flag_getter(i interface{}) bool {
	_, ok := i.(flag.Getter)
	return ok
}
