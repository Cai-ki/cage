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
