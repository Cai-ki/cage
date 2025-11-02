package media

import "errors"

var (
	// ErrNotImplemented is returned when a feature is not supported on the current platform.
	ErrNotImplemented = errors.New("media feature not implemented on this platform")

	// ErrSoxNotInstalled is returned when 'sox' is required but not found.
	ErrSoxNotInstalled = errors.New("audio recording requires 'sox' — install via 'brew install sox'")
)
