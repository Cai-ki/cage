//go:build darwin
// +build darwin

package media

func initPlatform() {
	screenshotCapturer = &macOSScreenshotCapturer{}
	audioRecorder = &macOSAudioRecorder{}
}
