package termcols

import (
	"errors"
	"regexp"
)

var layerMap map[string]layer = map[string]layer{"fg": FG, "bg": BG}

var colorMap map[string]SgrAttr = map[string]SgrAttr{
	"bold":         Bold,
	"faint":        Faint,
	"italic":       Italic,
	"underline":    Underline,
	"blink":        Blink,
	"reverse":      Reverse,
	"hide":         Hide,
	"strike":       Strike,
	"defaultstyle": DefaultStyle,

	"defaultfg": DefaultFg,
	"defaultbg": DefaultBg,

	"blackfg":  BlackFg,
	"blackbfg": BlackBfg,
	"blackbg":  BlackBg,
	"blackbbg": BlackBbg,

	"redfg":  RedFg,
	"redbfg": RedBfg,
	"redbg":  RedBg,
	"redbbg": RedBbg,

	"greenfg":  GreenFg,
	"greenbfg": GreenBfg,
	"greenbg":  GreenBg,
	"greenbbg": GreenBbg,

	"yellowfg":  YellowFg,
	"yellobfg":  YellowBfg,
	"yellowbg":  YellowBg,
	"yellowbbg": YellowBbg,

	"bluefg":  BlueFg,
	"bluebfg": BlueBfg,
	"bluebg":  BlueBg,
	"bluebbg": BlueBbg,

	"magentafg":  MagentaFg,
	"magentabfg": MagentaBfg,
	"magentabg":  MagentaBg,
	"magentabbg": MagentaBbg,

	"cyanfg":  CyanFg,
	"cyanbfg": CyanBfg,
	"cyanbg":  CyanBg,
	"cyanbbg": CyanBbg,

	"whitefg":  WhiteFg,
	"whitebfg": WhiteBfg,
	"whitebg":  WhiteBg,
	"whitebbg": WhiteBbg,
}

/*
PLAN:
=====

After the color map (which appears to be the more obvious choice for users)
returns false, I am going to use regex to do
	(1) do two regexpr matches for (a) rgb8 and (b) rgb24
	(2) in case one or the other passses, I will be using capturing groups
	    to get the paramters for (a) rgb8() or (b) rgb24() functions that
	    return SgrAttr like charm
	(3) in case map and two checks fail, I am going to return empty string and
	    false.
	(*) Eventually, I want to return an error (like a brand errors.New() one)
	    instead of a boolean.
*/

var errMap = errors.New("errMap")

func mapColor(s string) (SgrAttr, error) {
	re, err := regexp.Compile(`(?mi)^rgb(8|24)=(?:fg|bg)(?::\d{1,3}){1,3}$`)
	if err != nil {
		return "", errMap
	}

	if matchRegexp(re, s) {
		// Discern between 8 and 24
	}

	col, ok := colorMap[s]
	if !ok {
		return "", errMap
	}
	return col, nil
}

func matchRegexp(r *regexp.Regexp, val any) bool {
	valStr, ok := val.(string)
	if !ok {
		return false
	}
	return r.MatchString(valStr)
}
