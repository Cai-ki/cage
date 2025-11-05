package helper

import (
	"strings"

	"github.com/Cai-ki/cage/llm"
)

// DescribeToShellScript 根据自然语言描述生成安全的 shell 脚本
func DescribeToShellScript(description string) (string, error) {
	prompt := `
你是一个资深 DevOps 工程师，精通 Bash Shell 脚本编写。请根据用户的自然语言描述，生成一个**可直接执行、安全、带注释**的 shell 脚本。

要求：
1. 脚本必须以 #!/bin/bash 开头；
2. 必须包含 set -euo pipefail（严格模式）；
3. 所有变量引用必须用双引号（如 "$VAR"）；
4. 避免使用 cd 而不检查是否成功，应写成：cd /path || exit 1；
5. 如果涉及文件操作，先检查文件是否存在；
6. 如果涉及网络请求（如 curl），需处理失败情况；
7. 如果需要用户输入，使用 read -r 并加注释说明；
8. 脚本末尾应有明确的 exit 0；
9. **不要**使用 root 权限、不要删除系统关键目录（如 /）、不要写破坏性操作，除非用户明确要求且你加了确认提示；
10. 仅输出脚本内容，不要包含任何解释、标记（如 “bash”）或额外文本；
11. 用中文在关键步骤添加注释（以 # 开头）；
12. 如果描述模糊，优先生成通用、安全的版本，并在注释中说明假设。

现在，请根据以下描述生成脚本：

` + strings.TrimSpace(description)

	return llm.Completion(prompt)
}

// DescribeToRunnableShell 生成可直接粘贴到终端运行的 shell 命令片段（非脚本文件）
func DescribeToRunnableShell(description string) (string, error) {
	prompt := `
你是一位精通 Bash/Zsh 的开发者。请根据用户描述，生成一段**可直接粘贴到终端中运行**的 shell 命令。

要求：
1. 输出必须是**纯 shell 命令**，不包含 #!/bin/bash、不包含 exit、不包含文件头；
2. 如果是多行，请用换行分隔，确保粘贴后能逐行执行（或用 { ...; } 包裹）；
3. 所有操作必须安全：不删除系统文件，不修改 PATH，不后台运行危险进程；
4. 涉及文件路径时，优先使用相对路径或明确提示用户替换（如 ~/your_dir）；
5. 如果需要循环或条件，用单行形式，例如：for f in *.log; do echo "$f"; done
6. 如果需要临时变量，用小括号创建子 shell 避免污染环境，例如：(name="test"; echo "$name")
7. 涉及网络请求（如 curl）必须带错误处理（如 curl -f ... || echo "failed"）；
8. **不要**输出任何 Markdown、代码块符号、解释文字；
9. 用中文在行尾用 # 注释说明关键操作（但不要影响命令执行）；
10. 如果描述不明确，优先生成只读/查询类命令，避免写操作。

现在，请根据以下描述生成可粘贴运行的 shell 命令：

` + strings.TrimSpace(description)

	return llm.Completion(prompt)
}
