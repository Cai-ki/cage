package helper

import (
	_ "embed"
)

// DescribeToShellScriptParam 生成Shell脚本的参数
type DescribeToShellScriptParam struct {
	Input string
	Text  string
}

var _ Param = (*DescribeToShellScriptParam)(nil)

//go:embed prompts/DescribeToShellScriptPrompt.md
var DescribeToShellScriptPrompt string

func (DescribeToShellScriptParam) Prompt() string { return DescribeToShellScriptPrompt }

func (this *DescribeToShellScriptParam) Prepare() error {
	this.Text = this.Input
	return nil
}

func (this *DescribeToShellScriptParam) Do() (string, error) {
	return Do(this)
}

// DescribeToShellScript 根据自然语言描述生成安全的 shell 脚本
func DescribeToShellScript(description string) (string, error) {
	return (&DescribeToShellScriptParam{Input: description}).Do()
}

// DescribeToRunnableShellParam 生成可运行Shell命令的参数
type DescribeToRunnableShellParam struct {
	Input string
	Text  string
}

var _ Param = (*DescribeToRunnableShellParam)(nil)

//go:embed prompts/DescribeToRunnableShellPrompt.md
var DescribeToRunnableShellPrompt string

func (DescribeToRunnableShellParam) Prompt() string { return DescribeToRunnableShellPrompt }

func (this *DescribeToRunnableShellParam) Prepare() error {
	this.Text = this.Input
	return nil
}

func (this *DescribeToRunnableShellParam) Do() (string, error) {
	return Do(this)
}

// DescribeToRunnableShell 生成可直接粘贴到终端运行的 shell 命令片段
func DescribeToRunnableShell(description string) (string, error) {
	return (&DescribeToRunnableShellParam{Input: description}).Do()
}
