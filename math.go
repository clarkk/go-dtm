package dtm

import (
	"fmt"
	"math"
)

//	Multiply two integers with factor and round to nearest
func Factor_int64_multiply(a, b, factor int64) (int64, error){
	res, err := multiplication_overflow(a, b)
	if err != nil {
		return 0, err
	}
	
	//	Round half away from zero (round to nearest)
	half := factor / 2
	if res >= 0 {
		//	Check if adding half exceeds MaxInt64
		if res > math.MaxInt64 - half {
			return 0, fmt.Errorf("Multiplication rounding overflow: %d + %d", res, half)
		}
		res += half
	} else {
		//	Check if subtracting half exceeds MinInt64
		if res < math.MinInt64 + half {
			return 0, fmt.Errorf("Multiplication rounding overflow: %d - %d", res, half)
		}
		res -= half
	}
	
	return res / factor, nil
}

//	Calculate percentage with factor and round to nearest
func Factor_int64_percent(a, b, factor int64) (int64, error){
	res, err := multiplication_overflow(a, b)
	if err != nil {
		return 0, err
	}
	
	half := factor * 50
	if res >= 0 {
		//	Check if adding half exceeds MaxInt64
		if res > math.MaxInt64 - half {
			return 0, fmt.Errorf("rounding overflow: value too large")
		}
		res += half
	} else {
		//	Check if subtracting half exceeds MinInt64
		if res < math.MinInt64 + half {
			return 0, fmt.Errorf("rounding overflow: value too small")
		}
		res -= half
	}
	
	return res / (factor * 100), nil
}

//	Convert from one factor to another and round to nearest
func Factor_int64_convert(i, factor_from, factor_to int64) int64 {
	return (i * factor_to + factor_from / 2) / factor_from
}

func multiplication_overflow(a, b int64) (int64, error){
	if a == 0 || b == 0 {
		return 0, nil
	}
	res := a * b
	if res / b != a {
		return 0, fmt.Errorf("Multiplication overflow: %d * %d", a, b)
	}
	//	Handle the specific case where -1 * MinInt64 overflows to MinInt64
	if a == -1 && b == math.MinInt64 {
		return 0, fmt.Errorf("Multiplication overflow: %d * %d", a, b)
	}
	if b == -1 && a == math.MinInt64 {
		return 0, fmt.Errorf("Multiplication overflow: %d * %d", a, b)
	}
	return res, nil
}