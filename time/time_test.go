package time_test

import (
	"fmt"
	"testing"
	gotime "time"

	"github.com/mvndaai/go-tools/time"
	"github.com/stretchr/testify/assert"
)

func TestParseLongDuration(t *testing.T) {
	tests := []struct {
		input           string
		equivent        string
		invalidDuration string
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
		{input: "..2d", invalidDuration: "..2d"},
		{input: "..2m", equivent: "..2m"},
		{input: "..2,", equivent: "..2,"},
		{input: "1d2", equivent: "24h2"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ld, ldErr := time.ParseLongDuration(tt.input)
			n, nErr := gotime.ParseDuration(tt.equivent)
			assert.Equal(t, n, ld)

			if tt.invalidDuration != "" {
				nErr = fmt.Errorf(`time: invalid duration "%s"`, tt.invalidDuration)
			}
			assert.Equal(t, nErr, ldErr)
		})
	}
}
