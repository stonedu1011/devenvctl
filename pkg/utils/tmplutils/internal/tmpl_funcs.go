package internal

import (
	"encoding"
	"fmt"
	"math"
	"strconv"
	"strings"
	"text/template"
)

// Note: https://pkg.go.dev/text/template#hdr-Pipelines chainable argument should be the last parameter of any function

var (
	TmplFuncMap = template.FuncMap{
		"cap":   Capped,
		"pad":   Padding,
		"join":  Join,
	}
)

// Padding example: `{{padding -6 value}}` "{{padding 10 value}}"
func Padding(padding int, v interface{}) string {
	tag := "%" + strconv.Itoa(padding) + "v"
	return fmt.Sprintf(tag, v)
}

// Capped truncate given value to specified length
// if cap > 0: with tailing "..." if truncated
// if cap < 0: with middle "..." if truncated
func Capped(cap int, v interface{}) string {
	c := int(math.Abs(float64(cap)))
	s := Sprint(v)
	if len(s) <= c {
		return s
	}
	if cap > 0 {
		return fmt.Sprintf("%." + strconv.Itoa(c - 3) + "s...", s)
	} else if cap < 0 {
		lead := (c - 3) / 2
		tail := c - lead - 3
		return fmt.Sprintf("%." + strconv.Itoa(lead) + "s...%s", s, s[len(s)-tail:])
	} else {
		return ""
	}
}

func Join(sep string, values ...interface{}) string {
	strs := make([]string, 0, len(values))
	for _, v := range values {
		s := Sprint(v)
		if s != "" {
			strs = append(strs, s)
		}
	}
	str := strings.Join(strs, sep)
	return str
}

func Sprint(val interface{}) string {
	switch v := val.(type) {
	case nil:
		return ""
	case string:
		return v
	case []byte:
		return string(v)
	case fmt.Stringer:
		return v.String()
	case encoding.TextMarshaler:
		if s, e := v.MarshalText(); e == nil {
			return string(s)
		}
	}
	return fmt.Sprintf("%v", val)
}

