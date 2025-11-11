package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/Cai-ki/cage/quant"
	"github.com/Cai-ki/cage/sugar"
)

// 全局变量或持久化存储路径
const performanceFile = "trading_performance.json"

var recordMutex sync.RWMutex
var do sync.Once
var globalRecord *PerformanceRecord

// PerformanceRecord 记录每日/每次交易的绩效
type PerformanceRecord struct {
	Date           string  `json:"date"`            // 日期：2025-11-11
	InitialBalance float64 `json:"initial_balance"` // 初始余额（USDT）
	CurrentBalance float64 `json:"current_balance"` // 当前余额（USDT）
	CumulativePnL  float64 `json:"cumulative_pnl"`  // 累计盈亏（USDT）
	CumulativeROI  float64 `json:"cumulative_roi"`  // 累计收益率（%）
	TotalTrades    int     `json:"total_trades"`    // 总交易次数
	CurrentPrice   float64 `json:"current_price"`
}

func LoadOrCreatePerformanceRecord() (*PerformanceRecord, error) {
	recordMutex.Lock()
	defer recordMutex.Unlock()

	data, err := os.ReadFile(performanceFile)
	if err == nil {
		var record PerformanceRecord
		if json.Unmarshal(data, &record) == nil && record.Date == time.Now().Format("2006-01-02") {
			globalRecord = &record
			return &record, nil
		}
	}

	// 获取真实账户权益（USDT + 浮盈）
	initialEquity, err := getAccountEquity()
	if err != nil {
		return nil, err
	}

	price, _ := quant.FuturesGetTickerPrice("BTCUSDT")
	currentPrice := sugar.Must(sugar.StrToT[float64](price))

	record := &PerformanceRecord{
		Date:           time.Now().Format("2006-01-02"),
		InitialBalance: initialEquity,
		CurrentBalance: initialEquity,
		CumulativePnL:  0.0,
		CumulativeROI:  0.0,
		TotalTrades:    0, // 初始化为 0
		CurrentPrice:   currentPrice,
	}

	globalRecord = record
	savePerformanceRecordWithoutLock(record)

	return record, nil
}

// 记录交易并增加总交易次数
func RecordTrade() error {
	recordMutex.Lock()
	defer recordMutex.Unlock()

	record, err := LoadOrCreatePerformanceRecord()
	if err != nil {
		return err
	}

	// 增加总交易次数
	record.TotalTrades++

	// 重新计算当前权益（因为可能刚交易完）
	currentEquity, err := getAccountEquity()
	if err != nil {
		return err
	}

	record.CurrentBalance = currentEquity
	record.CumulativePnL = currentEquity - record.InitialBalance
	record.CumulativeROI = (record.CumulativePnL / record.InitialBalance) * 100

	price, _ := quant.FuturesGetTickerPrice("BTCUSDT")
	record.CurrentPrice = sugar.Must(sugar.StrToT[float64](price))

	globalRecord = record
	savePerformanceRecordWithoutLock(record)
	return nil
}

func GetPerformanceRecord() *PerformanceRecord {
	recordMutex.RLock()
	defer recordMutex.RUnlock()

	if globalRecord == nil {
		record, err := LoadOrCreatePerformanceRecord()
		if err != nil {
			return nil
		}
		return record
	}
	return globalRecord
}

func getAccountEquity() (float64, error) {
	// 1. 获取 USDT 可用余额
	usdtBalStr, err := quant.FuturesGetBalance("USDT")
	if err != nil {
		return 0, fmt.Errorf("failed to get USDT balance: %v", err)
	}
	availableUsdt, _ := strconv.ParseFloat(usdtBalStr, 64)

	// 2. 获取当前持仓的未实现盈亏（以 USDT 计）
	position, err := quant.FuturesGetPosition("BTCUSDT") // 替换为你的交易对
	if err != nil {
		return availableUsdt, nil // 如果无法获取持仓，至少返回可用余额
	}

	unrealizedPnL, _ := strconv.ParseFloat(position.UnRealizedProfit, 64)

	// 3. 总权益 = 可用 USDT + 未实现盈亏
	totalEquity := availableUsdt + unrealizedPnL

	return totalEquity, nil
}

func savePerformanceRecordWithoutLock(record *PerformanceRecord) {
	data, _ := json.MarshalIndent(record, "", "  ")
	os.WriteFile(performanceFile, data, 0644)
}

func savePerformanceRecord(record *PerformanceRecord) {
	recordMutex.Lock()
	defer recordMutex.Unlock()
	savePerformanceRecordWithoutLock(record)
}

func formatPerformanceSummary() string {
	record, err := LoadOrCreatePerformanceRecord()
	if err != nil {
		return "Performance data unavailable"
	}

	today := time.Now().Format("2006-01-02")
	if record.Date != today {
		// 应该不会发生，因为 LoadOrCreate 会处理
		return "Performance data outdated"
	}

	return fmt.Sprintf(
		"Strategy Performance (since start):\n"+
			"- Initial Capital: %.2f USDT\n"+
			"- Current Balance: %.2f USDT\n"+
			"- Cumulative PnL: %.2f USDT\n"+
			"- Cumulative ROI: %.2f%%\n"+
			"- Total Trades: %d",
		record.InitialBalance,
		record.CurrentBalance,
		record.CumulativePnL,
		record.CumulativeROI,
		record.TotalTrades,
	)
}
