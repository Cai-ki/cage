// capture/screen.go
package capture

import (
	"image"
	"sync"
)

var (
	defaultCapturer Capturer
	once            sync.Once
	initErr         error
)

// Screen captures the entire screen.
func Screen() (image.Image, error) {
	once.Do(func() {
		if defaultCapturer == nil {
			initErr = ErrNotImplemented
		}
	})
	if initErr != nil {
		return nil, initErr
	}
	return defaultCapturer.CaptureScreen()
}
