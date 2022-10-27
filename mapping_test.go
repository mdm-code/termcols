package termcols

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"testing"
)

func TestMapColors(t *testing.T) {
	cases := []struct {
		colors []string
		expOut error
	}{
		{[]string{}, nil},
		{[]string{"bluefg", "bold"}, nil},
		{[]string{"blink", "redbg"}, nil},
		{[]string{"italic", "strike", "rgb8=fg:255"}, nil},
		{[]string{"rgb24=bg:255:255:255", "blink", "magentafg"}, nil},

		// Failing colors
		{[]string{"italics", "strike"}, errMap},
		{[]string{"italic", "redgb"}, errMap},
		{[]string{"italic", "rgb8=gf:240"}, errMap},
	}
	for _, c := range cases {
		t.Run(strings.Join(c.colors, "-"), func(t *testing.T) {
			if _, out := MapColors(c.colors); out != c.expOut {
				t.Errorf("Have: %t, want %t", out, c.expOut)
			}
		})
	}
}

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

func TestCollateRgb8(t *testing.T) {
	re := regexp.MustCompile(
		`(?mi)^rgb8=(?P<layer>fg|bg):(?P<color>\d{1,3})$`,
	)
	cases := []struct {
		re    *regexp.Regexp
		str   string
		okExp bool
	}{
		{re, "rgb8=fg:126", true},
		{re, "", false},
		{re, "rgb8=222", false},
		{re, "rgb8=gb:255", false},
		{re, "rgb8=fg:", false},
		{re, "rgb8=bg:black", false},
		{re, "rgb8=fg:256", false},

		// NOTE: Made-up regular expressions to trigger collation errors
		{regexp.MustCompile(`.*`), "", false}, // missing layer field
		{
			regexp.MustCompile(`rgb8=(?P<layer>gb):(?P<color>red)`),
			"rgb8=gb:red",
			false,
		}, // `gb` is not a valid layer on the layerMap
		{
			regexp.MustCompile(`rgb8=(?P<layer>fg|bg):(?P<color>red)`),
			"rgb8=fg:red",
			false,
		}, // red cannot be converted to an int
		{
			regexp.MustCompile(`rgb8=(?P<layer>bg):.*`),
			"rgb8=bg:missing",
			false,
		}, // missing color field
	}
	for _, c := range cases {
		t.Run(c.str, func(t *testing.T) {
			if _, ok := collateRgb8(c.re, c.str); ok != c.okExp {
				t.Errorf("Have: %t; want %t", ok, c.okExp)
			}
		})
	}
}

func TestCollateRgb24(t *testing.T) {
	re := regexp.MustCompile(
		`(?mi)^rgb24=(?P<layer>fg|bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>\d{1,3})$`,
	)
	cases := []struct {
		re    *regexp.Regexp
		str   string
		okExp bool
	}{
		{re, "rgb24=fg:126:12:56", true},
		{re, "", false},
		{re, "rgb24=222:232:101", false},
		{re, "rgb24=gb:255:255:255", false},
		{re, "rgb24=fg:", false},
		{re, "rgb24=fg:312:120:2", false},
		{re, "rgb24=fg:120:120:257", false},

		// NOTE: Made-up regular expressions to trigger collation errors
		{regexp.MustCompile(`.*`), "", false}, // missing layer field
		{
			regexp.MustCompile(`rgb24=(?P<layer>gb):.*`),
			"rgb24=gb:red",
			false,
		}, // `gb` is not a valid layer on the layerMap
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>fg|bg):(?P<r>red)`,
			),
			"rgb24=fg:red",
			false,
		}, //  red cannot be converted to an int
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>fg|bg):(?P<g>green)`,
			),
			"rgb24=fg:green",
			false,
		}, // green cannot be converted to an int
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>fg|bg):(?P<b>blue)`,
			),
			"rgb24=fg:blue",
			false,
		}, // blue cannot be converted to an int
		{
			regexp.MustCompile(`rgb24=(?P<layer>bg):.*`),
			"rgb24=bg:255:255:255",
			false,
		}, // missing r field
		{
			regexp.MustCompile(`rgb24=(?P<layer>bg):(?P<r>\d{1,3}):.*`),
			"rgb24=bg:255:255:255",
			false,
		}, // missing g field
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):.*`,
			),
			"rgb24=bg:255:255:255",
			false,
		}, // missing b field
		{
			regexp.MustCompile(`rgb24=(?P<layer>bg):(?P<r>red):.*`),
			"rgb24=bg:red:255:255",
			false,
		}, // red cannot be converted to an int
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>bg):(?P<r>\d{1,3}):(?P<g>green):.*`,
			),
			"rgb24=bg:255:green:255",
			false,
		}, // green cannot be converted to an int
		{
			regexp.MustCompile(
				`rgb24=(?P<layer>bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>blue)`,
			),
			"rgb24=bg:255:255:blue",
			false,
		}, // blue cannot be converted to an int
	}
	for _, c := range cases {
		t.Run(c.str, func(t *testing.T) {
			if _, ok := collateRgb24(c.re, c.str); ok != c.okExp {
				t.Errorf("Have: %t; want %t", ok, c.okExp)
			}
		})
	}
}

func TestGetParams(t *testing.T) {
	re := regexp.MustCompile(
		`(?mi)^rgb8=(?P<layer>fg|bg):(?P<color>\d{1,3})$`,
	)
	cases := []struct {
		re     *regexp.Regexp
		str    string
		outExp map[string]string
	}{
		{re, "rgb8=fg:254", map[string]string{"layer": "fg", "color": "254"}},
		{re, "rgb8=bg:121", map[string]string{"layer": "bg", "color": "121"}},
		{re, "RGB8=FG:101", map[string]string{"layer": "FG", "color": "101"}},
		{re, "RGB8=fg:142", map[string]string{"layer": "fg", "color": "142"}},
		{re, "rgb8=BG:236", map[string]string{"layer": "BG", "color": "236"}},
	}
	for _, c := range cases {
		t.Run(c.str, func(t *testing.T) {
			out := getParams(c.re, c.str)
			if !reflect.DeepEqual(c.outExp, out) {
				t.Errorf("Have %v; want %v", out, c.outExp)
			}
		})
	}

}

func TestValidUint(t *testing.T) {
	cases := []struct {
		i     int
		str   string
		okExp bool
	}{
		{0, "0", true},
		{255, "255", true},
		{56, "56", true},
		{212, "212", true},
		{78, "78", true},
		{193, "193", true},
		{22, "22", true},

		{-1, "-1", false},
		{-3942, "-3942", false},
		{256, "256", false},
		{905, "905", false},
		{1293, "1293", false},
	}
	for _, c := range cases {
		t.Run(c.str, func(t *testing.T) {
			if ok := validUint8(c.i); ok != c.okExp {
				t.Errorf("Have %t; want: %t", ok, c.okExp)
			}
		})
	}
}
