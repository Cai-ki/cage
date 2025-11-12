# 包功能说明

quant 包是一个专为量化交易和自动化策略设计的加密货币交易所接口库，提供对 Binance 现货和合约交易的简化访问。该包封装了复杂的交易所 API 调用，提供统一的函数接口，支持市场数据获取、账户管理、订单执行等核心交易功能。通过环境变量配置 API 密钥和网络设置，开发者可以快速集成现货和合约交易能力到自己的交易系统中，特别适合用于算法交易、套利策略和风险管理等场景。

## 结构体与接口

```go
type CommissionRateResponse struct {
    Symbol              string `json:"symbol"`
    MakerCommissionRate string `json:"makerCommissionRate"`
    TakerCommissionRate string `json:"takerCommissionRate"`
}
```

CommissionRateResponse 结构体用于存储币安合约交易手续费率查询的响应数据。Symbol 字段表示交易对符号，如 "BTCUSDT"；MakerCommissionRate 字段表示挂单手续费率，例如 "0.000200" 表示 0.02%；TakerCommissionRate 字段表示吃单手续费率，例如 "0.000400" 表示 0.04%。

## 函数

```go
func BinanceFuturesRequest(method, endpoint string, params url.Values, result interface{}) error
```

BinanceFuturesRequest 函数执行通用的币安合约 API 请求，支持 GET、POST 等 HTTP 方法。method 参数指定 HTTP 请求方法，endpoint 参数指定 API 端点路径，params 参数包含请求参数，result 参数用于存储解析后的响应数据。函数会自动添加时间戳和签名，处理认证头信息，并返回请求执行结果。

```go
func BuyMarket(symbol string, quoteAmount float64) (*binance.CreateOrderResponse, error)
```

BuyMarket 函数在现货市场使用市价单买入指定交易对。symbol 参数指定交易对符号，quoteAmount 参数指定计价货币的数量。函数返回订单创建结果，包含订单ID、状态等信息。

```go
func CreateSignature(queryString, secretKey string) string
```

CreateSignature 函数使用 HMAC SHA256 算法创建请求签名。queryString 参数是需要签名的查询参数字符串，secretKey 参数是 API 密钥。函数返回十六进制格式的签名字符串，用于 API 请求的身份验证。

```go
func CreateSignatureFromParams(params url.Values, secretKey string) (signature, queryString string)
```

CreateSignatureFromParams 函数从 url.Values 参数创建签名。params 参数包含请求参数，secretKey 参数是 API 密钥。函数返回签名和编码后的查询字符串，便于构建完整的 API 请求。

```go
func FetchHistoricalPrices(symbol, interval string, limit int) ([]*binance.Kline, error)
```

FetchHistoricalPrices 函数获取现货市场的历史K线数据。symbol 参数指定交易对符号，interval 参数指定K线时间间隔，limit 参数指定返回的数据条数。函数返回K线数据数组，包含开盘价、最高价、最低价、收盘价等信息。

```go
func FetchLatestPrice(symbol string) (*binance.SymbolPrice, error)
```

FetchLatestPrice 函数获取现货市场指定交易对的最新价格。symbol 参数指定交易对符号。函数返回包含最新价格信息的 SymbolPrice 结构体。

```go
func FuturesBuyMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error)
```

FuturesBuyMarket 函数在合约市场使用市价单开多仓。symbol 参数指定交易对符号，quantity 参数指定合约数量。函数返回订单创建结果，用于确认开仓操作是否成功。

```go
func FuturesClosePosition(symbol string) (*futures.CreateOrderResponse, error)
```

FuturesClosePosition 函数自动检测并关闭指定交易对的持仓。symbol 参数指定交易对符号。函数会自动判断持仓方向并执行相应的平仓操作，使用 ReduceOnly 标记确保只减少持仓。

```go
func FuturesGetAccountDefaultFeeRate() (makerFeeRate string, takerFeeRate string, err error)
```

FuturesGetAccountDefaultFeeRate 函数获取账户的默认交易手续费率。函数返回挂单手续费率和吃单手续费率，通过查询 BTCUSDT 交易对的费率来获取默认值。

```go
func FuturesGetBalance(asset string) (string, error)
```

FuturesGetBalance 函数获取 USDT 保证金合约账户的余额。asset 参数指定资产类型，如 "USDT"。函数返回指定资产的钱包余额字符串。

```go
func FuturesGetCurrentFundingRate(symbol string) (fundingRate string, nextFundingTime int64, err error)
```

FuturesGetCurrentFundingRate 函数获取指定交易对的当前资金费率。symbol 参数指定交易对符号。函数返回资金费率字符串和下次资金费率更新的时间戳，资金费率如 "0.0001" 表示 0.01%。

```go
func FuturesGetFeeRateForSymbol(symbol string) (makerFeeRate string, takerFeeRate string, err error)
```

FuturesGetFeeRateForSymbol 函数获取指定交易对的当前交易手续费率。symbol 参数指定交易对符号。函数返回该交易对的挂单手续费率和吃单手续费率。

```go
func FuturesGetKlines(symbol, interval string, limit int) ([]*futures.Kline, error)
```

FuturesGetKlines 函数获取合约市场的历史K线数据。symbol 参数指定交易对符号，interval 参数指定K线时间间隔，limit 参数指定返回的数据条数。函数返回合约K线数据数组。

```go
func FuturesGetPosition(symbol string) (*futures.PositionRisk, error)
```

FuturesGetPosition 函数获取指定交易对的持仓信息。symbol 参数指定交易对符号。函数返回 PositionRisk 结构体，包含持仓数量、杠杆、未实现盈亏等详细信息。

```go
func FuturesGetTickerPrice(symbol string) (string, error)
```

FuturesGetTickerPrice 函数获取合约市场指定交易对的标记价格。symbol 参数指定交易对符号。函数返回标记价格字符串，用于合约估值和保证金计算。

```go
func FuturesSellMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error)
```

FuturesSellMarket 函数在合约市场使用市价单开空仓或平多仓。symbol 参数指定交易对符号，quantity 参数指定合约数量。函数返回订单创建结果，用于确认交易操作是否成功。

```go
func GetAccountBalance(currency string) (string, error)
```

GetAccountBalance 函数获取现货账户指定币种的可用余额。currency 参数指定币种符号，如 "BTC" 或 "USDT"。函数返回该币种的可用余额字符串。

```go
func ListOpenOrders(symbol string) ([]*binance.Order, error)
```

ListOpenOrders 函数获取现货市场指定交易对的所有未成交订单。symbol 参数指定交易对符号。函数返回订单数组，包含每个订单的详细信息如订单ID、数量、价格等。

```go
func SellMarket(symbol string, baseAmount float64) (*binance.CreateOrderResponse, error)
```

SellMarket 函数在现货市场使用市价单卖出指定交易对。symbol 参数指定交易对符号，baseAmount 参数指定基础货币的数量。函数返回订单创建结果，用于确认卖出操作是否成功。