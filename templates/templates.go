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
	"encoding/csv"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
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
	headers := make([]string, 0)

	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	reader := csv.NewReader(f)

	reader.Comma = comma

	var i int
	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		record := make(map[string]string)

		for j, elem := range line {
			if i == 0 {
				headers = append(headers, elem)
			} else {
				record[headers[j]] = elem
			}
		}

		if i == 0 {
			i++
			continue
		}

		m = append(m, record)

	}

	return m, nil
}

// func Execute(tmpls []*Template) error {
// 	templateNames := make(map[string]*template.Template)

// 	var err error
// 	var errCount int
// 	for _, t := range tmpls {
// 		if _, ok := templateNames[t.TemplateName]; !ok {
// 			var templateFile *os.File
// 			var templateContentBytes []byte
// 			var templateContent string

// 			templateFile, err = os.Open(t.TemplateName)
// 			if err != nil {
// 				fmt.Printf("!!! %s\n", err)
// 				errCount++
// 				continue
// 			}
// 			defer templateFile.Close()

// 			templateContentBytes, err = ioutil.ReadAll(templateFile)
// 			if err != nil {
// 				fmt.Printf("!!! %s\n", err)
// 				errCount++
// 				continue
// 			}

// 			templateContent = string(templateContentBytes)
// 			templateNames[t.TemplateName], err = template.New(t.TemplateName).Funcs(templateFuncs).Parse(templateContent)
// 			if err != nil {
// 				fmt.Printf("!!! %s\n", err)
// 				errCount++
// 				continue
// 			}

// 			flags := os.O_WRONLY | os.O_CREATE

// 			err = templateNames[t.TemplateName].Execute(file, t.Data)
// 			if err != nil {
// 				fmt.Printf("!!! %s\n", err)
// 				errCount++
// 				continue
// 			}

// 			if _, ok := t.Data["hostname"]; ok {
// 				fmt.Printf(">>> generated output for %s using %s\n", t.OutputFileName, t.TemplateName)
// 			} else {
// 				fmt.Printf(">>> generated output using %s\n", t.TemplateName)
// 			}
// 		}
// 	}

// 	if errCount > 0 {
// 		return fmt.Errorf("there had been an error(s) while generating output")
// 	}

// 	return nil
// }
