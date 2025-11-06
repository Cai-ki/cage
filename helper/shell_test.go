package helper_test

import (
	"fmt"
	"testing"

	"github.com/Cai-ki/cage/helper"
)

func TestDescribeToShellScript(t *testing.T) {
	// param := helper.DescribeToShellScriptParam{Input: "创建一个备份目录的脚本"}
	// txt, err := param.Do()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(param.Prompt())
	// t.Log(txt)

	t.Log(helper.DescribeToShellScript("创建一个备份目录的脚本"))
}

func TestDescribeToRunnableShell(t *testing.T) {
	desc := "列出当前目录下以及其子目录下所有 .go 文件，并统计总行数"
	cmd, _ := helper.DescribeToRunnableShell(desc)
	fmt.Println("=== 统计 Go 行数 ===")
	fmt.Println(cmd)
	return

	// 	desc1 := "列出当前目录下所有 .go 文件，并统计总行数"
	// 	cmd1, _ := helper.DescribeToRunnableShell(desc1)
	// 	fmt.Println("=== 统计 Go 行数 ===")
	// 	fmt.Println(cmd1)

	// 	desc2 := "检查 8080 端口是否被占用，如果被占用就显示占用进程"
	// 	cmd2, _ := helper.DescribeToRunnableShell(desc2)
	// 	fmt.Println("\n=== 检查端口占用 ===")
	// 	fmt.Println(cmd2)

	// 	desc3 := "从 https://httpbin.org/json 下载 JSON，提取其中的 “slideshow” 字段并格式化输出"
	// 	cmd3, _ := helper.DescribeToRunnableShell(desc3)
	// 	fmt.Println("\n=== 下载并解析 JSON ===")
	// 	fmt.Println(cmd3)

	// 	desc4 := "批量给当前目录所有 .sh 文件加可执行权限"
	// 	cmd4, _ := helper.DescribeToRunnableShell(desc4)
	// 	fmt.Println("\n=== 批量加权限（安全版） ===")
	// 	fmt.Println(cmd4)

	// 	desc5 := "显示最近 5 个 CPU 占用最高的进程"
	// 	cmd5, _ := helper.DescribeToRunnableShell(desc5)
	// 	fmt.Println("\n=== 高 CPU 进程 ===")
	// 	fmt.Println(cmd5)

	// param := helper.DescribeToRunnableShellParam{Input: "查找当前目录下所有.log文件并统计行数"}
	// txt, err := param.Do()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(param.Prompt())
	// t.Log(helper.Parse(param))
	// t.Log(txt)

	// t.Log(helper.DescribeToRunnableShell("查找当前目录下所有.log文件并统计行数"))
}
