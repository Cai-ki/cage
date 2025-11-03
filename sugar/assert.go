package sugar

import "fmt"

func Assert(cond bool, msg string) {
	if !cond {
		panic(msg)
	}
}

func Assertf(cond bool, format string, args ...any) {
	if !cond {
		panic(fmt.Sprintf(format, args...))
	}
}
