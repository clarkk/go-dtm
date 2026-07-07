package dtm

/*
	Test
	# go test . -v
*/

import "testing"

func Test_math(t *testing.T){
	t.Run("Factor_int64_multiply", func(t *testing.T){
		a		:= int64(123456)
		b 		:= int64(345678)
		factor	:= int64(100)
		
		want	:= int64(426760232)
		
		res, _ := Factor_int64_multiply(a, b, factor)
		if want != res {
			t.Fatalf("Want:\n%d\nGot:\n%d", want, res)
		}
	})
	
	t.Run("Factor_int64_percent", func(t *testing.T){
		i				:= int64(123456)
		percent 		:= int64(2595)		// 25.95%
		factor_percent	:= int64(100)
		
		want			:= int64(32037)
		
		res, _ := Factor_int64_percent(i, percent, factor_percent)
		if want != res {
			t.Fatalf("Want:\n%d\nGot:\n%d", want, res)
		}
	})
	
	t.Run("Factor_int64_convert", func(t *testing.T){
		i			:= int64(12345678)
		factor_from	:= int64(10000)
		factor_to	:= int64(100)
		
		want	:= int64(123457)
		
		res := Factor_int64_convert(i, factor_from, factor_to)
		if want != res {
			t.Fatalf("Want:\n%d\nGot:\n%d", want, res)
		}
	})
}