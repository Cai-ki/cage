//go:build windows
// +build windows

package record

import "io"

type windowsRecorder struct{}

func (w *windowsRecorder) Record(durationSeconds int) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}

func init() {
	defaultRecorder = &windowsRecorder{}
}
