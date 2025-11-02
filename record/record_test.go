package record_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Cai-ki/cage/record"
)

func TestRecord(t *testing.T) {
	// Skip if sox not installed (e.g., on CI or Linux)
	if _, err := record.Start(0); err != nil {
		if err == record.ErrNotImplemented || err == record.ErrSoxNotInstalled {
			t.Skipf("Recording not available: %v", err)
		}
		t.Fatalf("Unexpected error: %v", err)
	}

	audio, err := record.Start(3) // record 3 seconds
	if err != nil {
		t.Fatalf("Recording failed: %v", err)
	}
	defer audio.Close()

	// Save to testdata/
	testdataDir := filepath.Join("..", "testdata")
	if err := os.MkdirAll(testdataDir, 0755); err != nil {
		t.Fatal(err)
	}

	now := time.Now().Format("20060102_150405")
	outPath := filepath.Join(testdataDir, "recording_"+now+".wav")
	out, err := os.Create(outPath)
	if err != nil {
		t.Fatal(err)
	}
	defer out.Close()

	if _, err := io.Copy(out, audio); err != nil {
		t.Fatal(err)
	}

	t.Logf("Audio saved to: %s", outPath)
}
