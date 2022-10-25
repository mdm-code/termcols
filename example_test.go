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

func ExampleRgb8() {
	attr := termcols.Rgb8(termcols.FG, 12)
	fmt.Printf("%sColorized text!%s", attr, termcols.Reset)
	// Output: [38;5;12mColorized text![0m
}

func ExampleRgb24() {
	attr := termcols.Rgb24(termcols.BG, 221, 42, 89)
	fmt.Printf("%sColorized text!%s", attr, termcols.Reset)
	// Output: [48;2;221;42;89mColorized text![0m
}

func ExampleMapColor() {
	attr, _ := termcols.MapColor("bluefg")
	fmt.Printf("%sColorized text!%s", attr, termcols.Reset)
	// Output: [34mColorized text![0m
}

func ExampleMapColors() {
	attrs, _ := termcols.MapColors([]string{"bluefg", "yellowbg", "strike"})
	fmt.Printf("%s%s%sColorized text!%s", attrs[0], attrs[1], attrs[2], termcols.Reset)
	// Output: [34m[43m[9mColorized text![0m
}
