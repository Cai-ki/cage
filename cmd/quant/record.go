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
var records []PerformanceRecord

// PerformanceRecord 记录每日/每次交易的绩效
type PerformanceRecord struct {
	Date           string  `json:"date"`            // 日期：2025-11-11
	InitialBalance float64 `json:"initial_balance"` // 初始余额（USDT）
	CurrentBalance float64 `json:"current_balance"` // 当前余额（USDT）
	CumulativePnL  float64 `json:"cumulative_pnl"`  // 累计盈亏（USDT）
	CumulativeROI  float64 `json:"cumulative_roi"`  // 累计收益率（%）
	TotalTrades    int     `json:"total_trades"`    // 总交易次数
	CurrentPrice   float64 `json:"current_price"`

	Prompt    string `json:"prompt"`
	Decision  string `json:"decision"`
	ToolCalls string `json:"tool_calls"`
}

func LoadOrCreatePerformanceRecord() ([]PerformanceRecord, error) {
	recordMutex.Lock()
	defer recordMutex.Unlock()

	data, err := os.ReadFile(performanceFile)
	if err == nil {
		var loadedRecords []PerformanceRecord
		if json.Unmarshal(data, &loadedRecords) == nil && len(loadedRecords) > 0 {
			records = loadedRecords
			return loadedRecords, nil
		}
	}

	// 获取真实账户权益（USDT + 浮盈）
	initialEquity, err := getAccountEquity()
	if err != nil {
		return nil, err
	}

	price, _ := quant.FuturesGetTickerPrice("BTCUSDT")
	currentPrice := sugar.Must(sugar.StrToT[float64](price))

	firstRecord := PerformanceRecord{
		Date:           time.Now().Format("2006-01-02 15:04:05"),
		InitialBalance: initialEquity,
		CurrentBalance: initialEquity,
		CumulativePnL:  0.0,
		CumulativeROI:  0.0,
		TotalTrades:    0, // 初始化为 0
		CurrentPrice:   currentPrice,
		Prompt:         "",
		Decision:       "",
		ToolCalls:      "",
	}

	records = []PerformanceRecord{firstRecord}
	savePerformanceRecordsWithoutLock(records)

	return records, nil
}

// 记录交易并增加总交易次数
func RecordTrade(prompt, decision, toolcalls string, add bool) error {
	recordMutex.Lock()
	defer recordMutex.Unlock()

	// 加载现有记录
	data, err := os.ReadFile(performanceFile)
	if err == nil {
		var loadedRecords []PerformanceRecord
		if json.Unmarshal(data, &loadedRecords) == nil && len(loadedRecords) > 0 {
			records = loadedRecords
		}
	}

	var lastRecord PerformanceRecord
	if len(records) > 0 {
		lastRecord = records[len(records)-1]
	} else {
		// 如果没有记录，创建初始记录
		initialEquity, err := getAccountEquity()
		if err != nil {
			return err
		}
		price, _ := quant.FuturesGetTickerPrice("BTCUSDT")
		currentPrice := sugar.Must(sugar.StrToT[float64](price))

		lastRecord = PerformanceRecord{
			Date:           time.Now().Format("2006-01-02 15:04:05"),
			InitialBalance: initialEquity,
			CurrentBalance: initialEquity,
			CumulativePnL:  0.0,
			CumulativeROI:  0.0,
			TotalTrades:    0,
			CurrentPrice:   currentPrice,
			Prompt:         "",
			Decision:       "",
			ToolCalls:      "",
		}
		records = []PerformanceRecord{lastRecord}
	}

	// 创建新记录，基于最新记录
	newRecord := lastRecord // 复制最新记录
	newRecord.Date = time.Now().Format("2006-01-02 15:04:05")

	// 增加总交易次数
	if add {
		newRecord.TotalTrades++
	}

	// 重新计算当前权益（因为可能刚交易完）
	currentEquity, err := getAccountEquity()
	if err != nil {
		return err
	}

	newRecord.CurrentBalance = currentEquity
	newRecord.CumulativePnL = currentEquity - newRecord.InitialBalance
	newRecord.CumulativeROI = (newRecord.CumulativePnL / newRecord.InitialBalance) * 100

	price, _ := quant.FuturesGetTickerPrice("BTCUSDT")
	newRecord.CurrentPrice = sugar.Must(sugar.StrToT[float64](price))
	newRecord.Prompt = prompt
	newRecord.Decision = decision
	newRecord.ToolCalls = toolcalls

	// 添加新记录到数组
	records = append(records, newRecord)

	savePerformanceRecordsWithoutLock(records)
	return nil
}

func GetPerformanceRecord() *PerformanceRecord {
	recordMutex.RLock()
	defer recordMutex.RUnlock()

	// 尝试从文件加载最新记录
	data, err := os.ReadFile(performanceFile)
	if err == nil {
		var loadedRecords []PerformanceRecord
		if json.Unmarshal(data, &loadedRecords) == nil && len(loadedRecords) > 0 {
			records = loadedRecords
		}
	}

	if len(records) > 0 {
		return &records[len(records)-1] // 返回最后一个记录（最新的）
	}
	return nil
}

func GetAllPerformanceRecords() []PerformanceRecord {
	recordMutex.RLock()
	defer recordMutex.RUnlock()

	// 尝试从文件加载最新记录
	data, err := os.ReadFile(performanceFile)
	if err == nil {
		var loadedRecords []PerformanceRecord
		if json.Unmarshal(data, &loadedRecords) == nil {
			records = loadedRecords
		}
	}

	return records
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

func savePerformanceRecordsWithoutLock(records []PerformanceRecord) {
	data, _ := json.MarshalIndent(records, "", "  ")
	os.WriteFile(performanceFile, data, 0644)
}

func savePerformanceRecords(records []PerformanceRecord) {
	recordMutex.Lock()
	defer recordMutex.Unlock()
	savePerformanceRecordsWithoutLock(records)
}

func formatPerformanceSummary() string {
	record := GetPerformanceRecord()
	if record == nil {
		return "Performance data unavailable"
	}

	return fmt.Sprintf(
		"Strategy Performance (latest record):\n"+
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
