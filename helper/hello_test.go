package helper_test

import (
	"testing"

	"github.com/Cai-ki/cage/helper"
)

func TestHello(t *testing.T) {
	param := helper.HelloParam{Input: "1 + 1 = ?"}
	txt, err := param.Do()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(param.Prompt())
	param.Prepare()
	t.Log(param)

	t.Log(txt)
	t.Log(helper.Hello("1 + 1 = ?"))
}
