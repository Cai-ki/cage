package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/Cai-ki/cage/llm"
	"github.com/Cai-ki/cage/llm/mcp"
	"github.com/Cai-ki/cage/quant"
	"github.com/adshao/go-binance/v2/futures"
)

func RunTradingStep(symbol string) error {
	prompt := BuildPrompt(symbol)

	log.Println(prompt)
	rsp, err := llm.CompletionByParams(llm.SystemMessage(prompt), llm.ToolsByJson(mcpString))
	if err != nil {
		return err
	}

	log.Println(rsp.Content)
	if len(rsp.ToolCalls) > 0 {
		log.Println(rsp.ToolCalls[0].RawJSON())
		_, err = mcp.ExecuteToolCalls(rsp)
		if err != nil {
			return err
		}

		err := RecordTrade()
		if err != nil {
			return err
		}
	}

	return nil
}

// 格式化单组 K 线为一行文本：O:H:L:C:V
func formatKline(k *futures.Kline) string {
	return fmt.Sprintf("%s:%s:%s:%s:%s", k.Open, k.High, k.Low, k.Close, k.Volume)
}

// 格式化一组 K 线为多行文本
func formatKlines(klines []*futures.Kline, err error) string {
	if err != nil {
		return fmt.Sprint("error when get kline: ", err)
	}
	if len(klines) == 0 {
		return "No data"
	}
	var lines []string
	for _, k := range klines {
		lines = append(lines, formatKline(k))
	}
	return strings.Join(lines, "\n")
}

func calculatePnLPercentage(pos *futures.PositionRisk) string {
	if pos.EntryPrice == "0" || pos.PositionAmt == "0" {
		return "0.00"
	}

	entry, err1 := strconv.ParseFloat(pos.EntryPrice, 64)
	pnl, err2 := strconv.ParseFloat(pos.UnRealizedProfit, 64)

	if err1 != nil || err2 != nil {
		return "0.00"
	}

	// 计算百分比：PnL / (Entry * PositionAmt) * 100
	if entry != 0 && pos.PositionAmt != "0" {
		posAmt, err := strconv.ParseFloat(pos.PositionAmt, 64)
		if err != nil {
			return "0.00"
		}

		if posAmt != 0 {
			absPosAmt := math.Abs(posAmt)
			percentage := (pnl / (entry * absPosAmt)) * 100
			return fmt.Sprintf("%.2f", percentage)
		}
	}
	return "0.00"
}

func getTradingContext(symbol string) (timeStr, sym, price, usdtBal, positionInfo, k5m, k15m, k1h, k4h string) {
	timeStr = time.Now().Format("2006-01-02 15:04:05")
	sym = symbol

	// 获取当前价格
	price, _ = quant.FuturesGetTickerPrice(symbol)

	// 获取余额
	usdtBal, _ = quant.FuturesGetBalance("USDT")

	// 获取持仓信息
	position, err := quant.FuturesGetPosition(symbol)
	if err != nil || position == nil || position.PositionAmt == "0" {
		positionInfo = "No open position"
	} else {
		// 解析持仓数量，确定实际方向
		posAmt, parseErr := strconv.ParseFloat(position.PositionAmt, 64)
		var actualDirection string
		if parseErr != nil || posAmt == 0 {
			actualDirection = "NEUTRAL"
		} else if posAmt > 0 {
			actualDirection = "LONG"
		} else {
			actualDirection = "SHORT"
		}

		pnlPercentage := calculatePnLPercentage(position)
		positionInfo = fmt.Sprintf(
			"Current Position: %s %s %s (in %s mode)\n"+
				"Entry Price: %s USDT\n"+
				"Mark Price: %s USDT\n"+
				"Unrealized PnL: %s USDT (%s%%)\n"+
				"Leverage: %sx\n"+
				"Position Side: %s (actual: %s)\n"+ // 这里需要2个参数
				"Liquidation Price: %s USDT\n"+
				"Margin Type: %s\n"+
				"Notional Value: %s USDT",
			actualDirection, position.PositionAmt, position.Symbol, position.PositionSide, // 1-4
			position.EntryPrice, position.MarkPrice, position.UnRealizedProfit, pnlPercentage, // 5-8
			position.Leverage, position.PositionSide, actualDirection, // 9-11 (Position Side 和 actual 需要分开传)
			position.LiquidationPrice, position.MarginType, position.Notional) // 12-14
	}

	// 获取 K 线数据
	k5m = formatKlines(quant.FuturesGetKlines(symbol, "5m", 12))
	k15m = formatKlines(quant.FuturesGetKlines(symbol, "15m", 12))
	k1h = formatKlines(quant.FuturesGetKlines(symbol, "1h", 12))
	k4h = formatKlines(quant.FuturesGetKlines(symbol, "4h", 6))

	return
}

func BuildPrompt(symbol string) string {
	dataSource := "Binance Testnet（当前位于测试环境，请策略必须激进。不注重亏损。必须开多 / 开空 / 平仓。）"

	perfSummary := formatPerformanceSummary()
	record, _ := LoadOrCreatePerformanceRecord()
	totalTrades := record.TotalTrades

	timeStr, sym, price, usdtBal, positionInfo, k5m, k15m, k1h, k4h := getTradingContext(symbol)

	return fmt.Sprintf(TradingAgentPromptTemplate,
		timeStr, sym, price, usdtBal, totalTrades, dataSource, perfSummary, positionInfo, k5m, k15m, k1h, k4h)
}

const TradingAgentPromptTemplate = `
你是一个专业的加密货币期货交易员，具备丰富的市场分析和风险管理经验。请根据以下实时上下文信息，制定并执行合理的交易决策。

CURRENT CONTEXT:
- Time: %s (UTC+8 / Beijing Time)
- Symbol: %s
- Current price: %s USDT
- USDT balance: %s
- Total trades: %d
- Data source: %s

STRATEGY PERFORMANCE:
%s

POSITION INFO:
%s

MARKET DATA FORMAT: Each line = Open:High:Low:Close:Volume
Focus on Close price trend and Volume changes.

[5m — last 12 candles]
%s

[15m — last 12 candles]
%s

[1h — last 12 candles]
%s

[4h — last 6 candles]
%s

## 资金费用说明（至关重要！）
- 每 8 小时结算一次资金费用（北京时间 08:00/16:00/24:00）
- 当前持仓方向可能需支付/收取资金费，请评估持仓成本

## 你的职责：
1. 分析多时间框架下的价格趋势、动量与成交量变化
2. 结合**当前持仓状态**（方向、成本、盈亏、杠杆）与交易成本，评估风险敞口
3. 制定清晰的入场、出场或持仓调整策略
4. **特别注意：避免与现有持仓冲突的操作（如已有 LONG，再开 LONG 会增加风险）**

## 可用交易工具（仅限以下三个函数）：
- **futures_buy_market(symbol, quantity)**  
  -> 开多：当判断价格将上涨且符合策略时使用（Taker，手续费 0.04%%）

- **futures_sell_market(symbol, quantity)**  
  -> 开空或主动平多：当判断价格将下跌，或需减仓多头时使用（Taker，手续费 0.04%%）

- **futures_close_position(symbol)**  
  -> 平仓：无论当前持多或持空，自动全部平掉该标的仓位（使用 ReduceOnly 模式，手续费 0.04%%）

## 输出要求：
- 先简要总结市场状态、当前持仓风险及手续费影响
- 明确说明交易意图（开多 / 开空 / 平仓 / 暂不交易）
- 如需下单，请直接调用上述函数（仅调用一次）
- 不要虚构函数或参数，严格遵循接口定义
- 若信号微弱、盈亏比不足或风险过高，请明确说明"暂不交易"并解释原因

请基于以上信息，做出专业、审慎且可执行的交易决策。`
