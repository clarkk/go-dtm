package dtp

import (
	"math"
	"strconv"
	"testing"
)

func Test_from_int64(t *testing.T){
	tests := []struct {
		name	string
		value	int64
		factor	int64
		want	string
	}{
		{"positive", 199, 100, "1.99"},
		{"positive_less", 5, 100, "0.05"},
		{"negative", -199, 100, "-1.99"},
		{"negative_less", -5, 100, "-0.05"},
		{"positive_exact_multiple", 500, 100, "5.00"},
		{"negative_exact_multiple", -500, 100, "-5.00"},
		{"zero_value", 0, 100, "0.00"},
		{"zero_value_large_factor", 0, 10000, "0.0000"},
		{"positive_factor_one", 123, 1, "123"},
		{"positive_factor_ten", 5, 10, "0.5"},
		{"negative_factor_one", -123, 1, "-123"},
		{"negative_factor_ten", -5, 10, "-0.5"},
		{"positive_high_precision", 1, 1000000, "0.000001"},
		{"negative_high_precision", -1, 1000000, "-0.000001"},
		{"max_int64", math.MaxInt64, 100, "92233720368547758.07"},
		{"min_int64", math.MinInt64, 100, "-92233720368547758.08"},
		{"large_factor", 1, 1000000000000000000, "0.000000000000000001"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Decimal
			
			d.From_int64(tt.value, tt.factor)
			got := string(d)
			if got != tt.want {
				t.Fatalf("From_int64(%d, %d) = %q; want %q", tt.value, tt.factor, got, tt.want)
			}
			
			got_int, err := d.Int64(tt.factor)
			if err != nil {
				t.Fatalf("Int64(%d) error: %v", tt.factor, err)
			}
			if got_int != tt.value {
				t.Fatalf("Int64(%d) integer = %d; want %d", tt.factor, got_int, tt.value)
			}
			
			err = d.Invert(tt.factor)
			
			if tt.value == math.MinInt64 {
				if err != strconv.ErrRange {
					t.Errorf("Invert(%d) = %v; want ErrRange", tt.factor, err)
				}
				return
			}
			
			if err != nil {
				t.Fatalf("Invert(%d) error: %v", tt.factor, err)
			}
			
			inverted_int, err := d.Int64(tt.factor)
			if err != nil {
				t.Fatalf("Int64(%d) error: %v", tt.factor, err)
			}
			
			if inverted_int != -tt.value {
				t.Errorf("Invert resulted in %d; want %d", inverted_int, -tt.value)
			}
		})
	}
}

func Test_int64(t *testing.T) {
	tests := []struct {
		value		Decimal
		factor		int64
		want		int64
		want_err	bool
	}{
		{"123.456", 100, 12346, false},
		{"123.454", 100, 12345, false},
		{"123.455", 100, 12346, false},
		{"-123.456", 100, -12346, false},
		{"-123.454", 100, -12345, false},
		{"-123.455", 100, -12346, false},
		{"0.005", 100, 1, false},
		{"-0.005", 100, -1, false},
		{"0.004", 100, 0, false},
		{"-0.004", 100, 0, false},
		{"9.999", 100, 1000, false},
		{"-9.999", 100, -1000, false},
		{"12.3", 1000, 12300, false},
		{"12", 100, 1200, false},
		{"1.5", 1, 2, false},
		{"-1.5", 1, -2, false},
		{"", 100, 0, false},
		{"0", 100, 0, false},
		{"90071992547409.93", 100, 9007199254740993, false},
		{"9223372036854775807", 1, math.MaxInt64, false},
		{"-9223372036854775808", 1, math.MinInt64, false},
		{"9223372036854775807.5", 1, 0, true},
		{"-9223372036854775808.5", 1, 0, true},
		{"9223372036854775808", 1, 0, true},
		{"-9223372036854775809", 1, 0, true},
		{"abc", 100, 0, true},
		{"1.2.3", 100, 0, true},
	}
	
	for _, tt := range tests {
		t.Run(string(tt.value), func(t *testing.T) {
			got, err := tt.value.Int64(tt.factor)
			
			if (err != nil) != tt.want_err {
				t.Fatalf("Int64(%d) unexpected error: %v",
					tt.factor, err)
			}
			
			if !tt.want_err && got != tt.want {
				t.Errorf("Int64(%d) = %d; want %d",
					tt.factor, got, tt.want)
			}
		})
	}
}

func Test_decimal_places(t *testing.T) {
	tests := []struct {
		factor	int64
		want	int
	}{
		{1, 0},
		{10, 1},
		{100, 2},
		{1000, 3},
		{10000, 4},
		{1000000000000000000, 18},
	}
	
	for _, tt := range tests {
		if got := decimal_places(tt.factor); got != tt.want {
			t.Errorf("decimal_places(%d) = %d; want %d",
				tt.factor, got, tt.want)
		}
	}
}

func Test_invalid_decimal_factor(t *testing.T) {
	for _, factor := range []int64{-100, -1, 0, 2, 5, 15, 50, 101} {
		func() {
			defer func() {
				if recover() == nil {
					t.Errorf("expected panic for factor %d", factor)
				}
			}()
			
			decimal_places(factor)
		}()
	}
}