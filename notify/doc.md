# 包功能说明

notify 包提供了一个简单统一的接口来发送通知消息，目前主要支持通过 SMTP 协议发送电子邮件。该包的设计目标是简化通知系统的集成，让开发者能够通过环境变量配置邮件服务参数，无需复杂的初始化代码即可发送通知。典型使用场景包括系统监控告警、任务执行结果通知、用户活动提醒等场景，特别适合在需要轻量级通知功能的微服务和自动化脚本中使用。

## 结构体与接口

```go
type Notifier interface {
	Send(subject, message string) error
}
```

Notifier 接口定义了通知发送器的基本契约，任何实现了 Send 方法的类型都可以作为通知器使用。Send 方法接收主题和消息内容两个参数，返回可能的错误信息。

```go
type EmailNotifier struct {
	Service  string
	Email    string
	Password string
	Name     string
}
```

EmailNotifier 结构体实现了通过 SMTP 协议发送电子邮件的功能。Service 字段指定邮件服务提供商，支持 QQ、GMAIL、163、OUTLOOK 等常见服务；Email 字段设置发件人和收件人邮箱地址；Password 字段存储 SMTP 密码或应用令牌；Name 字段定义发件人显示名称，如未设置则使用邮箱地址作为默认值。

```go
func (e *EmailNotifier) Send(subject, body string) error
```

Send 方法是 EmailNotifier 的核心方法，用于发送电子邮件到指定的邮箱地址。该方法会根据配置的服务商自动选择正确的 SMTP 服务器地址、端口和加密方式，支持 SSL/TLS 直连和 STARTTLS 两种安全连接方式。参数 subject 设置邮件主题，body 设置邮件正文内容，返回发送过程中可能出现的错误。

## 函数

```go
func NewEmailNotifier() (*EmailNotifier, error)
```

NewEmailNotifier 函数从环境变量创建并初始化 EmailNotifier 实例。该函数会读取 SMTP_SERVICE（服务商）、SMTP_EMAIL（邮箱地址）、SMTP_PASSWORD（密码）和可选的 SMTP_NAME（显示名称）环境变量。如果必需的环境变量缺失，会返回相应的错误信息。

```go
func Send(subject, body string) error
```

Send 函数是包的便捷入口，使用默认的通知通道发送消息。该函数采用懒加载方式初始化默认的邮件通知器，只需在首次调用时读取环境变量配置。参数 subject 设置通知主题，body 设置通知正文，返回发送过程中可能出现的错误。如果初始化失败，后续调用都会返回相同的初始化错误。

## 变量与常量

包中未定义公开的变量与常量。