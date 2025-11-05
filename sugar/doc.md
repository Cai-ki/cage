# 包功能说明

sugar 包是一个 Go 语言工具包，提供了一系列语法糖和便捷函数，旨在简化日常开发中的常见操作。该包的设计目标是通过提供类型安全的泛型函数和错误处理辅助工具，减少样板代码的编写，提高开发效率。典型使用场景包括条件断言、错误处理、类型转换和零值处理等，特别适合在需要快速原型开发或简化复杂条件判断的项目中使用。

## 结构体与接口

本包未定义公开的结构体与接口。

## 函数

```go
func Assert(cond bool, msg string)
```
当条件 `cond` 为 false 时，该函数会触发 panic，并使用指定的消息 `msg` 作为错误信息。主要用于在开发阶段进行条件检查，确保程序状态符合预期。

```go
func Assertf(cond bool, format string, args ...any)
```
与 `Assert` 功能类似，但支持格式化错误消息。当条件 `cond` 为 false 时，会根据 `format` 和 `args` 生成格式化字符串，并以此触发 panic。

```go
func Coalsece[T comparable](values ...T) T
```
该泛型函数接收多个同类型参数，返回第一个非零值。如果所有值都是零值，则返回该类型的零值。适用于在多个可能的值中选择第一个有效值的场景。

```go
func ExitIfNot(cond bool, msg ...string)
```
当条件 `cond` 为 false 时，程序会立即退出。如果提供了 `msg` 参数，会将第一个消息字符串写入标准错误输出。常用于命令行工具中的前置条件检查。

```go
func ExitIfNotf(cond bool, format string, args ...any)
```
与 `ExitIfNot` 功能类似，但支持格式化错误消息。当条件不满足时，会根据 `format` 和 `args` 生成格式化字符串并输出到标准错误，然后退出程序。

```go
func ExitIfErr(err error, msg ...string)
```
当错误 `err` 不为 nil 时，程序会立即退出。如果提供了 `msg` 参数，会将第一个消息字符串写入标准错误输出。用于简化错误处理流程。

```go
func ExitIfErrf(err error, format string, args ...any)
```
与 `ExitIfErr` 功能类似，但支持格式化错误消息。当错误不为 nil 时，会根据 `format` 和 `args` 生成格式化字符串并输出到标准错误，然后退出程序。

```go
func Must[T any](v T, err error) T
```
该泛型函数接收一个值 `v` 和一个错误 `err`。如果错误不为 nil，会触发 panic；否则返回原始值 `v`。常用于简化必须成功操作的错误处理，如文件打开、数据解析等。

```go
func StrToT[T any](str string) (T, error)
```
将字符串 `str` 转换为指定的类型 `T`。支持基本数据类型包括整数、浮点数、布尔值和字符串。如果转换失败，返回零值和错误信息。使用反射实现类型安全的转换。

```go
func StrToTWithDefault[T any](str string, def T) T
```
将字符串 `str` 转换为指定的类型 `T`，如果转换失败则返回默认值 `def`。支持的基本数据类型与 `StrToT` 相同，但在转换失败时不会返回错误，而是直接返回提供的默认值。

## 变量与常量

本包未定义公开的变量与常量。