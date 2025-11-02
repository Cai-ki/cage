package capture_test

import (
	"image/png"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Cai-ki/cage/capture"
)

func TestScreen(t *testing.T) {
	img, err := capture.Screen()
	if err != nil {
		t.Fatalf("Screen capture failed: %v", err)
	}
	// 确保 testdata 目录存在（项目根目录下）
	testdataDir := filepath.Join("..", "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatalf("Failed to create testdata dir: %v", err)
	}

	// 生成带时间戳的文件名，避免覆盖
	now := time.Now().Format("20060102_150405")
	outputPath := filepath.Join(testdataDir, "screenshot_"+now+".png")

	f, err := os.Create(outputPath)
	if err != nil {
		t.Fatalf("Failed to create output file: %v", err)
	}
	defer f.Close()

	if err := png.Encode(f, img); err != nil {
		t.Fatalf("Failed to encode PNG: %v", err)
	}

	t.Logf("Screenshot saved to: %s", outputPath)
}
