package dtp

import "testing"

func Test_from_int64(t *testing.T){
	tests := []struct {
		name	string
		value	int64
		factor	float64
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
		{"positive_factor_one", -123, 1, "-123"},
		{"positive_factor_ten", -5, 10, "-0.5"},
		{"positive_high_precision", 1, 1000000, "0.000001"},
		{"negative_high_precision", -1, 1000000, "-0.000001"},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var d Decimal
			
			d.From_int64(tt.value, tt.factor)
			got := string(d)
			if got != tt.want {
				t.Errorf("From_int64(%d, %f) = %q; want %q", tt.value, tt.factor, got, tt.want)
			}
			
			got_int, err := d.Int64(tt.factor)
			if err != nil {
				t.Errorf("Int64(%f) error: %v", tt.factor, err)
			}
			if got_int != tt.value {
				t.Errorf("Int64(%f) integer = %d; want %d", tt.factor, got_int, tt.value)
			}
			
			if tt.value != 0 {
				err := d.Invert(tt.factor)
				if err != nil {
					t.Errorf("Invert(%f) error: %v", tt.factor, err)
				}
				inverted_int, _ := d.Int64(tt.factor)
				if inverted_int != tt.value * -1 {
					t.Errorf("Invert resulted in %d; want %d", inverted_int, tt.value * -1)
				}
			}
		})
	}
}