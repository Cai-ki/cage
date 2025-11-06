package helper

import (
	"reflect"
	"strings"
	"text/template"

	"github.com/Cai-ki/cage/llm"
)

var templateCache = make(map[string]*template.Template)

type Param interface {
	Prompt() string
	Prepare() error
	Do() (string, error)
}

func Parse(param Param) (string, error) {
	if err := param.Prepare(); err != nil {
		return "", err
	}

	typeName := reflect.TypeOf(param).Name()
	tmpl, ok := templateCache[typeName]
	if !ok {
		var err error
		tmpl, err = template.New(typeName).Parse(strings.TrimSpace(param.Prompt()))
		if err != nil {
			return "", err
		}
		templateCache[typeName] = tmpl
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, param); err != nil {
		return "", err
	}

	return buf.String(), nil
}

func Do(param Param) (string, error) {
	prompt, err := Parse(param)
	if err != nil {
		return "", err
	}

	return llm.CompletionBySystem(prompt)
}
