package sugar_test

import (
	"testing"

	"github.com/Cai-ki/cage/sugar"
)

func TestStrToT(t *testing.T) {
	// === 成功转换测试 ===
	t.Run("string", func(t *testing.T) {
		got, err := sugar.StrToT[string]("hello")
		if err != nil {
			t.Error(err)
		}
		want := "hello"
		if got != want {
			t.Errorf("sugar.StrToT[string] = %q, want %q", got, want)
		}
	})

	t.Run("int", func(t *testing.T) {
		got, err := sugar.StrToT[int]("42")
		if err != nil {
			t.Error(err)
		}
		want := 42
		if got != want {
			t.Errorf("sugar.StrToT[int] = %d, want %d", got, want)
		}
	})

	t.Run("int64", func(t *testing.T) {
		got, err := sugar.StrToT[int64]("9223372036854775807")
		if err != nil {
			t.Error(err)
		}
		want := int64(9223372036854775807)
		if got != want {
			t.Errorf("sugar.StrToT[int64] = %d, want %d", got, want)
		}
	})

	t.Run("float64", func(t *testing.T) {
		got, err := sugar.StrToT[float64]("3.14159")
		if err != nil {
			t.Error(err)
		}
		want := 3.14159
		if got != want {
			t.Errorf("sugar.StrToT[float64] = %f, want %f", got, want)
		}
	})

	t.Run("bool_true", func(t *testing.T) {
		got, err := sugar.StrToT[bool]("true")
		if err != nil {
			t.Error(err)
		}
		want := true
		if got != want {
			t.Errorf("sugar.StrToT[bool] = %t, want %t", got, want)
		}
	})

	t.Run("bool_false", func(t *testing.T) {
		got, err := sugar.StrToT[bool]("false")
		if err != nil {
			t.Error(err)
		}
		want := false
		if got != want {
			t.Errorf("sugar.StrToT[bool] = %t, want %t", got, want)
		}
	})

	t.Run("uint", func(t *testing.T) {
		got, err := sugar.StrToT[uint]("100")
		if err != nil {
			t.Error(err)
		}
		want := uint(100)
		if got != want {
			t.Errorf("sugar.StrToT[uint] = %d, want %d", got, want)
		}
	})

	t.Run("invalid_int", func(t *testing.T) {
		_, err := sugar.StrToT[int]("not_a_number")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid_float", func(t *testing.T) {
		_, err := sugar.StrToT[float64]("not.a.number")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("invalid_bool", func(t *testing.T) {
		_, err := sugar.StrToT[bool]("maybe")
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("unsupported_type", func(t *testing.T) {
		_, err := sugar.StrToT[[]byte]("hello")
		if err != nil {
			t.Error(err)
		}
	})
}

func TestStrToTWithDefault(t *testing.T) {
	// === 成功转换 ===
	t.Run("int_success", func(t *testing.T) {
		got := sugar.StrToTWithDefault("42", 0)
		want := 42
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("float64_success", func(t *testing.T) {
		got := sugar.StrToTWithDefault("3.14", 0.0)
		want := 3.14
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("bool_true", func(t *testing.T) {
		got := sugar.StrToTWithDefault("true", false)
		want := true
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	t.Run("string", func(t *testing.T) {
		got := sugar.StrToTWithDefault("hello", "default")
		want := "hello"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	// === 转换失败，返回默认值 ===
	t.Run("int_invalid", func(t *testing.T) {
		got := sugar.StrToTWithDefault("not_a_number", -1)
		want := -1
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("float_invalid", func(t *testing.T) {
		got := sugar.StrToTWithDefault("xyz", 999.0)
		want := 999.0
		if got != want {
			t.Errorf("got %f, want %f", got, want)
		}
	})

	t.Run("bool_invalid", func(t *testing.T) {
		got := sugar.StrToTWithDefault("maybe", true)
		want := true
		if got != want {
			t.Errorf("got %t, want %t", got, want)
		}
	})

	// === 不支持的类型（如 struct）返回默认值 ===
	type MyStruct struct{ Name string }
	t.Run("unsupported_type", func(t *testing.T) {
		def := MyStruct{Name: "default"}
		got := sugar.StrToTWithDefault("anything", def)
		if got != def {
			t.Error("expected default value for unsupported type")
		}
	})

	// === 边界值测试 ===
	t.Run("int64_max", func(t *testing.T) {
		got := sugar.StrToTWithDefault("9223372036854775807", int64(0))
		want := int64(9223372036854775807)
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})

	t.Run("uint64_large", func(t *testing.T) {
		got := sugar.StrToTWithDefault("18446744073709551615", uint64(0))
		want := uint64(18446744073709551615)
		if got != want {
			t.Errorf("got %d, want %d", got, want)
		}
	})
}
