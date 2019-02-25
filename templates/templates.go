package templates

import (
	"encoding/csv"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"text/template"
)

// Template stores exactly one row and related to it template of a data read from CSV file
type Template struct {
	Data         map[string]string
	TemplateName string

	OutputPath       string
	OutputFileName   string
	SingleFileOutput bool
}

// New creates and returns pointer to the Template
func New(data map[string]string, templateName string, path string, filename string, singleFileOutput bool) (*Template, error) {
	t := &Template{}

	t.Data = data
	t.TemplateName = templateName
	t.OutputPath = path
	t.OutputFileName = filename
	t.SingleFileOutput = singleFileOutput

	return t, nil
}

func (t *Template) SetGlobalVars(m map[string]string) {
	for k, v := range m {
		t.Data[k] = v
	}
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

func Execute(tmpls []*Template) error {
	templateNames := make(map[string]*template.Template)

	var err error
	var errCount int
	for _, t := range tmpls {
		if _, ok := templateNames[t.TemplateName]; !ok {
			var templateFile *os.File
			var templateContentBytes []byte
			var templateContent string

			templateFile, err = os.Open(t.TemplateName)
			if err != nil {
				fmt.Printf("!!! %s\n", err)
				errCount++
				continue
			}
			defer templateFile.Close()

			templateContentBytes, err = ioutil.ReadAll(templateFile)
			if err != nil {
				fmt.Printf("!!! %s\n", err)
				errCount++
				continue
			}

			templateContent = string(templateContentBytes)
			templateNames[t.TemplateName], err = template.New(t.TemplateName).Funcs(templateFuncs).Parse(templateContent)
			if err != nil {
				fmt.Printf("!!! %s\n", err)
				errCount++
				continue
			}

			flags := os.O_WRONLY | os.O_CREATE
			if t.SingleFileOutput == true {
				flags = flags | os.O_APPEND
			}

			file, err := os.OpenFile(t.OutputPath+t.OutputFileName, flags, 0644)
			if err != nil {
				fmt.Printf("!!! %s\n", err)
				errCount++
				continue
			}
			defer file.Close()

			err = templateNames[t.TemplateName].Execute(file, t.Data)
			if err != nil {
				fmt.Printf("!!! %s\n", err)
				errCount++
				continue
			}

			if _, ok := t.Data["hostname"]; ok {
				fmt.Printf(">>> generated output for %s using %s\n", t.OutputFileName, t.TemplateName)
			} else {
				fmt.Printf(">>> generated output using %s\n", t.TemplateName)
			}
		}
	}

	if errCount > 0 {
		return fmt.Errorf("there had been an error(s) while generating output")
	}

	return nil
}
