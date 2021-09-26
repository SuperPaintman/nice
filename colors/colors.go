// This package is inspired by the chalk's ansi-styles and supports-color
// (JavaScript).
//
// See: https://en.wikipedia.org/wiki/ANSI_escape_code .
// See: https://misc.flogisoft.com/bash/tip_colors_and_formatting .
// See: https://www.npmjs.com/package/ansi-styles .
// See: https://www.npmjs.com/package/supports-color .
// See: https://github.com/termstandard/colors .

package colors

import (
	"os"
	"strconv"

	"github.com/mattn/go-isatty"
)

type termSupports uint8

const (
	supportsColor termSupports = 1 << iota
	supportsANSI256
	supportsTrueColor
)

func terminalSupports(lookup func(key string) (string, bool)) (s termSupports) {
	// Terminals.
	if v, ok := lookup("TERM"); ok && v == "dumb" {
		return
	}

	if colorTerm, ok := lookup("COLORTERM"); ok {
		s |= supportsColor

		if colorTerm == "truecolor" {
			s |= supportsANSI256
			s |= supportsTrueColor
		}
		return
	}

	if termProg, ok := lookup("TERM_PROGRAM"); ok {
		if termProg == "Apple_Terminal" {
			s |= supportsColor
			s |= supportsANSI256
			s |= supportsTrueColor
			return
		}

		// TODO(SuperPaintman): add iTerm.app
	}

	// TODO(SuperPaintman): /-256(color)?$/i

	// TODO(SuperPaintman): /^screen|^xterm|^vt100|^vt220|^rxvt|color|ansi|cygwin|linux/i

	// TODO(SuperPaintman): add win32 checker.

	// CI.
	if _, ok := lookup("CI"); ok {
		cis := [...]string{
			"TRAVIS",
			"CIRCLECI",
			"APPVEYOR",
			"GITLAB_CI",
			"GITHUB_ACTIONS",
			"BUILDKITE",
			"DRONE",
		}

		for _, name := range cis {
			if _, ok := lookup(name); ok {
				s |= supportsColor
				return
			}
		}

		if v, ok := lookup("CI_NAME"); ok && v == "codeship" {
			s |= supportsColor
			return
		}
	}

	// TODO(SuperPaintman): add TeamCity checker.

	return
}

func SupportsColor() bool     { return supports&supportsColor != 0 }
func SupportsANSI256() bool   { return supports&supportsANSI256 != 0 }
func SupportsTrueColor() bool { return supports&supportsTrueColor != 0 }

type Mode uint8

const (
	Auto Mode = 1 << iota >> 1
	Never
	Always
	ForceANSI256
	ForceTrueColor
)

var (
	supports                                              termSupports
	mode                                                  Mode
	shouldUseColors, shouldUseANSI256, shouldUseTrueColor bool
)

func init() {
	// Check is TTY.
	isTTY := isatty.IsTerminal(os.Stdout.Fd()) ||
		isatty.IsCygwinTerminal(os.Stdout.Fd())

	if isTTY {
		supports = terminalSupports(os.LookupEnv)
	}

	shouldUseColors, shouldUseANSI256, shouldUseTrueColor = computeShouldUse(mode, supports)
}

func SetMode(m Mode) {
	mode = m
	shouldUseColors, shouldUseANSI256, shouldUseTrueColor = computeShouldUse(mode, supports)
}

func computeShouldUse(m Mode, s termSupports) (colors bool, ansi256 bool, trueColor bool) {
	if m&Never == 0 {
		colors = m&Always != 0 || s&supportsColor != 0

		ansi256 = colors &&
			(s&supportsANSI256 != 0 || m&ForceANSI256 != 0)

		trueColor = colors &&
			(s&supportsTrueColor != 0 || m&ForceTrueColor != 0)
	}

	return
}

type Attribute uint8

func (a Attribute) String() string {
	if !shouldUseColors {
		return ""
	}

	return attributeToString(uint8(a))
}

//go:generate python ./generate_reset_attributes.py

func (a Attribute) Reset() Attribute {
	return resetAttributes[a]
}

const (
	Reset              Attribute = 0
	ResetBold          Attribute = 22 // 21 isn't widely supported and 22 does the same thing.
	ResetDim           Attribute = 22
	ResetItalic        Attribute = 23
	ResetUnderline     Attribute = 24
	ResetInverse       Attribute = 27
	ResetHidden        Attribute = 28
	ResetStrikethrough Attribute = 29
	ResetOverline      Attribute = 55

	Bold          Attribute = 1
	Dim           Attribute = 2
	Italic        Attribute = 3
	Underline     Attribute = 4
	Inverse       Attribute = 7
	Hidden        Attribute = 8
	Strikethrough Attribute = 9
	Overline      Attribute = 53
)

const (
	ResetColor Attribute = 39

	Black   Attribute = 30
	Red     Attribute = 31
	Green   Attribute = 32
	Yellow  Attribute = 33
	Blue    Attribute = 34
	Magenta Attribute = 35
	Cyan    Attribute = 36
	White   Attribute = 37

	BlackBright   Attribute = 90
	RedBright     Attribute = 91
	GreenBright   Attribute = 92
	YellowBright  Attribute = 93
	BlueBright    Attribute = 94
	MagentaBright Attribute = 95
	CyanBright    Attribute = 96
	WhiteBright   Attribute = 97

	// Aliases.
	Gray Attribute = BlackBright
)

const (
	ResetBgColor Attribute = 49

	BgBlack   Attribute = 40
	BgRed     Attribute = 41
	BgGreen   Attribute = 42
	BgYellow  Attribute = 43
	BgBlue    Attribute = 44
	BgMagenta Attribute = 45
	BgCyan    Attribute = 46
	BgWhite   Attribute = 47

	BgBlackBright   Attribute = 100
	BgRedBright     Attribute = 101
	BgGreenBright   Attribute = 102
	BgYellowBright  Attribute = 103
	BgBlueBright    Attribute = 104
	BgMagentaBright Attribute = 105
	BgCyanBright    Attribute = 106
	BgWhiteBright   Attribute = 107

	// Aliases.
	BgGray Attribute = BgBlackBright
)

//go:generate python ./generate_ansi_attribute_string.py

// attributeToString converts Attribute to string.
//
// A hack for fast and inlinable Attribute to string converion (like what the
// stringer does).
//
// Unfortunately Go can't inline strconv.Itoa (at least Go 1.16) and we can't
// write a regular function. So we need some evil slice hacks :(.
//
// If you look in the git log you will find a faster version of this function
// but it has fewer ways to be inlined.
//
// NOTE(SuperPaintman): I would like to remove it in the future.
func attributeToString(i uint8) string {
	return ansiAttributeString[ansiAttributeIndex[uint16(i)]:ansiAttributeIndex[uint16(i)+1]]
}

// TODO(SuperPaintman): make it inlinable.
func ANSI256(color uint8) string {
	if !shouldUseANSI256 {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[38;5;" + strconv.Itoa(int(color)) + "m"
}

// TODO(SuperPaintman): make it inlinable.
func BgANSI256(color uint8) string {
	if !shouldUseANSI256 {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[48;5;" + strconv.Itoa(int(color)) + "m"
}

// 24-bit or truecolor or ANSI 16 millions.
func TrueColor(r, g, b uint8) string {
	if !shouldUseTrueColor {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice.
	return "\x1b[38;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m"
}

func BgTrueColor(r, g, b uint8) string {
	if !shouldUseTrueColor {
		return ""
	}

	// TODO(SuperPaintman): optimize it with a preassembled slice for uint8.
	return "\x1b[48;2;" +
		strconv.Itoa(int(r)) + ";" +
		strconv.Itoa(int(g)) + ";" +
		strconv.Itoa(int(b)) + "m"
}

type RGB struct {
	R, G, B uint8
}

func TrueColorRGB(color RGB) string {
	return TrueColor(color.R, color.G, color.B)
}

func BgTrueColorRGB(color RGB) string {
	return BgTrueColor(color.R, color.G, color.B)
}
