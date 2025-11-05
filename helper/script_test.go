package helper_test

import (
	"fmt"
	"testing"

	"github.com/Cai-ki/cage/helper"
)

func TestDescribeToShellScript(t *testing.T) {
	// 场景1：日常开发常用
	desc1 := "每天凌晨2点备份 ~/projects 目录到 /backup，保留7天，用 tar.gz 压缩，并记录日志"
	script1, _ := helper.DescribeToShellScript(desc1)
	fmt.Println("=== 场景1：自动备份脚本 ===")
	fmt.Println(script1)

	// 场景2：跨平台兼容（macOS / Linux）
	desc2 := "检查系统是 macOS 还是 Linux，然后安装 jq：mac 用 brew，linux 用 apt"
	script2, _ := helper.DescribeToShellScript(desc2)
	fmt.Println("\n=== 场景2：跨平台安装 jq ===")
	fmt.Println(script2)

	// 场景3：API 调用 + 错误处理
	desc3 := "调用 https://api.example.com/health，如果返回 200 则输出 'OK'，否则重试3次，最后失败则告警"
	script3, _ := helper.DescribeToShellScript(desc3)
	fmt.Println("\n=== 场景3：带重试的健康检查 ===")
	fmt.Println(script3)

	// 场景4：你的实际痛点（权限修复）
	desc4 := "给当前目录下所有 .sh 文件添加可执行权限，并确保它们有 #!/bin/bash 头"
	script4, _ := helper.DescribeToShellScript(desc4)
	fmt.Println("\n=== 场景4：修复 shell 脚本权限 ===")
	fmt.Println(script4)
}

func TestDescribeToRunnableShell(t *testing.T) {
	desc1 := "列出当前目录下所有 .go 文件，并统计总行数"
	cmd1, _ := helper.DescribeToRunnableShell(desc1)
	fmt.Println("=== 统计 Go 行数 ===")
	fmt.Println(cmd1)

	desc2 := "检查 8080 端口是否被占用，如果被占用就显示占用进程"
	cmd2, _ := helper.DescribeToRunnableShell(desc2)
	fmt.Println("\n=== 检查端口占用 ===")
	fmt.Println(cmd2)

	desc3 := "从 https://httpbin.org/json 下载 JSON，提取其中的 “slideshow” 字段并格式化输出"
	cmd3, _ := helper.DescribeToRunnableShell(desc3)
	fmt.Println("\n=== 下载并解析 JSON ===")
	fmt.Println(cmd3)

	desc4 := "批量给当前目录所有 .sh 文件加可执行权限"
	cmd4, _ := helper.DescribeToRunnableShell(desc4)
	fmt.Println("\n=== 批量加权限（安全版） ===")
	fmt.Println(cmd4)

	desc5 := "显示最近 5 个 CPU 占用最高的进程"
	cmd5, _ := helper.DescribeToRunnableShell(desc5)
	fmt.Println("\n=== 高 CPU 进程 ===")
	fmt.Println(cmd5)
}
