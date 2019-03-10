package templates

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"strings"
	"unicode/utf8"
)

var charMap = map[rune]rune{
	'Ą': 'A',
	'ą': 'a',
	'Ć': 'C',
	'ć': 'c',
	'Ę': 'E',
	'ę': 'e',
	'Ł': 'L',
	'ł': 'l',
	'Ń': 'N',
	'ń': 'n',
	'Ó': 'O',
	'ó': 'o',
	'Ś': 'S',
	'ś': 's',
	'Ź': 'Z',
	'ź': 'z',
	'Ż': 'Z',
	'ż': 'z',
}

// isUTF8 checks if given byte's slice is correctly encoded with UTF-8
func isUTF8(b []byte) bool {
	return utf8.Valid(b)
}

// normUTF8 checks if given byte's slice begins with BOM (Byte Order Mark) and if so, truncates it and returns plain UTF-8
func normUTF8(b []byte) []byte {
	if bytes.Compare(b[:3], []byte{0xef, 0xbb, 0xbf}) == 0 {
		return b[3:]
	}

	return b
}

// Normalize replaces accent characters in a given string with an ASCII counterpart and returns it
func Normalize(s string) string {
	var nstr strings.Builder

	for _, ch := range s {
		if ch > 0x80 {
			if v, ok := charMap[ch]; ok {
				nstr.WriteRune(v)
			}
			continue
		} else {
			nstr.WriteRune(ch)
		}
	}

	return nstr.String()
}

// ReadCSV reads from r and returns data arranged in slice of maps
func ReadCSV(r io.Reader, comma rune) ([]map[string]string, error) {
	m := make([]map[string]string, 0)

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}

	if !isUTF8(b) {
		return nil, fmt.Errorf("csv data file is not encoded in utf-8 or ascii")
	}

	b = normUTF8(b)

	reader := csv.NewReader(bytes.NewReader(b))
	reader.Comma = comma

	csvContent, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	// when we have a header separated...
	header := make([]string, 0)
	header = append(header, csvContent[0]...)

	// ...we need to omit it
	csvContent = csvContent[1:]

	for _, line := range csvContent {
		record := make(map[string]string)

		for j, field := range line {
			colName := header[j]
			record[colName] = Normalize(field)
		}

		m = append(m, record)

	}

	return m, nil
}
