package colors

import (
	"strconv"
	"testing"
)

func TestAttribute_Reset(t *testing.T) {
	tt := []struct {
		name      string
		attribute Attribute
		want      Attribute
	}{
		{
			name:      "Reset",
			attribute: Reset,
			want:      Reset,
		},
		{
			name:      "Bold",
			attribute: Bold,
			want:      ResetBold,
		},
		{
			name:      "Red",
			attribute: Red,
			want:      ResetColor,
		},
		{
			name:      "BgRed",
			attribute: BgRed,
			want:      ResetBgColor,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.attribute.Reset()
			if got != tc.want {
				t.Errorf("Reset(): got = %q, want = %q", got, tc.want)
			}
		})
	}
}

func TestAttributeToString(t *testing.T) {
	const maxUint8 = int(^uint8(0))

	for i := 0; i < maxUint8+1; i++ {
		got := attributeToString(uint8(i))

		want := "\x1b[" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("%d: got = %q, want = %q", i, got, want)
		}
	}
}
