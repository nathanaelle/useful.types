package types // import "github.com/nathanaelle/useful.types"

import (
	"testing"
)

func Test_IPAddr(t *testing.T) {
	d := new(IpAddr)
	Has_All_Interfaces(t, d)
}
