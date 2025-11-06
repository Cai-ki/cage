package helper

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AnalyzeProjectParam struct {
	Dir  string
	Code string
}

var _ Param = (*AnalyzeProjectParam)(nil)

//go:embed prompts/AnalyzeProjectPrompt.md
var analyzeProjectPrompt string

func (AnalyzeProjectParam) Prompt() string { return analyzeProjectPrompt }

func (this *AnalyzeProjectParam) Prepare() error {
	files, err := findAllDocMD(this.Dir)
	if err != nil {
		return fmt.Errorf("读取源码文件失败: %w", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("目录中未找到任何 .go 源文件（非测试）")
	}

	var sb strings.Builder

	// 按文件名排序，保证顺序稳定
	fileNames := make([]string, 0, len(files))
	for name := range files {
		fileNames = append(fileNames, name)
	}
	sort.Strings(fileNames)

	for _, name := range fileNames {
		content := files[name]
		sb.WriteString(fmt.Sprintf("=== FILE: %s ===\n", name))
		sb.WriteString(content)
		sb.WriteString("\n\n")
	}

	this.Code = sb.String()
	return nil
}

func (this *AnalyzeProjectParam) Do() (string, error) {
	return Do(this)
}

func AnalyzeProject(dir string) (string, error) {
	return (&AnalyzeProjectParam{Dir: dir}).Do()
}

// GenerateProjectDoc 调用doc.md分析并写入README
func GenerateProjectDoc(dir, outputPath string) error {
	param := AnalyzeProjectParam{Dir: dir}
	doc, err := param.Do()
	if err != nil {
		return err
	}

	outputDir := filepath.Dir(outputPath)
	if outputDir != "." {
		if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
			return fmt.Errorf("创建输出目录失败: %w", err)
		}
	}

	doc = "```\n以下内容由 AI 脚本总结\n```\n" + doc

	if err := os.WriteFile(outputPath, []byte(doc), 0644); err != nil {
		return fmt.Errorf("写入文档失败: %w", err)
	}

	return nil
}

// findAllDocMD 从环境变量 ROOT_DIR 指定的根目录开始，
// 递归查找所有名为 doc.md 的文件，并返回文件路径到内容的映射。
func findAllDocMD(rootDir string) (map[string]string, error) {
	if rootDir == "" {
		return nil, fmt.Errorf("dir is not set")
	}

	result := make(map[string]string)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err // 跳过无法访问的路径
		}
		if !info.IsDir() && info.Name() == "doc.md" {
			content, readErr := os.ReadFile(path)
			if readErr != nil {
				return readErr // 或者可以选择记录错误并继续
			}
			relPath, err := filepath.Rel(rootDir, path)
			if err != nil {
				relPath = path // fallback
			}
			result[relPath] = string(content)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("error walking the path %s: %w", rootDir, err)
	}

	return result, nil
}
