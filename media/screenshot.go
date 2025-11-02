package media

import "image"

// ScreenshotCapturer captures static screen images.
type ScreenshotCapturer interface {
	CaptureScreen() (image.Image, error)
}
