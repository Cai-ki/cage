# 包功能说明

media 包提供了跨平台的媒体捕获功能，主要包括屏幕截图和音频录制。该包通过平台特定的实现来支持不同操作系统，目前完整支持 macOS 平台，在 Linux 和 Windows 平台上返回未实现错误。包内部使用延迟初始化机制，确保平台检测和资源初始化只执行一次。典型使用场景包括应用程序需要捕获用户屏幕内容或录制音频输入的场景，如远程协助、语音笔记、教学演示等应用。

## 结构体与接口

```go
type AudioRecorder interface {
	Record(durationSeconds int) (io.ReadCloser, error)
}
```

AudioRecorder 接口定义了音频录制功能，包含一个 Record 方法，用于录制指定时长的音频数据并返回可读取的音频流。

```go
type ScreenshotCapturer interface {
	CaptureScreen() (image.Image, error)
}
```

ScreenshotCapturer 接口定义了屏幕截图功能，包含一个 CaptureScreen 方法，用于捕获当前屏幕图像并返回图像对象。

## 函数

```go
func Screenshot() (image.Image, error)
```

Screenshot 函数捕获整个屏幕的图像，返回一个 image.Image 对象。在 macOS 平台上实际执行截图操作，在 Linux 和 Windows 平台上返回 ErrNotImplemented 错误。

```go
func RecordAudio(durationSeconds int) (io.ReadCloser, error)
```

RecordAudio 函数录制麦克风音频指定时长（以秒为单位），返回一个 WAV 格式的音频流（16kHz 采样率、16 位深度、单声道）。在 macOS 平台上需要 sox 工具支持，在 Linux 和 Windows 平台上返回 ErrNotImplemented 错误。

## 变量与常量

```go
var ErrNotImplemented = errors.New("media feature not implemented on this platform")
```

ErrNotImplemented 错误表示当前平台不支持该媒体功能，当在 Linux 或 Windows 平台上调用截图或录音功能时会返回此错误。

```go
var ErrSoxNotInstalled = errors.New("audio recording requires 'sox' — install via 'brew install sox'")
```

ErrSoxNotInstalled 错误表示在 macOS 平台上进行音频录制时缺少必需的 sox 工具，需要用户通过 brew install sox 命令安装 sox 后才能使用音频录制功能。