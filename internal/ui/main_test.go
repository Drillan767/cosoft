package ui

import (
	"testing"
	"time"
)

func Test_roundHourToQuarter(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		t    time.Time
		want time.Time
	}{
		{
			name: "first_quarter",
			t:    time.Date(2026, 1, 27, 13, 8, 0, 0, time.Local),
			want: time.Date(2026, 1, 27, 13, 15, 0, 0, time.Local),
		},
		{
			name: "second_quarter",
			t:    time.Date(2026, 1, 27, 13, 23, 0, 0, time.Local),
			want: time.Date(2026, 1, 27, 13, 30, 0, 0, time.Local),
		},
		{
			name: "third_quarter",
			t:    time.Date(2026, 1, 27, 13, 38, 0, 0, time.Local),
			want: time.Date(2026, 1, 27, 13, 45, 0, 0, time.Local),
		},
		{
			name: "next_hour",
			t:    time.Date(2026, 1, 27, 13, 53, 0, 0, time.Local),
			want: time.Date(2026, 1, 27, 14, 0, 0, 0, time.Local),
		},
		{
			name: "next_quarter",
			t:    time.Date(2026, 1, 27, 14, 0, 0, 0, time.Local),
			want: time.Date(2026, 1, 27, 14, 15, 0, 0, time.Local),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := roundHourToQuarter(tt.t)
			if got != tt.want {
				t.Errorf("roundHourToQuarter() = %v, want %v", got, tt.want)
			}
		})
	}
}
