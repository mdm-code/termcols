package main

import (
	"errors"
	"io"
	"os"
	"testing"
)

type mockWriter struct {
	buff []byte
}

type mockReader struct {
	text []byte
}

func (w *mockWriter) Write(p []byte) (int, error) {
	w.buff = append(w.buff, p...)
	return 0, nil
}

func (w *mockWriter) String() string {
	return string(w.buff)
}

func (r *mockReader) Read(p []byte) (int, error) {
	p = append(p, []byte("Colorize me!")...)
	return 0, io.EOF
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

		{"fail-01", []string{}, errParsing},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			_, _, err := parse(c.args, func([]string) ([]io.Reader, func(), error) {
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
	}
	for _, c := range cases {
		err := pipe(c.r, c.w, c.s)
		if !errors.Is(err, c.err) {
			t.Errorf("Have %T; want %T", err, c.err)
		}
	}
}
