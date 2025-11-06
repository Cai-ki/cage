// package helper

package helper

import (
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type AnalyzePackageParam struct {
	Dir  string
	Code string
}

var _ Param = (*AnalyzePackageParam)(nil)

//go:embed prompts/AnalyzePackagePrompt.md
var analyzePackagePrompt string

func (AnalyzePackageParam) Prompt() string { return analyzePackagePrompt }

func (this *AnalyzePackageParam) Prepare() error {
	files, err := readGoFilesRecursively(this.Dir)
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

func (this *AnalyzePackageParam) Do() (string, error) {
	return Do(this)
}

func AnalyzePackage(dir string) (string, error) {
	return (&AnalyzePackageParam{Dir: dir}).Do()
}

// GeneratePackageDoc 调用源码分析并写入文档（保持原有功能）
func GeneratePackageDoc(dir, outputPath string) error {
	param := AnalyzePackageParam{Dir: dir}
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

	if err := os.WriteFile(outputPath, []byte(doc), 0644); err != nil {
		return fmt.Errorf("写入文档失败: %w", err)
	}

	return nil
}

// readGoFilesRecursively 递归读取目录下所有 .go 文件（排除 *_test.go）
func readGoFilesRecursively(dir string) (map[string]string, error) {
	files := make(map[string]string)

	err := filepath.WalkDir(dir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil // 继续遍历子目录
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		if strings.HasSuffix(path, "_test.go") {
			return nil // 跳过测试文件
		}

		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %w", path, err)
		}
		// 使用相对于 dir 的路径作为 key，更清晰
		relPath, err := filepath.Rel(dir, path)
		if err != nil {
			relPath = path // fallback
		}
		files[relPath] = string(content)
		return nil
	})

	if err != nil {
		return nil, err
	}
	return files, nil
}

func packageNameFromDir(dir string) string {
	base := filepath.Base(filepath.Clean(dir))
	if base == "." || base == "/" {
		return "root"
	}
	return base
}
