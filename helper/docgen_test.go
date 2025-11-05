package helper_test

import (
	"path"
	"testing"

	"github.com/Cai-ki/cage/helper"
)

func TestGeneratePackageDoc(t *testing.T) {
	pkgPath := "."
	err := helper.GeneratePackageDoc(pkgPath, path.Join(pkgPath, "doc.md"))
	if err != nil {
		t.Fatal(err)
	}
}
