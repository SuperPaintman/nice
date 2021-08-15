// This package is inspired by the chalk's ansi-styles (JavaScript).
//
// See: https://en.wikipedia.org/wiki/ANSI_escape_code .
// See: https://misc.flogisoft.com/bash/tip_colors_and_formatting .
// See: https://www.npmjs.com/package/ansi-styles .

package colors

type Mode uint8

const (
	Auto Mode = iota
	Never
	Always
)

var Colors Mode = Auto

type Attribute uint8

func (a Attribute) String() string {
	// TODO(SuperPaintman): check TTY.
	return attributeToString(uint8(a))
}

//go:generate python ./generate_reset_attributes.py

func (s Attribute) Reset() Attribute {
	return resetAttributes[s]
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
// NOTE(SuperPaintman): I would like to remove it in the future.
func attributeToString(i uint8) string {
	const (
		size1d = 4 // '\x1b' '[' '0' 'm'
		size2d = 5 // '\x1b' '[' '0' '0' 'm'
		size3d = 6 // '\x1b' '[' '0' '0' '0' 'm'

		offset2d = size1d * 10                // 0 ... 9
		offset3d = offset2d + size2d*(100-10) // 0 ... 9 + 10 ... 99
	)

	if i < 10 {
		return ansiAttributeString[int(i)*size1d : (int(i)+1)*size1d]
	}

	if i < 100 {
		return ansiAttributeString[offset2d+(int(i)-10)*size2d : offset2d+(int(i)-10+1)*size2d]
	}

	return ansiAttributeString[offset3d+(int(i)-100)*size3d : offset3d+(int(i)-100+1)*size3d]
}
