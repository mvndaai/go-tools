package redact

import (
	"testing"
)

func TestWords(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{in: "", expected: ""},
		{in: "Bob", expected: "B**"},
		{in: "Bob Jones", expected: "B** J****"},
		{in: "Bob K Jones", expected: "B** K J****"},
		{in: "The 🐶🪵 is brown.", expected: "T** 🐶* i* b*****"},
		{in: "many   spaces", expected: "m***   s*****"},
		{in: "123 w 450 e", expected: "1** w 4** e"},
		{in: "220 Main Street", expected: "2** M*** S*****"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if actual := Words(tt.in); tt.expected != actual {
				t.Errorf("'%s' did not match expected '%s'", actual, tt.expected)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{in: "", expected: ""},
		{in: "example@example.com", expected: "e******@example.com"},
		{in: "🐶🪵@b.com", expected: "🐶*@b.com"},
		{in: "joe+s@m@gmail.com", expected: "j****@m@gmail.com"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if actual := Email(tt.in); tt.expected != actual {
				t.Errorf("'%s' did not match expected '%s'", actual, tt.expected)
			}
		})
	}
}

func TestPhone(t *testing.T) {
	tests := []struct {
		in       string
		expected string
	}{
		{in: "", expected: ""},
		{in: "1", expected: "*"},
		{in: "1234", expected: "****"},
		{in: "12345", expected: "1****"},
		{in: "1-2-34", expected: "1-2-**"},
		{in: "801-123-1234", expected: "801-123-****"},
		{in: "801.123.1234", expected: "801.123.****"},
		{in: "801.123.123", expected: "801.123****"},
		{in: "12🐶", expected: "***"},
	}

	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			if actual := Phone(tt.in); tt.expected != actual {
				t.Errorf("'%s' did not match expected '%s'", actual, tt.expected)
			}
		})
	}
}
