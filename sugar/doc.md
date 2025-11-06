# 包功能说明

sugar 包是一个 Go 语言工具包，提供了一系列便捷的辅助函数，旨在简化日常开发中的常见操作。该包的核心设计目标是减少样板代码，提高开发效率，通过类型安全的泛型实现和简洁的 API 设计，为开发者提供可靠的编程工具。典型使用场景包括条件断言、错误处理、零值判断、字符串转换等常见编程任务，特别适合在需要快速原型开发或代码简化的项目中应用。

## 结构体与接口

本包未定义公开的结构体与接口。

## 函数

```go
func Assert(cond bool, msg string)
```
Assert 函数用于条件断言，当条件 cond 为 false 时会触发 panic，并使用提供的消息 msg 作为错误信息。该函数适用于在开发阶段验证程序逻辑的正确性。

```go
func Assertf(cond bool, format string, args ...any)
```
Assertf 函数是 Assert 的格式化版本，当条件 cond 为 false 时会触发 panic，并使用 format 和 args 参数格式化错误信息。支持动态生成详细的错误消息。

```go
func Coalsece[T comparable](values ...T) T
```
Coalesce 函数接受一个可变参数列表 values，返回第一个非零值。如果所有值都是零值，则返回类型 T 的零值。该函数使用泛型约束 comparable，适用于任何可比较的类型。

```go
func ExitIfNot(cond bool, msg ...string)
```
ExitIfNot 函数检查条件 cond，如果为 false 则退出程序。可选的 msg 参数会写入标准错误输出，支持提供多个字符串但只使用第一个。

```go
func ExitIfNotf(cond bool, format string, args ...any)
```
ExitIfNotf 函数是 ExitIfNot 的格式化版本，当条件 cond 为 false 时退出程序，并使用 format 和 args 参数格式化错误消息到标准错误输出。

```go
func ExitIfErr(err error, msg ...string)
```
ExitIfErr 函数检查错误 err，如果不为 nil 则退出程序。可选的 msg 参数会写入标准错误输出，支持提供多个字符串但只使用第一个。

```go
func ExitIfErrf(err error, format string, args ...any)
```
ExitIfErrf 函数是 ExitIfErr 的格式化版本，当错误 err 不为 nil 时退出程序，并使用 format 和 args 参数格式化错误消息到标准错误输出。

```go
func Must[T any](v T, err error) T
```
Must 函数接受一个值 v 和错误 err，如果 err 不为 nil 则触发 panic，否则返回 v。该函数常用于简化错误处理，特别是在初始化操作中。

```go
func StrToT[T any](str string) (T, error)
```
StrToT 函数将字符串 str 转换为指定的类型 T。支持整数、浮点数、布尔值和字符串类型转换。如果转换失败，返回类型 T 的零值和错误信息。

```go
func StrToTWithDefault[T any](str string, def T) T
```
StrToTWithDefault 函数将字符串 str 转换为指定的类型 T，如果转换失败则返回默认值 def。支持整数、浮点数、布尔值和字符串类型转换，提供安全的转换机制。

## 变量与常量

本包未定义公开的变量与常量。