package dtp

import (
	"math"
	"strconv"
	"strings"
)

type Decimal string

func (d *Decimal) From_int64(value int64, factor float64){
	if factor <= 1 {
		*d = Decimal(strconv.FormatInt(value, 10))
		return
	}
	
	//	Determine the number of decimal places
	exp := int(math.Round(math.Log10(factor)))
	
	abs_value := value
	var sign string
	if value < 0 {
		abs_value = -value
		sign = "-"
	}
	
	//	Convert absolute value to string and pad with zeros if it's smaller than the factor
	s := strconv.FormatInt(abs_value, 10)
	if len(s) <= exp {
		s = strings.Repeat("0", exp - len(s) + 1) + s
	}
	
	//	Manually insert the decimal point to avoid float64 errors
	cut := len(s) - exp
	*d = Decimal(sign+s[:cut]+"."+s[cut:])
}

func (d *Decimal) Int64(factor float64) (int64, error){
	if *d == "" {
		return 0, nil
	}
	
	f, err := strconv.ParseFloat(string(*d), 64)
	if err != nil {
		return 0, err
	}
	return int64(math.Round(f * factor)), nil
}

func (d *Decimal) Invert(factor float64) error {
	i, err := d.Int64(factor)
	if err != nil {
		return err
	}
	d.From_int64(i * -1, factor)
	return nil
}