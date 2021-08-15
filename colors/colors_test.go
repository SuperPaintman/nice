package colors

import (
	"strconv"
	"testing"
)

func TestStyle_Reset(t *testing.T) {
	tt := []struct {
		name  string
		style Style
		want  Style
	}{
		{
			name:  "Reset",
			style: Reset,
			want:  Reset,
		},
		{
			name:  "Bold",
			style: Bold,
			want:  ResetBold,
		},
		{
			name:  "Red",
			style: Red,
			want:  ResetColor,
		},
		{
			name:  "BgRed",
			style: BgRed,
			want:  ResetBgColor,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.style.Reset()
			if got != tc.want {
				t.Errorf("Reset(): got = %q, want = %q", got, tc.want)
			}
		})
	}
}

func TestStyleToString(t *testing.T) {
	const maxUint8 = int(^uint8(0))

	for i := 0; i < maxUint8+1; i++ {
		got := styleToString(uint8(i))

		want := "\x1b[" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("%d: got = %q, want = %q", i, got, want)
		}
	}
}
