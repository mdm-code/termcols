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
	// ErrMap indicates that there were issues with disambiguating color names.
	ErrMap = errors.New("Color mapping error")
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
			return []SgrAttr{}, ErrMap
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
			return "", ErrMap
		}
		return col, nil
	}
	re24 := regexp.MustCompile(
		`(?mi)^rgb24=(?P<layer>fg|bg):(?P<r>\d{1,3}):(?P<g>\d{1,3}):(?P<b>\d{1,3})$`,
	)
	if matchRegexp(re24, s) {
		col, ok := collateRgb24(re24, s)
		if !ok {
			return "", ErrMap
		}
		return col, nil
	}
	return "", ErrMap
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
	l, ok := getLayer(params)
	if !ok {
		return "", false
	}
	c, ok := getColor(params, "color")
	if !ok {
		return "", false
	}
	result := Rgb8(l, c)
	return result, true
}

// CollateRgb24 parses string s into SgrAttr using the provided regex r.
func collateRgb24(r *regexp.Regexp, s string) (SgrAttr, bool) {
	params := getParams(r, s)
	l, ok := getLayer(params)
	if !ok {
		return "", false
	}
	red, ok := getColor(params, "r")
	if !ok {
		return "", false
	}
	green, ok := getColor(params, "g")
	if !ok {
		return "", false
	}
	blue, ok := getColor(params, "b")
	if !ok {
		return "", false
	}
	result := Rgb24(l, red, green, blue)
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

// GetColor returns the color based on the key in the params map.
func getColor(params map[string]string, key string) (uint8, bool) {
	val, ok := params[key]
	if !ok {
		return 0, false
	}
	col, err := strconv.Atoi(val)
	if err != nil {
		return 0, false
	}
	if ok := validUint8(col); !ok {
		return 0, false
	}
	return uint8(col), true
}

// ValidUint8 verifies if the integer i falls in range [0, 255] of uint8.
func validUint8(i int) bool {
	if i >= 0 && i <= 255 {
		return true
	}
	return false
}
