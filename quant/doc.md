# 包功能说明

本包提供了一个简化的 AI 友好型加密货币交易所接口，目前支持通过环境变量配置的币安（现货）交易所。该包封装了交易所的核心功能，包括价格查询、历史数据获取、市价买卖交易、账户余额查询和订单管理，旨在为量化交易策略和自动化交易系统提供简洁易用的编程接口。典型使用场景包括加密货币交易机器人、价格监控系统和量化分析工具。

## 结构体与接口

本包未定义公开的结构体与接口，所有功能通过函数提供。

## 函数

```go
func FetchLatestPrice(symbol string) (*binance.SymbolPrice, error)
```

获取指定交易对的最新价格信息。参数 symbol 为交易对符号（如 "BTC/USDT"），函数会自动标准化符号格式。返回币安的 SymbolPrice 结构体指针，包含价格和符号信息。如果查询失败或找不到价格数据，会返回相应的错误。

```go
func FetchHistoricalPrices(symbol, interval string, limit int) ([]*binance.Kline, error)
```

获取指定交易对的历史 K 线数据。参数 symbol 为交易对符号，interval 为 K 线间隔（如 "1m", "1h", "1d"），limit 为获取的数据条数限制。返回币安 Kline 结构体指针的切片，包含开盘价、最高价、最低价、收盘价等历史价格信息。

```go
func BuyMarket(symbol string, quoteAmount float64) (*binance.CreateOrderResponse, error)
```

使用计价货币（如 USDT）进行市价买入操作。参数 symbol 为交易对符号，quoteAmount 为计价货币金额。返回币安的 CreateOrderResponse 结构体指针，包含订单执行结果。需要正确设置 EXCHANGE_API_KEY 和 EXCHANGE_API_SECRET 环境变量才能使用交易功能。

```go
func SellMarket(symbol string, baseAmount float64) (*binance.CreateOrderResponse, error)
```

卖出指定数量的基础资产。参数 symbol 为交易对符号，baseAmount 为基础资产数量。返回币安的 CreateOrderResponse 结构体指针，包含订单执行结果。需要正确设置交易所 API 密钥环境变量才能执行交易操作。

```go
func GetAccountBalance(currency string) (string, error)
```

查询指定币种的可用余额。参数 currency 为币种代码（如 "BTC", "USDT"）。返回该币种的可用余额字符串，如果找不到对应币种则返回 "0"。需要正确设置交易所 API 密钥环境变量才能访问账户信息。

```go
func ListOpenOrders(symbol string) ([]*binance.Order, error)
```

列出指定交易对的所有未成交订单。参数 symbol 为交易对符号。返回币安 Order 结构体指针的切片，包含所有未完成订单的详细信息。需要正确设置交易所 API 密钥环境变量才能查询订单状态。

## 变量与常量

本包未定义公开的变量与常量。