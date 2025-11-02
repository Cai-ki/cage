//go:build linux
// +build linux

package record

import "io"

type linuxRecorder struct{}

func (l *linuxRecorder) Record(durationSeconds int) (io.ReadCloser, error) {
	return nil, ErrNotImplemented
}

func init() {
	defaultRecorder = &linuxRecorder{}
}
