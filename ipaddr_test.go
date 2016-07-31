package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_IPAddr(t *testing.T) {
	d := new(IpAddr)
	Has_All_Interfaces(t, d)
}

func Benchmark_IpAddr_Set(b *testing.B) {
	d := new(IpAddr)
	for i := 0; i < b.N; i++ {
		d.Set("192.168.0.1")
	}
}

func Benchmark_IpAddr_String(b *testing.B) {
	d := new(StoreSize)
	d.Set("192.168.0.1")
	for i := 0; i < b.N; i++ {
		d.String()
	}
}
