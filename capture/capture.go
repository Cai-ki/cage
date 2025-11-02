// Package capture provides cross-platform screen capturing utilities.
package capture

import "image"

// Capturer is the interface for screen capture implementations.
type Capturer interface {
	CaptureScreen() (image.Image, error)
}
