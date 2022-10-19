package termcols

import "testing"

func TestMapping(t *testing.T) {
	cases := []struct {
		color  string
		expOut error
	}{
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

		{"purplefg", errMap},
		{"greybg", errMap},
		{"pinkbbg", errMap},
		{"orangebfg", errMap},

		{"rgb8=fg:25", errMap},
		{"rgb24=bg:255:24:123", errMap},
	}
	for _, c := range cases {
		t.Run(c.color, func(t *testing.T) {
			if _, out := mapColor(c.color); out != c.expOut {
				t.Errorf("Have: %t, want %t", out, c.expOut)
			}
		})
	}
}
