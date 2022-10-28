/*
tcols - add color to text on the terminal

Tcols reads text from a file and writes the colorized text to the standard
output.

Usage:
	tcols [-s|--style arg...] [file...]

Options:
	-h, --help   show this help message and exit
	-s, --style  list of styles and colors to apply to text

Example:
	tcols -style 'bold bluefg' < <(echo -n 'Hello, world!')

Output:
	Raw: \033[1m\033[34mHello, World!\033[0m

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
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/mdm-code/termcols"
)

const (
	exitSuccess exitCode = iota
	exitFailure
)

var (
	styles []string
	usage  = fmt.Sprintf(`tcols - add color to text on the terminal

Tcols reads text from a file and writes the colorized text to the standard
output.

Usage:
	tcols [-s|--style arg...] [file...]

Options:
	-h, --help   show this help message and exit
	-s, --style  list of styles and colors to apply to text

Example:
	tcols -style 'bold bluefg' < <(echo -n 'Hello, world!')

Output:
	[1m[34mHello, World![0m

The program returns text read from a file with Select Graphic Rendition control
sequences prepended and the reset control sequence appended at the end. The
sequence of attributes passed to the --style flag of the command is preserved,
so colors and styles can (un)intentionally cancel out one another.

Styles:
	%s %s %s %s
	%s %s %s %s
	%s

Colors:
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s %s %s
	%s %s

Rgb8:
	%s %s
Rgb24:
	%s %s
`,
		`[1mbold[0m`,
		`[2mfaint[0m`,
		`[3mitalic[0m`,
		`[4munderline[0m`,
		`[5mblink[0m`,
		`[7mreverse[0m`,
		`[8mhide[0m`,
		`[9mstrike[0m`,
		`[10mdefaultstyle[0m`,
		`[30mblackfg[0m`,
		`[90mblackbfg[0m`,
		`[40mblackbg[0m`,
		`[100mblackbbg[0m`,
		`[31mredfg[0m`,
		`[91mredbbfg[0m`,
		`[41mredbg[0m`,
		`[101mredbbbg[0m`,
		`[32mgreenfg[0m`,
		`[92mgreenbfg[0m`,
		`[42mgreenbg[0m`,
		`[102mgreenbbg[0m`,
		`[33myellowfg[0m`,
		`[93myellowbfg[0m`,
		`[43myellowbg[0m`,
		`[103myellowbbg[0m`,
		`[34mbluefg[0m`,
		`[94mbluebfg[0m`,
		`[44mbluebg[0m`,
		`[104mbluebbg[0m`,
		`[35mmagentafg[0m`,
		`[95mmagentabfg[0m`,
		`[45mmagentabg[0m`,
		`[105mmagentabbg[0m`,
		`[36mcyanfg[0m`,
		`[96mcyanbfg[0m`,
		`[46mcyanbg[0m`,
		`[106mcyanbbg[0m`,
		`[37mwhitefg[0m`,
		`[97mwhitebfg[0m`,
		`[47mwhitebg[0m`,
		`[107mwhitebbg[0m`,
		`[39mdefaultfg[0m`,
		`[49mdefaultbg[0m`,
		`[38;5;178mrgb8=fg:178[0m`,
		`[48;5;57mrgb8=bg:57[0m`,
		`[38;2;178;12;240mrgb24=fg:178:12:240[0m`,
		`[48;2;57;124;12mrgb24=bg:57:124:12[0m`,
	)
)

type (
	exitCode = int
	exitFunc func(exitCode)
)

type failer struct {
	w    io.Writer
	fn   exitFunc
	code exitCode
	mu   sync.Locker
}

func newFailer(w io.Writer, fn exitFunc, code exitCode) failer {
	return failer{w, fn, code, &sync.Mutex{}}
}

func (f *failer) fail(e error) (exitFunc, exitCode) {
	f.mu.Lock()
	fmt.Fprintf(f.w, e.Error())
	f.mu.Unlock()
	return f.fn, f.code
}

func args() {
	for _, flagName := range []string{"s", "style"} {
		flag.Func(
			flagName,
			"list of styles and colors to apply to text",
			func(v string) error {
				styles = strings.Fields(v)
				return nil
			},
		)
	}
	flag.Usage = func() { fmt.Print(usage) }
	flag.Parse()
}

func readText(bb *[]byte, rr ...io.Reader) error {
	if len(rr) == 0 {
		return nil
	}
	for _, r := range rr {
		b, err := ioutil.ReadAll(r)
		if err != nil {
			return err
		}
		*bb = append(*bb, b...)
	}
	return nil
}

func argsFiles() ([]io.Reader, func(), error) {
	var files []io.Reader
	closer := func() {
		for _, f := range files {
			switch t := f.(type) {
			case io.Closer:
				t.Close()
			}
		}
	}
	for _, fname := range flag.Args() {
		f, err := os.OpenFile(fname, os.O_RDONLY, 0)
		if err != nil {
			return files, closer, err
		}
		files = append(files, f)
	}
	return files, closer, nil
}

func main() {
	args()
	f := newFailer(os.Stderr, os.Exit, exitFailure)
	text := make([]byte, 0, 64)
	if len(flag.Args()) > 0 {
		files, closer, err := argsFiles()
		defer closer()
		if err != nil {
			exit, code := f.fail(err)
			exit(code)
		}
		err = readText(&text, files...)
		if err != nil {
			exit, code := f.fail(err)
			exit(code)
		}
	} else {
		err := readText(&text, os.Stdin)
		if err != nil {
			exit, code := f.fail(err)
			exit(code)
		}
	}
	out := bufio.NewWriter(os.Stdout)
	colors, err := termcols.MapColors(styles)
	if err != nil {
		exit, code := f.fail(err)
		exit(code)
	}
	output := termcols.Colorize(string(text), colors...)
	_, err = out.WriteString(output)
	if err != nil {
		exit, code := f.fail(err)
		exit(code)
	}
	if err := out.Flush(); err != nil {
		exit, code := f.fail(err)
		exit(code)
	}
	os.Exit(exitSuccess)
}
