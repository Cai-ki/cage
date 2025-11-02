// Package media provides cross-platform media capture: screenshots and audio recording.
package media

import (
	"image"
	"io"
	"sync"
)

var (
	screenshotCapturer ScreenshotCapturer
	audioRecorder      AudioRecorder
	once               sync.Once
)

// Screenshot captures the entire screen and returns an image.
func Screenshot() (image.Image, error) {
	once.Do(initPlatform)
	if screenshotCapturer == nil {
		return nil, ErrNotImplemented
	}
	return screenshotCapturer.CaptureScreen()
}

// RecordAudio records microphone audio for the given duration (seconds).
// Returns a WAV audio stream (16kHz, 16-bit, mono).
func RecordAudio(durationSeconds int) (io.ReadCloser, error) {
	once.Do(initPlatform)
	if audioRecorder == nil {
		return nil, ErrNotImplemented
	}
	return audioRecorder.Record(durationSeconds)
}
