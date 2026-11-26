package dtp

import (
	"math"
	"strconv"
	"strings"
)

type Decimal string

func (d *Decimal) From_int64(value, factor int64){
	places := decimal_places(factor)
	s := strconv.FormatInt(value, 10)
	if places == 0 {
		*d = Decimal(s)
		return
	}
	
	digits := strings.TrimPrefix(s, "-")
	if len(digits) <= places {
		digits = strings.Repeat("0", places+1-len(digits)) + digits
	}
	
	point := len(digits) - places
	s = digits[:point]+"."+digits[point:]
	if value < 0 {
		s = "-"+s
	}
	*d = Decimal(s)
}

func (d *Decimal) Int64(factor int64) (int64, error){
	places := decimal_places(factor)
	s := string(*d)
	if s == "" {
		return 0, nil
	}
	
	negative := s[0] == '-'
	if negative || s[0] == '+' {
		s = s[1:]
	}
	
	whole, fraction, _ := strings.Cut(s, ".")
	if whole == "" && fraction == "" {
		return 0, strconv.ErrSyntax
	}
	
	for i := 0; i < len(whole); i++ {
		if whole[i] < '0' || whole[i] > '9' {
			return 0, strconv.ErrSyntax
		}
	}
	
	for i := 0; i < len(fraction); i++ {
		if fraction[i] < '0' || fraction[i] > '9' {
			return 0, strconv.ErrSyntax
		}
	}
	
	round_up := false
	if len(fraction) > places {
		round_up = fraction[places] >= '5'
		fraction = fraction[:places]
	}
	
	if whole == "" {
		whole = "0"
	}
	digits := whole+fraction+strings.Repeat("0", places - len(fraction))
	if negative {
		digits = "-"+digits
	}
	
	value, err := strconv.ParseInt(digits, 10, 64)
	if err != nil || !round_up {
		return value, err
	}
	
	if negative {
		if value == math.MinInt64 {
			return value, strconv.ErrRange
		}
		return value - 1, nil
	}
	if value == math.MaxInt64 {
		return value, strconv.ErrRange
	}
	return value + 1, nil
}

func (d *Decimal) Invert(factor int64) error {
	value, err := d.Int64(factor)
	if err != nil {
		return err
	}
	if value == math.MinInt64 {
		return strconv.ErrRange
	}
	
	d.From_int64(-value, factor)
	return nil
}

func decimal_places(factor int64) int {
	places := 0
	for factor > 1 && factor % 10 == 0 {
		factor /= 10
		places++
	}
	if factor != 1 {
		panic("Invalid decimal factor")
	}
	return places
}