package termcols

import (
	"fmt"
	"strings"
	"unsafe"
)

const (
	// Esc octal sequence \033 works with Bash, Zsh and Dash (appears not to in
	// Ksh and Csh). There are other escape sequence code representations: the
	// ctrl-key ^[ sequence, unicode \u001b, hexadecimal \x1b, \0x1B, decimal
	// 27 and \e are the ones that I came across.
	Esc = "\033"

	// Csi stands for Control Sequence Introducer
	Csi = Esc + "["
)

// Layer
const (
	FG Layer = Csi + "38"
	BG Layer = Csi + "48"
)

// Reset control sequence
const Reset SgrAttr = Csi + "0m"

// Style
const (
	Bold         SgrAttr = Csi + "1m"
	Faint        SgrAttr = Csi + "2m"
	Italic       SgrAttr = Csi + "3m"
	Underline    SgrAttr = Csi + "4m"
	Blink        SgrAttr = Csi + "5m"
	Reverse      SgrAttr = Csi + "7m"
	Hide         SgrAttr = Csi + "8m"
	Strike       SgrAttr = Csi + "9m"
	DefaultStyle SgrAttr = Csi + "10m"
)

// Normal foreground
const (
	BlackFg   SgrAttr = Csi + "30m"
	RedFg     SgrAttr = Csi + "31m"
	GreenFg   SgrAttr = Csi + "32m"
	YellowFg  SgrAttr = Csi + "33m"
	BlueFg    SgrAttr = Csi + "34m"
	MagentaFg SgrAttr = Csi + "35m"
	CyanFg    SgrAttr = Csi + "36m"
	WhiteFg   SgrAttr = Csi + "37m"
	DefaultFg SgrAttr = Csi + "39m"
)

// Bright foreground
const (
	BlackBfg   SgrAttr = Csi + "90m"
	RedBfg     SgrAttr = Csi + "91m"
	GreenBfg   SgrAttr = Csi + "92m"
	YellowBfg  SgrAttr = Csi + "93m"
	BlueBfg    SgrAttr = Csi + "94m"
	MagentaBfg SgrAttr = Csi + "95m"
	CyanBfg    SgrAttr = Csi + "96m"
	WhiteBfg   SgrAttr = Csi + "97m"
)

// Normal background
const (
	BlackBg   SgrAttr = Csi + "40m"
	RedBg     SgrAttr = Csi + "41m"
	GreenBg   SgrAttr = Csi + "42m"
	YellowBg  SgrAttr = Csi + "43m"
	BlueBg    SgrAttr = Csi + "44m"
	MagentaBg SgrAttr = Csi + "45m"
	CyanBg    SgrAttr = Csi + "46m"
	WhiteBg   SgrAttr = Csi + "47m"
	DefaultBg SgrAttr = Csi + "49m"
)

// Bright background
const (
	BlackBbg   SgrAttr = Csi + "100m"
	RedBbg     SgrAttr = Csi + "101m"
	GreenBbg   SgrAttr = Csi + "102m"
	YellowBbg  SgrAttr = Csi + "103m"
	BlueBbg    SgrAttr = Csi + "104m"
	MagentaBbg SgrAttr = Csi + "105m"
	CyanBbg    SgrAttr = Csi + "106m"
	WhiteBbg   SgrAttr = Csi + "107m"
)

// SgrAttr corresponds to a SGR (Select Graphic Rendition) control sequence
// sets display attributes. Each SGR parameter remains active until it is reset
// with the `CSI 0m` [RESET] control sequence.
type SgrAttr string

// Layer is used to specify whether the color should be applied to either
// foreground or background. The default format for the RGB set
// foreground/background color control sequence for 24-bit colors is
// {Layer};2;{R};{G};{B}m, and for 8-bit colors this is {Layer};5;{Color}m as
// implemented in [Rgb8] and [Rgb24] public functions respectively.
type Layer string

// Colorize returns a string literal s with attrs SGR control sequences
// prepended and the reset control sequence appended at the end. The sequence
// of attrs passed to the function call is preserved, so colors and styles can
// (un)intentionally cancel out one another.
func Colorize(s string, attrs ...SgrAttr) string {
	if len(attrs) == 0 {
		return s
	}
	prefix := *(*[]string)(unsafe.Pointer(&attrs))
	return strings.Join(prefix, "") + s + string(Reset)
}

// Rgb8 returns the set foreground/background 8-bit color control sequence. It
// accepts the target layer l parameter that can either be set to foreground or
// background. The c parameter stands for the color. It corresponds to one of
// the colors from a 256-color lookup table, hence the parameter should be in
// the range [0, 255].
func Rgb8(l Layer, c uint8) SgrAttr {
	seq := fmt.Sprintf("%s;5;%dm", l, c)
	return SgrAttr(seq)
}

// Rgb24 returns the set foreground/background 24-bit color control sequence.
// It accepts the target layer l parameter that can either be set to foreground
// or background. The next three r, g, b parameters correspond to a 24-bit
// color sequence split into three 8-bit sets. RGB parameters should be in the
// range [0, 255].
func Rgb24(l Layer, r, g, b uint8) SgrAttr {
	seq := fmt.Sprintf("%s;2;%d;%d;%dm", l, r, g, b)
	return SgrAttr(seq)
}
