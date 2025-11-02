#!/bin/bash

# 创建 media 目录
mkdir -p media

# media/media.go
cat > media/media.go << 'EOF'
// Package media provides cross-platform media capture: screenshots and audio recording.
package media

import (
	"io"
	"sync"
)

var (
	screenshotCapturer ScreenshotCapturer
	audioRecorder      AudioRecorder
	once               sync.Once
)

// Screenshot captures the entire screen and returns an image.
func Screenshot() (Image, error) {
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
EOF

# media/screenshot.go
cat > media/screenshot.go << 'EOF'
package media

import "image"

// Image is an alias for image.Image for cleaner API.
type Image = image.Image

// ScreenshotCapturer captures static screen images.
type ScreenshotCapturer interface {
	CaptureScreen() (Image, error)
}
EOF

# media/audio.go
cat > media/audio.go << 'EOF'
package media

import "io"

// AudioRecorder captures audio streams.
type AudioRecorder interface {
	Record(durationSeconds int) (io.ReadCloser, error)
}
EOF

# media/errors.go
cat > media/errors.go << 'EOF'
package media

import "errors"

var (
	// ErrNotImplemented is returned when a feature is not supported on the current platform.
	ErrNotImplemented = errors.New("media feature not implemented on this platform")

	// ErrSoxNotInstalled is returned when 'sox' is required but not found.
	ErrSoxNotInstalled = errors.New("audio recording requires 'sox' — install via 'brew install sox'")
)
EOF

# media/macos_screenshot.go
cat > media/macos_screenshot.go << 'EOF'
//go:build darwin
// +build darwin

package media

import (
	"image"
	"image/png"
	"os"
	"os/exec"
)

type macOSScreenshotCapturer struct{}

func (m *macOSScreenshotCapturer) CaptureScreen() (Image, error) {
	f, err := os.CreateTemp("", "screenshot_*.png")
	if err != nil {
		return nil, err
	}
	tmpPath := f.Name()
	f.Close()
	defer os.Remove(tmpPath)

	cmd := exec.Command("screencapture", "-x", "-t", "png", tmpPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	f, err = os.Open(tmpPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}
EOF

# media/macos_audio.go
cat > media/macos_audio.go << 'EOF'
//go:build darwin
// +build darwin

package media

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type macOSAudioRecorder struct{}

func (r *macOSAudioRecorder) Record(durationSeconds int) (io.ReadCloser, error) {
	if _, err := exec.LookPath("sox"); err != nil {
		return nil, ErrSoxNotInstalled
	}

	tmpFile, err := os.CreateTemp("", "recording_*.wav")
	if err != nil {
		return nil, err
	}
	tmpPath := tmpFile.Name()
	tmpFile.Close()
	defer os.Remove(tmpPath)

	cmd := exec.Command(
		"sox",
		"-d",
		"-r", "16000",
		"-b", "16",
		"-c", "1",
		tmpPath,
		"trim", "0", fmt.Sprintf("%d", durationSeconds),
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("sox failed: %w", err)
	}

	return os.Open(tmpPath)
}
EOF

# media/macos.go
cat > media/macos.go << 'EOF'
//go:build darwin
// +build darwin

package media

func initPlatform() {
	screenshotCapturer = &macOSScreenshotCapturer{}
	audioRecorder = &macOSAudioRecorder{}
}
EOF

# media/linux.go
cat > media/linux.go << 'EOF'
//go:build linux
// +build linux

package media

func initPlatform() {
	// Leave capturers nil → ErrNotImplemented
}
EOF

# media/windows.go
cat > media/windows.go << 'EOF'
//go:build windows
// +build windows

package media

func initPlatform() {
	// Leave capturers nil → ErrNotImplemented
}
EOF

# media/media_test.go
cat > media/media_test.go << 'EOF'
package media

import (
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func saveTestFile(t *testing.T, name string, data interface{}) string {
	testdataDir := filepath.Join("..", "testdata")
	os.MkdirAll(testdataDir, 0755)

	now := time.Now().Format("20060102_150405")
	path := filepath.Join(testdataDir, name+"_"+now)

	switch v := data.(type) {
	case Image:
		f, _ := os.Create(path + ".png")
		defer f.Close()
		png.Encode(f, v)
	case io.Reader:
		f, _ := os.Create(path + ".wav")
		defer f.Close()
		io.Copy(f, v)
	}
	t.Logf("Saved to: %s", path)
	return path
}

func TestScreenshot(t *testing.T) {
	img, err := Screenshot()
	if err != nil {
		if err == ErrNotImplemented {
			t.Skip("Screenshot not available on this platform")
		}
		t.Fatal(err)
	}
	saveTestFile(t, "screenshot", img)
}

func TestRecordAudio(t *testing.T) {
	audio, err := RecordAudio(3)
	if err != nil {
		if err == ErrNotImplemented || err == ErrSoxNotInstalled {
			t.Skipf("Audio recording not available: %v", err)
		}
		t.Fatal(err)
	}
	defer audio.Close()
	saveTestFile(t, "recording", audio)
}
EOF

echo "✅ Done! Created media/ directory with all Go files for macOS screenshot and audio recording."
echo "💡 Remember to run: brew install sox"
echo "🧪 Test with: cd media && go test -v"