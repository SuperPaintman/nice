package colors

import (
	"strconv"
	"testing"
)

func init() {
	// TODO(SuperPaintman): add tests for cases when terminal does not support colors.
	SetMode(Always | ForceANSI256 | ForceTrueColor)
}

func TestComputeShouldUse(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name          string
		mode          Mode
		supports      termSupports
		wantColors    bool
		wantANSI256   bool
		wantTrueColor bool
	}{
		{
			name:          "mode auto and term supports nothing",
			mode:          Auto,
			supports:      0,
			wantColors:    false,
			wantANSI256:   false,
			wantTrueColor: false,
		},
		{
			name:          "mode auto and term supports colors",
			mode:          Auto,
			supports:      supportsColor,
			wantColors:    true,
			wantANSI256:   false,
			wantTrueColor: false,
		},
		{
			name:          "mode auto and term supports ANSI256",
			mode:          Auto,
			supports:      supportsColor | supportsANSI256,
			wantColors:    true,
			wantANSI256:   true,
			wantTrueColor: false,
		},
		{
			name:          "mode auto and term supports true colors",
			mode:          Auto,
			supports:      supportsColor | supportsTrueColor,
			wantColors:    true,
			wantANSI256:   false,
			wantTrueColor: true,
		},
		{
			name:          "mode auto and term supports all",
			mode:          Auto,
			supports:      supportsColor | supportsANSI256 | supportsTrueColor,
			wantColors:    true,
			wantANSI256:   true,
			wantTrueColor: true,
		},
		{
			name:          "mode auto and term supports all",
			mode:          Auto,
			supports:      supportsColor | supportsANSI256 | supportsTrueColor,
			wantColors:    true,
			wantANSI256:   true,
			wantTrueColor: true,
		},
		{
			name:          "mode never and term supports all",
			mode:          Never,
			supports:      supportsColor | supportsANSI256 | supportsTrueColor,
			wantColors:    false,
			wantANSI256:   false,
			wantTrueColor: false,
		},
		{
			name:          "mode always and term supports nothing",
			mode:          Always,
			supports:      0,
			wantColors:    true,
			wantANSI256:   false,
			wantTrueColor: false,
		},
		{
			name:          "mode always and term supports nothing and force ANSI256",
			mode:          Always | ForceANSI256,
			supports:      0,
			wantColors:    true,
			wantANSI256:   true,
			wantTrueColor: false,
		},
		{
			name:          "mode always and term supports nothing and force true color",
			mode:          Always | ForceTrueColor,
			supports:      0,
			wantColors:    true,
			wantANSI256:   false,
			wantTrueColor: true,
		},
		{
			name:          "mode always and term supports nothing and force all",
			mode:          Always | ForceANSI256 | ForceTrueColor,
			supports:      0,
			wantColors:    true,
			wantANSI256:   true,
			wantTrueColor: true,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			gotColors, gotANSI256, gotTrueColor := computeShouldUse(tc.mode, tc.supports)

			if gotColors != tc.wantColors {
				t.Errorf("computeShouldUse(%04b, %03b): gotColors: got = %v, want = %v",
					tc.mode, tc.supports, gotColors, tc.wantColors,
				)
			}

			if gotANSI256 != tc.wantANSI256 {
				t.Errorf("computeShouldUse(%04b, %03b): gotANSI256: got = %v, want = %v",
					tc.mode, tc.supports, gotANSI256, tc.wantANSI256,
				)
			}

			if gotTrueColor != tc.wantTrueColor {
				t.Errorf("computeShouldUse(%04b, %03b): gotTrueColor: got = %v, want = %v",
					tc.mode, tc.supports, gotTrueColor, tc.wantTrueColor,
				)
			}
		})
	}
}

type testEnv map[string]string

func (e testEnv) String() string {
	buf := ""
	for k, v := range e {
		if buf != "" {
			buf += " "
		}

		buf += k + "=" + v
	}

	return buf
}

func (e testEnv) LookupEnv(key string) (string, bool) {
	v, ok := e[key]
	return v, ok
}

func TestTerminalSupports(t *testing.T) {
	tt := []struct {
		name string
		env  testEnv
		want termSupports
	}{
		{
			name: "TERM=dumb",
			env: testEnv{
				"TERM": "dumb",
			},
			want: 0,
		},
		{
			name: "TERM=dumb with COLORTERM",
			env: testEnv{
				"TERM":      "dumb",
				"COLORTERM": "truecolor",
			},
			want: 0,
		},
		{
			name: "COLORTERM=<any>",
			env: testEnv{
				"COLORTERM": "test",
			},
			want: supportsColor,
		},
		{
			name: "COLORTERM=truecolor",
			env: testEnv{
				"COLORTERM": "truecolor",
			},
			want: supportsColor | supportsANSI256 | supportsTrueColor,
		},
		{
			name: "TERM_PROGRAM=Apple_Terminal",
			env: testEnv{
				"TERM_PROGRAM": "Apple_Terminal",
			},
			want: supportsColor | supportsANSI256 | supportsTrueColor,
		},
		{
			name: "CI=<any> TRAVIS=<any>",
			env: testEnv{
				"CI":     "test",
				"TRAVIS": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> CIRCLECI=<any>",
			env: testEnv{
				"CI":       "test",
				"CIRCLECI": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> APPVEYOR=<any>",
			env: testEnv{
				"CI":       "test",
				"APPVEYOR": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> GITLAB_CI=<any>",
			env: testEnv{
				"CI":        "test",
				"GITLAB_CI": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> GITHUB_ACTIONS=<any>",
			env: testEnv{
				"CI":             "test",
				"GITHUB_ACTIONS": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> BUILDKITE=<any>",
			env: testEnv{
				"CI":        "test",
				"BUILDKITE": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> DRONE=<any>",
			env: testEnv{
				"CI":    "test",
				"DRONE": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> TRAVIS=<any>",
			env: testEnv{
				"CI":     "test",
				"TRAVIS": "test",
			},
			want: supportsColor,
		},
		{
			name: "CI=<any> CI_NAME=codeship",
			env: testEnv{
				"CI":      "test",
				"CI_NAME": "codeship",
			},
			want: supportsColor,
		},
	}

	for _, tc := range tt {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := terminalSupports(tc.env.LookupEnv)
			if got != tc.want {
				t.Errorf("terminalSupports(%s): got = %03b, want = %03b",
					tc.env, got, tc.want,
				)
			}
		})
	}
}

func TestAttribute_Reset(t *testing.T) {
	t.Parallel()

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
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := tc.attribute.Reset()
			if got != tc.want {
				t.Errorf("Reset(): got = %q, want = %q", got, tc.want)
			}
		})
	}
}

const maxUint8 = int(^uint8(0))

func TestAttributeToString(t *testing.T) {
	t.Parallel()

	for i := 0; i < maxUint8+1; i++ {
		got := attributeToString(uint8(i))

		want := "\x1b[" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("attributeToString(%d): got = %q, want = %q", i, got, want)
		}
	}
}

func TestANSI256(t *testing.T) {
	t.Parallel()

	for i := 0; i < maxUint8+1; i++ {
		got := ANSI256(uint8(i))

		want := "\x1b[38;5;" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("ANSI256(%d): got = %q, want = %q", i, got, want)
		}
	}
}

func TestBgANSI256(t *testing.T) {
	t.Parallel()

	for i := 0; i < maxUint8+1; i++ {
		got := BgANSI256(uint8(i))

		want := "\x1b[48;5;" + strconv.Itoa(i) + "m"

		if got != want {
			t.Errorf("BgANSI256(%d): got = %q, want = %q", i, got, want)
		}
	}
}

const (
	trueColorMax  = 1<<24 - 1
	trueColorStep = 256 - 9 // We use this step to make move overlaps.
)

func TestTrueColor(t *testing.T) {
	t.Parallel()

	for i := 0; i < trueColorMax+1; i += trueColorStep {
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
			fn := t.Fatalf
			if testing.Verbose() {
				fn = t.Errorf
			}

			fn("TrueColor(%d, %d, %d): got = %q, want = %q", r, g, b, got, want)
		}
	}
}

func TestBgTrueColor(t *testing.T) {
	t.Parallel()

	for i := 0; i < trueColorMax+1; i += trueColorStep {
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
			fn := t.Fatalf
			if testing.Verbose() {
				fn = t.Errorf
			}

			fn("BgTrueColor(%d, %d, %d): got = %q, want = %q", r, g, b, got, want)
		}
	}
}

func TestTrueColorRGB(t *testing.T) {
	t.Parallel()

	for i := 0; i < trueColorMax+1; i += trueColorStep {
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
			fn := t.Fatalf
			if testing.Verbose() {
				fn = t.Errorf
			}

			fn("TrueColorRGB(%v): got = %q, want = %q", rgb, got, want)
		}
	}
}

func TestBgTrueColorRGB(t *testing.T) {
	t.Parallel()

	for i := 0; i < trueColorMax+1; i += trueColorStep {
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
			fn := t.Fatalf
			if testing.Verbose() {
				fn = t.Errorf
			}

			fn("BgTrueColorRGB(%v): got = %q, want = %q", rgb, got, want)
		}
	}
}

var attributeStringRes string

func BenchmarkAttribute_String(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for c := 0; c <= 255; c++ {
			attributeStringRes = Attribute(c).String()
		}
	}
}
