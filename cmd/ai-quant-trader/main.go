package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Cai-ki/cage/llm"
	"github.com/Cai-ki/cage/quant"
	"github.com/adshao/go-binance/v2"
)

type Decision struct {
	Action string  `json:"action"` // "buy", "sell", "hold"
	Amount float64 `json:"amount"`
	Reason string  `json:"reason"`

	AnalystView string `json:"analyst_view"`
	RiskView    string `json:"risk_view"`
}

func formatKlines(klines []*binance.Kline, limit int) string {
	n := len(klines)
	start := 0
	if n > limit {
		start = n - limit
	}
	var b strings.Builder
	for _, k := range klines[start:] {
		b.WriteString(fmt.Sprintf("O:%s H:%s L:%s C:%s V:%s\n", k.Open, k.High, k.Low, k.Close, k.Volume))
	}
	return b.String()
}

func getCurrentPosition(symbol string, latest *binance.SymbolPrice) (holding bool, avgCost float64) {
	// 简化：从历史记录推算（实际可用 quant.GetAccountBalance + 订单记录）
	history := loadAllHistory(symbol)
	btcBalStr, _ := quant.GetAccountBalance("BTC")
	btcBal, _ := strconv.ParseFloat(btcBalStr, 64)
	currentPrice, _ := strconv.ParseFloat(latest.Price, 64) // 需要传入 latest.Price
	if currentPrice <= 0 {
		currentPrice = 1.0
	}
	if btcBal*currentPrice > 1.0 { // 持仓价值 > 1 USDT
		// 找最近一次买入均价（简化版）
		for i := len(history) - 1; i >= 0; i-- {
			if history[i].Action == "buy" {
				avgCost, _ = strconv.ParseFloat(history[i].Price, 64)
				break
			}
		}
		return true, avgCost
	}
	return false, 0
}

// ====== 在 main.go 顶部已有的 Decision 保持不变 ======

func RunTradingStep(symbol string) error {
	// Ensure init record exists
	if _, _, err := ensureInitRecord(symbol); err != nil {
		return fmt.Errorf("failed to ensure init record: %w", err)
	}

	// 1. Fetch data
	latest, err := quant.FetchLatestPrice(symbol)
	if err != nil {
		return fmt.Errorf("fetch latest price: %w", err)
	}

	usdtBalStr, _ := quant.GetAccountBalance("USDT")
	btcBalStr, _ := quant.GetAccountBalance("BTC")
	usdtBal, _ := strconv.ParseFloat(usdtBalStr, 64)
	btcBal, _ := strconv.ParseFloat(btcBalStr, 64)

	k5m, _ := quant.FetchHistoricalPrices(symbol, "5m", 48)
	k15m, _ := quant.FetchHistoricalPrices(symbol, "15m", 48)
	k1h, _ := quant.FetchHistoricalPrices(symbol, "1h", 24)
	k4h, _ := quant.FetchHistoricalPrices(symbol, "4h", 12)

	recentHist := loadRecentHistory(symbol, 3)
	histStr := formatRecentHistory(recentHist)
	summary := generateWeeklySummary(symbol)
	holding, avgCost := getCurrentPosition(symbol, latest)
	positionStr := fmt.Sprintf("Currently holding: %s BTC (est. avg cost: %.2f USDT)", btcBalStr, avgCost)
	if !holding {
		positionStr = "Currently not holding any BTC."
	}

	dataSource := "DATA SOURCE: Binance Testnet (prices may be synthetic)"

	// 3. Prompt
	prompt := fmt.Sprintf(`
You are an AI quant team with three roles:

[1] Market Analyst:
- Analyze price action, indicators, and K-line patterns across 5m/15m/1h/4h.
- Output: a concise technical assessment.

[2] Risk Manager:
- Check current balances, position size, recent trades, and exposure limits.
- Output: risk constraints and max allowable order size.

[3] Trading Executor:
- Combine [1] and [2] to decide final action.
- Output: {"action":"...", "amount":..., "reason":"..."}

CRITICAL RULES:
- NEVER contradict your recent actions (see "Recent decisions" below).
- NEVER exceed available balance.
- Be conservative: never risk >20%% of available capital in one trade.
- If uncertain, HOLD.

CURRENT CONTEXT:
- Time: %s
- Symbol: %s
- Current price: %s USDT
- USDT balance: %s
- BTC balance: %s
- %s
- %s

STRATEGY PERFORMANCE (last 7 days):
%s

RECENT DECISIONS (last 3 actions — AVOID CONFLICTS):
%s

MARKET DATA (OHLCV — focus on CLOSE and VOLUME trends):

[5m — last 12 candles]
%s

[15m — last 12 candles]
%s

[1h — last 12 candles]
%s

[4h — last 6 candles]
%s

YOUR TASK:
1. Decide action: "buy", "sell", or "hold"
2. Specify amount:
   - If "buy": USDT amount (e.g., 25.0)
   - If "sell": BTC amount (e.g., 0.0015)
   - If "hold": 0

YOUR RESPONSE FORMAT:
{
  "analyst_view": "string",
  "risk_view": "string",
  "action": "buy/sell/hold",
  "amount": number,
  "reason": "string"
}
DO NOT include any other text, markdown, or explanation.
`,
		time.Now().UTC().Format("2006-01-02 15:04:05 UTC"),
		symbol,
		latest.Price,
		usdtBalStr,
		btcBalStr,
		positionStr,
		dataSource,
		summary,
		histStr,
		formatKlines(k5m, 12),
		formatKlines(k15m, 12),
		formatKlines(k1h, 12),
		formatKlines(k4h, 6),
	)

	// 4. LLM call
	resp, err := llm.Completion(prompt)
	if err != nil {
		return fmt.Errorf("LLM call failed: %w", err)
	}

	// 5. Parse
	var dec Decision
	if err := json.Unmarshal([]byte(resp), &dec); err != nil {
		dec = Decision{Action: "hold", Amount: 0, Reason: "LLM parse error"}
	}

	// === 新增：记录完整 prompt + response ===
	appendPromptLog(PromptLogEntry{
		Timestamp:   time.Now().UTC(),
		Symbol:      symbol,
		Prompt:      prompt,
		LLMResponse: resp,
		Decision:    dec,
		Price:       latest.Price,
	})

	// 6. Execute
	subject := fmt.Sprintf("[AI Quant] %s %.6f %s", dec.Action, dec.Amount, symbol)
	var body strings.Builder
	body.WriteString(fmt.Sprintf("Reason: %s\nPrice: %s\nTime: %s\n",
		dec.Reason, latest.Price, time.Now().UTC().Format("15:04:05")))

	switch dec.Action {
	case "buy":
		if dec.Amount > 0 && dec.Amount <= usdtBal*0.95 {
			_, err := quant.BuyMarket(symbol, dec.Amount)
			if err != nil {
				body.WriteString(fmt.Sprintf("⚠️ BUY FAILED: %v", err))
			} else {
				body.WriteString(fmt.Sprintf("✅ BOUGHT %.2f USDT", dec.Amount))
			}
		} else {
			body.WriteString("⚠️ Invalid or insufficient buy amount")
		}
	case "sell":
		if dec.Amount > 0 && dec.Amount <= btcBal*0.99 {
			_, err := quant.SellMarket(symbol, dec.Amount)
			if err != nil {
				body.WriteString(fmt.Sprintf("⚠️ SELL FAILED: %v", err))
			} else {
				body.WriteString(fmt.Sprintf("✅ SOLD %.6f BTC", dec.Amount))
			}
		} else {
			body.WriteString("⚠️ Invalid or insufficient sell amount")
		}
	default:
		body.WriteString("⏸️ HOLD")
	}

	// 7. Log & Notify
	record := TradeRecord{
		Timestamp:   time.Now().UTC(),
		Symbol:      symbol,
		Action:      dec.Action,
		Price:       latest.Price,
		Reason:      dec.Reason,
		AnalystView: dec.AnalystView,
		RiskView:    dec.RiskView,
	}
	if dec.Action == "buy" {
		record.AmountUSDT = dec.Amount
	} else if dec.Action == "sell" {
		record.AmountBTC = dec.Amount
	}
	appendToHistory(record)

	exportNetValueCSV(symbol)

	// notify.Send(subject, body.String())
	fmt.Println(subject)
	fmt.Println(body.String())

	return nil
}

func main() {
	symbol := "BTC/USDT"

	os.Setenv("RUN_LOOP", "true")

	const t = 5

	// 初始化 server
	server := NewHTMLServer("testdata/net_value.csv", "testdata/trading_history.json")
	server.Start("127.0.0.1:9999")

	if os.Getenv("RUN_LOOP") == "true" {
		for {
			fmt.Println("🚀 Starting trading step...")
			if err := RunTradingStep(symbol); err != nil {
				fmt.Printf("❌ Error: %v\n", err)
			}
			fmt.Printf("⏳ Sleeping %d minutes...\n", t)
			time.Sleep(1 * time.Second)
			server.Update()
			time.Sleep(t * time.Minute)
		}
	} else {
		RunTradingStep(symbol)
	}
}
