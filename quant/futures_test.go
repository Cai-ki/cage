package quant_test

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	_ "github.com/Cai-ki/cage/config"
	"github.com/Cai-ki/cage/quant"
	"github.com/Cai-ki/cage/sugar"
)

func TestFuturesGetBalance(t *testing.T) {
	bal := sugar.Must(quant.FuturesGetBalance("USDT"))
	t.Log("Futures USDT balance:", bal)
}

func TestFuturesGetPosition(t *testing.T) {
	pos := sugar.Must(quant.FuturesGetPosition("BTCUSDT"))
	data := sugar.Must(json.Marshal(pos))
	t.Log("Position:", string(data))
}

func TestFuturesGetKlines(t *testing.T) {
	klines := sugar.Must(quant.FuturesGetKlines("BTCUSDT", "1h", 5))
	data := sugar.Must(json.Marshal(klines))
	t.Log("Klines:", string(data))
}

// ⚠️ DANGEROUS: Only run on futures testnet with small amounts
func TestFuturesBuySellAndClose(t *testing.T) {
	// Check balance
	bal := sugar.Must(quant.FuturesGetBalance("USDT"))
	usdt, _ := strconv.ParseFloat(bal, 64)
	if usdt < 10.0 {
		t.Skip("Insufficient USDT balance on futures testnet")
	}

	symbol := "BTCUSDT"

	// Buy (open long)
	order1 := sugar.Must(quant.FuturesBuyMarket(symbol, 0.001))
	t.Logf("Opened long: %+v", order1)

	// Check position
	pos := sugar.Must(quant.FuturesGetPosition(symbol))
	t.Logf("Position after buy: %+v", pos)

	// Close position
	closeOrder := sugar.Must(quant.FuturesClosePosition(symbol))
	t.Logf("Closed position: %+v", closeOrder)

	// Verify closed
	pos2 := sugar.Must(quant.FuturesGetPosition(symbol))
	amt, err := strconv.ParseFloat(pos2.PositionAmt, 64)
	if err != nil {
		t.Fatalf("Failed to parse position amount: %v", err)
	}

	// 判断是否接近 0（考虑浮点精度）
	if math.Abs(amt) > 1e-8 {
		t.Errorf("Position not fully closed: %s", pos2.PositionAmt)
	}
}

func TestFuturesGetCurrentFundingRate(t *testing.T) {
	symbol := "BTC/USDT"
	fundingRate, nextFundingTime, err := quant.FuturesGetCurrentFundingRate(symbol)

	if err != nil {
		t.Fatalf("Failed to get funding rate for %s: %v", symbol, err)
	}

	// 打印结果供参考
	fmt.Printf("Symbol: %s\n", symbol)
	fmt.Printf("Current Funding Rate: %s\n", fundingRate)
	fmt.Printf("Next Funding Time (Unix Timestamp): %d\n", nextFundingTime)
	fmt.Printf("Next Funding Time (Human Readable): %s\n", time.Unix(nextFundingTime/1000, 0).UTC().Format(time.RFC3339))

	// 验证返回值类型和基本逻辑
	// 资金费率通常是一个小数字符串，如 "0.0001" 或 "-0.00005"
	// nextFundingTime 应该是一个正整数时间戳
	if nextFundingTime <= 0 {
		t.Errorf("Expected a positive NextFundingTime, got %d", nextFundingTime)
	}

	// 简单检查 fundingRate 是否为一个看起来像数字的字符串
	// 更严格的验证可以使用 strconv.ParseFloat
	// 例如，资金费率可能为 "0" (在特定时间点)
	if fundingRate == "" {
		t.Errorf("Expected a non-empty FundingRate string, got empty")
	}

	t.Logf("Successfully retrieved funding rate for %s: %s", symbol, fundingRate)
}

// 添加这个函数来检查服务器时间
func checkServerTime() error {
	resp, err := http.Get("https://testnet.binancefuture.com/fapi/v1/time")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var timeResponse struct {
		ServerTime int64 `json:"serverTime"`
	}

	if err := json.Unmarshal(body, &timeResponse); err != nil {
		return err
	}

	localTime := time.Now().UnixMilli()
	serverTime := timeResponse.ServerTime
	diff := localTime - serverTime

	fmt.Printf("Local time:  %d (%s)\n", localTime, time.Unix(localTime/1000, 0).Format(time.RFC3339))
	fmt.Printf("Server time: %d (%s)\n", serverTime, time.Unix(serverTime/1000, 0).Format(time.RFC3339))
	fmt.Printf("Time difference: %d ms (%.2f seconds)\n", diff, float64(diff)/1000.0)

	if abs(diff) > 60000 { // 超过1分钟
		return fmt.Errorf("time difference too large: %d ms", diff)
	}

	return nil
}

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func TestFuturesGetFeeRateForSymbol(t *testing.T) {
	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")

	if apiKey == "" || secretKey == "" {
		t.Skip("Skipping test: API credentials not set")
	}

	// 首先检查时间同步
	t.Log("Checking server time synchronization...")
	if err := checkServerTime(); err != nil {
		t.Fatalf("Time synchronization check failed: %v", err)
	}

	makerRate, takerRate, err := quant.FuturesGetFeeRateForSymbol("BTCUSDT")

	if err != nil {
		t.Fatalf("Failed to get fee rate: %v", err)
	}

	t.Logf("Successfully retrieved fee rates - Maker: %s, Taker: %s", makerRate, takerRate)

	// 验证费率
	if makerRate == "" || takerRate == "" {
		t.Errorf("Expected non-empty fee rates")
	}
}

func TestFuturesGetAccountDefaultFeeRate(t *testing.T) {
	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")

	if apiKey == "" || secretKey == "" {
		t.Skip("Skipping test: API credentials not set in environment variables (EXCHANGE_FUTURES_API_KEY, EXCHANGE_FUTURES_API_SECRET)")
	}

	makerRate, takerRate, err := quant.FuturesGetAccountDefaultFeeRate()

	if err != nil {
		t.Fatalf("Failed to get default account fee rate: %v", err)
	}

	fmt.Printf("Account Default Fee Rate (via BTCUSDT):\n")
	fmt.Printf("  Maker Fee Rate: %s\n", makerRate)
	fmt.Printf("  Taker Fee Rate: %s\n", takerRate)

	// Similar validation as above
	if makerRate == "" || takerRate == "" {
		t.Errorf("Expected non-empty maker and taker fee rates, got maker='%s', taker='%s'", makerRate, takerRate)
	}

	makerFloat, _ := strconv.ParseFloat(makerRate, 64)
	takerFloat, _ := strconv.ParseFloat(takerRate, 64)
	if takerFloat < makerFloat {
		t.Logf("Warning: Taker fee (%f) is less than Maker fee (%f). This can happen based on account tier or BNB holdings.", takerFloat, makerFloat)
	}

	t.Logf("Successfully retrieved default account fee rates - Maker: %s (%f), Taker: %s (%f)", makerRate, makerFloat, takerRate, takerFloat)
}

// 测试基础连接和可用接口
func TestBinanceFuturesAPIConnection(t *testing.T) {
	// 测试1: Ping 接口
	resp, err := http.Get("https://testnet.binancefuture.com/fapi/v1/ping")
	if err != nil {
		t.Fatalf("Ping failed: %v", err)
	}
	defer resp.Body.Close()
	t.Logf("Ping response status: %d", resp.StatusCode)

	// 测试2: 服务器时间接口
	resp2, err := http.Get("https://testnet.binancefuture.com/fapi/v1/time")
	if err != nil {
		t.Fatalf("Time API failed: %v", err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp2.Body)
	t.Logf("Time API response: %s", string(body))

	// 测试3: 交易所信息
	resp3, err := http.Get("https://testnet.binancefuture.com/fapi/v1/exchangeInfo")
	if err != nil {
		t.Fatalf("ExchangeInfo failed: %v", err)
	}
	defer resp3.Body.Close()
	t.Logf("ExchangeInfo status: %d", resp3.StatusCode)
}

// 尝试不同的接口路径
func TestDifferentFeeEndpoints(t *testing.T) {
	endpoints := []string{
		"https://testnet.binancefuture.com/fapi/v1/commissionRate",
		"https://testnet.binancefuture.com/fapi/v2/commissionRate",
		"https://testnet.binancefuture.com/fapi/v1/account", // 账户信息可能包含费率
		"https://testnet.binancefuture.com/fapi/v1/makerCommission",
		"https://testnet.binancefuture.com/fapi/v1/takerCommission",
	}

	for _, endpoint := range endpoints {
		resp, err := http.Get(endpoint)
		if err != nil {
			t.Logf("Endpoint %s: ERROR - %v", endpoint, err)
			continue
		}
		t.Logf("Endpoint %s: Status %d", endpoint, resp.StatusCode)
		resp.Body.Close()
	}
}

func TestBinanceAPIDiagnostic(t *testing.T) {
	// 1. 测试基础连接
	t.Log("=== Testing Basic Connectivity ===")
	resp, err := http.Get("https://testnet.binancefuture.com/fapi/v1/ping")
	if err != nil {
		t.Fatalf("Cannot connect to testnet: %v", err)
	}
	t.Logf("Ping successful: %d", resp.StatusCode)
	resp.Body.Close()

	// 2. 测试服务器时间
	t.Log("=== Testing Server Time ===")
	resp2, err := http.Get("https://testnet.binancefuture.com/fapi/v1/time")
	if err != nil {
		t.Fatalf("Time API failed: %v", err)
	}
	body, _ := io.ReadAll(resp2.Body)
	t.Logf("Server time response: %s", string(body))
	resp2.Body.Close()

	// 3. 测试认证请求（如果有API密钥）
	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	if apiKey != "" {
		t.Log("=== Testing Authenticated Request ===")
		// 测试账户信息接口
		endpoint := "https://testnet.binancefuture.com/fapi/v2/account"
		req, _ := http.NewRequest("GET", endpoint, nil)
		req.Header.Set("X-MBX-APIKEY", apiKey)

		client := &http.Client{Timeout: 10 * time.Second}
		resp3, err := client.Do(req)
		if err != nil {
			t.Logf("Authenticated request error: %v", err)
		} else {
			t.Logf("Authenticated request status: %d", resp3.StatusCode)
			if resp3.Body != nil {
				resp3.Body.Close()
			}
		}
	}
}

// 专门测试签名的函数
func TestSignatureGeneration(t *testing.T) {
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")
	if secretKey == "" {
		t.Skip("Secret key not set")
	}

	// 测试用例
	testCases := []struct {
		name   string
		params url.Values
	}{
		{
			name: "Basic symbol query",
			params: url.Values{
				"symbol":    []string{"BTCUSDT"},
				"timestamp": []string{"1762941810113"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			queryString := tc.params.Encode()
			t.Logf("Query string: %s", queryString)

			mac := hmac.New(sha256.New, []byte(secretKey))
			mac.Write([]byte(queryString))
			signature := hex.EncodeToString(mac.Sum(nil))

			t.Logf("Generated signature: %s", signature)

			// 验证签名长度
			if len(signature) != 64 {
				t.Errorf("Expected signature length 64, got %d", len(signature))
			}
		})
	}
}

// 使用 Binance 官方示例验证签名
func TestOfficialSignatureExample(t *testing.T) {
	// 官方示例：https://binance-docs.github.io/apidocs/futures/cn/#user-data-32
	secretKey := "NhqPtmdSJYdKjVHjA7PZj4Mge3R5YNiP1e3UZjInClVN65XAbvqqM6A7H5fATj0j"
	queryString := "symbol=LTCBTC&side=BUY&type=LIMIT&timeInForce=GTC&quantity=1&price=0.1&recvWindow=5000&timestamp=1499827319559"

	expectedSignature := "c8db56825ae71d6d79447849e617115f4a920fa2acdcab2b053c4b2838bd6b71"

	mac := hmac.New(sha256.New, []byte(secretKey))
	mac.Write([]byte(queryString))
	actualSignature := hex.EncodeToString(mac.Sum(nil))

	t.Logf("Expected: %s", expectedSignature)
	t.Logf("Actual:   %s", actualSignature)

	if actualSignature != expectedSignature {
		t.Errorf("Signature mismatch! Expected %s, got %s", expectedSignature, actualSignature)
	} else {
		t.Log("✓ Signature algorithm is correct!")
	}
}
