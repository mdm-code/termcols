package termcols

import (
	"fmt"
	"regexp"
	"testing"
)

func TestMapColor(t *testing.T) {
	cases := []struct {
		color  string
		expOut error
	}{
		// Existing colors and styles
		{"bold", nil},
		{"faint", nil},
		{"italic", nil},
		{"blink", nil},
		{"reverse", nil},
		{"strike", nil},
		{"defaultstyle", nil},
		{"defaultbg", nil},
		{"blackfg", nil},
		{"redbfg", nil},
		{"greenfg", nil},
		{"yellowbg", nil},
		{"bluebbg", nil},
		{"magentafg", nil},
		{"cyanbg", nil},
		{"whitebfg", nil},

		// Not-implemented colors and styles
		{"purplefg", errMap},
		{"greybg", errMap},
		{"pinkbbg", errMap},
		{"orangebfg", errMap},

		// Passing RGB patterns
		{"RGB8=fg:25", nil},
		{"rgb8=bg:240", nil},
		{"rgb8=fg:123", nil},
		{"rgb24=BG:8:246:22", nil},
		{"RGB24=bg:123:22:40", nil},
		{"rgb24=fg:0:12:255", nil},

		// Failing RGB patterns
		{"", errMap},                         // empty string
		{"rgb24", errMap},                    // missing parameters
		{"rgb8=gb:227", errMap},              // unknown layer (8)
		{"rgb24=gf:227:12:142", errMap},      // unknown layer (24)
		{"rgb9=bg:227", errMap},              // unknown bit size
		{"rgb8=bg:255:255", errMap},          // too many color values (8)
		{"rgb24=fg:255:255:255:255", errMap}, // too many color values (24)
		{"rgb24=gf:12:245:0", errMap},        // unknown layer
		{"rgb8=bg:256", errMap},              // 256 > uint8 255 cap (8)
		{"rgb24=bg:255:256:123", errMap},     // 256 > uint8 255 cap (24)
	}
	for _, c := range cases {
		t.Run(c.color, func(t *testing.T) {
			if _, out := MapColor(c.color); out != c.expOut {
				t.Errorf("Have: %t, want %t", out, c.expOut)
			}
		})
	}
}

func TestMatchRegexp(t *testing.T) {
	cases := []struct {
		re    *regexp.Regexp
		s     interface{}
		okExp bool
	}{
		{
			regexp.MustCompile(
				`(?mi)^rgb8=(?P<layer>fg|bg):(?P<color>\d{1,3})$`,
			),
			"rgb8=fg:255",
			true,
		},
		{
			regexp.MustCompile(
				`(?mi)^rgb24=(?P<layer>fg|bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>\d{1,3})$`,
			),
			"rgb24=fg:255:255:255",
			true,
		},
		{regexp.MustCompile(`rgb:.*`), "RGB24", false},
		{regexp.MustCompile(`.*`), struct{}{}, false},
		{regexp.MustCompile(`.*`), 10_000, false},
		{regexp.MustCompile(`.*`), 3.14, false},
	}
	for _, c := range cases {
		t.Run(fmt.Sprintf("%v", c.s), func(t *testing.T) {
			if ok := matchRegexp(c.re, c.s); ok != c.okExp {
				t.Errorf("Have: %t; want %t", ok, c.okExp)
			}
		})
	}
}

func TestCollateRgb8(t *testing.T) {}

func TestCollateRgb24(t *testing.T) {}

func TestGetParams(t *testing.T) {}

func TestValidUint(t *testing.T) {}
