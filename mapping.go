package termcols

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
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

var errMap = errors.New("Color mapping error")

func mapColor(s string) (SgrAttr, error) {
	col, ok := colorMap[s]
	if ok {
		return col, nil
	}
	// TODO (michal): remove redundant top-level regex and just do re8 and re24
	p := `(?mi)^rgb(?:8|24)=(?:fg|bg)(?::\d{1,3}){1,3}$`
	re, err := regexp.Compile(p)
	if err != nil {
		return "", errMap
	}
	if matchRegexp(re, s) {
		p8 := `(?mi)^rgb8=(?P<layer>fg|bg):(?P<color>\d{1,3})$`
		re8, err := regexp.Compile(p8)
		if err != nil {
			return "", errMap
		}
		if matchRegexp(re8, s) {
			var c8 SgrAttr
			if c8, ok := collateRgb8(re8, s); !ok {
				return c8, errMap
			}
			return c8, nil
		}
		p24 := `(?mi)^rgb24=(?P<layer>fg|bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>\d{1,3})$`
		re24, err := regexp.Compile(p24)
		if err != nil {
			return "", errMap
		}
		if matchRegexp(re24, s) {
			var c24 SgrAttr
			if c24, ok := collateRgb24(re24, s); !ok {
				return c24, errMap
			}
			return c24, nil
		}
	}
	return "", errMap
}

func matchRegexp(r *regexp.Regexp, val any) bool {
	valStr, ok := val.(string)
	if !ok {
		return false
	}
	return r.MatchString(valStr)
}

func collateRgb8(r *regexp.Regexp, s string) (SgrAttr, bool) {
	params := getParams(r, s)
	lr, ok := params["layer"]
	if !ok {
		return "", ok
	}
	lr = strings.ToLower(lr)
	l, ok := layerMap[lr]
	if !ok {
		return "", ok
	}

	c, ok := params["color"]
	if !ok {
		return "", ok
	}
	col, err := strconv.Atoi(c)
	if err != nil {
		return "", false
	}
	if ok := validUint(col); !ok {
		return "", ok
	}

	result := Rgb8(l, uint8(col))
	return result, true
}

func collateRgb24(r *regexp.Regexp, s string) (SgrAttr, bool) {
	params := getParams(r, s)
	lr, ok := params["layer"]
	if !ok {
		return "", ok
	}
	lr = strings.ToLower(lr)
	l, ok := layerMap[lr]
	if !ok {
		return "", ok
	}

	// Red
	cr, ok := params["r"]
	if !ok {
		return "", ok
	}
	rcol, err := strconv.Atoi(cr)
	if err != nil {
		return "", false
	}
	if ok := validUint(rcol); !ok {
		return "", ok
	}

	// Green
	cg, ok := params["g"]
	if !ok {
		return "", ok
	}
	gcol, err := strconv.Atoi(cg)
	if err != nil {
		return "", false
	}
	if ok = validUint(gcol); !ok {
		return "", ok
	}

	// Blue
	cb, ok := params["b"]
	if !ok {
		return "", ok
	}
	bcol, err := strconv.Atoi(cb)
	if err != nil {
		return "", false
	}
	if ok := validUint(bcol); !ok {
		return "", ok
	}

	result := Rgb24(l, uint8(rcol), uint8(gcol), uint8(bcol))
	return result, true
}

func getParams(r *regexp.Regexp, s string) map[string]string {
	match := r.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			result[name] = match[i]
		}
	}
	return result
}

func validUint(i int) bool {
	if i >= 0 && i <= 255 {
		return true
	}
	return false
}
