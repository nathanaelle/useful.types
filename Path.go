package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"os"
)




type	Path	string

func (d *Path)Get() interface{} {
	return string(*d)
}

func (d *Path) String() string {
	return string(*d)
}

func (d *Path)UnmarshalJSON(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *Path)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}

func (d *Path) UnmarshalTOML(data []byte) error  {
	return d.Set(string(bytes.Trim(data,"\"")))
}


func (d *Path) Set(data string) (err error) {
	_, err = os.Stat(data)
	if err == nil {
		*d = Path(data)
	}
	return
}




/*
type	PathList	[]Path


func (d *PathList) Set(data string) (err error) {
	_, err = os.Stat(data)
	if err == nil {
		*d = append( *d, Path(data) )
	}
	return
}


func (d *PathList) String() string {
	return fmt.Sprintf("%s", *d)
}
*/
