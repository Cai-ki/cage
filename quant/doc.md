# 包功能说明

本包提供了一个面向 AI 友好的加密货币交易所接口封装，目前支持 Binance 现货和合约交易。通过环境变量配置交易所认证信息，包内部自动初始化客户端连接。设计目标是简化量化交易策略的实现，提供统一的函数接口进行行情获取、账户查询和订单操作，适用于自动化交易系统和量化策略研究场景。

## 结构体与接口

本包未定义公开的结构体与接口，所有功能通过函数提供。

## 函数

```go
func FetchLatestPrice(symbol string) (*binance.SymbolPrice, error)
```

获取指定交易对的最新价格信息。参数 symbol 为交易对符号（如 "BTC/USDT"），函数内部会自动标准化格式。返回包含最新价格信息的 SymbolPrice 结构体指针，如果获取失败或交易对不存在则返回错误。

```go
func FetchHistoricalPrices(symbol, interval string, limit int) ([]*binance.Kline, error)
```

获取指定交易对的历史 K 线数据。参数 symbol 为交易对符号，interval 为 K 线间隔（如 "1m", "1h"），limit 为获取的数据条数。返回 K 线数据切片，包含开盘价、最高价、最低价、收盘价等信息。

```go
func BuyMarket(symbol string, quoteAmount float64) (*binance.CreateOrderResponse, error)
```

使用市价单买入交易对，以计价货币（如 USDT）数量作为下单金额。参数 symbol 为交易对符号，quoteAmount 为计价货币数量。返回订单创建响应，包含订单 ID、状态等信息。需要正确设置 EXCHANGE_SPOT_API_KEY 和 EXCHANGE_SPOT_API_SECRET 环境变量。

```go
func SellMarket(symbol string, baseAmount float64) (*binance.CreateOrderResponse, error)
```

使用市价单卖出交易对，以基础货币数量作为下单数量。参数 symbol 为交易对符号，baseAmount 为基础货币数量。返回订单创建响应。需要正确设置交易所 API 密钥环境变量。

```go
func GetAccountBalance(currency string) (string, error)
```

查询指定币种的可用余额。参数 currency 为币种符号（如 "BTC", "USDT"）。返回该币种的可用余额字符串，如果币种不存在则返回 "0"。需要正确设置交易所 API 密钥环境变量。

```go
func ListOpenOrders(symbol string) ([]*binance.Order, error)
```

列出指定交易对的所有未成交订单。参数 symbol 为交易对符号。返回订单信息切片，包含每个订单的详细信息。需要正确设置交易所 API 密钥环境变量。

```go
func FuturesGetBalance(asset string) (string, error)
```

获取 USDT 保证金合约账户的余额。参数 asset 为资产符号（如 "USDT"）。返回该资产的钱包余额字符串。需要正确设置 EXCHANGE_FUTURES_API_KEY 和 EXCHANGE_FUTURES_API_SECRET 环境变量。

```go
func FuturesGetPosition(symbol string) (*futures.PositionRisk, error)
```

获取指定交易对的持仓风险信息。参数 symbol 为交易对符号。返回包含持仓数量、杠杆、未实现盈亏等信息的 PositionRisk 结构体指针。如果该交易对没有持仓则返回错误。

```go
func FuturesBuyMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error)
```

使用市价单开立多头仓位。参数 symbol 为交易对符号，quantity 为下单数量。返回合约订单创建响应。需要正确设置合约交易 API 密钥环境变量。

```go
func FuturesSellMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error)
```

使用市价单开立空头仓位或平仓多头仓位。参数 symbol 为交易对符号，quantity 为下单数量。返回合约订单创建响应。需要正确设置合约交易 API 密钥环境变量。

```go
func FuturesClosePosition(symbol string) (*futures.CreateOrderResponse, error)
```

自动检测并平掉指定交易对的所有持仓。参数 symbol 为交易对符号。函数会自动检测持仓方向和数量，并使用减仓单（ReduceOnly）进行平仓操作。返回订单创建响应。

```go
func FuturesGetKlines(symbol, interval string, limit int) ([]*futures.Kline, error)
```

获取合约交易的历史 K 线数据。参数 symbol 为交易对符号，interval 为 K 线间隔，limit 为获取的数据条数。返回合约 K 线数据切片。

## 变量与常量

本包未定义公开的变量与常量。