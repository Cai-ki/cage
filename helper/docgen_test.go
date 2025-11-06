package helper_test

import (
	"path"
	"testing"

	"github.com/Cai-ki/cage/helper"
)

func TestAnalyzePackage(t *testing.T) {
	// param := helper.AnalyzePackageParam{Dir: "/Users/caiki/Code/project/cage/config"}
	// txt, err := param.Do()
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// t.Log(param.Prompt())
	// t.Log(param)

	// t.Log(txt)
	t.Log(helper.AnalyzePackage("/Users/caiki/Code/project/cage/config"))
}

func TestGeneratePackageDoc(t *testing.T) {
	p := "."
	err := helper.GeneratePackageDoc(p, path.Join(p, "doc.md"))
	if err != nil {
		t.Fatal(err)
	}
}
