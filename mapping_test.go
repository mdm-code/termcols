package termcols

import "testing"

func TestMapping(t *testing.T) {
	cases := []struct {
		color  string
		expOut bool
	}{
		{"bold", true},
		{"faint", true},
		{"italic", true},
		{"blink", true},
		{"reverse", true},
		{"strike", true},
		{"defaultstyle", true},
		{"defaultbg", true},
		{"blackfg", true},
		{"redbfg", true},
		{"greenfg", true},
		{"yellowbg", true},
		{"bluebbg", true},
		{"magentafg", true},
		{"cyanbg", true},
		{"whitebfg", true},

		{"purplefg", false},
		{"greybg", false},
		{"pinkbbg", false},
		{"orangebfg", false},
	}
	for _, c := range cases {
		t.Run(c.color, func(t *testing.T) {
			if _, out := mapColor(c.color); out != c.expOut {
				t.Errorf("Have: %t, want %t", out, c.expOut)
			}
		})
	}
}
