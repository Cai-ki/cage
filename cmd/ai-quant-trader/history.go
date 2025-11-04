package main

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/Cai-ki/cage/quant"
)

type TradeRecord struct {
	Timestamp  time.Time `json:"timestamp"`
	Symbol     string    `json:"symbol"`
	Action     string    `json:"action"` // "buy", "sell", "hold"
	AmountUSDT float64   `json:"amount_usdt,omitempty"`
	AmountBTC  float64   `json:"amount_btc,omitempty"`
	Price      string    `json:"price"`
	Reason     string    `json:"reason"`

	AnalystView string `json:"analyst_view"`
	RiskView    string `json:"risk_view"`
}

const historyFile = "testdata/trading_history.json"

func loadAllHistory(symbol string) []TradeRecord {
	data, err := os.ReadFile(historyFile)
	if err != nil {
		return nil
	}
	var all []TradeRecord
	json.Unmarshal(data, &all)

	var filtered []TradeRecord
	for _, r := range all {
		if r.Symbol == symbol {
			filtered = append(filtered, r)
		}
	}
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].Timestamp.Before(filtered[j].Timestamp)
	})
	return filtered
}

func loadRecentHistory(symbol string, maxEntries int) []TradeRecord {
	all := loadAllHistory(symbol)
	if len(all) == 0 {
		return nil
	}
	start := len(all) - maxEntries
	if start < 0 {
		start = 0
	}
	return all[start:]
}

func appendToHistory(record TradeRecord) {
	var all []TradeRecord
	if data, err := os.ReadFile(historyFile); err == nil {
		json.Unmarshal(data, &all)
	}
	all = append(all, record)
	data, _ := json.MarshalIndent(all, "", "  ")
	os.WriteFile(historyFile, data, 0644)
}

func formatRecentHistory(records []TradeRecord) string {
	if len(records) == 0 {
		return "No prior decisions.\n"
	}
	var b strings.Builder
	b.WriteString("Recent decisions (last 3):\n")
	for _, r := range records {
		if r.Action == "init" { // 跳过 init 记录
			continue
		}
		ts := r.Timestamp.Format("15:04")
		switch r.Action {
		case "buy":
			b.WriteString(fmt.Sprintf("- [%s] BUY %.2f USDT at %s (%s)\n", ts, r.AmountUSDT, r.Price, r.Reason))
		case "sell":
			b.WriteString(fmt.Sprintf("- [%s] SELL %.6f BTC at %s (%s)\n", ts, r.AmountBTC, r.Price, r.Reason))
		default:
			b.WriteString(fmt.Sprintf("- [%s] HOLD at %s (%s)\n", ts, r.Price, r.Reason))
		}
	}
	return b.String()
}

func generateWeeklySummary(symbol string) string {
	history := loadAllHistory(symbol)
	if len(history) == 0 {
		return "No historical trades."
	}

	cutoff := time.Now().UTC().Add(-7 * 24 * time.Hour)
	var recent []TradeRecord
	for _, r := range history {
		if r.Timestamp.After(cutoff) {
			recent = append(recent, r)
		}
	}

	if len(recent) == 0 {
		return "No trades in the past 7 days."
	}

	buys, sells, wins := 0, 0, 0
	var lastBuyPrice float64

	for _, r := range recent {
		if r.Action == "buy" {
			buys++
			lastBuyPrice, _ = strconv.ParseFloat(r.Price, 64)
		} else if r.Action == "sell" {
			sells++
			if lastBuyPrice > 0 {
				sellPrice, _ := strconv.ParseFloat(r.Price, 64)
				if sellPrice > lastBuyPrice {
					wins++
				}
			}
		}
	}

	if buys == 0 && sells == 0 {
		return "No trades in the past 7 days."
	}

	wins = 0
	for _, r := range recent {
		if r.Action == "buy" {
			buys++
			lastBuyPrice, _ = strconv.ParseFloat(r.Price, 64)
		} else if r.Action == "sell" {
			sells++
			if lastBuyPrice > 0 {
				sellPrice, _ := strconv.ParseFloat(r.Price, 64)
				if sellPrice > lastBuyPrice {
					wins++
				}
			}
		}
	}

	if sells > 0 {
		winRate := float64(wins) / float64(sells) * 100
		return fmt.Sprintf("[7-Day Summary] Trades: %d buy, %d sell | Win rate: %.1f%%", buys, sells, winRate)
	} else {
		return fmt.Sprintf("[7-Day Summary] Trades: %d buy, %d sell", buys, sells)
	}
}

// PromptLogEntry records the full context sent to LLM and its response
type PromptLogEntry struct {
	Timestamp   time.Time `json:"timestamp"`
	Symbol      string    `json:"symbol"`
	Prompt      string    `json:"prompt"`
	LLMResponse string    `json:"llm_response"`
	Decision    Decision  `json:"decision"`
	Price       string    `json:"price_at_time"`
}

const promptLogFile = "testdata/prompt_log.json"

// appendPromptLog saves the full prompt and LLM interaction
func appendPromptLog(entry PromptLogEntry) {
	var all []PromptLogEntry
	if data, err := os.ReadFile(promptLogFile); err == nil {
		json.Unmarshal(data, &all)
	}
	all = append(all, entry)
	data, _ := json.MarshalIndent(all, "", "  ")
	os.WriteFile(promptLogFile, data, 0644)
}

// ensureInitRecord checks if an "init" record exists in history.
// If not, it creates one using current API balance and appends it.
// It returns the initial USDT and BTC balances.
func ensureInitRecord(symbol string) (initUSDT float64, initBTC float64, err error) {
	history := loadAllHistory(symbol)

	// Check if "init" record already exists
	for _, r := range history {
		if r.Action == "init" {
			return r.AmountUSDT, r.AmountBTC, nil
		}
	}

	// No "init" record found: create one from current exchange balance
	usdtBalStr, err1 := quant.GetAccountBalance("USDT")
	btcBalStr, err2 := quant.GetAccountBalance("BTC")
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("failed to fetch balance for init: USDT err=%v, BTC err=%v", err1, err2)
	}

	initUSDT, err1 = strconv.ParseFloat(usdtBalStr, 64)
	initBTC, err2 = strconv.ParseFloat(btcBalStr, 64)
	if err1 != nil || err2 != nil {
		return 0, 0, fmt.Errorf("failed to parse balance: USDT=%s (%v), BTC=%s (%v)", usdtBalStr, err1, btcBalStr, err2)
	}

	// Get current price to fill the "Price" field (required by struct, though unused for init)
	latest, err := quant.FetchLatestPrice(symbol)
	priceStr := "1.0"
	if err == nil {
		priceStr = latest.Price
	}

	// Create and append the "init" record
	initRecord := TradeRecord{
		Timestamp:  time.Now().UTC(),
		Symbol:     symbol,
		Action:     "init",
		AmountUSDT: initUSDT,
		AmountBTC:  initBTC,
		Price:      priceStr,
		Reason:     "Initial account snapshot",
	}

	appendToHistory(initRecord)
	fmt.Printf("✅ Initialized trading history with USDT=%.2f, BTC=%.6f\n", initUSDT, initBTC)
	return initUSDT, initBTC, nil
}
