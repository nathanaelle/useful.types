package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"net/url"
)


// Wrapper Type for url.URL providing missing interfaces
// All the parsing and the validation are done by url.Parse
type	URL	url.URL


func (d *URL)Set(data string) (err error) {
	dest, err := url.Parse(data)
	if err == nil {
		*d = URL(*dest)
	}
	return err
}

func (d *URL)Get() interface{} {
	return url.URL(*d)
}

func (d *URL)UnmarshalTOML(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *URL)String() string {
	var u_d	*url.URL
	*u_d = url.URL(*d)
	return u_d.String()
}

func (d *URL)UnmarshalJSON(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *URL)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}
