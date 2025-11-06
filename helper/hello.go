package helper

import (
	_ "embed"
)

type HelloParam struct {
	Input string
	Text  string
}

var _ Param = (*HelloParam)(nil)

//go:embed prompts/HelloPrompt.md
var helloPrompt string

func (HelloParam) Prompt() string { return helloPrompt }

func (this *HelloParam) Prepare() error {
	this.Text = "|" + this.Input + "|"
	return nil
}

func (this *HelloParam) Do() (string, error) {
	return Do(this)
}

func Hello(input string) (string, error) {
	return (&HelloParam{Input: input}).Do()
}
