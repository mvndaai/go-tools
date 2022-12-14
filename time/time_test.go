package time_test

import (
	"fmt"
	"testing"
	gotime "time"

	"github.com/mvndaai/go-tools/time"
	"github.com/stretchr/testify/assert"
)

func TestConvertLongDuration(t *testing.T) {
	defaultMultiplier := map[string]float64{
		"d": 24,
		"w": 168,
	}

	tests := []struct {
		in         string
		out        string
		multiplier map[string]float64
		err        error
	}{
		{in: "1h", out: "1h"},
		{in: "1d1h", out: "24h1h"},
		{in: "--1d", out: "--24h"},
		{in: "-+-1d", out: "-+-24h"},
		{in: "+-1d", out: "+-24h"},
		{in: "-+1d", out: "-+24h"},
		{in: "", out: ""},
		{in: "5", out: "5"},
		{in: "1fortnight", out: "336h", multiplier: map[string]float64{"fortnight": 24 * 14, "w": 168}},
		{in: "..5d", err: fmt.Errorf(`time: invalid duration "..5d"`)},
		{in: "1f", out: "1f", multiplier: map[string]float64{"f|g": 1}},
		{in: "1g", out: "1g", multiplier: map[string]float64{"f|g": 1}},
		{in: "1f|g", out: "1h", multiplier: map[string]float64{"f|g": 1}},
		{in: "1f", out: "1f", multiplier: map[string]float64{".*": 3}},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			m := defaultMultiplier
			if tt.multiplier != nil {
				m = tt.multiplier
			}

			d, err := time.ConvertLongDuration(tt.in, m)
			assert.Equal(t, tt.out, d)
			assert.Equal(t, tt.err, err)
		})
	}
}

func TestParseLongDuration(t *testing.T) {
	tests := []struct {
		input    string
		equivent string
		err      error
	}{
		{input: "1h", equivent: "1h"},
		{input: "1d", equivent: "24h"},
		{input: "1w", equivent: fmt.Sprintf("%dh", 7*24)},
		{input: "1d1w1h", equivent: fmt.Sprintf("%dh", 7*24+24+1)},
		{input: "1plug", equivent: "1plug"},
		{input: "1d2m", equivent: "24h2m"},
		{input: "5", equivent: "5"},
		{input: "1d2h", equivent: "24h2h"},
		{input: "1d1d", equivent: "48h"},
		{input: ".5h", equivent: ".5h"},
		{input: "1d2h", equivent: "24h2h"},
		{input: "-1d", equivent: "-24h"},
		{input: "--2m", equivent: "--2m"},
		{input: "..2d", err: fmt.Errorf(`time: invalid duration "..2d"`)},
		{input: "..2m", equivent: "..2m"},
		{input: "..2,", equivent: "..2,"},
		{input: "1d2", equivent: "24h2"},
		{input: "+-2h", equivent: "+-2h"},
		{input: "", equivent: ""},
		{input: "0", equivent: "0"},
		{input: "-0", equivent: "-0"},
		{input: "+0", equivent: "+0"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ld, ldErr := time.ParseLongDuration(tt.input)
			n, nErr := gotime.ParseDuration(tt.equivent)
			assert.Equal(t, n, ld)
			if tt.err != nil {
				nErr = tt.err
			}
			assert.Equal(t, nErr, ldErr)
		})
	}
}

func TestParseDurationWithUnits(t *testing.T) {
	tests := []struct {
		input    string
		equivent string
		err      error
	}{
		{input: "1h", equivent: "1h"},
		{input: "1d", equivent: "24h"},
		{input: "1w", equivent: fmt.Sprintf("%dh", 7*24)},
		{input: "1d1w1h", equivent: fmt.Sprintf("%dh", 7*24+24+1)},
		{input: "1plug", equivent: "1plug"},
		{input: "1d2m", equivent: "24h2m"},
		{input: "5", equivent: "5"},
		{input: "1d2h", equivent: "24h2h"},
		{input: "1d1d", equivent: "48h"},
		{input: ".5h", equivent: ".5h"},
		{input: "1d2h", equivent: "24h2h"},
		{input: "-1d", equivent: "-24h"},
		{input: "--2m", equivent: "--2m"},
		{input: "..2d", err: fmt.Errorf(`time: invalid duration "..2d"`)},
		{input: "..2m", equivent: "..2m"},
		{input: "..2,", equivent: "..2,"},
		{input: "1d2", err: fmt.Errorf(`time: missing unit in duration "1d2"`)},
		{input: "+-2h", equivent: "+-2h"},
		{input: "", equivent: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			unitMap := map[string]uint64{
				"ns": uint64(gotime.Nanosecond),
				"us": uint64(gotime.Microsecond),
				"??s": uint64(gotime.Microsecond), // U+00B5 = micro symbol
				"??s": uint64(gotime.Microsecond), // U+03BC = Greek letter mu
				"ms": uint64(gotime.Millisecond),
				"s":  uint64(gotime.Second),
				"m":  uint64(gotime.Minute),
				"h":  uint64(gotime.Hour),
				"d":  uint64(gotime.Hour * 24),
				"w":  uint64(gotime.Hour * 168),
			}

			ld, ldErr := time.ParseDurationWithUnits(tt.input, unitMap)
			n, nErr := gotime.ParseDuration(tt.equivent)
			assert.Equal(t, n, ld)
			assert.Equal(t, n, ld)
			if tt.err != nil {
				nErr = tt.err
			}
			assert.Equal(t, nErr, ldErr)
		})
	}
}

func TestCompareParseDurations(t *testing.T) {
	tests := []struct {
		input       string
		errorPrefix string
	}{
		{input: "1h"},
		{input: "1d"},
		{input: "1w"},
		{input: "1d1w1h"},
		{input: "1plug"},
		{input: "1d2m"},
		{input: "5"},
		{input: "1d2h"},
		{input: "1d1d"},
		{input: ".5h"},
		{input: "1d2h"},
		{input: "-1d"},
		{input: "--2m"},
		{input: "..2d"},
		{input: "..2m"},
		{input: "..2,"},
		{input: "1d2", errorPrefix: "time: missing unit in duration"},
		{input: "+-2h", errorPrefix: "time: invalid duration"},
		{input: ""},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {

			unitMap := map[string]uint64{
				"ns": uint64(gotime.Nanosecond),
				"us": uint64(gotime.Microsecond),
				"??s": uint64(gotime.Microsecond), // U+00B5 = micro symbol
				"??s": uint64(gotime.Microsecond), // U+03BC = Greek letter mu
				"ms": uint64(gotime.Millisecond),
				"s":  uint64(gotime.Second),
				"m":  uint64(gotime.Minute),
				"h":  uint64(gotime.Hour),
				"d":  uint64(time.Day),
				"w":  uint64(time.Week),
			}

			wu, uwErr := time.ParseDurationWithUnits(tt.input, unitMap)
			ld, ldErr := time.ParseLongDuration(tt.input)
			assert.Equal(t, wu, ld)

			if uwErr != ldErr {
				if tt.errorPrefix != "" {
					assert.Contains(t, uwErr.Error(), tt.errorPrefix)
					assert.Contains(t, ldErr.Error(), tt.errorPrefix)
				} else {
					assert.Equal(t, uwErr, ldErr)
				}
			}
		})
	}
}

func BenchmarkTestParseDurationWithUnits(b *testing.B) {
	const ds = "2w3d5h2m"
	unitMap := map[string]uint64{
		"ns": uint64(gotime.Nanosecond),
		"us": uint64(gotime.Microsecond),
		"??s": uint64(gotime.Microsecond), // U+00B5 = micro symbol
		"??s": uint64(gotime.Microsecond), // U+03BC = Greek letter mu
		"ms": uint64(gotime.Millisecond),
		"s":  uint64(gotime.Second),
		"m":  uint64(gotime.Minute),
		"h":  uint64(gotime.Hour),
		"d":  uint64(time.Day),
		"w":  uint64(time.Week),
	}

	for i := 0; i < b.N; i++ {
		time.ParseDurationWithUnits(ds, unitMap)
	}
}

func BenchmarkTestParseLongDuration(b *testing.B) {
	const ds = "2w3d5h2m"

	for i := 0; i < b.N; i++ {
		time.ParseLongDuration(ds)
	}
}

// BenchmarkTestParseDurationWithUnits-10	17278627	67.03 ns/op	   0 B/op	 0 allocs/op
// BenchmarkTestParseLongDuration-10     	  464401	2515  ns/op	3171 B/op	39 allocs/op
