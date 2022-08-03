package time

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	Day  time.Duration = 24 * time.Hour
	Week               = 168 * time.Hour
)

func ConvertLongDuration(s string, hourMultiplier map[string]float64) (string, error) {
	if len(s) == 0 {
		return s, nil
	}

	var converted = s

	var operator string
	if s[0] == '+' || s[0] == '-' {
		operator = string(s[0])
		converted = converted[1:]
	}

	var units []string
	for unit, _ := range hourMultiplier {
		units = append(units, unit)
	}

	regexStr := fmt.Sprintf(`(?m)([0-9.]+)([%s]+)`, strings.Join(units, "|"))
	var re = regexp.MustCompile(regexStr)
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		if m, ok := hourMultiplier[match[2]]; ok {
			v, err := strconv.ParseFloat(match[1], 64)
			if err != nil {
				return "", fmt.Errorf(`time: invalid duration "%s"`, s)
			}
			hours := v * m
			h := strconv.FormatFloat(hours, 'f', -1, 64) + "h"
			converted = strings.Replace(converted, match[0], h, 1)
		}
	}

	converted = string(operator) + converted

	return converted, nil
}

// ParseLongDuration allows w,d in time.ParseDuration with a d=24h and w=7d
func ParseLongDuration(s string) (time.Duration, error) {
	multiplier := map[string]float64{
		"d": 24,
		"w": 168,
	}
	converted, err := ConvertLongDuration(s, multiplier)
	if err != nil {
		return 0, err
	}
	return time.ParseDuration(converted)
}

func ParseDurationWithUnits(s string, unitMap map[string]uint64) (time.Duration, error) {
	// [-+]?([0-9]*(\.[0-9]*)?[a-z]+)+
	orig := s
	var d uint64
	neg := false

	// Consume [-+]?
	if s != "" {
		c := s[0]
		if c == '-' || c == '+' {
			neg = c == '-'
			s = s[1:]
		}
	}
	// Special case: if all that is left is "0", this is zero.
	if s == "0" {
		return 0, nil
	}
	if s == "" {
		return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
	}
	for s != "" {
		var (
			v, f  uint64      // integers before, after decimal point
			scale float64 = 1 // value = v + f/scale
		)

		var err error

		// The next character must be [0-9.]
		if !(s[0] == '.' || '0' <= s[0] && s[0] <= '9') {
			return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
		}
		// Consume [0-9]*
		pl := len(s)
		v, s, err = leadingInt(s)
		if err != nil {
			return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
		}
		pre := pl != len(s) // whether we consumed anything before a period

		// Consume (\.[0-9]*)?
		post := false
		if s != "" && s[0] == '.' {
			s = s[1:]
			pl := len(s)
			f, scale, s = leadingFraction(s)
			post = pl != len(s)
		}
		if !pre && !post {
			// no digits (e.g. ".s" or "-.s")
			return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
		}

		// Consume unit.
		i := 0
		for ; i < len(s); i++ {
			c := s[i]
			if c == '.' || '0' <= c && c <= '9' {
				break
			}
		}
		if i == 0 {
			return 0, fmt.Errorf(`time: missing unit in duration "%s"`, orig)
		}
		u := s[:i]
		s = s[i:]
		unit, ok := unitMap[u]
		if !ok {
			return 0, fmt.Errorf(`time: unknown unit "%v" in duration "%s"`, u, orig)
		}
		if v > 1<<63/unit {
			// overflow
			return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
		}
		v *= unit
		if f > 0 {
			// float64 is needed to be nanosecond accurate for fractions of hours.
			// v >= 0 && (f*unit/scale) <= 3.6e+12 (ns/h, h is the largest unit)
			v += uint64(float64(f) * (float64(unit) / scale))
			if v > 1<<63 {
				// overflow
				return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
			}
		}
		d += v
		if d > 1<<63 {
			return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
		}
	}
	if neg {
		return -time.Duration(d), nil
	}
	if d > 1<<63-1 {
		return 0, fmt.Errorf(`time: invalid duration "%s"`, orig)
	}
	return time.Duration(d), nil

}

var errLeadingInt = errors.New("time: bad [0-9]*") // never printed
func leadingInt(s string) (x uint64, rem string, err error) {
	i := 0
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if x > 1<<63/10 {
			// overflow
			return 0, "", errLeadingInt
		}
		x = x*10 + uint64(c) - '0'
		if x > 1<<63 {
			// overflow
			return 0, "", errLeadingInt
		}
	}
	return x, s[i:], nil
}

// leadingFraction consumes the leading [0-9]* from s.
// It is used only for fractions, so does not return an error on overflow,
// it just stops accumulating precision.
func leadingFraction(s string) (x uint64, scale float64, rem string) {
	i := 0
	scale = 1
	overflow := false
	for ; i < len(s); i++ {
		c := s[i]
		if c < '0' || c > '9' {
			break
		}
		if overflow {
			continue
		}
		if x > (1<<63-1)/10 {
			// It's possible for overflow to give a positive number, so take care.
			overflow = true
			continue
		}
		y := x*10 + uint64(c) - '0'
		if y > 1<<63 {
			overflow = true
			continue
		}
		x = y
		scale *= 10
	}
	return x, scale, s[i:]
}
