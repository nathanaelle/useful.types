package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"errors"
	"strconv"
)



var	storesize_validchar []byte =  []byte{ 'o','B','i','k','K','M','G','T','P','E','Z','Y' }


func storesize_is_validchar(b byte) (bool,int) {
	for i,valid := range storesize_validchar {
		if b == valid {
			return true,i
		}
	}
	return false,0
}



/*
 *	type:		Duration
 *	content:	time duration aka intergers with time units
 */
type	StoreSize	int64

func (d *StoreSize)Set(data string) (err error) {
	return d.byte_set([]byte(data))
}

func (d *StoreSize)Get() interface{} {
	return int64(*d)
}

func (d *StoreSize)UnmarshalTOML(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *StoreSize)String() string {
	return strconv.FormatInt(int64(*d),10)
}

func (d *StoreSize)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *StoreSize)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}

func (d *StoreSize)byte_set(data []byte) (err error) {
	max		:= len(data)-1
	digit_only	:= true
	binary_unit	:= false
	power		:= 0
	last_num	:= 0

	for i,b	:= range data {
		if !digit_only {
			if b == 'i' {
				if i != max-1 {
					return errors.New("invalid StoreSize : "+string(data))
				}
				binary_unit = true
				continue
			}

			if b == 'o' || b == 'B' {
				if i != max {
					return errors.New("invalid StoreSize : "+string(data))
				}
				continue
			}

			return errors.New("invalid StoreSize : "+string(data))
			continue
		}

		if b >= '0' && b <= '9' {
			continue
		}
		last_num	= i
		digit_only	= false

		ok,pos := storesize_is_validchar(b)
		if  !ok {
			return errors.New("invalid StoreSize : "+string(data))
		}

		if i < max-2 {
			return errors.New("invalid StoreSize : "+string(data))
		}

		if b == 'i' {
			return errors.New("invalid StoreSize : "+string(data))
		}

		if b == 'o' || b == 'B' {
			if i != max {
				return errors.New("invalid StoreSize : "+string(data))
			}
			continue
		}

		if pos > 3 {
			power = pos-3
			continue
		}
		power=1
	}

	if digit_only {
		v,err := strconv.ParseInt(string(data), 10, 64 )
		if err != nil {
			return err
		}
		*d = StoreSize(v)
		return nil
	}


	v,err := strconv.ParseInt(string(data[0:last_num]), 10, 64 )
	if err != nil {
		return err
	}


	factor	:= int64(1000)
	if binary_unit {
		factor = 1024
	}

	for power > 0 {
		v = v*factor
		power--
	}

	*d = StoreSize(v)

	return nil
}
