package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"errors"
	"strconv"
	_ "math"

)



var	storesize_validchar []byte =  []byte{ 'o','B','i','k','K','M','G','T','P','E','Z','Y' }


func storesize_is_validchar(b byte) (bool,int) {
	switch b {
	case	'o','B':
		return true,0
	case	'i':
		return true,2
	case	'k','K':
		return true,3
	case	'M':
		return true,5
	case	'G':
		return true,6
	case	'T':
		return true,7
	case	'P':
		return true,8
	case	'E':
		return true,9
	case	'Z':
		return true,10
	case	'Y':
		return true,11
	}

	for i,valid := range storesize_validchar {
		if b == valid {
			return true,i
		}
	}
	return false,0
}


/*
 *	type:		int64
 *	content:	size in bytes / octet
 */
type	StoreSize	int64


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


/*
func (d *StoreSize)byte_text() []byte {
	l_text	:= 0
	v_i	:= int64(*d)
	neg	:= false

	if v == 0 {
		return	[]byte{ '0' }
	}

	if v < 0 {
		l_text	= 1
		neg	= true
		v	= - v
	}
	v_f	:= float64(v_i)

	ln10	:= math.Log10(v_f)/3
	ln2	:= math.Log2(v_f)/10



	b_ok	:= true
	d_ok	:= true
	d10_3	:= int64(*d)
	d2_10	:= int64(*d)


}
*/


func (d *StoreSize)Set(data string) (err error) {
	return d.byte_set([]byte(data))
}

func (d *StoreSize)Get() interface{} {
	return int64(*d)
}

func (d *StoreSize)UnmarshalTOML(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *StoreSize)UnmarshalText(data []byte) (err error) {
	return d.byte_set(data)
}


func (d *StoreSize)MarshalText() ([]byte,error) {
	return []byte(d.String()),nil
}



func (d StoreSize)String() string {
	return strconv.FormatInt(int64(d),10)
}

func (d *StoreSize)UnmarshalJSON(data []byte) (err error) {
	return d.byte_set(bytes.Trim(data,"\""))
}

func (d *StoreSize)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}
