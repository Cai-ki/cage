# 包功能说明

quant 包提供了一个简化的 AI 友好接口，用于访问加密货币交易所。目前支持通过环境变量配置的币安（Binance）现货交易。该包封装了常见的交易操作，包括获取实时价格、历史 K 线数据、执行市价买卖订单、查询账户余额和查看未成交订单等，旨在为量化交易策略和自动化交易系统提供便捷的接入方式。

## 结构体与接口

本包未定义公开的结构体与接口，主要依赖于第三方库 `github.com/adshao/go-binance/v2` 中的类型。

## 函数

```go
func FetchLatestPrice(symbol string) (*binance.SymbolPrice, error)
```

FetchLatestPrice 获取指定交易对的最新价格。参数 symbol 为交易对符号（如 "BTC/USDT"），函数内部会进行标准化处理。返回值为币安 SymbolPrice 结构体指针，包含价格信息。如果获取失败或未找到价格数据，会返回相应的错误。

```go
func FetchHistoricalPrices(symbol, interval string, limit int) ([]*binance.Kline, error)
```

FetchHistoricalPrices 获取指定交易对的历史 K 线数据。参数 symbol 为交易对符号，interval 为 K 线间隔（如 "1h"），limit 为获取的数据条数限制。返回值为币安 Kline 结构体指针的切片，包含开盘价、最高价、最低价、收盘价等历史数据。

```go
func BuyMarket(symbol string, quoteAmount float64) (*binance.CreateOrderResponse, error)
```

BuyMarket 使用计价货币（如 USDT）执行市价买入订单。参数 symbol 为交易对符号，quoteAmount 为计价货币金额。返回值为币安 CreateOrderResponse 结构体指针，包含订单执行结果。如果客户端未初始化（缺少 API 密钥）或下单失败，会返回相应的错误。

```go
func SellMarket(symbol string, baseAmount float64) (*binance.CreateOrderResponse, error)
```

SellMarket 执行市价卖出订单，卖出指定数量的基础资产。参数 symbol 为交易对符号，baseAmount 为基础资产数量。返回值为币安 CreateOrderResponse 结构体指针，包含订单执行结果。如果客户端未初始化（缺少 API 密钥）或下单失败，会返回相应的错误。

```go
func GetAccountBalance(currency string) (string, error)
```

GetAccountBalance 查询指定币种的可用余额。参数 currency 为币种代码（如 "BTC"）。返回值为可用余额的字符串表示。如果客户端未初始化（缺少 API 密钥）或查询失败，会返回相应的错误。如果未找到指定币种的余额，则返回 "0"。

```go
func ListOpenOrders(symbol string) ([]*binance.Order, error)
```

ListOpenOrders 获取指定交易对的未成交订单列表。参数 symbol 为交易对符号。返回值为币安 Order 结构体指针的切片，包含订单详细信息。如果客户端未初始化（缺少 API 密钥）或查询失败，会返回相应的错误。

## 变量与常量

本包未定义公开的变量与常量。