//go:build darwin
// +build darwin

package record

import (
	"fmt"
	"io"
	"os"
	"os/exec"
)

type macOSRecorder struct{}

func (r *macOSRecorder) Record(durationSeconds int) (io.ReadCloser, error) {
	// Check if sox is available
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

	// Record from default mic: 16kHz, 16-bit, mono, WAV
	cmd := exec.Command(
		"sox",
		"-d",          // default audio device (mic)
		"-r", "16000", // sample rate
		"-b", "16", // bit depth
		"-c", "1", // channels (mono)
		tmpPath,                                         // output file
		"trim", "0", fmt.Sprintf("%d", durationSeconds), // duration
	)

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("sox recording failed: %w", err)
	}

	return os.Open(tmpPath)
}

func init() {
	defaultRecorder = &macOSRecorder{}
}
