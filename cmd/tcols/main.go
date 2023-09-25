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

	tcols --style 'bold bluefg' < <(echo -n 'Hello, world!')
	echo -n Hello, world\! | tcols -s bold -s bluefg

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
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/mdm-code/termcols"
	"golang.org/x/term"
)

const (
	exitSuccess exitCode = iota
	exitFailure
)

var (
	styles     []string
	errPiping  error = errors.New("cannot read/write on nil interfaces")
	usageAttrs       = [...][2]string{
		{"Hello, world!", string(termcols.Bold) + string(termcols.BlueFg) + "%s" + string(termcols.Reset)},
		{"bold", string(termcols.Bold) + "%s" + string(termcols.Reset)},
		{"faint", string(termcols.Faint) + "%s" + string(termcols.Reset)},
		{"italic", string(termcols.Italic) + "%s" + string(termcols.Reset)},
		{"underline", string(termcols.Underline) + "%s" + string(termcols.Reset)},
		{"blink", string(termcols.Blink) + "%s" + string(termcols.Reset)},
		{"reverse", string(termcols.Reverse) + "%s" + string(termcols.Reset)},
		{"hide", string(termcols.Hide) + "%s" + string(termcols.Reset)},
		{"strike", string(termcols.Strike) + "%s" + string(termcols.Reset)},
		{"defaultstyle", string(termcols.DefaultStyle) + "%s" + string(termcols.Reset)},
		{"blackfg", string(termcols.BlackFg) + "%s" + string(termcols.Reset)},
		{"blackbfg", string(termcols.BlackBfg) + "%s" + string(termcols.Reset)},
		{"blackbg", string(termcols.BlackBg) + "%s" + string(termcols.Reset)},
		{"blackbbg", string(termcols.BlackBbg) + "%s" + string(termcols.Reset)},
		{"redfg", string(termcols.RedFg) + "%s" + string(termcols.Reset)},
		{"redbfg", string(termcols.RedBfg) + "%s" + string(termcols.Reset)},
		{"redbg", string(termcols.RedBg) + "%s" + string(termcols.Reset)},
		{"redbbg", string(termcols.RedBbg) + "%s" + string(termcols.Reset)},
		{"greenfg", string(termcols.GreenFg) + "%s" + string(termcols.Reset)},
		{"greenbfg", string(termcols.GreenBfg) + "%s" + string(termcols.Reset)},
		{"greenbg", string(termcols.GreenBg) + "%s" + string(termcols.Reset)},
		{"greenbbg", string(termcols.GreenBbg) + "%s" + string(termcols.Reset)},
		{"yellowfg", string(termcols.YellowFg) + "%s" + string(termcols.Reset)},
		{"yellowbfg", string(termcols.YellowBfg) + "%s" + string(termcols.Reset)},
		{"yellowbg", string(termcols.YellowBg) + "%s" + string(termcols.Reset)},
		{"yellowbbg", string(termcols.YellowBbg) + "%s" + string(termcols.Reset)},
		{"bluefg", string(termcols.BlueFg) + "%s" + string(termcols.Reset)},
		{"bluebfg", string(termcols.BlueBfg) + "%s" + string(termcols.Reset)},
		{"bluebg", string(termcols.BlueBg) + "%s" + string(termcols.Reset)},
		{"bluebbg", string(termcols.BlueBbg) + "%s" + string(termcols.Reset)},
		{"magentafg", string(termcols.MagentaFg) + "%s" + string(termcols.Reset)},
		{"magentabfg", string(termcols.MagentaBfg) + "%s" + string(termcols.Reset)},
		{"magentabg", string(termcols.MagentaBg) + "%s" + string(termcols.Reset)},
		{"magentabbg", string(termcols.MagentaBbg) + "%s" + string(termcols.Reset)},
		{"cyanfg", string(termcols.CyanFg) + "%s" + string(termcols.Reset)},
		{"cyanbfg", string(termcols.CyanBfg) + "%s" + string(termcols.Reset)},
		{"cyanbg", string(termcols.CyanBg) + "%s" + string(termcols.Reset)},
		{"cyanbbg", string(termcols.CyanBbg) + "%s" + string(termcols.Reset)},
		{"whitefg", string(termcols.WhiteFg) + "%s" + string(termcols.Reset)},
		{"whitebfg", string(termcols.WhiteBfg) + "%s" + string(termcols.Reset)},
		{"whitebg", string(termcols.WhiteBg) + "%s" + string(termcols.Reset)},
		{"whitebbg", string(termcols.WhiteBbg) + "%s" + string(termcols.Reset)},
		{"defaultfg", string(termcols.DefaultFg) + "%s" + string(termcols.Reset)},
		{"defaultbg", string(termcols.DefaultBg) + "%s" + string(termcols.Reset)},
		{"rgb8=fg:178", `[38;5;178m%s[0m`},
		{"rgb8=bg:57", `[48;5;57m%s[0m`},
		{"rgb24=fg:178:12:240", `[38;2;178;12;240m%s[0m`},
		{"rgb24=fg:57:124:12", `[48;2;57;124;12m%s[0m`},
	}
	usage = `tcols - add color to text on the terminal

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
	%s

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
`
)

type (
	exitCode = int
	openFn   = func([]string, func(string) (*os.File, error)) ([]io.Reader, func(), error)
	exitFunc func(exitCode)

	failer struct {
		w    io.Writer
		fn   exitFunc
		code exitCode
		mu   sync.Locker
	}

	concurrentWriter struct {
		w *bufio.Writer
		sync.Mutex
	}
)

func (f *failer) fail(e error) (exitFunc, exitCode) {
	f.mu.Lock()
	fmt.Fprintf(f.w, e.Error())
	f.mu.Unlock()
	return f.fn, f.code
}

func (cw *concurrentWriter) Write(p []byte) (n int, err error) {
	cw.Lock()
	n, err = cw.w.Write(p)
	cw.Unlock()
	return
}

func (cw *concurrentWriter) Flush() error {
	return cw.w.Flush()
}

func newFailer(w io.Writer, fn exitFunc, code exitCode) failer {
	return failer{w, fn, code, &sync.Mutex{}}
}

func newConcurrentWriter(w io.Writer) *concurrentWriter {
	return &concurrentWriter{w: bufio.NewWriter(w)}
}

// PrepUsageAttrs collates colored or non-colored usage attributes.
func prepUsageAttrs(colored bool) []any {
	result := make([]any, 0, len(usageAttrs))
	for _, v := range usageAttrs {
		if colored {
			parsed := fmt.Sprintf(v[1], v[0])
			result = append(result, parsed)
			continue
		}
		result = append(result, v[0])
	}
	return result
}

func parse(args []string, open openFn) ([]io.Reader, func(), error) {
	fs := flag.NewFlagSet("tcols", flag.ExitOnError)
	for _, fName := range []string{"s", "styles"} {
		fs.Func(
			fName,
			"list of styles and colors to apply to text",
			func(v string) error {
				styles = append(styles, strings.Fields(v)...)
				return nil
			},
		)
	}
	fs.Usage = func() {
		usageOut := os.Stdout
		if term.IsTerminal(int(usageOut.Fd())) {
			colored := prepUsageAttrs(true)
			fmt.Fprintf(usageOut, fmt.Sprintf(usage, colored...))
			return
		}
		bw := prepUsageAttrs(false)
		fmt.Fprintf(usageOut, fmt.Sprintf(usage, bw...))
	}
	err := fs.Parse(args)
	if err != nil {
		return []io.Reader{}, func() {}, err
	}
	if len(fs.Args()) > 0 {
		return open(fs.Args(), os.Open)
	}
	return []io.Reader{os.Stdin}, func() {}, nil
}

// Open opens files to have their contents read. The function f serves as the
// main callable responsible for opening files.
func open(fnames []string, f func(string) (*os.File, error)) ([]io.Reader, func(), error) {
	var files []io.Reader
	closer := func() {
		for _, f := range files {
			switch t := f.(type) {
			case io.Closer:
				t.Close()
			}
		}
	}
	for _, fname := range fnames {
		f, err := f(fname)
		if err != nil {
			return files, closer, err
		}
		files = append(files, f)
	}
	return files, closer, nil
}

// Pipe transfers the input text colorized according to the provided styles
// from the r reader to the w writer. The colorize parameter controls if the
// text should be colorized.
func pipe(r io.Reader, w io.Writer, styles []string, colorize bool) error {
	if r == nil && w == nil {
		return errPiping
	}
	text, err := ioutil.ReadAll(r)
	if err != nil {
		return errPiping
	}
	colors, err := termcols.MapColors(styles)
	if err != nil {
		return err
	}
	if colorize {
		colored := termcols.Colorize(string(text), colors...)
		_, err = io.WriteString(w, colored)
		if err != nil {
			return errPiping
		}
		return nil
	}
	_, err = w.Write(text)
	if err != nil {
		return errPiping
	}
	return nil
}

func run(args []string, fn openFn) error {
	files, closer, err := parse(args, fn)
	defer closer()
	if err != nil {
		return err
	}

	out := newConcurrentWriter(os.Stdout)

	var colorize bool
	if term.IsTerminal(int(os.Stdout.Fd())) {
		colorize = true
	}

	var wg sync.WaitGroup
	wg.Add(len(files))

	done := make(chan struct{})
	fail := make(chan error)

	for _, f := range files {
		go func(r io.Reader) {
			defer wg.Done()
			err := pipe(r, out, styles, colorize)
			if err != nil {
				fail <- err
			}
		}(f)
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		break
	case err := <-fail:
		return err
	}

	if err := out.Flush(); err != nil {
		return err
	}
	return nil
}

func main() {
	f := newFailer(os.Stderr, os.Exit, exitFailure)
	err := run(os.Args[1:], open)
	if err != nil {
		exit, code := f.fail(err)
		exit(code)
	}
	os.Exit(exitSuccess)
}
