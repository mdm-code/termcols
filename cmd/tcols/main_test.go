package main

import (
	"os"
	"testing"
)

// Unit tests
// ==========
// 1. Add tests for readText with reflections on the buffer pointed here
// 2. Check how to mock OpenFile or Open in Go
// 3. Test the fail function with ExampleTests
// 4. Create Run command -- call it for each file in argsFiles or for os.Stdin
// 5. This might help me run colorize with goroutines

type mockWriter struct {
	buff []byte
}

func (m *mockWriter) Write(p []byte) (int, error) {
	m.buff = append(m.buff, p...)
	return 0, nil
}

func (m *mockWriter) String() string {
	return string(m.buff)
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
			err := parse(c.args)
			if !errors.Is(err, c.err) {
				t.Errorf("Have %T; want %T", err, c.err)
			}
		})
	}
}

}
