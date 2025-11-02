package record

import "errors"

// ErrNotImplemented is returned when audio recording is not supported on the current platform.
var ErrNotImplemented = errors.New("audio recording not implemented on this platform")

// ErrSoxNotInstalled is returned when 'sox' is required but not found in PATH.
var ErrSoxNotInstalled = errors.New("audio recording requires 'sox' — install via 'brew install sox'")
