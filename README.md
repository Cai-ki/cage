```
以下内容由 AI 脚本总结
```
# Go 项目工具集

## 项目概述

本项目是一个综合性的 Go 工具集合，提供了多个功能模块，涵盖环境配置、AI 辅助开发、大语言模型集成、媒体捕获、通知服务和量化交易等领域。各模块设计独立且可复用，能够显著提升开发效率，简化常见开发任务。

## 包功能说明

### config - 环境配置管理

**核心功能**：自动从项目根目录加载 `.env` 文件到环境变量。

**公开接口**：
- `LoadEnvFromRoot()` - 自动查找项目根目录并加载环境变量文件

**使用场景**：规范化项目配置管理，避免硬编码敏感信息，支持多环境配置。

### helper - AI 智能转换工具

**核心功能**：基于大语言模型的数据格式转换和代码生成。

**主要转换功能**：
- JSON ↔ SQL/Go Struct/Protobuf
- SQL → JSON Schema
- CSV → SQL
- 自然语言描述 → SQL/Shell 脚本
- Go 包源码分析生成文档

**公开接口**：
- `JsonToSql()`, `JsonToGoStruct()`, `JsonToProto()`
- `SqlToJSONSchema()`, `CsvToSql()`, `DescriptionToSql()`
- `AnalyzePackage()`, `GeneratePackageDoc()`
- `DescribeToShellScript()`, `DescribeToRunnableShell()`

**使用场景**：快速原型开发、数据建模、API 设计、自动化脚本生成。

### llm - 大语言模型客户端

**核心功能**：统一的 LLM API 客户端，支持文本生成、视觉理解和文本嵌入。

**公开接口**：
- `Completion()`, `CompletionBySystem()` - 文本生成
- `Vision()`, `VisionWithPrompt()` - 图像分析
- `Embedding()`, `EmbeddingWithDim()` - 文本向量化

**配置支持**：通过环境变量配置 API 密钥、模型参数和基础 URL，兼容 OpenAI 和 Ollama 等服务。

### media - 跨平台媒体捕获

**核心功能**：屏幕截图和音频录制。

**公开接口**：
- `Screenshot()` - 捕获屏幕图像
- `RecordAudio()` - 录制音频流

**平台支持**：完整支持 macOS，Linux 和 Windows 返回未实现错误。

### notify - 统一通知服务

**核心功能**：通过 SMTP 发送电子邮件通知。

**公开接口**：
- `Send()` - 发送通知消息
- `NewEmailNotifier()` - 创建邮件通知器

**服务支持**：QQ、Gmail、163、Outlook 等主流邮件服务商。

### quant - 加密货币交易接口

**核心功能**：简化的币安交易所 API 封装。

**公开接口**：
- `FetchLatestPrice()`, `FetchHistoricalPrices()` - 获取价格数据
- `BuyMarket()`, `SellMarket()` - 执行交易订单
- `GetAccountBalance()`, `ListOpenOrders()` - 账户管理

**使用场景**：量化交易策略、自动化交易系统。

### sugar - 便捷开发工具包

**核心功能**：简化常见编程任务的辅助函数。

**主要功能**：
- `Assert()`, `Assertf()` - 条件断言
- `Coalesce()` - 非零值选择
- `ExitIfNot()`, `ExitIfErr()` - 错误处理
- `Must()` - 简化错误检查
- `StrToT()` - 类型安全转换

**使用场景**：减少样板代码，提高开发效率。

## 快速开始

### 环境配置

在项目根目录创建 `.env` 文件：

```env
LLM_API_KEY=your_openai_api_key
LLM_BASE_URL=https://api.openai.com/v1
SMTP_SERVICE=GMAIL
SMTP_EMAIL=your_email@gmail.com
SMTP_PASSWORD=your_app_password
```

### 基础使用示例

```go
package main

import (
    "fmt"
    "github.com/your-username/your-project/helper"
    "github.com/your-username/your-project/llm"
    "github.com/your-username/your-project/notify"
)

func main() {
    // AI 文本生成
    response, err := llm.Completion("Hello, how are you?")
    if err != nil {
        panic(err)
    }
    fmt.Println(response)

    // JSON 转 Go 结构体
    jsonStr := `{"name": "John", "age": 30}`
    goStruct, err := helper.JsonToGoStruct(jsonStr)
    if err != nil {
        panic(err)
    }
    fmt.Println(goStruct)

    // 发送通知
    err = notify.Send("系统通知", "任务执行完成")
    if err != nil {
        panic(err)
    }
}
```

### 构建和运行

```bash
# 克隆项目
git clone https://github.com/your-username/your-project.git
cd your-project

# 安装依赖
go mod tidy

# 构建
go build -o myapp ./cmd/your-app

# 运行
./myapp
```

## 项目结构概览

```
.
├── config/          # 环境配置管理
├── helper/          # AI 智能转换工具
├── llm/            # 大语言模型客户端
├── media/          # 媒体捕获功能
├── notify/         # 通知服务
├── quant/          # 量化交易接口
├── sugar/          # 便捷开发工具
└── cmd/            # 可执行程序入口
```

## 依赖要求

- Go 1.18+（支持泛型）
- 各包可能有额外的系统依赖（如 macOS 音频录制需要 sox）
- 网络连接（用于 LLM API 调用和交易所接口）

## 许可证

[在此添加项目许可证信息]