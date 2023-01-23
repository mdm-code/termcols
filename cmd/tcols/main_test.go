package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/mdm-code/termcols"
)

type (
	mockWriter struct{ buff []byte }
	failWriter struct{}
	mockReader struct{ text []byte }
	failReader struct{}
)

func (w *mockWriter) Write(p []byte) (int, error) {
	w.buff = append(w.buff, p...)
	return 0, nil
}

func (w *failWriter) Write(p []byte) (int, error) {
	return 0, fmt.Errorf("errored")
}

func (w *mockWriter) String() string {
	return string(w.buff)
}

func (r *mockReader) Read(p []byte) (int, error) {
	p = append(p, []byte("Colorize me!")...)
	return 0, io.EOF
}

func (r *failReader) Read(p []byte) (int, error) {
	return 0, fmt.Errorf("errored")
}

func TestFail(t *testing.T) {
	cases := []struct {
		w    *mockWriter
		fn   exitFunc
		code exitCode
		err  error
	}{
		{&mockWriter{}, func(int) {}, exitFailure, os.ErrNotExist},
	}
	for _, c := range cases {
		f := newFailer(c.w, c.fn, c.code)
		exit, code := f.fail(c.err)
		defer exit(code)

		if c.w.String() != c.err.Error() {
			t.Errorf("Have %s; want %s", c.w.String(), c.err.Error())
		}
	}
}

func TestCliParse(t *testing.T) {
	cases := []struct {
		name string
		args []string
		err  error
	}{
		{"pass-01", []string{"-styles", "redfg bluefg"}, nil},
		{"pass-02", []string{"-s", "strike rgb24=fg:242:121:64"}, nil},
		{"pass-03", []string{"-s", "yellowbg", "--styles", "bluefg"}, nil},
		{"pass-04", []string{}, nil},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, err := parse(c.args, func([]string, func(string) (*os.File, error)) ([]io.Reader, func(), error) {
				return []io.Reader{}, func() {}, nil
			})
			if !errors.Is(err, c.err) {
				t.Errorf("Have %T; want %T", err, c.err)
			}
		})
	}
}

// TestPipeText tests a single, single-threaded pass of text data.
func TestPipeText(t *testing.T) {
	cases := []struct {
		r   io.Reader
		w   io.Writer
		s   []string
		err error
	}{
		{&mockReader{}, &mockWriter{}, []string{}, nil},
		{nil, nil, []string{}, errPiping},
		{&mockReader{}, &mockWriter{}, []string{"blue"}, termcols.ErrMap},
		{&failReader{}, &mockWriter{}, []string{}, errPiping},
		{&mockReader{}, &failWriter{}, []string{}, errPiping},
	}
	for _, c := range cases {
		err := pipe(c.r, c.w, c.s)
		if !errors.Is(err, c.err) {
			t.Errorf("Have %T; want %T", err, c.err)
		}
	}
}

func TestOpen(t *testing.T) {
	errOpen := errors.New("open error")
	cases := []struct {
		name   string
		fnames []string
		f      func(string) (*os.File, error)
		err    error
	}{
		{
			"pass",
			[]string{"one.txt", "two.md", "three.html"},
			func(fname string) (*os.File, error) {
				return &os.File{}, nil
			},
			nil,
		},
		{
			"fail-open",
			[]string{"one.txt", "two.md", "three.html"},
			func(fname string) (*os.File, error) {
				return nil, errOpen
			},
			errOpen,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, closer, err := open(c.fnames, c.f)
			defer closer()
			if err != c.err {
				t.Errorf("Have %T; want %T", err, c.err)
			}
		})
	}
}

func TestRun(t *testing.T) {
	f := func(fname []string, f func(string) (*os.File, error)) ([]io.Reader, func(), error) {
		return []io.Reader{strings.NewReader("hello")}, func() {}, nil
	}
	cases := []struct {
		name string
		args []string
		fn   openFn
		err  error
	}{
		{"pass-01", []string{"-s", "greenbg yellowfg bold", "1.pyc", "2.c"}, f, nil},
		{"pass-02", []string{}, f, nil},
		{"fail-01", []string{"--styles", "wacky", "hello.py"}, f, termcols.ErrMap},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := run(c.args, c.fn)
			if !errors.Is(err, c.err) {
				t.Errorf("Have %v; want %v", err, c.err)
			}
		})
	}
}
