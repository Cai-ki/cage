package media

import "io"

// AudioRecorder captures audio streams.
type AudioRecorder interface {
	Record(durationSeconds int) (io.ReadCloser, error)
}
