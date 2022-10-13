package termcols

import (
	"fmt"
	"strings"
	"unsafe"
)

const (
	// ESC sequence \033 works with Bash, Zsh and Dash (appears not to in Ksh
	// and Csh). There are other escape sequence code representations: \x1b,
	// \0x1b, \u001b and \e are the ones that I came across. All of these
	// appear not to work in Bash, and yet they do work in Zsh.
	ESC = `\033`
	CSI = ESC + `[`
)

// Layer
const (
	FG layer = CSI + `38`
	BG layer = CSI + `48`
)

const RESET sgrAttr = CSI + `0m` // Reset control sequence

// Style
const (
	BOLD          sgrAttr = CSI + `1m`
	FAINT         sgrAttr = CSI + `2m`
	ITALIC        sgrAttr = CSI + `3m`
	UNDEDRLINE    sgrAttr = CSI + `4m`
	BLINK         sgrAttr = CSI + `5m`
	REVERSE       sgrAttr = CSI + `7m`
	HIDE          sgrAttr = CSI + `8m`
	STRIKE        sgrAttr = CSI + `9m`
	DEFAULT_STYLE sgrAttr = CSI + `10m`
)

// Normal foreground
const (
	BLACK_FG   sgrAttr = CSI + `30m`
	RED_FG     sgrAttr = CSI + `31m`
	GREEN_FG   sgrAttr = CSI + `32m`
	YELLOW_FG  sgrAttr = CSI + `33m`
	BLUE_FG    sgrAttr = CSI + `34m`
	MAGENTA_FG sgrAttr = CSI + `35m`
	CYAN_FG    sgrAttr = CSI + `36m`
	WHITE_FG   sgrAttr = CSI + `37m`
	DEFAULT_FG sgrAttr = CSI + `39m`
)

// Bright foreground
const (
	BLACK_BFG   sgrAttr = CSI + `90m`
	RED_BFG     sgrAttr = CSI + `91m`
	GREEN_BFG   sgrAttr = CSI + `92m`
	YELLOW_BFG  sgrAttr = CSI + `93m`
	BLUE_BFG    sgrAttr = CSI + `94m`
	MAGENTA_BFG sgrAttr = CSI + `95m`
	CYAN_BFG    sgrAttr = CSI + `96m`
	WHITE_BFG   sgrAttr = CSI + `97m`
)

// Normal background
const (
	BLACK_BG   sgrAttr = CSI + `40m`
	RED_BG     sgrAttr = CSI + `41m`
	GREEN_BG   sgrAttr = CSI + `42m`
	YELLOW_BG  sgrAttr = CSI + `43m`
	BLUE_BG    sgrAttr = CSI + `44m`
	MAGENTA_BG sgrAttr = CSI + `45m`
	CYAN_BG    sgrAttr = CSI + `46m`
	WHITE_BG   sgrAttr = CSI + `47m`
	DEFAULT_BG sgrAttr = CSI + `49m`
)

// Bright background
const (
	BLACK_BBG   sgrAttr = CSI + `100m`
	RED_BBG     sgrAttr = CSI + `101m`
	GREEN_BBG   sgrAttr = CSI + `102m`
	YELLOW_BBG  sgrAttr = CSI + `103m`
	BLUE_BBG    sgrAttr = CSI + `104m`
	MAGENTA_BBG sgrAttr = CSI + `105m`
	CYAN_BBG    sgrAttr = CSI + `106m`
	WHITE_BBG   sgrAttr = CSI + `107m`
)

// SGR (Select Graphic Rendition) control sequence sets display attributes.
// Each SGR parameter remains active until it is reset with the `CSI 0m`
// [RESET] control sequence.
type sgrAttr string

// Layer is used to specify whether the color should be applied to either
// foreground or background. The default format for the RGB set
// foreground/background color control sequence for 24-bit colors is
// {layer};2;{R};{G};{B}m, and for 8-bit colors this is {layer};5;{Color}m as
// implemented in [Rgb8] and [Rgb24] public functions respectively.
type layer string

// Colorize returns a string literal s with attrs SGR control sequences
// prepended and the reset control sequence appended at the end. The sequence
// of attrs passed to the function call is preserved, so colors and styles can
// (un)intentionally cancel out one another.
func Colorize(s string, attrs ...sgrAttr) string {
	if len(attrs) == 0 {
		return s
	}
	prefix := *(*[]string)(unsafe.Pointer(&attrs))
	return strings.Join(prefix, "") + s + string(RESET)
}

// Rgb8 returns the set foreground/background 8-bit color control sequence. It
// accepts the target layer l parameter that can either be set to foreground or
// background. The c parameter stands for the color. It corresponds to one of
// the colors from a 256-color lookup table, hence the parameter should be in
// the range [0, 255].
func Rgb8(l layer, c uint8) sgrAttr {
	seq := fmt.Sprintf("%s;5;%d", l, c)
	return sgrAttr(seq)
}

// Rgb24 returns the set foreground/background 24-bit color control sequence.
// It accepts the target layer l parameter that can either be set to foreground
// or background. The next three r, g, b parameters correspond to a 24-bit
// color sequence split into three 8-bit sets. RGB parameters should be in the
// range [0, 255].
func Rgb24(l layer, r, g, b uint8) sgrAttr {
	seq := fmt.Sprintf("%s;2;%d;%d;%dm", l, r, g, b)
	return sgrAttr(seq)
}
