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

func (m *macOSScreenshotCapturer) CaptureScreen() (image.Image, error) {
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
