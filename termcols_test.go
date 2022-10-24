package termcols

import (
	"testing"
)

// BenchmarkColorize was arranged to test the performance of the slice type
// conversion using unsafe.Pointer as opposed to looping over the source
// slice and appending elements one after another to the target slice.
//
// In my testing slice type conversion with unsafe.Pointer is almost four
// times faster than the memory-safe alternative.
func BenchmarkColorize(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Colorize("Colorize me!", BlackFg)
	}
}

func TestColorize(t *testing.T) {
	cases := []struct {
		attrs  []SgrAttr
		expOut string
	}{
		{
			[]SgrAttr{},
			" Colorize me! ",
		},
		{
			[]SgrAttr{Bold, BlackFg, WhiteBbg},
			"\033[1m\033[30m\033[107m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Faint, RedFg, CyanBbg},
			"\033[2m\033[31m\033[106m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Italic, GreenFg, MagentaBbg},
			"\033[3m\033[32m\033[105m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Underline, YellowFg, BlueBbg},
			"\033[4m\033[33m\033[104m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Blink, BlueFg, YellowBbg},
			"\033[5m\033[34m\033[103m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Reverse, MagentaFg, GreenBbg},
			"\033[7m\033[35m\033[102m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Hide, CyanFg, RedBbg},
			"\033[8m\033[36m\033[101m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Strike, WhiteFg, BlackBbg},
			"\033[9m\033[37m\033[100m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{DefaultStyle, DefaultFg, BlackBbg},
			"\033[10m\033[39m\033[100m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Bold, BlackBfg, DefaultBg},
			"\033[1m\033[90m\033[49m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Faint, RedBfg, WhiteBg},
			"\033[2m\033[91m\033[47m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Italic, GreenBfg, CyanBg},
			"\033[3m\033[92m\033[46m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Underline, YellowBfg, MagentaBg},
			"\033[4m\033[93m\033[45m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Blink, BlueBfg, YellowBg},
			"\033[5m\033[94m\033[43m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Reverse, MagentaBfg, BlueBg},
			"\033[7m\033[95m\033[44m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Hide, CyanBfg, GreenBg},
			"\033[8m\033[96m\033[42m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Strike, WhiteBfg, RedBg},
			"\033[9m\033[97m\033[41m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{DefaultStyle, WhiteBfg, BlackBg},
			"\033[10m\033[97m\033[40m Colorize me! \033[0m",
		},
	}
	for _, c := range cases {
		t.Run(c.expOut, func(t *testing.T) {
			if out := Colorize(" Colorize me! ", c.attrs...); out != c.expOut {
				t.Errorf("Have: %s, want: %s", out, c.expOut)
			}
		})
	}
}

func TestRgb8(t *testing.T) {
	cases := []struct {
		l      Layer
		col    uint8
		expOut SgrAttr
	}{
		{FG, 12, SgrAttr("\033[38;5;12m")},
		{BG, 82, SgrAttr("\033[48;5;82m")},
		{FG, 32, SgrAttr("\033[38;5;32m")},
		{BG, 180, SgrAttr("\033[48;5;180m")},
		{FG, 255, SgrAttr("\033[38;5;255m")},
		{BG, 234, SgrAttr("\033[48;5;234m")},
	}
	for _, c := range cases {
		t.Run(string(c.expOut), func(t *testing.T) {
			if out := Rgb8(c.l, c.col); out != c.expOut {
				t.Errorf("Have: %s, want: %s", out, c.expOut)
			}
		})
	}
}

func TestRgb24(t *testing.T) {
	cases := []struct {
		l       Layer
		r, g, b uint8
		expOut  SgrAttr
	}{
		{FG, 12, 145, 67, SgrAttr("\033[38;2;12;145;67m")},
		{BG, 200, 247, 23, SgrAttr("\033[48;2;200;247;23m")},
		{FG, 89, 12, 2, SgrAttr("\033[38;2;89;12;2m")},
		{BG, 150, 73, 0, SgrAttr("\033[48;2;150;73;0m")},
		{FG, 0, 255, 255, SgrAttr("\033[38;2;0;255;255m")},
		{BG, 12, 59, 90, SgrAttr("\033[48;2;12;59;90m")},
	}
	for _, c := range cases {
		t.Run(string(c.expOut), func(t *testing.T) {
			if out := Rgb24(c.l, c.r, c.g, c.b); out != c.expOut {
				t.Errorf("Have: %s, want: %s", out, c.expOut)
			}
		})
	}
}

func TestCombination(t *testing.T) {
	cases := []struct {
		attrs  []SgrAttr
		expOut string
	}{
		{
			[]SgrAttr{Italic, GreenFg, Rgb8(BG, 176)},
			"\033[3m\033[32m\033[48;5;176m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Underline, Rgb8(FG, 44), MagentaBg},
			"\033[4m\033[38;5;44m\033[45m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Reverse, MagentaBfg, Rgb8(BG, 255)},
			"\033[7m\033[95m\033[48;5;255m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Bold, Rgb24(FG, 12, 45, 62), Rgb8(BG, 89)},
			"\033[1m\033[38;2;12;45;62m\033[48;5;89m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Strike, Rgb24(FG, 78, 22, 0), BlackBbg},
			"\033[9m\033[38;2;78;22;0m\033[100m Colorize me! \033[0m",
		},
		{
			[]SgrAttr{Hide, CyanBfg, Rgb24(BG, 255, 255, 255)},
			"\033[8m\033[96m\033[48;2;255;255;255m Colorize me! \033[0m",
		},
	}
	for _, c := range cases {
		t.Run(c.expOut, func(t *testing.T) {
			if out := Colorize(" Colorize me! ", c.attrs...); out != c.expOut {
				t.Errorf("Have: %s, want: %s", out, c.expOut)
			}
		})
	}
}
