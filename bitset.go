package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"encoding/binary"
	"encoding/base64"
)


type BitSet struct {
  loglength	uint8
  set		[]uint64
}


// Create a bitset that contains 2^loglength bits
// Caution ! a loglength of 33 means 8Gibits so it uses 1Go of RAM
func NewBitSet(loglength uint8) *BitSet {
	if loglength <7 {
		return	&BitSet{ loglength, make([]uint64, 1) }
	}
	return	&BitSet{ loglength, make([]uint64, 1<<(loglength-6)) }
}


func (bs *BitSet)Set(data string) error {
	return bs.byte_set([]byte(data))
}


func (bs *BitSet)byte_set(data64 []byte) error {
	data	:= make([]byte,base64.URLEncoding.DecodedLen(len(data64)))
	_,err	:= base64.URLEncoding.Decode(data,data64)
	if err != nil {
		return nil
	}

	if data[0] <7 {
		bs.set	= make([]uint64, 1)
	} else {
		bs.set	= make([]uint64, 1<<(data[0]-6))
	}
	return	binary.Read(bytes.NewReader(data), binary.BigEndian, bs)
}


func (bs *BitSet)byte_get() ([]byte,error) {
	data	:= bytes.NewBuffer(make([]byte,0,1+(1<<(bs.loglength-3))))
	err	:= binary.Write(data, binary.BigEndian, bs)
	if err != nil {
		return []byte{},err
	}
	data64	:= make([]byte,base64.URLEncoding.EncodedLen(data.Len()))
	base64.URLEncoding.Encode(data64,data.Bytes())

	return data64,nil
}


func (bs *BitSet)UnmarshalTOML(data []byte) (err error) {
	return bs.byte_set(bytes.Trim(data,"\""))
}


func (bs BitSet)Get() interface{} {
	d,_ := bs.byte_get()
	return d
}


func (bs BitSet)String() string {
	d,_ := bs.byte_get()
	return string(d)
}


func (bs *BitSet)UnmarshalJSON(data []byte) (err error) {
	return bs.byte_set(bytes.Trim(data,"\""))
}


func (bs *BitSet)MarshalJSON() (data []byte,err error) {
	return []byte("\""+bs.String()+"\""),nil
}
