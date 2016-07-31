package	types	// import "github.com/nathanaelle/useful.types"

import (
	"bytes"
	"errors"
	"unicode"
)


/*
 *	type:		FQDN
 *	content:	full qualified domain name
 *	pitfall:	don't verify if the fqdn exists
 *			don't verify if the tld exists
 *			don't handle conversion from unicode to punycode
 */
type	FQDN	string


func (d *FQDN)Get() interface{} {
	return string(*d)
}

func (d *FQDN)UnmarshalJSON(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}

func (d *FQDN)MarshalJSON() (data []byte,err error) {
	return []byte("\""+d.String()+"\""),nil
}

func (d *FQDN)UnmarshalText(data []byte) (err error) {
	return d.Set(string(data))
}

func (d *FQDN)MarshalText() (data []byte,err error) {
	return []byte(d.String()),nil
}



func (d FQDN) String() string {
	return string(d)
}

func (d *FQDN) UnmarshalTOML(data []byte) (err error) {
	return d.Set(string(bytes.Trim(data,"\"")))
}


func (d *FQDN) Set(t_fqdn string) (err error) {
	if len(t_fqdn)< 1 || len(t_fqdn)>253 {
		err	= errors.New("invalid FQDN : "+t_fqdn )
		return
	}

	dot	:= false
	hyphen	:= false
	for pos,char := range t_fqdn {
		switch	{
		case	char == '\\' :
			dot	= false
			hyphen	= false

		case	char == '-' :
			if pos == 0 {
				return	errors.New("begin with hyphen is forbidden for FQDN ["+t_fqdn+"]")
			}
			if dot {
				return	errors.New("hyphen after dot is forbidden for FQDN ["+t_fqdn+"]")
			}
			dot	= false
			hyphen	= true

		case	char == '.' :
			if dot {
				return	errors.New("double dot is forbidden for FQDN ["+t_fqdn+"]")
			}
			if hyphen {
				return	errors.New("dot after hyphen is forbidden for FQDN ["+t_fqdn+"]")
			}
			dot	= true
			hyphen	= false

		case	unicode.IsNumber(char) || unicode.IsLetter(char):
			dot	= false
			hyphen	= false

		default:
			return	errors.New("invalid char ["+string(char)+"] for FQDN ["+t_fqdn+"]")
		}
	}

	*d = FQDN(t_fqdn)

	return nil
}





/*
 *	Explode a FQDN to a slice of strings
 */
func (d *FQDN) Split() []string {
	res	:= make([]string, 1)
	fqdn	:= string(*d)
	begin	:= 0

	quote	:= false
	for pos,char := range fqdn {
		switch	char {
		case	'\\':
			quote = !quote

		case	'.':
			if !quote {
				if len(fqdn[begin:pos]) > 0 {
					res	= append(res, fqdn[begin:pos] )
				}
				begin	= pos+1
			}

		default:
			quote	= false
		}
	}

	if begin < len(fqdn) {
		res	= append(res, fqdn[begin:len(fqdn)] )
	}

	return res
}


/*
 *	Explode a FQDN to a slice of strings
 */
func (d *FQDN) PathToRoot() []string {
	res	:= make([]string, 1)
	fqdn	:= string(*d)
	end	:= len(fqdn)
	last	:= 0
	res[0]	= fqdn

	quote	:= false
	for pos,char := range fqdn {
		switch	char {
		case	'\\':
			quote = !quote

		case	'.':
			if !quote {
				res	= append(res, fqdn[pos:end] )
				last	= pos
			}
			quote = false

		default:
			quote = false
		}
	}

	if last < end-1 {
		res	= append(res, "." )
	}

	return res
}


func (d *FQDN) ToPunny() (*FQDN) {
	return d
}



func (d *FQDN) FromPunny() (*FQDN) {
	return d
}
