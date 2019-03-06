// Copyright © 2019 Pawel Potrykus <pawel.potrykus@gmail.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package templates

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
	"unicode/utf8"
)

// Template stores exactly one row and related to it template of a data read from CSV file
type Template struct {
	Data            map[string]string
	TemplateName    string
	TemplateContent string
}

// New creates and returns pointer to the Template
func New(data map[string]string, templateName string, templateReader io.Reader) (*Template, error) {
	t := &Template{}

	t.Data = data
	t.TemplateName = templateName

	b, err := ioutil.ReadAll(templateReader)
	if err != nil {
		return nil, err
	}

	if !isUTF8(b) {
		return nil, fmt.Errorf("template file is not encoded in utf-8 or ascii")
	}

	b = normUTF8(b)

	t.TemplateContent = string(b)

	return t, nil
}

// SetGlobalVars sets additional variables to use while generating output from template
func (t *Template) SetGlobalVars(m map[string]string) {
	for k, v := range m {
		t.Data[k] = v
	}
}

// Fprintt fills template with data and write the results to 'w'. It returns number of characters written and an error (if any)
func Fprintt(w io.Writer, tplContent string, tplData map[string]string) (int, error) {
	tt, err := template.New("").Funcs(templateFuncs).Parse(tplContent)
	if err != nil {
		return 0, err
	}

	strWriter := &strings.Builder{}
	tt.Execute(strWriter, tplData)

	strReader := strings.NewReader(strWriter.String())
	io.Copy(w, strReader)

	return len(strWriter.String()), nil
}

// Sprintt fills template with data and returns it
func Sprintt(tplContent string, tplData map[string]string) string {
	strWriter := &strings.Builder{}

	_, err := Fprintt(strWriter, tplContent, tplData)
	if err != nil {
		return ""
	}

	return strWriter.String()
}

// Printt prints template filled with data to 'stdout'
func Printt(tplContent string, tplData map[string]string) {
	Fprintt(os.Stdout, tplContent, tplData)
}

// Execute executes template and outputs to 'w'
func (t *Template) Execute(w io.Writer) error {
	tt, err := template.New(t.TemplateName).Funcs(templateFuncs).Parse(t.TemplateContent)
	if err != nil {
		return err
	}

	return tt.Execute(w, t.Data)
}

// ReadCSV reads CSV file and returns data arranged in slice of maps
func ReadCSV(filename string, comma rune) ([]map[string]string, error) {
	m := make([]map[string]string, 0)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
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
			record[colName] = field
		}

		m = append(m, record)

	}

	return m, nil
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
