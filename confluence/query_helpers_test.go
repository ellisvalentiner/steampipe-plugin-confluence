package confluence

import (
	"testing"
)

func TestRequestPageSize(t *testing.T) {
	tests := []struct {
		name            string
		maxItems        int
		streamed        int
		defaultPageSize int
		want            int
	}{
		{"no limit returns default", 0, 0, 100, 100},
		{"no limit with streamed returns default", 0, 50, 100, 100},
		{"exact limit returns remaining", 10, 0, 100, 10},
		{"partial streamed returns remaining", 10, 3, 100, 7},
		{"fully streamed returns zero", 10, 10, 100, 0},
		{"over-streamed returns zero", 10, 15, 100, 0},
		{"remaining smaller than default", 5, 2, 100, 3},
		{"negative maxItems treated as no limit", -1, 0, 100, 100},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := requestPageSize(tt.maxItems, tt.streamed, tt.defaultPageSize)
			if got != tt.want {
				t.Errorf("requestPageSize(%d, %d, %d) = %d; want %d",
					tt.maxItems, tt.streamed, tt.defaultPageSize, got, tt.want)
			}
		})
	}
}
