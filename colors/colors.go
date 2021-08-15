// This package is inspired by the chalk's ansi-styles (JavaScript).
//
// See: https://www.npmjs.com/package/ansi-styles .
// See: https://misc.flogisoft.com/bash/tip_colors_and_formatting .

package colors

type Mode uint8

const (
	Auto Mode = iota
	Never
	Always
)

var Colors Mode = Auto

type Style uint8

func (s Style) String() string {
	// TODO(SuperPaintman): check TTY.
	return styleToString(uint8(s))
}

//go:generate python ./generate_reset_colors.py

func (s Style) Reset() Style {
	return resetStyles[s]
}

const (
	Reset              Style = 0
	ResetBold          Style = 22 // 21 isn't widely supported and 22 does the same thing.
	ResetDim           Style = 22
	ResetItalic        Style = 23
	ResetUnderline     Style = 24
	ResetInverse       Style = 27
	ResetHidden        Style = 28
	ResetStrikethrough Style = 29
	ResetOverline      Style = 55

	Bold          Style = 1
	Dim           Style = 2
	Italic        Style = 3
	Underline     Style = 4
	Inverse       Style = 7
	Hidden        Style = 8
	Strikethrough Style = 9
	Overline      Style = 53
)

const (
	ResetColor Style = 39

	Black   Style = 30
	Red     Style = 31
	Green   Style = 32
	Yellow  Style = 33
	Blue    Style = 34
	Magenta Style = 35
	Cyan    Style = 36
	White   Style = 37

	BlackBright   Style = 90
	RedBright     Style = 91
	GreenBright   Style = 92
	YellowBright  Style = 93
	BlueBright    Style = 94
	MagentaBright Style = 95
	CyanBright    Style = 96
	WhiteBright   Style = 97

	// Aliases.
	Gray Style = BlackBright
)

const (
	ResetBgColor Style = 49

	BgBlack   Style = 40
	BgRed     Style = 41
	BgGreen   Style = 42
	BgYellow  Style = 43
	BgBlue    Style = 44
	BgMagenta Style = 45
	BgCyan    Style = 46
	BgWhite   Style = 47

	BgBlackBright   Style = 100
	BgRedBright     Style = 101
	BgGreenBright   Style = 102
	BgYellowBright  Style = 103
	BgBlueBright    Style = 104
	BgMagentaBright Style = 105
	BgCyanBright    Style = 106
	BgWhiteBright   Style = 107

	// Aliases.
	BgGray Style = BgBlackBright
)

//go:generate python ./generate_ansi_style_string.py

// styleToString converts Style to string.
//
// A hack for fast and inlinable Style to string converion (like what the
// stringer does).
//
// Unfortunately Go can't inline strconv.Itoa (at least Go 1.16) and we can't
// write a regular function. So we need some evil slice hacks :(.
//
// NOTE(SuperPaintman): I would like to remove it in the future.
func styleToString(i uint8) string {
	const (
		size1d = 4 // '\x1b' '[' '0' 'm'
		size2d = 5 // '\x1b' '[' '0' '0' 'm'
		size3d = 6 // '\x1b' '[' '0' '0' '0' 'm'

		offset2d = size1d * 10                // 0 ... 9
		offset3d = offset2d + size2d*(100-10) // 0 ... 9 + 10 ... 99
	)

	if i < 10 {
		return ansiStyleString[int(i)*size1d : (int(i)+1)*size1d]
	}

	if i < 100 {
		return ansiStyleString[offset2d+(int(i)-10)*size2d : offset2d+(int(i)-10+1)*size2d]
	}

	return ansiStyleString[offset3d+(int(i)-100)*size3d : offset3d+(int(i)-100+1)*size3d]
}
