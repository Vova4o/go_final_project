package nextdate

import (
	"testing"
	"time"
)

func TestNextDate(t *testing.T) {
	tests := []struct {
		date   string
		repeat string
		want   string
	}{
		{"20240126", "", ""},
		{"20240126", "k 34", ""},
		{"20240126", "ooops", ""},
		{"15000156", "y", ""},
		{"ooops", "y", ""},
		{"16890220", "y", `20240220`},
		{"20250701", "y", `20260701`},
		{"20240101", "y", `20250101`},
		{"20231231", "y", `20241231`},
		{"20240229", "y", `20250301`},
		{"20240301", "y", `20250301`},
		{"20240113", "d", ""},
		{"20240113", "d 7", `20240127`},
		{"20240120", "d 20", `20240209`},
		{"20240202", "d 30", `20240303`},
		{"20240320", "d 401", ""},
		{"20231225", "d 12", `20240130`},
		{"20240228", "d 1", "20240229"},
		{"20231225", "d 12", `20240130`},
	}

	for _, tt := range tests {
		nowStr := "20240126"
		now, _ := time.Parse("20060102", nowStr)
		got, _ := NextDate(now, tt.date, tt.repeat)
		if got != tt.want {
			t.Errorf("NextDate(%v, %v, %v) = %v; want %v", now, tt.date, tt.repeat, got, tt.want)
		}
	}
}
