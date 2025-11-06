package config

import (
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

func init() {
	LoadEnvFromRoot()
}

// LoadEnvFromRoot 尝试从项目根目录加载 .env 文件
// 项目根目录定义为包含 go.mod 的最近祖先目录
func LoadEnvFromRoot() error {
	// 从当前工作目录向上查找 go.mod
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	for {
		// 检查当前目录是否有 go.mod
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			// 找到了项目根目录，加载 .env
			return godotenv.Load(filepath.Join(dir, ".env"))
		}

		// 向上一层
		parent := filepath.Dir(dir)
		if parent == dir {
			// 到达文件系统根目录，未找到 go.mod
			return os.ErrNotExist
		}
		dir = parent
	}
}
