package colors

import (
	"strconv"
	"testing"
)

// TODO(SuperPaintman): change it on the CI.
const trueColorTestStep = 1<<8 - 1

func init() {
	// TODO(SuperPaintman): add tests for cases when terminal does not support colors.
	supportsColor = true
	supportsANSI256 = true
	supportsTrueColor = true
}

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

const maxUint8 = int(^uint8(0))

func TestAttributeToString(t *testing.T) {
	for i := 0; i < maxUint8+1; i++ {
		got := attributeToString(uint8(i))

		want := "\x1b[" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("attributeToString(%d): got = %q, want = %q", i, got, want)
		}
	}
}

func TestANSI256(t *testing.T) {
	for i := 0; i < maxUint8+1; i++ {
		got := ANSI256(uint8(i))

		want := "\x1b[38;5;" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("ANSI256(%d): got = %q, want = %q", i, got, want)
		}
	}
}

func TestBgANSI256(t *testing.T) {
	for i := 0; i < maxUint8+1; i++ {
		got := BgANSI256(uint8(i))

		want := "\x1b[48;5;" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("BgANSI256(%d): got = %q, want = %q", i, got, want)
		}
	}
}

const trueColorMax = 1<<24 - 1

func TestTrueColor(t *testing.T) {
	for i := 0; i < trueColorMax+1; i += trueColorTestStep {
		r := uint8((i >> 0) & 255)
		g := uint8((i >> 8) & 255)
		b := uint8((i >> 16) & 255)

		got := TrueColor(r, g, b)

		want := "\x1b[38;2;" +
			strconv.Itoa(int(r)) + ";" +
			strconv.Itoa(int(g)) + ";" +
			strconv.Itoa(int(b)) + "m"

		if got != want {
			// Fatal because we don't want to see ~16 million errors.
			t.Fatalf("TrueColor(%d, %d, %d): got = %q, want = %q", r, g, b, got, want)
		}
	}
}

func TestBgTrueColor(t *testing.T) {
	for i := 0; i < trueColorMax+1; i += trueColorTestStep {
		r := uint8((i >> 0) & 255)
		g := uint8((i >> 8) & 255)
		b := uint8((i >> 16) & 255)

		got := BgTrueColor(r, g, b)

		want := "\x1b[48;2;" +
			strconv.Itoa(int(r)) + ";" +
			strconv.Itoa(int(g)) + ";" +
			strconv.Itoa(int(b)) + "m"

		if got != want {
			// Fatal because we don't want to see ~16 million errors.
			t.Fatalf("BgTrueColor(%d, %d, %d): got = %q, want = %q", r, g, b, got, want)
		}
	}
}

func TestTrueColorRGB(t *testing.T) {
	for i := 0; i < trueColorMax+1; i += trueColorTestStep {
		rgb := RGB{
			R: uint8((i >> 0) & 255),
			G: uint8((i >> 8) & 255),
			B: uint8((i >> 16) & 255),
		}

		got := TrueColorRGB(rgb)

		want := "\x1b[38;2;" +
			strconv.Itoa(int(rgb.R)) + ";" +
			strconv.Itoa(int(rgb.G)) + ";" +
			strconv.Itoa(int(rgb.B)) + "m"

		if got != want {
			// Fatal because we don't want to see ~16 million errors.
			t.Fatalf("TrueColorRGB(%v): got = %q, want = %q", rgb, got, want)
		}
	}
}

func TestBgTrueColorRGB(t *testing.T) {
	for i := 0; i < trueColorMax+1; i += trueColorTestStep {
		rgb := RGB{
			R: uint8((i >> 0) & 255),
			G: uint8((i >> 8) & 255),
			B: uint8((i >> 16) & 255),
		}

		got := BgTrueColorRGB(rgb)

		want := "\x1b[48;2;" +
			strconv.Itoa(int(rgb.R)) + ";" +
			strconv.Itoa(int(rgb.G)) + ";" +
			strconv.Itoa(int(rgb.B)) + "m"

		if got != want {
			// Fatal because we don't want to see ~16 million errors.
			t.Fatalf("BgTrueColorRGB(%v): got = %q, want = %q", rgb, got, want)
		}
	}
}
