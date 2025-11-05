package helper

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/Cai-ki/cage/llm"
)

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

// AnalyzePackageBySourceCode 读取整个目录的 Go 源码并交由大模型总结可导出功能
func AnalyzePackageBySourceCode(dir string) (string, error) {
	files, err := readGoFilesRecursively(dir)
	if err != nil {
		return "", fmt.Errorf("读取源码文件失败: %w", err)
	}

	if len(files) == 0 {
		return "", fmt.Errorf("目录中未找到任何 .go 源文件（非测试）")
	}

	var sb strings.Builder
	sb.WriteString("以下是该 Go 包的全部源代码（已排除 _test.go 文件），每个文件以「=== FILE: 文件名 ===」开头：\n\n")

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

	sourceCode := sb.String()

	prompt := "你是一位精通 Go 语言的系统架构师，请根据以下完整的 Go 包源代码（包含所有非测试 .go 文件），生成一份**结构清晰、格式规范的 Markdown 文档**，用于说明该包的公开接口与功能。\n\n" +
		"输出要求如下：\n\n" +
		"1. **文档必须是标准 Markdown 格式**，可直接保存为 .md 文件并在 GitHub、VS Code 等环境中正确渲染。\n\n" +
		"2. **整体结构**：\n" +
		"   - 第一行：`# 包功能说明`\n" +
		"   - 接着一段 100–200 字的中文概览，说明包的核心用途、设计目标和典型使用场景。\n" +
		"   - 然后按类别分节：`## 结构体与接口`、`## 函数`、`## 变量与常量`\n\n" +
		"3. **每一项的格式**：\n" +
		"   - 先用 **Go 代码块**（带语言标识）展示其完整定义或签名，例如：\n" +
		"     ```go\n" +
		"     func NewClient(apiKey string) *Client\n" +
		"     ```\n" +
		"   - 紧接着用一段**中文段落**解释其功能、参数、返回值、使用注意事项等。\n" +
		"   - 每项之间空一行。\n\n" +
		"4. **覆盖范围**（仅限可被外部包访问的项）：\n" +
		"   - 所有首字母大写的：函数、结构体、接口、变量、常量\n" +
		"   - 结构体的公开方法（如 `(c *Client) DoSomething()`）\n" +
		"   - 如果某个结构体有公开字段（如 `type Config struct { Timeout int }`），请在结构体代码块下方说明各字段含义。\n\n" +
		"5. **禁止行为**：\n" +
		"   - 不得编造代码中不存在的函数、字段或行为。\n" +
		"   - 不得使用非 Markdown 的格式（如纯文本星号列表、无语言标识的代码块）。\n" +
		"   - 不要包含“根据代码”“如上所示”等冗余引导语。\n\n" +
		"6. **语言**：全文使用简体中文，技术术语准确。\n\n" +
		"以下是该包的完整源码（每个文件以「=== FILE: 文件名 ===」开头）：\n" +
		sourceCode

	result, err := llm.Completion(prompt)
	if err != nil {
		return "", fmt.Errorf("调用大模型失败: %w", err)
	}

	// 可选：简单后处理，去除可能的多余引号或格式
	result = strings.Trim(result, " \n\"`")
	return result, nil
}

// GeneratePackageDoc 调用源码分析并写入文档
func GeneratePackageDoc(dir, outputPath string) error {
	doc, err := AnalyzePackageBySourceCode(dir)
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

func packageNameFromDir(dir string) string {
	base := filepath.Base(filepath.Clean(dir))
	if base == "." || base == "/" {
		return "root"
	}
	return base
}
