package termcols_test

import (
	"fmt"

	"github.com/mdm-code/termcols"
)

// ExampleColorize shows how to use termcols public API Colorize function.
// It uses a combination of style, foreground and colors to stylize the
// ``Colorize text\!'' string.
func ExampleColorize() {
	s := termcols.Colorize(
		"Colorized text!",
		termcols.RedFg,
		termcols.Underline,
		termcols.Rgb24(termcols.BG, 120, 255, 54),
	)
	fmt.Println(s)
	// Output: [31m[4m[48;2;120;255;54mColorized text![0m
}
