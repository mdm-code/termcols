package termcols

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var layerMap map[string]Layer = map[string]Layer{"fg": FG, "bg": BG}

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

var (
	errMap = errors.New("Color mapping error")
)

// MapColors attempts to interpret string elements of the ss slice as a set of
// predefined colors/styles or an RGB8/24 string pattern that is expected to
// come in one of the case-insensitive patterns listed below. Otherwise the
// function returns an empty slice and errMap.
//
//   RGB 8  : rgb8=[fg|bg]:[0-255]
//   RGB 24 : rgb24=[fg|bg]:[0-255]:[0-255]:[0-255]
func MapColors(ss []string) ([]SgrAttr, error) {
	result := make([]SgrAttr, 0, 3)
	for _, s := range ss {
		attr, err := MapColor(s)
		if err != nil {
			return []SgrAttr{}, errMap
		}
		result = append(result, attr)
	}
	return result, nil
}

// MapColor attempts to interpret the string s as either one of the predefined
// colors/styles or an RGB8/24 string pattern that is expected to come in the
// one of the case-insensitive patterns listed below. Otherwise the function
// returns an empty string of type SgrAttr and errMap.
//
//   RGB 8  : rgb8=[fg|bg]:[0-255]
//   RGB 24 : rgb24=[fg|bg]:[0-255]:[0-255]:[0-255]
func MapColor(s string) (SgrAttr, error) {
	col, ok := colorMap[strings.ToLower(s)]
	if ok {
		return col, nil
	}
	re8 := regexp.MustCompile(
		`(?mi)^rgb8=(?P<layer>fg|bg):(?P<color>\d{1,3})$`,
	)
	if matchRegexp(re8, s) {
		col, ok := collateRgb8(re8, s)
		if !ok {
			return "", errMap
		}
		return col, nil
	}
	re24 := regexp.MustCompile(
		`(?mi)^rgb24=(?P<layer>fg|bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>\d{1,3})$`,
	)
	if matchRegexp(re24, s) {
		col, ok := collateRgb24(re24, s)
		if !ok {
			return "", errMap
		}
		return col, nil
	}
	return "", errMap
}

// MatchRegexp checks if val matches the provided regex r.
func matchRegexp(r *regexp.Regexp, val any) bool {
	valStr, ok := val.(string)
	if !ok {
		return false
	}
	return r.MatchString(valStr)
}

// CollateRgb8 parses string s into SgrAttr using the provided regex r.
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

// CollateRgb24 parses string s into SgrAttr using the provided regex r.
func collateRgb24(r *regexp.Regexp, s string) (SgrAttr, bool) {
	params := getParams(r, s)

	// TODO (michal): move the layer check to a separate function
	l, ok := getLayer(params)
	if !ok {
		return "", false
	}

	// TODO (michal): move color check to a separate function
	// NOTE: This should make it easier to target unit tests
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

// GetParams extracts regex named capturing group names and values.
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

// GetLayer returns Layer based on the key in the params map.
func getLayer(params map[string]string) (Layer, bool) {
	lr, ok := params["layer"]
	if !ok {
		return "", false
	}
	lr = strings.ToLower(lr)
	l, ok := layerMap[lr]
	if !ok {
		return "", false
	}
	return l, true
}

// ValidUint verifies if the integer i falls in range [0, 255] of uint8.
func validUint(i int) bool {
	if i >= 0 && i <= 255 {
		return true
	}
	return false
}
