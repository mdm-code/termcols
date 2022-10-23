/*
Package termcols implements ANSI color codes that can be used to color text on
the terminal. Different styles and foreground/background colors can be chained
together through an intuitive package API to arrive at some cool visual
effects.

The selection of style and color control sequences implemented by the package
was largely based on an exhaustive list of Select Graphic Rendition (SGR)
control sequences available at [Wikipedia ANSI]. It is a great resource in case
one or more elements appear not to be supported in a given terminal.

The Escape sequence for termcols is set to \033, which means that it should
work without any issues with Bash, Zsh or Dash. Other shells might not support
it.

The same applies to 8-bit and 24-bit colors: there is no guarantee that these
escape sequences are supported will be rendered properly on some terminals.
Results may vary, so it is good practice to test it first for compatibility.

The package has two public functions MapColor and MapColors that accept string
values to try and map it onto a valid SgrAttr, however, it has been made
implemented to simplify the terminal tcols command.

# Usage

  package main

  import (
  	"fmt"

	"github.com/mdm-code/termcols"
  )

  func main() {
	s := termcols.Colorize(
		"Colorized text!",
		termcols.RedFg,
		termcols.Underline,
		termcols.Rgb24(termcols.BG, 120, 255, 54),
	)
  	fmt.Println(s)
  }

[Wikipedia ANSI]: https://en.wikipedia.org/wiki/ANSI_escape_code
*/
package termcols
