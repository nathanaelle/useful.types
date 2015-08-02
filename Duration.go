package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"time"
)





/*
 *	type:		Duration
 *	content:	time duration aka intergers with time units
 */
type Duration time.Duration

func (d *Duration)Set(data string) (err error) {
	return d.byte_set([]byte(data))
}

func (d *Duration)byte_set(data []byte) (err error) {
	tmp_d, err := time.ParseDuration(string(data))
	if err == nil {
		*d = Duration(tmp_d)
	}
	return err
}


func (d *Duration)Get() interface{} {
	return time.Duration(*d)
}

func (d *Duration)UnmarshalTOML(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *Duration)String() string {
	return time.Duration(*d).String()
}

func (d *Duration)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *Duration)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}
