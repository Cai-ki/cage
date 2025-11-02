// Package record provides cross-platform audio recording from the microphone.
package record

import (
	"io"
	"sync"
)

// Recorder is the interface for audio recording implementations.
type Recorder interface {
	Record(durationSeconds int) (io.ReadCloser, error)
}

var (
	defaultRecorder Recorder
	once            sync.Once
)

// Start records audio from the default microphone for the given duration (in seconds).
// It returns a ReadCloser that yields WAV audio data (16kHz, 16-bit, mono, PCM).
func Start(durationSeconds int) (io.ReadCloser, error) {
	once.Do(func() {
		// defaultRecorder is set by platform-specific init()
	})
	if defaultRecorder == nil {
		return nil, ErrNotImplemented
	}
	return defaultRecorder.Record(durationSeconds)
}
