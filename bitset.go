package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"errors"
	"encoding/binary"
	"encoding/base64"
)


type BitSet struct {
  length	uint64
  set		[]uint64
}



func u64len(l uint64) uint64 {
	if l%64 == 0 {
		return uint64(l/64)
	}
	return uint64(l/64)+1
}



func NewBitSet(length uint64) *BitSet {
	l64 := u64len(length)

	return	&BitSet{ length, make([]uint64, l64) }
}


func (bs *BitSet)Set(data string) error {
	return bs.byte_set([]byte(data))
}


func (bs *BitSet)byte_set(data64 []byte) (err error) {
	data	:= make([]byte,base64.URLEncoding.DecodedLen(len(data64)))
	_,err	= base64.URLEncoding.Decode(data,data64)
	if err != nil {
		return
	}

	if len(data) == 0 {
		bs.length = 0
		bs.set = []uint64{}
		return nil
	}

	reader	:= bytes.NewReader(data)
	length	:= uint64(0)

	err	= binary.Read(reader, binary.BigEndian, &length)
	if err != nil {
		return
	}

	if uint64(len(data)) < (8*(length+1)) {
		return errors.New("Not a BitSet")
	}

	bs.length = length
	bs.set	= make([]uint64, length)

	return	binary.Read(reader, binary.BigEndian, bs.set)
}


func (bs BitSet)byte_get() (buff []byte, err error) {
	data	:= bytes.NewBuffer(make([]byte,0,8*(bs.length+1)))
	err	= binary.Write(data, binary.BigEndian, bs.length)
	if err != nil {
		return []byte{},err
	}

	err	= binary.Write(data, binary.BigEndian, bs.set[0:bs.length])
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
	d,err := bs.byte_get()
	if err != nil {
		return err.Error()
	}
	return string(d)
}


func (bs *BitSet)UnmarshalJSON(data []byte) (err error) {
	return bs.byte_set(bytes.Trim(data,"\""))
}


func (bs BitSet)MarshalJSON() (data []byte,err error) {
	return []byte("\""+bs.String()+"\""),nil
}


func (bs BitSet) Bit(pos uint64, value bool) BitSet {
	var ret	*BitSet

	mod	:= uint(pos%uint64(64))
	idx	:= uint64(pos/uint64(64))
	target	:= uint64(uint64(1)<<mod)

	switch idx < bs.length {
	case	true:
		ret = NewBitSet(bs.length)
	case	false:
		ret = NewBitSet(idx+1)
	}
	copy(ret.set[0:bs.length], bs.set[0:bs.length])

	switch value {
	case	true:
		ret.set[idx] |= target
	case	false:
		ret.set[idx] &^= target
	}

	return *ret
}


func (bs BitSet) Union(b2 BitSet) BitSet {
	b_a, b_b := bs, b2

	if b_a.length < b_b.length {
		b_a, b_b = b_b, b_a
	}

	ret := NewBitSet(b_a.length)

	i := 0
	for i < len(b_b.set) {
		ret.set[i] = b_a.set[i] | b_b.set[i]
		i++
	}

	for i < len(b_a.set) {
		ret.set[i] = b_a.set[i]
		i++
	}

	return *ret
}


func (bs BitSet) Intersection(b2 BitSet) BitSet {
	b_a, b_b := bs, b2

	if b_a.length < b_b.length {
		b_a, b_b = b_b, b_a
	}

	ret := NewBitSet(b_a.length)

	i := 0
	for i < len(b_b.set) {
		ret.set[i] = b_a.set[i] & b_b.set[i]
		i++
	}

	return *ret
}
