package	types	// import "github.com/nathanaelle/useful.types"


import	(
	"testing"
	"flag"
	"fmt"
	"encoding/json"
)



type	valid_t struct {
	str	string
	data	interface{}
}





func Has_All_Interfaces(t *testing.T,usefultype interface{})  {
	if !Is_fmt_stringer(usefultype) {
		t.Errorf("not %20s  : %#v", "fmt.Stringer", usefultype)
	}
	if !Is_flag_value(usefultype) {
		t.Errorf("not %20s  : %#v", "flag.Value", usefultype)
	}
	if !Is_flag_getter(usefultype) {
		t.Errorf("not %20s  : %#v", "flag.Getter", usefultype)
	}
	if !Is_json_marshaler(usefultype) {
		t.Errorf("not %20s  : %#v", "json.Marshaler", usefultype)
	}
	if !Is_json_unmarshaler(usefultype) {
		t.Errorf("not %20s  : %#v", "json.Unmarshaler", usefultype)
	}
}




func Is_fmt_stringer(i interface{}) bool {
	_,ok := i.(fmt.Stringer)
	return ok
}

func Is_flag_value(i interface{}) bool {
	_,ok := i.(flag.Value)
	return ok
}

func Is_flag_getter(i interface{}) bool {
	_,ok := i.(flag.Getter)
	return ok
}

func Is_json_marshaler(i interface{}) bool {
	_,ok := i.(json.Marshaler)
	return ok
}

func Is_json_unmarshaler(i interface{}) bool {
	_,ok := i.(json.Unmarshaler)
	return ok
}
