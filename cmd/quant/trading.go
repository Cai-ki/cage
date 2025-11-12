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
	_, _ = LoadOrCreatePerformanceRecord()
	do.Do(StartServer)

	prompt := BuildPrompt(symbol)

	log.Println(prompt)
	rsp, err := llm.CompletionByParams(llm.SystemMessage(prompt), llm.ToolsByJson(mcpString))
	if err != nil {
		return err
	}

	log.Println(rsp.Content)
	if len(rsp.ToolCalls) > 0 {
		_, err = mcp.ExecuteToolCalls(rsp)
		if err != nil {
			return err
		}
		log.Println("Execute tool calls success")

		toolCallsStr := ""
		for _, v := range rsp.ToolCalls {
			toolCallsStr += v.RawJSON() + "\n"
		}

		log.Println(toolCallsStr)
		err := RecordTrade(prompt, rsp.Content, toolCallsStr, globalMemory, len(rsp.ToolCalls) > 1)
		if err != nil {
			return err
		}

	} else {
		err := RecordTrade(prompt, rsp.Content, "", globalMemory, false)
		if err != nil {
			return err
		}
	}
	log.Println("Record trade success")
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

func getTradingContext(symbol string) (timeStr, sym, price, usdtBal, positionInfo, multiIndicator, rate string) {
	timeStr = getTime()
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

	// // 平衡的中频配置 - 推荐使用
	// config := &quant.IndicatorConfig{
	// 	EMAs:       []int{12, 26, 50}, // 短中结合
	// 	MAs:        []int{20, 60},     // 实用周期
	// 	RSI:        []int{14},         // 标准RSI
	// 	MACD:       true,
	// 	Stochastic: []int{14, 3}, // 标准随机
	// 	ATR:        []int{14},    // 标准波动率
	// 	Bollinger:  []int{20, 2}, // 标准布林带
	// }

	// // K线数量保持你的原配置
	// k5m, err := quant.FuturesGetKlines(symbol, "5m", 150)
	// k15m, err := quant.FuturesGetKlines(symbol, "15m", 120)
	// k1h, err := quant.FuturesGetKlines(symbol, "1h", 100)

	// 激进中高频配置 - 极速响应
	config := &quant.IndicatorConfig{
		EMAs:       []int{5, 13, 21}, // 超短期EMA
		MAs:        []int{8, 21},     // 快速移动平均
		RSI:        []int{6, 14},     // 超快速RSI
		MACD:       true,
		Stochastic: []int{5, 3},  // 超敏感随机
		ATR:        []int{7, 14}, // 快速波动率
		Bollinger:  []int{13, 2}, // 极窄布林带
	}

	// K线数据更少，只关注最近变化
	k5m, err := quant.FuturesGetKlines(symbol, "5m", 50)
	k15m, err := quant.FuturesGetKlines(symbol, "15m", 40)
	k1h, err := quant.FuturesGetKlines(symbol, "1h", 30)

	calculator := quant.NewIndicatorCalculator(config)

	timeframeData := map[string][]*futures.Kline{
		"5m":  k5m,
		"15m": k15m,
		"1h":  k1h,
	}

	multiIndicator = calculator.CalculateMultiTimeframe("BTCUSDT", timeframeData).ToSimpleString()

	fundingRate, nextFundingTime, err := quant.FuturesGetCurrentFundingRate(sym)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	nextFundingTimeStr := time.Unix(nextFundingTime/1000, 0).In(loc).Format("2006-01-02 15:04:05")
	// 输出：2024-01-16T00:00:00+08:00 （北京时间，比UTC快8小时）
	makerRate, takerRate, err := quant.FuturesGetFeeRateForSymbol("BTCUSDT")
	rate = fmt.Sprintf("Fee rates - Maker: %s, Taker: %s\nCurrent Funding Rate: %s\nNext Funding Time: %s\n", makerRate, takerRate, fundingRate, nextFundingTimeStr)
	return
}

func BuildPrompt(symbol string) string {
	tag := "当前位于测试环境，策略允许激进，允许高风险高收益操作。"

	perfSummary := formatPerformanceSummary()

	record := GetPerformanceRecord()
	totalTrades := record.TotalTrades

	timeStr, sym, price, usdtBal, positionInfo, multiIndicator, rate := getTradingContext(symbol)

	return fmt.Sprintf(TradingAgentPromptTemplate,
		tag, timeStr, sym, price, usdtBal, totalTrades, perfSummary, positionInfo, multiIndicator, rate)
}

const TradingAgentPromptTemplate = `
[%s]

你是一个专业的加密货币期货交易员，具备丰富的市场分析和风险管理经验。请根据以下实时上下文信息，制定并执行合理的交易决策。

CURRENT CONTEXT:
- Time: %s (UTC+8 / Beijing Time)
- Symbol: %s
- Current price: %s USDT
- USDT balance: %s
- Total trades: %d

STRATEGY PERFORMANCE:

%s

当前持仓状态:

%s

当前市场:

%s

资金费用以及手续费说明:

%s

## 你的职责：
1. 分析多时间框架下的价格趋势、动量与成交量变化
2. 结合**当前持仓状态**（方向、成本、盈亏、杠杆）与交易成本，评估风险敞口
3. 制定清晰的入场、出场或持仓调整策略

## 可用函数：
- **futures_buy_market(symbol, quantity)**  
  -> 开多：当判断价格将上涨且符合策略时使用

- **futures_sell_market(symbol, quantity)**  
  -> 开空或主动平多：当判断价格将下跌，或需减仓多头时使用

- **futures_close_position(symbol)**  
  -> 平仓：无论当前持多或持空，自动全部平掉该标的仓位（使用 ReduceOnly 模式)

- **save_memory(memory)**  
  -> 记忆化：将需要持久化的记忆存储下来，记忆会传入下次调用时的上下文中（必须调用）

## 输出要求：
- 先简要总结市场状态、当前持仓风险及手续费影响
- 明确说明交易意图（开多 / 开空 / 平仓 / 暂不交易）
- 若信号微弱、盈亏比不足或风险过高，请明确说明"暂不交易"并解释原因

## 调用要求：
- 无论是否进行交易，你**必须调用一次** save_memory 函数存储记忆，记忆可以包含此时市场摘要，此次决策的原因，对当前整体局势的分析，长远的战略等（根据需要可调整），以自然段格式展现。
- 如需下单，请直接调用上述函数（不要调用逻辑矛盾）
- 不要虚构函数或参数，严格遵循接口定义

请基于以上信息，做出专业、审慎且可执行的交易决策。
`
