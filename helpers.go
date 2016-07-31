package	types	// import "github.com/nathanaelle/useful.types"

import	(
	"math"
	"strconv"
)


var 	humanscale_exp_byte	[]byte = []byte { 'y','z','a','f','p','n','µ','m',' ','k','M','G','T','P','E','Z','Y' }


func Round(x float64) int {
	switch	{
	case	x < 0:
		return int(x-0.5)
	case	x > 0:
		return int(x+0.5)
	default:
		return 0
	}
}


func MetricPower(value float64, base float64, max_precision int) (neg bool, scaled float64, precision int, power int) {
	power_f	:= math.Floor(math.Log2(value)/math.Log2(base))
	power	= int(power_f)
	scaled	= value / math.Pow(base, power_f)
	if value < 0 {
		neg	= true
	}

	if max_precision <= 0 {
		return
	}

	_,fracf	:= math.Modf(scaled)
	frac	:= Round(fracf*math.Pow(10,float64(max_precision)))
	if frac == 0 {
		return
	}
	precision= max_precision

	for precision > 0 {
		if (frac % 10) != 0 {
			break
		}
		frac = Round(float64(frac)/10)
		precision--
	}

	return
}


func HumanScaleString(value float64, base float64, unit string) string {
	if value == 0 {
		return	"0"+unit
	}

	exp	:= []string { "y","z","a","f","p","n","µ","m","","k","M","G","T","P","E","Z","Y" }
	_, h_v, p, s	:= MetricPower(value, base, 3)

	if s > -9 && s < 9 {
		return	strconv.FormatFloat(h_v,'f',p,64)+exp[int(s)+8]+unit
	}

	return strconv.FormatFloat(value,'E',6,64)+unit
}

func HumanScaleBytes(value float64, base float64, unit []byte) []byte {
	if value == 0 {
		return	append( append(make([]byte,0,1+len(unit)), '0'), unit...)
	}

	neg, h_v, p, s	:= MetricPower(value, base, 3)
	sign	:= 0
	if neg {
		sign	= 1
	}

	switch {
	case	s == 0:
		return	append( strconv.AppendFloat(make([]byte,0,sign+4+p+len(unit)), h_v,'f',p,64), unit...)
	case s > -9 && s < 9 :
		return	append(	append(
			strconv.AppendFloat(make([]byte,0,sign+4+p+1+len(unit)), h_v,'f',p,64),
			humanscale_exp_byte[int(s)+8]),
			unit...)
	default:
		return	append( strconv.AppendFloat(make([]byte,0,sign+12+len(unit)), value,'E',6,64), unit...)
	}
}





func FieldsFuncN(s string, hope int, f func(rune) bool) []string {
	p_is_sep:= true
	is_sep	:= true
	begin	:= -1
	end	:= -1
	res	:= make( []string, 0, hope )

	for i,rune := range s {
		p_is_sep = is_sep
		is_sep = f(rune)
		switch {
			case is_sep && !p_is_sep:
				end = i
				res = append( res, s[begin:end] )

			case !is_sep && p_is_sep:
				begin	= i
		}
	}

	if(begin>-1 && begin>end ) {
		res = append( res, s[begin:len(s)] )
	}

	return res
}
