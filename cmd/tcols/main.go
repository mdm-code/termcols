/*
tcols - add color to text on the terminal

Tcols reads text from a file and writes the colorized text to the standard
output.

Usage:
	tcols [options] [file]

Options:
	-h, --help   show this help message and exit
	-s, --style  specify styles and colors to apply

Example:
	tcols bold < <(echo 'Hello, world!')

Output:
	\033[1mHello, World!\033[0m

The program returns text read from a file with Select Graphic Rendition control
sequences prepended and the reset control sequence appended at the end. The
sequence of attributes passed to the --style flag of the command is preserved,
so colors and styles can (un)intentionally cancel out one another.
*/
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mdm-code/termcols"
)

const usage = `tcols - add color to text on the terminal

Tcols reads text from a file and writes the colorized text to the standard
output.

Usage:
	tcols [options] [file]

Options:
	-h, --help   show this help message and exit
	-s, --style  specify styles and colors to apply

Example:
	tcols --style bold < <(echo 'Hello, world!')

Output:
	\033[1mHello, World!\033[0m

The program returns text read from a file with Select Graphic Rendition control
sequences prepended and the reset control sequence appended at the end. The
sequence of attributes passed to the --style flag of the command is preserved,
so colors and styles can (un)intentionally cancel out one another.
`

const (
	exitSuccess = iota
	exitFailure
)

var (
	style string
)

// TODO: Add a reference on `go doc tcols` in the `README.md` file.
// TODO: Add run function with string return so that it can be Example tested
// TODO Style should be evaluated with a map string:SgrAttr and a function
func main() {
	flag.StringVar(&style, "s", "", "")
	flag.StringVar(&style, "style", "", "")
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()

	out := bufio.NewWriter(os.Stdout)

	text, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(exitFailure)
	}

	output := termcols.Colorize(string(text), termcols.Bold)
	_, err = out.WriteString(output)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(exitFailure)
	}

	if err := out.Flush(); err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(exitFailure)
	}
	os.Exit(exitSuccess)
}
