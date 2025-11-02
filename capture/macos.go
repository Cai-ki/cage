//go:build darwin
// +build darwin

package capture

import (
	"image"
	"image/png"
	"os"
	"os/exec"
)

type macOSCapturer struct{}

func (m *macOSCapturer) CaptureScreen() (image.Image, error) {
	// 创建临时文件（自动命名）
	f, err := os.CreateTemp("", "screenshot_*.png")
	if err != nil {
		return nil, err
	}
	tmpPath := f.Name()
	f.Close() // screencapture 需要写入该文件，所以先关闭
	defer os.Remove(tmpPath)

	// 调用 macOS screencapture
	cmd := exec.Command("screencapture", "-x", "-t", "png", tmpPath)
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	// 读取截图
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

func init() {
	defaultCapturer = &macOSCapturer{}
}
