//go:build windows
// +build windows

package capture

import "image"

type windowsCapturer struct{}

func (w *windowsCapturer) CaptureScreen() (image.Image, error) {
	// TODO: 实现 Windows 截图（如 PowerShell + Add-Type）
	panic("Windows screen capture not implemented yet")
}

func init() {
	defaultCapturer = &windowsCapturer{}
}
