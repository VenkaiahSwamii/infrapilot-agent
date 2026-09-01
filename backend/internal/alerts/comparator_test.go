package alerts

import (
	"testing"
)

func TestCompare(t *testing.T) {
	tests := []struct {
		name      string
		value     float64
		operator  string
		threshold float64
		want      bool
	}{
		// Prompt mandatory test cases (Phase 13)
		{
			name:      "92 > 90",
			value:     92.0,
			operator:  ">",
			threshold: 90.0,
			want:      true,
		},
		{
			name:      "50 > 90",
			value:     50.0,
			operator:  ">",
			threshold: 90.0,
			want:      false,
		},
		{
			name:      "95 >= 95",
			value:     95.0,
			operator:  ">=",
			threshold: 95.0,
			want:      true,
		},
		{
			name:      "20 != 30",
			value:     20.0,
			operator:  "!=",
			threshold: 30.0,
			want:      true,
		},
		{
			name:      "85 < 90",
			value:     85.0,
			operator:  "<",
			threshold: 90.0,
			want:      true,
		},
		// Additional operator edge cases
		{
			name:      "90 <= 90",
			value:     90.0,
			operator:  "<=",
			threshold: 90.0,
			want:      true,
		},
		{
			name:      "100 == 100",
			value:     100.0,
			operator:  "==",
			threshold: 100.0,
			want:      true,
		},
		{
			name:      "100 = 100",
			value:     100.0,
			operator:  "=",
			threshold: 100.0,
			want:      true,
		},
		{
			name:      "30 != 30",
			value:     30.0,
			operator:  "!=",
			threshold: 30.0,
			want:      false,
		},
		{
			name:      "invalid operator",
			value:     50.0,
			operator:  "INVALID",
			threshold: 50.0,
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Compare(tt.value, tt.operator, tt.threshold)
			if got != tt.want {
				t.Errorf("Compare(%.1f, %q, %.1f) = %v, want %v",
					tt.value, tt.operator, tt.threshold, got, tt.want)
			}
		})
	}
}
