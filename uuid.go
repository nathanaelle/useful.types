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
	UUIDv1		byte	= 0x10
	UUIDv2			= 0x20
	UUIDv3			= 0x30
	UUIDv4			= 0x40
	UUIDv5			= 0x50
	UUIDv1MacRand		= 0xe0
	UUIDv1_timestamp	= 0xf0
)


const	(
	UUID_NCS	byte	= 0x00
	UUID_RFC		= 0x80
	UUID_MS			= 0xc0
	UUID_UNUSED		= 0xe0

)


var 	monotonic_v1	uint32	= 0

// HardWareAddress must be set with the desired MAC Address before the generation of a RFC UUIDv1
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

	switch	version {
		case	UUIDv1:
			var	hw	[6]byte
			buffer	:= bytes.NewBuffer(make([]byte,0,16))

			copy(hw[:],HardWareAddress[0:6])
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			err	= binary.Write(buffer, binary.BigEndian, uuidv1{ uint32(now), uint16(now>>32), uint16(now>>48), uint16(seq), hw })
			if err !=nil {
				return
			}
			t_uuid	:= buffer.Bytes()
			t_uuid[6]= (t_uuid[6]&0x0f)|(UUIDv1)
			t_uuid[8]= (t_uuid[8]&0x3f)|(UUID_RFC)
			copy(uuid[:],t_uuid)
			return


		case	UUIDv1MacRand:
			var	hw	[6]byte
			buffer	:= bytes.NewBuffer(make([]byte,0,16))

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
			t_uuid[6]= (t_uuid[6]&0x0f)|(UUIDv1)
			t_uuid[8]= (t_uuid[8]&0x3f)|(UUID_RFC)
			copy(uuid[:],t_uuid)
			return

		case	UUIDv1_timestamp:
			buffer	:= bytes.NewBuffer(make([]byte,0,16))

			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			err	= binary.Write(buffer, binary.BigEndian, uuidv1{ uint32(now>>32), uint16(now>>16), uint16(now>>4), uint16(seq), [6]byte{0,0,0,0,0,0} })
			if err !=nil {
				return
			}
			t_uuid	:= buffer.Bytes()
			t_uuid[6]= (t_uuid[6]&0x0f)|(UUIDv1)
			t_uuid[8]= (t_uuid[8]&0x3f)|(UUID_RFC)
			copy(uuid[:],t_uuid)
			return


		case	UUIDv4:
			_, err = rand.Read(uuid[:])
			if err !=nil {
				return
			}
			uuid[6]= (uuid[6]&0x0f)|(UUIDv4)
			uuid[8]= (uuid[8]&0x3f)|(UUID_RFC)
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
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d UUID)String() string {
	return fmt.Sprintf("%x-%x-%x-%x-%x",d[0:4],d[4:6],d[6:8],d[8:10],d[10:16])
}

func (d *UUID)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *UUID)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}

func (d *UUID)Set(data string) (err error) {
	return d.byte_set([]byte(data))
}


func (d *UUID) byte_set(uuid []byte) (err error) {
	t_uuid	:= [16]byte{}
	t_d	:= bytes.Split(uuid,[]byte{'-'})

	if len(t_d) != 5 {
		return errors.New("not a uuid : ["+ string(uuid) +"]")
	}

	if len(t_d[0]) != 8 {
		return errors.New("not a uuid : ["+ string(uuid) +"]")
	}
	_,err	= hex.Decode(t_uuid[0:4],t_d[0])
	if err	!= nil {
		return
	}


	if len(t_d[1]) != 4 || len(t_d[2]) != 4 || len(t_d[3]) != 4 {
		return errors.New("not a uuid : ["+ string(uuid) +"]")
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
		return errors.New("not a uuid : ["+ string(uuid) +"]")
	}
	_,err	= hex.Decode(t_uuid[10:16],t_d[4])
	if err	!= nil {
		return
	}

	switch t_uuid[6]&0xf0 {
		case UUIDv1,UUIDv2,UUIDv3,UUIDv4,UUIDv5:
		default:
			return errors.New("unknown version : ["+ string(uuid) +"]")
	}

	switch {
		case (t_uuid[8]&0x80)==UUID_NCS:
		case (t_uuid[8]&0xc0)==UUID_RFC:
		case (t_uuid[8]&0xe0)==UUID_MS:
		case (t_uuid[8]&0xe0)==UUID_UNUSED:
		default:
			err= errors.New("unknown variant : ["+ string(uuid) +"]")
			return
	}
	*d = UUID(t_uuid)

	return nil
}


func (u UUID)IsValid() bool {
	t_d := [16]byte(u)

	switch t_d[6]&0xf0 {
		case UUIDv1,UUIDv2,UUIDv3,UUIDv4,UUIDv5:
		default:
			return false
	}

	switch {
		case (t_d[8]&0x80)==UUID_NCS:
		case (t_d[8]&0xc0)==UUID_RFC:
		case (t_d[8]&0xe0)==UUID_MS:
		case (t_d[8]&0xe0)==UUID_UNUSED:

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
