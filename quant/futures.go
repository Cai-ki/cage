package quant

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	futures "github.com/adshao/go-binance/v2/futures"
)

var (
	futuresClient *futures.Client
)

func init() {
	exchange := os.Getenv("EXCHANGE_NAME")
	if exchange == "" {
		exchange = "binance"
	}

	if strings.ToLower(exchange) != "binance" {
		return
	}

	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")
	if apiKey == "" || secretKey == "" {
		return
	}

	// Enable Testnet if needed
	if os.Getenv("BINANCE_TESTNET") == "true" {
		futures.UseTestnet = true
	}

	futuresClient = futures.NewClient(apiKey, secretKey)
}

// standardizeSymbolFutures converts "BTC/USDT" → "BTCUSDT"
// Same as spot, but kept separate for clarity
func standardizeSymbolFutures(symbol string) string {
	return standardizeSymbol(symbol)
}

// FuturesGetBalance returns USDT-margined futures account balance
func FuturesGetBalance(asset string) (string, error) {
	if futuresClient == nil {
		return "", errors.New("futures client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	resp, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", err
	}
	for _, b := range resp.Assets {
		if strings.EqualFold(b.Asset, asset) {
			return b.WalletBalance, nil
		}
	}
	return "0", nil
}

// FuturesGetPosition returns position info for a symbol
func FuturesGetPosition(symbol string) (*futures.PositionRisk, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	positions, err := futuresClient.NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("no position found for %s", symbol)
	}
	return positions[0], nil
}

// FuturesBuyMarket opens a long position with market order
func FuturesBuyMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	qStr := fmt.Sprintf("%.8f", quantity)
	return futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Type(futures.OrderTypeMarket).
		Side(futures.SideTypeBuy).
		Quantity(qStr).
		Do(context.Background())
}

// FuturesSellMarket opens a short position or closes long
func FuturesSellMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	qStr := fmt.Sprintf("%.8f", quantity)
	return futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Type(futures.OrderTypeMarket).
		Side(futures.SideTypeSell).
		Quantity(qStr).
		Do(context.Background())
}

// FuturesClosePosition closes current position (auto detect direction)
func FuturesClosePosition(symbol string) (*futures.CreateOrderResponse, error) {
	pos, err := FuturesGetPosition(symbol)
	if err != nil {
		return nil, err
	}
	qty := pos.PositionAmt
	if qty == "0" {
		return nil, fmt.Errorf("no open position for %s", symbol)
	}

	// Determine close side
	var side futures.SideType
	qtyAbs := qty
	if strings.HasPrefix(qty, "-") {
		side = futures.SideTypeBuy // short → buy to close
		qtyAbs = qty[1:]
	} else {
		side = futures.SideTypeSell // long → sell to close
	}

	return futuresClient.NewCreateOrderService().
		Symbol(standardizeSymbolFutures(symbol)).
		Type(futures.OrderTypeMarket).
		Side(side).
		Quantity(qtyAbs).
		ReduceOnly(true). // important: only reduce
		Do(context.Background())
}

// FuturesGetKlines returns historical futures klines
func FuturesGetKlines(symbol, interval string, limit int) ([]*futures.Kline, error) {
	symbol = standardizeSymbolFutures(symbol)
	return futuresClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
}

func FuturesGetTickerPrice(symbol string) (string, error) {
	if futuresClient == nil {
		return "", errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	resp, err := futuresClient.NewPremiumIndexService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	if len(resp) > 0 {
		return resp[0].MarkPrice, nil // 或者 resp[0].LastPrice
	}
	return "", errors.New("no price data")
}

// FuturesGetCurrentFundingRate 获取指定交易对的当前资金费率
// 返回资金费率 (例如 0.0001 表示 0.01%) 和下次更新时间戳
func FuturesGetCurrentFundingRate(symbol string) (fundingRate string, nextFundingTime int64, err error) {
	if futuresClient == nil {
		return "", 0, errors.New("futures client not initialized (missing EXCHANGE_FUTURES_API_KEY/SECRET)")
	}
	symbol = standardizeSymbolFutures(symbol)

	// 使用 PremiumIndexService 获取标记价格和资金费率信息
	// Binance API 的 Premium Index 端点包含了资金费率信息
	resp, err := futuresClient.NewPremiumIndexService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", 0, err
	}

	if len(resp) == 0 {
		return "", 0, fmt.Errorf("no funding rate data found for symbol %s", symbol)
	}

	// 检查是否获取到了资金费率
	// FundingRate 可能为 nil 或 "0" 如果当前是结算时刻或无费率
	if resp[0].LastFundingRate == "" {
		return "0", resp[0].NextFundingTime, nil // 或者返回错误，取决于您的需求
	}

	return resp[0].LastFundingRate, resp[0].NextFundingTime, nil
}

// CommissionRateResponse 定义 API 响应结构
type CommissionRateResponse struct {
	Symbol              string `json:"symbol"`
	MakerCommissionRate string `json:"makerCommissionRate"` // 例如 "0.000200" 表示 0.02%
	TakerCommissionRate string `json:"takerCommissionRate"` // 例如 "0.000400" 表示 0.04%
}

// FuturesGetFeeRateForSymbol 获取指定交易对的当前交易手续费率
func FuturesGetFeeRateForSymbol(symbol string) (makerFeeRate string, takerFeeRate string, err error) {
	symbol = standardizeSymbolFutures(symbol)

	params := url.Values{}
	params.Add("symbol", symbol)

	var commissionData CommissionRateResponse
	err = BinanceFuturesRequest("GET", "/fapi/v1/commissionRate", params, &commissionData)
	if err != nil {
		return "", "", err
	}

	return commissionData.MakerCommissionRate, commissionData.TakerCommissionRate, nil
}

// FuturesGetAccountDefaultFeeRate 获取账户的默认交易手续费率
func FuturesGetAccountDefaultFeeRate() (makerFeeRate string, takerFeeRate string, err error) {
	return FuturesGetFeeRateForSymbol("BTCUSDT")
}

// CreateSignature 创建 HMAC SHA256 签名
// queryString: 需要签名的查询参数字符串
// secretKey: API密钥
// 返回: 十六进制格式的签名字符串
func CreateSignature(queryString, secretKey string) string {
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(queryString))
	return hex.EncodeToString(h.Sum(nil))
}

// CreateSignatureFromParams 从 url.Values 创建签名
// params: 请求参数
// secretKey: API密钥
// 返回: 签名和编码后的查询字符串
func CreateSignatureFromParams(params url.Values, secretKey string) (signature, queryString string) {
	queryString = params.Encode()
	signature = CreateSignature(queryString, secretKey)
	return signature, queryString
}

// BinanceFuturesRequest 通用的 Binance Futures API 请求函数
// method: HTTP 方法 (GET, POST, etc.)
// endpoint: API 端点路径 (e.g., "/fapi/v1/commissionRate")
// params: 请求参数
// result: 用于存储响应数据的指针
func BinanceFuturesRequest(method, endpoint string, params url.Values, result interface{}) error {
	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")
	if apiKey == "" || secretKey == "" {
		return errors.New("API credentials not set")
	}

	baseURL := "https://testnet.binancefuture.com"

	// 添加时间戳（如果还没有）
	if params.Get("timestamp") == "" {
		params.Add("timestamp", strconv.FormatInt(time.Now().UnixMilli(), 10))
	}

	// 创建签名
	signature, queryString := CreateSignatureFromParams(params, secretKey)

	// 构建完整URL
	fullURL := baseURL + endpoint + "?" + queryString + "&signature=" + signature

	req, err := http.NewRequest(method, fullURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("X-MBX-APIKEY", apiKey)

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to perform request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal response JSON: %w", err)
	}

	return nil
}
