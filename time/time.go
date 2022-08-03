package time

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// ParseLongDuration allows w,d in time.ParseDuration with a d=24h and w=7d
func ParseLongDuration(s string) (time.Duration, error) {
	if len(s) == 0 {
		return time.ParseDuration(s)
	}

	var converted = s

	var operator string
	if converted[0] == '+' {
		operator = "+"
		converted = converted[1:]
	}
	if s[0] == '-' {
		operator = "-"
		converted = converted[1:]
	}

	var hours float64
	fmt.Println("original converted", converted)
	var re = regexp.MustCompile(`(?m)([0-9.]+)([a-z]+)`)
	for _, match := range re.FindAllStringSubmatch(s, -1) {
		var multiplier float64
		switch match[2] {
		case "h":
			multiplier = 1
		case "d":
			multiplier = 24
		case "w":
			multiplier = 168
		default:
			continue
		}
		v, err := strconv.ParseFloat(match[1], 64)
		if err != nil {
			return 0, fmt.Errorf(`time: invalid duration "%s"`, s)
		}
		hours += v * multiplier
		converted = strings.Replace(converted, match[0], "", 1)
	}

	fmt.Println("converted before hours", converted)
	if hours != 0 {
		h := fmt.Sprintf("%f", hours)
		h = strings.TrimRight(strings.TrimRight(h, "0"), ".")
		converted = fmt.Sprintf("%sh%s", h, converted)
	}
	fmt.Println("converted before operator", converted)
	converted = fmt.Sprintf("%s%s", string(operator), converted)

	return time.ParseDuration(converted)
}
