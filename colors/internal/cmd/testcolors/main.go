package main

import (
	"fmt"

	"github.com/SuperPaintman/nice/colors"
)

const trueColorEnabled = false

func main() {
	// Modifier.
	fmt.Printf(":%s[ X ]%s: ", colors.Bold, colors.Bold.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Dim, colors.Dim.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Italic, colors.Italic.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Underline, colors.Underline.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Inverse, colors.Inverse.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Hidden, colors.Hidden.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Strikethrough, colors.Strikethrough.Reset())
	fmt.Printf(":%s[ X ]%s:\n", colors.Overline, colors.Overline.Reset())

	// Color.
	fmt.Printf(":%s[ X ]%s: ", colors.Black, colors.Black.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Red, colors.Red.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Green, colors.Green.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Yellow, colors.Yellow.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Blue, colors.Blue.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Magenta, colors.Magenta.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.Cyan, colors.Cyan.Reset())
	fmt.Printf(":%s[ X ]%s:\n", colors.White, colors.White.Reset())

	// Color Bright.
	fmt.Printf(":%s[ X ]%s: ", colors.BlackBright, colors.BlackBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.RedBright, colors.RedBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.GreenBright, colors.GreenBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.YellowBright, colors.YellowBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BlueBright, colors.BlueBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.MagentaBright, colors.MagentaBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.CyanBright, colors.CyanBright.Reset())
	fmt.Printf(":%s[ X ]%s:\n", colors.WhiteBright, colors.WhiteBright.Reset())

	// BgColor.
	fmt.Printf(":%s[ X ]%s: ", colors.BgBlack, colors.BgBlack.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgRed, colors.BgRed.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgGreen, colors.BgGreen.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgYellow, colors.BgYellow.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgBlue, colors.BgBlue.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgMagenta, colors.BgMagenta.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgCyan, colors.BgCyan.Reset())
	fmt.Printf(":%s[ X ]%s:\n", colors.BgWhite, colors.BgWhite.Reset())

	// BgColor Bright.
	fmt.Printf(":%s[ X ]%s: ", colors.BgBlackBright, colors.BgBlackBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgRedBright, colors.BgRedBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgGreenBright, colors.BgGreenBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgYellowBright, colors.BgYellowBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgBlueBright, colors.BgBlueBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgMagentaBright, colors.BgMagentaBright.Reset())
	fmt.Printf(":%s[ X ]%s: ", colors.BgCyanBright, colors.BgCyanBright.Reset())
	fmt.Printf(":%s[ X ]%s:\n", colors.BgWhiteBright, colors.BgWhiteBright.Reset())

	// ANSI 256.
	for i := 0; i < 256; i++ {
		if i != 0 && i%16 == 0 {
			fmt.Print("\n")
		} else if i != 0 {
			fmt.Print(" ")
		}

		fmt.Printf(":%s[ %03d ]%s:", colors.ANSI256(uint8(i)), i, colors.Reset)
	}
	fmt.Print("\n")

	// BgANSI 256.
	for i := 0; i < 256; i++ {
		if i != 0 && i%16 == 0 {
			fmt.Print("\n")
		} else if i != 0 {
			fmt.Print(" ")
		}

		fmt.Printf(":%s[ %3d ]%s:", colors.BgANSI256(uint8(i)), i, colors.Reset)
	}
	fmt.Print("\n")

	// TrueColor.
	if trueColorEnabled {
		for i := 0; i < 1<<24; i++ {
			if i != 0 && i%16 == 0 {
				fmt.Print("\n")
			} else if i != 0 {
				fmt.Print(" ")
			}

			r := uint8((i >> 0) & 255)
			g := uint8((i >> 8) & 255)
			b := uint8((i >> 16) & 255)

			fmt.Printf(":%s[ %3d %3d %3d ]%s:", colors.TrueColor(r, g, b), r, g, b, colors.Reset)
		}
		fmt.Print("\n")
	}

	// BgTrueColor.
	if trueColorEnabled {
		for i := 0; i < 1<<24; i++ {
			if i != 0 && i%16 == 0 {
				fmt.Print("\n")
			} else if i != 0 {
				fmt.Print(" ")
			}

			r := uint8((i >> 0) & 255)
			g := uint8((i >> 8) & 255)
			b := uint8((i >> 16) & 255)

			fmt.Printf(":%s[ %3d %3d %3d ]%s:", colors.BgTrueColor(r, g, b), r, g, b, colors.Reset)
		}
		fmt.Print("\n")
	}
}
