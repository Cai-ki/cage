//go:build linux
// +build linux

package capture

import "image"

type linuxCapturer struct{}

func (l *linuxCapturer) CaptureScreen() (image.Image, error) {
	// TODO: 实现 Linux 截图（如调用 maim, gnome-screenshot, 或 xwd）
	panic("Linux screen capture not implemented yet")
}

func init() {
	defaultCapturer = &linuxCapturer{}
}
