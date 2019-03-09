package templates

import (
	"bytes"
	"testing"
)

var testCasesNormalize = []struct {
	input    string
	expected string
}{
	{
		input:    "ĄĆĘŁŃÓŚŻŹąćęłńóśźż",
		expected: "ACELNOSZZacelnoszz",
	},
}

var testCasesUTF8 = []struct {
	str    string
	isUTF8 bool
	bom    bool
}{
	{str: "♠ ♣ ♥ ♦", isUTF8: true},
	{str: string([]byte{0x41, 0x42, 0x43}), isUTF8: true},
	{str: string([]byte{0xef, 0xbb, 0xbf, 0x41, 0x42, 0x43}), isUTF8: true, bom: true},
	{str: string([]byte{0xd8, 0x01, 0xdc, 0x37}), isUTF8: false},
	{str: string([]byte{0x01, 0xd8, 0x37, 0xdc}), isUTF8: false},
}

func TestIsUTF8(t *testing.T) {
	for _, tc := range testCasesUTF8 {
		if isUTF8([]byte(tc.str)) != tc.isUTF8 {
			t.Errorf("expected to get %t, instead got %t", tc.isUTF8, !tc.isUTF8)
		}
	}
}

func TestNormUTF8(t *testing.T) {
	for _, tc := range testCasesUTF8 {
		if tc.isUTF8 {
			normalized := normUTF8([]byte(tc.str))
			if tc.bom == true {
				if bytes.Compare(normalized, []byte(tc.str[3:])) != 0 {
					t.Error("expected to get normalized UTF8 string, instead got something different")
				}
			} else {
				if bytes.Compare(normalized, []byte(tc.str)) != 0 {
					t.Error("expected to get normalized UTF8 string, instead got something different")
				}
			}
		}

	}
}

func TestNormalize(t *testing.T) {
	for _, tc := range testCasesNormalize {
		if Normalize(tc.input) != tc.expected {
			t.Error("expected to get normalized string, but something went wrong")
		}
	}
}

func BenchmarkNormalize(b *testing.B) {
	var str = "ĄĆĘŁŃÓŚŻŹąćęłńóśźż"

	for n := 0; n < b.N; n++ {
		Normalize(str)
	}
}
