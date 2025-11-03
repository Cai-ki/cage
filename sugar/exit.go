package sugar

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

func ExitIfErr(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			os.Stderr.WriteString(msg[0])
		}
		os.Exit(0)
	}
}

func ExitIfErrf(err error, format string, args ...any) {
	if err != nil {
		if format != "" {
			fmt.Fprintf(os.Stderr, format, args...)
		}
		os.Exit(0)
	}
}
