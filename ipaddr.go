package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"errors"
	"net"
)

// Wrapper Type for net.IP providing missing interfaces
// All the parsing and the validation are done by net.ParseIP
type	IpAddr	net.IP

func (d *IpAddr)byte_set(data []byte) (err error) {
	dest := IpAddr(net.ParseIP(string(data)))
	if dest == nil {
		return errors.New("invalid IpAddr : "+string(data))
	}
	*d = dest

	return nil
}

func (d *IpAddr)Get() interface{} {
	return net.IP(*d)
}

func (d *IpAddr)UnmarshalTOML(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *IpAddr)String() string {
	return net.IP(*d).String()
}

func (d *IpAddr)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *IpAddr)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}

func (d *IpAddr) ToTCPAddr(port string) (*net.TCPAddr, error)   {
	ip := net.IP(*d)
	if ip.To4() == nil {
		return net.ResolveTCPAddr( "tcp", "["+ip.String()+"]:"+port )
	}
	return net.ResolveTCPAddr( "tcp", ip.String()+":"+port )
}
