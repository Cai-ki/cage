package helper_test

import (
	"os"
	"path"
	"testing"

	_ "github.com/Cai-ki/cage/config"
	"github.com/Cai-ki/cage/helper"
)

func TestGenerateProjectDoc(t *testing.T) {
	p := os.Getenv("ROOT_DIR")

	if p == "" {
		t.Fatal(p)
	}

	err := helper.GenerateProjectDoc(p, path.Join(p, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
}
