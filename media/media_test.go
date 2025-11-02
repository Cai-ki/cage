package media_test

import (
	"image"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Cai-ki/cage/media"
)

func saveTestFile(t *testing.T, name string, data interface{}) string {
	testdataDir := filepath.Join("..", "testdata")
	os.MkdirAll(testdataDir, 0755)

	now := time.Now().Format("20060102_150405")
	path := filepath.Join(testdataDir, name+"_"+now)

	switch v := data.(type) {
	case image.Image:
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
	img, err := media.Screenshot()
	if err != nil {
		if err == media.ErrNotImplemented {
			t.Skip("Screenshot not available on this platform")
		}
		t.Fatal(err)
	}
	saveTestFile(t, "screenshot", img)
}

func TestRecordAudio(t *testing.T) {
	audio, err := media.RecordAudio(3)
	if err != nil {
		if err == media.ErrNotImplemented || err == media.ErrSoxNotInstalled {
			t.Skipf("Audio recording not available: %v", err)
		}
		t.Fatal(err)
	}
	defer audio.Close()
	saveTestFile(t, "recording", audio)
}
