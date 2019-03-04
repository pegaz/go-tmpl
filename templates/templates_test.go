package templates

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

var tplContent = "example template with {{.Name}}\n"
var tplData = map[string]string{"Name": "*name*"}

func redirectStdout() (func() string, error) {
	stdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		return nil, err
	}

	os.Stdout = w

	defered := func() string {
		w.Close()

		var buf bytes.Buffer
		io.Copy(&buf, r)

		os.Stdout = stdout

		return buf.String()
	}

	return defered, nil
}

func TestExecute(t *testing.T) {
	r := strings.NewReader(tplContent)

	tpl, err := New(tplData, "test_template", r)
	if err != nil {
		t.Error(err)
	}

	w := &strings.Builder{}

	err = tpl.Execute(w)
	if err != nil {
		t.Error(err)
	}

	referenceString := strings.Replace(tplContent, "{{.Name}}", tplData["Name"], 1)
	if w.String() != referenceString {
		t.Error("expected to get exactly the same string from template's output as the reference, instead it is different")
	}
}

func TestPrintt(t *testing.T) {
	defered, err := redirectStdout()
	if err != nil {
		t.Error(err)
	}

	Printt(tplContent, tplData)

	got := defered()

	referenceString := strings.Replace(tplContent, "{{.Name}}", tplData["Name"], 1)
	if got != referenceString {
		t.Error("expected to get exactly the same string from template's output as the reference, instead it is different")
	}
}

func TestSprintt(t *testing.T) {
	got := Sprintt(tplContent, tplData)

	referenceString := strings.Replace(tplContent, "{{.Name}}", tplData["Name"], 1)
	if got != referenceString {
		t.Error("expected to get exactly the same string from template's output as the reference, instead it is different")
	}
}

func TestFprintt(t *testing.T) {
	strWriter := &strings.Builder{}

	n, err := Fprintt(strWriter, tplContent, tplData)
	if err != nil {
		t.Error(err)
	}

	referenceString := strings.Replace(tplContent, "{{.Name}}", tplData["Name"], 1)
	if len(referenceString) != n {
		t.Errorf("expected to get the output of length %d, instead got %d", len(referenceString), n)
	}

	if strWriter.String() != referenceString {
		t.Error("expected to get exactly the same string from template's output as the reference, instead it is different")
	}
}
