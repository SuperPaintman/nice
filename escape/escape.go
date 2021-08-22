// See: https://en.wikipedia.org/wiki/ANSI_escape_code .

package escape

import "strconv"

const (
	ESC = "\x1b"    // Starts all the escape sequences.
	CSI = ESC + "[" // Control Sequence Introducer.
)

// CursorUp moves the cursor n cells in the given direction.
// If the cursor is already at the edge of the screen, this has no effect.
func CursorUp(n int) string {
	if n <= 1 {
		return CSI + "1A"
	}

	return CSI + strconv.Itoa(n) + "A"
}

// CursorDown moves the cursor n cells in the given direction.
// If the cursor is already at the edge of the screen, this has no effect.
func CursorDown(n int) string {
	if n <= 1 {
		return CSI + "1B"
	}

	return CSI + strconv.Itoa(n) + "B"
}

// CursorForward moves the cursor n cells in the given direction.
// If the cursor is already at the edge of the screen, this has no effect.
func CursorForward(n int) string {
	if n <= 1 {
		return CSI + "1C"
	}

	return CSI + strconv.Itoa(n) + "C"
}

// CursorBackward moves the cursor n cells in the given direction.
// If the cursor is already at the edge of the screen, this has no effect.
func CursorBackward(n int) string {
	if n <= 1 {
		return CSI + "1D"
	}

	return CSI + strconv.Itoa(n) + "D"
}

// CursorNextLine moves the cursor to beginning of the line n lines down.
func CursorNextLine(n int) string {
	if n <= 1 {
		return CSI + "1E"
	}

	return CSI + strconv.Itoa(n) + "E"
}

// CursorPreviousLine moves the cursor to beginning of the line n lines up.
func CursorPreviousLine(n int) string {
	if n <= 1 {
		return CSI + "1F"
	}

	return CSI + strconv.Itoa(n) + "F"
}

// CursorHorizontalAbsolute moves the cursor to column n.
func CursorHorizontalAbsolute(n int) string {
	if n <= 1 {
		return CSI + "1G"
	}

	return CSI + strconv.Itoa(n) + "G"
}

const (
	CursorLeft = CSI + "G" // Moves the cusor to beginning of the line.
)

// CursorPosition moves the cursor to row n, column m.
// Moves the cursor to row n, column m. The values are 1-based.
func CursorPosition(n, m int) string {
	if n < 1 {
		n = 1
	}
	if m < 1 {
		m = 1
	}

	return CSI + strconv.Itoa(n) + ";" + strconv.Itoa(m) + "H"
}

func CursorTo(x, y int) string {
	if y == 0 {
		return CursorHorizontalAbsolute(x + 1)
	}

	return CursorPosition(y+1, x+1)
}

// TODO(SuperPaintman): optimize it.
func CursorMove(x, y int) string {
	res := ""

	if x < 0 {
		res += CursorBackward(-1 * x)
	} else if x > 0 {
		res += CursorForward(x)
	}

	if y < 0 {
		res += CursorUp(-1 * y)
	} else if y > 0 {
		res += CursorDown(y)
	}

	return res
}

const (
	EraseDown      = CSI + "J"  // Clears from cursor to end of screen.
	EraseUp        = CSI + "1J" // Clears from cursor to beginning of the screen.
	EraseScreen    = CSI + "2J" // Clears entire screen.
	EraseLineEnd   = CSI + "K"  // Clears from cursor to the end of the line.
	EraseLineStart = CSI + "1K" // Clears from cursor to beginning of the line.
	EraseLine      = CSI + "2K" // Clears entire line. Cursor position does not change.
)

// TODO(SuperPaintman): add "S" and "T" for the scroll.

// Popular private sequences.

const (
	CursorShow = CSI + "?25h" // Shows the cursor, from the VT220.
	CursorHide = CSI + "?25l" // Hides the cursor.
	// CursorSave    = ESC + "7"
	// CursorRestore = ESC + "8"
)

const (
	ClearScreen = ESC + "c" // Triggers a full reset of the terminal to its original state.
)
