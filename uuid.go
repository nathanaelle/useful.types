package	types	// import "github.com/nathanaelle/useful.types"

import	(
	"net"
	"time"
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


var 	(
	ErrInvalidUUID		error	= errors.New("Invalid UUID")
	ErrUnknownVersionUUID	error	= errors.New("Uknown UUID Version")
	ErrUnknownVariantUUID	error	= errors.New("Uknown UUID Variant")
)


// 100ns ticks from 1582-10-15 to 1970-01-01 = 387 years + 97 leap years - 3 non leap centuries
const	from15821015to19700101 int64	= (387+97-3)*(24*60*60)*(1000*1000*10)

var 	monotonic_v1	uint32	= 0

// HardWareAddress must be set with the desired MAC Address before the generation of a RFC UUIDv1
var	HardWareAddress	net.HardwareAddr

type	(
	UUID	[16]byte

	uuidv1 struct {
		t1	uint32
		t2	uint16
		t3	uint16
		Seq	uint16
		HW	[6]byte
	}
)

func ts60bits(ts int64) uint64 {
	return	uint64(ts/100 + from15821015to19700101)&0x0fffffffffffffff
}


func NewUUID(version byte) (uuid UUID, err error)  {

	switch	version {
		case	UUIDv1:
			seq	:= atomic.AddUint32(&monotonic_v1, 1)
			now	:= ts60bits(time.Now().UnixNano())
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
			now	:= ts60bits(time.Now().UnixNano())
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
			now	:= ts60bits(time.Now().UnixNano())<<4
			binary.BigEndian.PutUint32(uuid[0:4], uint32(now>>32))
			binary.BigEndian.PutUint16(uuid[4:6], uint16(now>>16))
			binary.BigEndian.PutUint16(uuid[6:8], uint16(now)>>4)
			binary.BigEndian.PutUint16(uuid[8:10], uint16(seq))
			_, err = rand.Read(uuid[10:16])

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


func (d UUID)byte_text(t []byte) []byte {
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


func valid_hex(data []byte) bool {
	for _, h := range data {
		if !(h >= '0' && h <= '9') && !(h >= 'a' && h <= 'f') {
			return	false
		}
	}
	return	true
}

func uuidIsValidVersion(d [16]byte) bool {
	switch d[6]&0xf0 {
	case UUIDv1,UUIDv2,UUIDv3,UUIDv4,UUIDv5:
		return true
	}
	return	false
}

func uuidIsValidVariant(d [16]byte) bool {
	switch {
	case (d[8]&0x80)==UUID_NCS:
		return	true
	case (d[8]&0xc0)==UUID_RFC:
		return	true
	case (d[8]&0xe0)==UUID_MS:
		return	true
	case (d[8]&0xe0)==UUID_UNUSED:
		return	true
	}
	return	false
}

func uuidIsValidSlice(uuid []byte) bool {
	if len(uuid)!=36 {
		return	false
	}
	hexOK	:= valid_hex(uuid[0:8]) && valid_hex(uuid[9:13]) && valid_hex(uuid[14:18]) && valid_hex(uuid[19:23]) && valid_hex(uuid[24:36])
	dashOK	:= uuid[8] == '-' && uuid[13] == '-' && uuid[18] == '-' && uuid[23] == '-'
	return	hexOK && dashOK
}


func (d *UUID) byte_set(uuid []byte) error {
	t_uuid	:= [16]byte{}
	if !uuidIsValidSlice(uuid) {
		return ErrInvalidUUID
	}

	if _,err := hex.Decode(t_uuid[0:4],uuid[0:8]); err	!= nil {
		return	err
	}

	if _,err := hex.Decode(t_uuid[4:6],uuid[9:13]); err	!= nil {
		return	err
	}

	if _,err := hex.Decode(t_uuid[6:8],uuid[14:18]); err	!= nil {
		return	err
	}

	if _,err := hex.Decode(t_uuid[8:10],uuid[19:23]); err	!= nil {
		return	err
	}
	if _,err := hex.Decode(t_uuid[10:16],uuid[24:36]); err	!= nil {
		return	err
	}

	if !uuidIsValidVersion(t_uuid) {
		return ErrUnknownVersionUUID
	}

	if !uuidIsValidVariant(t_uuid) {
		return	ErrUnknownVariantUUID
	}

	*d = UUID(t_uuid)

	return nil
}



func (d UUID)Get() interface{} {
	ret := d
	return ret
}

func (d UUID)MarshalText() ([]byte,error) {
	return	d.byte_text(make([]byte,36)), nil
}

func (d *UUID)UnmarshalText(data[]byte) error {
	return	d.byte_set(data)
}

func (d UUID)MarshalBinary() ([]byte,error) {
	var ret [16]byte
	copy(ret[:], d[:])
	return	ret[:], nil
}

func (d *UUID)UnmarshalBinary(data[]byte) error {
	if len(data)!=16 {
		return ErrInvalidUUID
	}
	copy(d[:],data)
	return	nil
}

func (d UUID)String() string {
	return string(d.byte_text(make([]byte,36)))
}


func (d *UUID)Set(data string) (err error) {
	return d.byte_set([]byte(data))
}

func (u UUID)IsValid() bool {
	t_uuid := [16]byte(u)

	return uuidIsValidVersion(t_uuid) && uuidIsValidVariant(t_uuid)
}


func (u1 UUID)IsEqual(u2 UUID) bool {
	t_u1 := [16]byte(u1)
	t_u2 := [16]byte(u2)
	return subtle.ConstantTimeCompare(t_u1[:],t_u2[:]) == 1
}
