package	types	// import "github.com/nathanaelle/useful.types"

import	(
	"net"
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
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			binary.BigEndian.PutUint32(uuid[0:4], uint32(now))
			binary.BigEndian.PutUint16(uuid[4:6], uint16(now>>32))
			binary.BigEndian.PutUint16(uuid[6:8], uint16(now>>48))
			binary.BigEndian.PutUint16(uuid[8:10], uint16(seq))
			copy(uuid[10:16], HardWareAddress[0:6])
			uuid[6]= (uuid[6]&0x0f)|(UUIDv1)
			uuid[8]= (uuid[8]&0x3f)|(UUID_RFC)
			return

		case	UUIDv1MacRand:
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			binary.BigEndian.PutUint32(uuid[0:4], uint32(now))
			binary.BigEndian.PutUint16(uuid[4:6], uint16(now>>32))
			binary.BigEndian.PutUint16(uuid[6:8], uint16(now>>48))
			binary.BigEndian.PutUint16(uuid[8:10], uint16(seq))
			_, err = rand.Read(uuid[10:16])
			uuid[6]= (uuid[6]&0x0f)|(UUIDv1)
			uuid[8]= (uuid[8]&0x3f)|(UUID_RFC)
			return

		case	UUIDv1_timestamp:
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= uint64(time.Now().UnixNano()/100 + 12219292800000)
			binary.BigEndian.PutUint32(uuid[0:4], uint32(now))
			binary.BigEndian.PutUint16(uuid[4:6], uint16(now>>32))
			binary.BigEndian.PutUint16(uuid[6:8], uint16(now>>48))
			binary.BigEndian.PutUint16(uuid[8:10], uint16(seq))

			uuid[6]= (uuid[6]&0x0f)|(UUIDv1)
			uuid[8]= (uuid[8]&0x3f)|(UUID_RFC)
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


func (d *UUID)byte_text(t []byte) []byte {
	if len(t) < 36 {
		panic("UUID len too short")
	}
	hex.Encode(t[0:8], d[0:4])
	t[8] = '-'
	hex.Encode(t[9:13], d[4:6])
	t[13] = '-'
	hex.Encode(t[14:18], d[6:8])
	t[18] = '-'
	hex.Encode(t[19:23], d[8:10])
	t[23] = '-'
	hex.Encode(t[24:36], d[10:16])

	return	t
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



func (d *UUID)Get() interface{} {
	return [16]byte(*d)
}

func (d *UUID)UnmarshalTOML(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *UUID)MarshalText() ([]byte,error) {
	return	d.byte_text(make([]byte,36)), nil
}

func (d *UUID)UnmarshalText(data[]byte) error {
	return	d.byte_set(data)
}

func (d *UUID)String() string {
	return string(d.byte_text(make([]byte,36)))
}

func (d *UUID)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *UUID)MarshalJSON() (data []byte,err error) {
	t := make([]byte,38)
	t[0] = '"'
	t[37]= '"'

	d.byte_text(t[1:37])

	return t,nil
}

func (d *UUID)Set(data string) (err error) {
	return d.byte_set([]byte(data))
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
