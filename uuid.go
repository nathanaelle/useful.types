package	types	// import "github.com/nathanaelle/useful.types"

import	(
	"net"
	"fmt"
	"time"
	"bytes"
	"errors"
	"sync/atomic"
	"crypto/rand"
	"encoding/hex"
	"crypto/subtle"
	"encoding/binary"
)

const	(
	UUIDv1		byte	= iota+1
	UUIDv2
	UUIDv3
	UUIDv4
	UUIDv5
)


const	(
	UUIDRFC		byte	= (iota+8)<<4
	UUIDvar1
	UUIDvar2
	UUIDvar3
)

const	UUIDv1Rand	byte	= UUIDv1|UUIDvar1
const	UUIDv1SortRand	byte	= UUIDv1|UUIDvar2


var 	monotonic_v1	uint32	= 0

var	HardWareAddress	net.HardwareAddr

type	UUID	[16]byte

type	uuidv1 struct {
	T1	uint32
	T2	uint16
	T3	uint16
	Seq	uint16
	HW	[6]byte
}



func NewUUID(version byte) (uuid UUID, err error)  {
	ver	:= version&0x0f
	variant	:= version&0xf0

	switch	version {
		case	UUIDv1:
			var	hw	[6]byte
			buffer	:= new(bytes.Buffer)

			copy(hw[:],HardWareAddress[0:6])
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			err	= binary.Write(buffer, binary.BigEndian, uuidv1{ uint32(now), uint16(now>>32), uint16(now>>48), uint16(seq), hw })
			if err !=nil {
				return
			}
			t_uuid	:= buffer.Bytes()
			t_uuid[6]= (t_uuid[6]&0x0f)|(ver<<4)
			t_uuid[8]= (t_uuid[8]&0x0f)|(UUIDvar1)
			copy(uuid[:],t_uuid)
			return


		case	UUIDv1Rand:
			var	hw	[6]byte
			buffer	:= new(bytes.Buffer)

			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			_, err = rand.Read(hw[:])
			if err !=nil {
				return
			}
			err	= binary.Write(buffer, binary.BigEndian, uuidv1{ uint32(now), uint16(now>>32), uint16(now>>48), uint16(seq), hw })
			if err !=nil {
				return
			}
			t_uuid	:= buffer.Bytes()
			t_uuid[6]= (t_uuid[6]&0x0f)|(ver<<4)
			t_uuid[8]= (t_uuid[8]&0x0f)|(variant)
			copy(uuid[:],t_uuid)
			return

		case	UUIDv1SortRand:
			var	hw	[6]byte
			buffer	:= new(bytes.Buffer)

			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			_, err = rand.Read(hw[:])
			if err !=nil {
				return
			}
			err	= binary.Write(buffer, binary.BigEndian, uuidv1{ uint32(now>>32), uint16(now>>16), uint16(now>>4), uint16(seq), hw })
			if err !=nil {
				return
			}
			t_uuid	:= buffer.Bytes()
			t_uuid[6]= (t_uuid[6]&0x0f)|(ver<<4)
			t_uuid[8]= (t_uuid[8]&0x0f)|(variant)
			copy(uuid[:],t_uuid)
			return


		case	UUIDv4:
			_, err = rand.Read(uuid[:])
			if err !=nil {
				return
			}
			uuid[6]= (uuid[6]&0x0f)|(ver<<4)
			uuid[8]= (uuid[8]&0x0f)|(UUIDRFC)
			return

		default:
			err	= errors.New("cant generate this uuid")
			return
	}

}

func (d *UUID)Get() interface{} {
	return [16]byte(*d)
}

func (d *UUID)UnmarshalTOML(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *UUID)String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",d[0:4],d[4:6],d[6:8],d[8:10],d[10:16])
}

func (d *UUID)UnmarshalJSON(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *UUID)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}


func (d *UUID) Set(uuid string) (err error) {
	t_uuid	:= [16]byte{}
	t_d	:= bytes.Split([]byte(uuid),[]byte{'-'})

	if len(t_d) != 5 {
		return errors.New("not a uuid : ["+ uuid +"]")
	}

	if len(t_d[0]) != 8 {
		return errors.New("not a uuid : ["+ uuid +"]")
	}
	_,err	= hex.Decode(t_uuid[0:4],t_d[0])
	if err	!= nil {
		return
	}


	if len(t_d[1]) != 4 || len(t_d[2]) != 4 || len(t_d[3]) != 4 {
		return errors.New("not a uuid : ["+ uuid +"]")
	}
	_,err	= hex.Decode(t_uuid[4:6],t_d[1])
	if err	!= nil {
		return
	}
	_,err	= hex.Decode(t_uuid[6:8],t_d[2])
	if err	!= nil {
		return
	}
	_,err	= hex.Decode(t_uuid[8:10],t_d[3])
	if err	!= nil {
		return
	}

	if len(t_d[4]) != 12 {
		return errors.New("not a uuid : ["+ uuid +"]")
	}
	_,err	= hex.Decode(t_uuid[10:16],t_d[4])
	if err	!= nil {
		return
	}

	switch t_uuid[6]>>4 {
		case UUIDv1,UUIDv2,UUIDv3,UUIDv4,UUIDv5:
		default:
			return errors.New("unknown version : ["+ uuid +"]")
	}

	switch t_uuid[8]&0xf0 {
		case UUIDvar1,UUIDvar2,UUIDvar3,UUIDRFC:
		default:
			err= errors.New("unknown variant : ["+ uuid +"]")
			return
	}
	*d = UUID(t_uuid)

	return nil
}


func (u UUID)IsValid() bool {
	t_d := [16]byte(u)

	switch t_d[6]>>4 {
		case UUIDv1,UUIDv2,UUIDv3,UUIDv4,UUIDv5:
		default:
			return false
	}

	switch t_d[8]&0xf0 {
		case UUIDvar1,UUIDvar2,UUIDvar3,UUIDRFC:
		default:
			return false
	}

	return true
}


func (u1 UUID)IsEqual(u2 UUID) bool {
	t_u1 := [16]byte(u1)
	t_u2 := [16]byte(u2)
	return subtle.ConstantTimeCompare(t_u1[:],t_u2[:]) == 1
}
