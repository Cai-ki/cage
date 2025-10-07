package cage

import (
	"fmt"
	"os"
)

func ExitIfNot(cond bool, msg ...string) {
	if !cond {
		if len(msg) > 0 {
			os.Stderr.WriteString(msg[0])
		}
		os.Exit(0)
	}
}

func ExitIfNotf(cond bool, format string, args ...any) {
	if !cond {
		if format != "" {
			fmt.Fprintf(os.Stderr, format, args...)
		}
		os.Exit(0)
	}
}
