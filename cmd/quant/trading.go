package main

import (
	"fmt"
	"log"
	"math"
	"strconv"
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
		err := RecordTrade(prompt, rsp.Content, toolCallsStr, len(rsp.ToolCalls) > 1)
		if err != nil {
			return err
		}

	} else {
		err := RecordTrade(prompt, rsp.Content, "", false)
		if err != nil {
			return err
		}
	}
	log.Println("Record trade success")
	return nil
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
			actualDirection = "中性"
		} else if posAmt > 0 {
			actualDirection = "多头"
		} else {
			actualDirection = "空头"
		}

		pnlPercentage := calculatePnLPercentage(position)
		// positionInfo = fmt.Sprintf(
		// 	"当前持仓状态: %s %s %s (%s 模式)\n"+
		// 		"- 开仓均价: %s USDT\n"+
		// 		"- 标记价格: %s USDT\n"+
		// 		"- 未实现盈亏: %s USDT (%s%%)\n"+
		// 		"- 杠杆倍数: %s倍\n"+
		// 		"- 持仓方向: %s (实际: %s)\n"+
		// 		"- 强平价格: %s USDT\n"+
		// 		"- 保证金模式: %s\n"+
		// 		"- 名义价值: %s USDT",
		// 	actualDirection, position.PositionAmt, position.Symbol, position.PositionSide, position.EntryPrice, position.MarkPrice, position.UnRealizedProfit, pnlPercentage, position.Leverage, position.PositionSide, actualDirection, position.LiquidationPrice, position.MarginType, position.Notional)

		// 计算指标
		// 计算指标
		notional, _ := strconv.ParseFloat(position.Notional, 64)
		leverage, _ := strconv.ParseFloat(position.Leverage, 64)
		markPrice, _ := strconv.ParseFloat(position.MarkPrice, 64)
		liquidationPrice, _ := strconv.ParseFloat(position.LiquidationPrice, 64)
		unrealizedProfit, _ := strconv.ParseFloat(position.UnRealizedProfit, 64)
		pnlPercent, _ := strconv.ParseFloat(pnlPercentage, 64)

		// 占用保证金（正确）
		marginUsed := notional / leverage

		// 总资金 = USDT余额（修正这里！）
		usdtBalance, _ := strconv.ParseFloat(usdtBal, 64)
		// totalEquity := usdtBalance + notional + unrealizedProfit  // ❌ 错误的

		// 仓位占比（修正）
		positionRatio := (marginUsed / usdtBalance) * 100 // ✅ 用总资金，不是总权益

		// 距离强平百分比（保持原样，这个计算正确）
		liquidationDistance := 0.0
		if actualDirection == "多头" {
			liquidationDistance = ((markPrice - liquidationPrice) / markPrice) * 100
		} else if actualDirection == "空头" {
			liquidationDistance = ((liquidationPrice - markPrice) / markPrice) * 100
		} else {
			liquidationDistance = 100.0 // 无持仓
		}

		// 风险等级（保持原样）
		riskLevel := "低风险"
		if liquidationDistance < 5 {
			riskLevel = "极高风险"
		} else if liquidationDistance < 10 {
			riskLevel = "高风险"
		} else if liquidationDistance < 15 {
			riskLevel = "中等风险"
		}

		// 持仓方向映射
		var positionSideMap = map[string]string{
			"BOTH":  "双向",
			"LONG":  "多头",
			"SHORT": "空头",
		}

		// 保证金模式映射
		var marginTypeMap = map[string]string{
			"cross":    "全仓",
			"isolated": "逐仓",
		}

		// 在使用前先进行映射转换
		mappedPositionSide := positionSideMap[position.PositionSide] // "LONG" -> "多头"
		mappedMarginType := marginTypeMap[position.MarginType]       // "cross" -> "全仓"

		positionInfo = fmt.Sprintf(
			"- 持仓状态: %s %s %s (%s 模式)\n"+
				"- 开仓均价: %s USDT\n"+
				"- 标记价格: %s USDT\n"+
				"- 未实现盈亏: %.2f USDT (%.2f%%)\n"+ // 改为%.2f格式，更整洁
				"- 杠杆倍数: %.0f倍\n"+
				"- 持仓方向: %s (实际: %s)\n"+
				"- 强平价格: %s USDT\n"+
				"- 保证金模式: %s\n"+
				"- 名义价值: %s USDT\n"+
				"- 风险与仓位\n"+
				"- 占用保证金: %.2f USDT\n"+
				"- 总资金: %.2f USDT\n"+
				"- 仓位占比: %.1f%%\n"+
				"- 距离强平: %.1f%%\n"+
				"- 风险等级: %s",
			actualDirection, position.PositionAmt, position.Symbol, mappedPositionSide,
			position.EntryPrice, position.MarkPrice, unrealizedProfit, pnlPercent, // 使用解析后的变量
			leverage, mappedPositionSide, actualDirection,
			position.LiquidationPrice, mappedMarginType, position.Notional,
			marginUsed, usdtBalance, positionRatio, liquidationDistance, riskLevel)
	}

	// 平衡的中频配置 - 推荐使用
	config := &quant.IndicatorConfig{
		EMAs:       []int{12, 26, 50}, // 短中结合
		MAs:        []int{20, 60},     // 实用周期
		RSI:        []int{14},         // 标准RSI
		MACD:       true,
		Stochastic: []int{14, 3}, // 标准随机
		ATR:        []int{14},    // 标准波动率
		Bollinger:  []int{20, 2}, // 标准布林带
	}

	// K线数量保持你的原配置
	k3m, err := quant.FuturesGetKlines(symbol, "3m", 150)
	k5m, err := quant.FuturesGetKlines(symbol, "5m", 150)
	k15m, err := quant.FuturesGetKlines(symbol, "15m", 120)
	k1h, err := quant.FuturesGetKlines(symbol, "1h", 100)

	// // 激进中高频配置 - 极速响应
	// config := &quant.IndicatorConfig{
	// 	EMAs:       []int{5, 13, 21}, // 超短期EMA
	// 	MAs:        []int{8, 21},     // 快速移动平均
	// 	RSI:        []int{6, 14},     // 超快速RSI
	// 	MACD:       true,
	// 	Stochastic: []int{5, 3},  // 超敏感随机
	// 	ATR:        []int{7, 14}, // 快速波动率
	// 	Bollinger:  []int{13, 2}, // 极窄布林带
	// }

	// // K线数据更少，只关注最近变化
	// k5m, err := quant.FuturesGetKlines(symbol, "5m", 50)
	// k15m, err := quant.FuturesGetKlines(symbol, "15m", 40)
	// k1h, err := quant.FuturesGetKlines(symbol, "1h", 30)

	calculator := quant.NewIndicatorCalculator(config)

	timeframeData := map[string][]*futures.Kline{
		"3m":  k3m,
		"5m":  k5m,
		"15m": k15m,
		"1h":  k1h,
	}

	// 按时间周期顺序输出（从大到小）
	timeframeOrder := []string{"1h", "15m", "5m", "3m"}

	multiIndicator = calculator.CalculateMultiTimeframe("BTCUSDT", timeframeData).ToSimpleString(timeframeOrder)

	fundingRate, nextFundingTime, err := quant.FuturesGetCurrentFundingRate(sym)
	loc, _ := time.LoadLocation("Asia/Shanghai")
	nextFundingTimeStr := time.Unix(nextFundingTime/1000, 0).In(loc).Format("2006-01-02 15:04:05")
	makerRate, takerRate, err := quant.FuturesGetFeeRateForSymbol("BTCUSDT")
	rate = fmt.Sprintf("- 手续费率 - Maker(挂单): %s, Taker(吃单): %s\n- 当前资金费率: %s\n- 下次资金费率时间: %s\n",
		makerRate, takerRate, fundingRate, nextFundingTimeStr)
	return
}

func BuildPrompt(symbol string) string {
	perfSummary := formatPerformanceSummary()

	record := GetPerformanceRecord()
	totalTrades := record.TotalTrades

	timeStr, sym, price, usdtBal, positionInfo, multiIndicator, rate := getTradingContext(symbol)

	return fmt.Sprintf(TradingAgentPromptTemplate,
		timeStr, sym, price, usdtBal, totalTrades, perfSummary, positionInfo, multiIndicator, rate)
}

const TradingAgentPromptTemplate = `
---
**你是一个激进但专业的加密货币期货交易员，擅长抓住市场机会并主动管理风险。你的目标是最大化资金利用率，在控制风险的前提下积极交易。**
---

## 📊 账户与市场概览
- **当前时间**: %s (UTC+8 / Beijing Time)
- **交易标的**: %s
- **实时价格**: %s USDT
- **可用保证金**: %s USDT
- **历史交易次数**: %d

## 🔍 上次决策回顾
%s

## 💼 当前持仓状态
%s

## 📈 技术指标分析
%s

## 💰 费率信息
%s

## ⚖️ 仓位管理规则

### 计算公式
- **仓位占比** = (持仓保证金 ÷ 总权益) × 100%%
- **持仓保证金** = 持仓价值 ÷ 杠杆倍数
- **持仓价值** = 持仓数量 × 当前价格
- **总权益** = 可用保证金 + 持仓价值 + 未实现盈亏

### 仓位目标
- 🟢 **理想仓位**: 30%%-60%%
- 🟡 **最小仓位**: 10%%
- 🟠 **警戒仓位**: >70%%
- 🔴 **当前评估**: 如仓位<5%%，必须积极加仓

## 🎯 交易策略指导

### 趋势确认原则
- ✅ **3个周期指标一致** → 重仓参与（40%%-60%%）
- ☑️ **2个周期指标一致** → 中等仓位（25%%-40%%）
- ⚠️ **仅1个周期有信号** → 轻仓试探（15%%-25%%）
- ❌ **无明确信号** → 保持现有仓位或微调

### 持仓管理纪律
- 📈 **盈利<2.5%%时**: 坚决持有，不加不减
- 💰 **盈利2.5%%-5%%时**: 可考虑部分止盈，但保留至少50%%仓位
- 🎉 **盈利>5%%时**: 逐步止盈，但保持20%%以上仓位参与趋势
- 📉 **浮亏时**: 基于技术指标判断是否加仓摊薄成本，不要恐慌性平仓

## 📝 输出要求（严格执行）
1. **本次决策的技术分析依据**
2. **仓位管理逻辑和风险评估**
3. **明确说明交易意图**（开多 / 开空 / 平仓 / 暂不交易）
4. **若信号微弱、盈亏比不足或风险过高**，请明确说明"暂不交易"并解释原因

## 🔧 调用要求（严格执行）

### 函数调用规则
- 🚨 **必须且仅能**通过函数调用来执行操作
- 🔄 **每次响应必须包含 exactly two 工具调用**：
  1. 一个交易类操作（"futures_buy_market" / "futures_sell_market" / "futures_close_position"）
  2. **必须调用** "save_memory(memory)"，传入完整分析与记忆

### 记忆保存要求（零容忍）
- 💾 **必须调用 save_memory**，否则系统无法学习优化
- 📋 **记忆内容必须包含**：
  1. 后续价格预期和具体操作计划
  2. 对上次决策的反思（如有）

### 具体调用场景

| 场景 | 操作流程 |
|------|----------|
| **需要交易时** | 1. 交易操作函数<br>2. "save_memory" (必须) |
| **无需交易时** | 1. 跳过交易操作<br>2. "save_memory" (必须) |

### 🚫 禁止行为
- ❌ 只调用交易操作不调用 "save_memory"
- ❌ 不调用任何函数
- ❌ 记忆内容空泛（如仅"看好上涨"）

> ⚠️ **重要提醒**: 
> - 如果你认为无需交易，请**仅调用 "save_memory"**
> - 如果你需要交易，请**同时调用交易函数和 "save_memory"**
`

// const TradingAgentPromptTemplate = `
// ---
// **你是一个激进但专业的加密货币期货交易员，擅长抓住市场机会并主动管理风险。你的目标是最大化资金利用率，在控制风险的前提下积极交易。**
// ---

// # 账户与市场概览:
// - 当前时间: %s (UTC+8 / Beijing Time)
// - 交易标的: %s
// - 实时价格: %s USDT
// - 可用保证金: %s USDT
// - 历史交易次数: %d

// # 上次决策回顾:
// %s

// # 当前持仓状态:
// %s

// # 技术指标分析:
// %s

// # 费率信息:
// %s

// # 仓位管理规则：
// **计算公式**：
// - 仓位占比 = (持仓保证金 ÷ 总权益) × 100%%
// - 持仓保证金 = 持仓价值 ÷ 杠杆倍数
// - 持仓价值 = 持仓数量 × 当前价格
// - 总权益 = 可用保证金 + 持仓价值 + 未实现盈亏

// **仓位目标**：
// - 理想仓位：30%%-60%%
// - 最小仓位：10%%
// - 警戒仓位：>70%%
// - 当前评估：如仓位<5%%，必须积极加仓

// # 交易策略指导：

// **趋势确认原则：**
// - 3个周期指标一致 → 重仓参与（40%%-60%%）
// - 2个周期指标一致 → 中等仓位（25%%-40%%）
// - 仅1个周期有信号 → 轻仓试探（15%%-25%%）
// - 无明确信号 → 保持现有仓位或微调

// **持仓管理纪律：**
// - 盈利<2.5%%时：坚决持有，不加不减
// - 盈利2.5%%-5%%时：可考虑部分止盈，但保留至少50%%仓位
// - 盈利>5%%时：逐步止盈，但保持20%%以上仓位参与趋势
// - 浮亏时：基于技术指标判断是否加仓摊薄成本，不要恐慌性平仓

// # 输出要求（严格执行）：
// - 本次决策的技术分析依据
// - 仓位管理逻辑和风险评估
// - 明确说明交易意图（开多 / 开空 / 平仓 / 暂不交易）
// - 若信号微弱、盈亏比不足或风险过高，请明确说明"暂不交易"并解释原因

// # 调用要求（严格执行）：

// ## 函数调用规则
// - 你**必须且仅能**通过函数调用来执行操作
// - **每次响应必须包含 exactly two 工具调用**：
//   1. 一个交易类操作（futures_buy_market / futures_sell_market / futures_close_position）
//   2. **必须调用 save_memory(memory)**，传入完整分析与记忆

// ## 记忆保存要求（零容忍）
// - **必须调用 save_memory**，否则系统无法学习优化
// - **记忆内容必须包含**：
//   1. 后续价格预期和具体操作计划
//   2. 对上次决策的反思（如有）

// ## 具体调用场景：

// **需要交易时**：
// 1. futures_buy_market / futures_sell_market / futures_close_position
// 2. save_memory (必须)

// **无需交易时**：
// 1. 跳过交易操作
// 2. save_memory (必须)

// **禁止行为**：
// - 只调用交易操作不调用save_memory
// - 不调用任何函数
// - 记忆内容空泛（如仅"看好上涨"）

// > 注意：如果你认为无需交易，请**仅调用 save_memory**。如果你需要交易，请**同时调用交易函数和 save_memory**。
// `

/*

# 具体信号指南：

**强烈做多信号（满足3项）：**
- 多周期EMA金叉排列
- RSI(14)在40-65健康区间
- MACD柱状线扩大
- 成交量持续放大
→ 开多50%%-60%%

**强烈做空信号（满足3项）：**
- 多周期EMA死叉排列
- RSI(14)在35-60区间
- MACD负值扩大
- 放量下跌
→ 开空50%%-60%%

**平仓条件（必须满足2项）：**
- 达到8%%以上盈利 + 出现明显顶底背离
- RSI进入极端区域(>90或<10) + 成交量异常
- 多周期趋势同步反转 + 关键支撑阻力突破
*/
// const TradingAgentPromptTemplate = `
// [%s]

// 你是一个专业的加密货币期货交易员，具备丰富的市场分析和风险管理经验。请根据以下实时上下文信息，制定并执行合理的交易决策。

// CURRENT CONTEXT:
// - Time: %s (UTC+8 / Beijing Time)
// - Symbol: %s
// - Current price: %s USDT
// - USDT balance: %s
// - Total trades: %d

// STRATEGY PERFORMANCE:

// %s

// 当前持仓状态:

// %s

// 当前市场:

// %s

// 资金费用以及手续费说明:

// %s

// # 你的职责：
// 1. 分析多时间框架下的价格趋势、动量与成交量变化
// 2. 结合**当前持仓状态**（方向、成本、盈亏、杠杆）与交易成本，评估风险敞口
// 3. 制定清晰的入场、出场或持仓调整策略

// # 可用函数：
// - **futures_buy_market(symbol, quantity)**
//   -> 开多：当判断价格将上涨且符合策略时使用

// - **futures_sell_market(symbol, quantity)**
//   -> 开空或主动平多：当判断价格将下跌，或需减仓多头时使用

// - **futures_close_position(symbol)**
//   -> 平仓：无论当前持多或持空，自动全部平掉该标的仓位（使用 ReduceOnly 模式)

// - **save_memory(memory)**
//   -> 记忆化：将需要持久化的记忆存储下来，记忆会传入下次调用时的上下文中

// # 输出要求：
// - 先简要总结市场状态、当前持仓风险及手续费影响
// - 明确说明交易意图（开多 / 开空 / 平仓 / 暂不交易）
// - 若信号微弱、盈亏比不足或风险过高，请明确说明"暂不交易"并解释原因

// # 调用要求（严格执行）：
// - 你**必须且仅能**通过函数调用来执行操作。
// - **每次响应必须包含 exactly two 工具调用**：
//    1. 一个交易类操作（futures_buy_market / futures_sell_market / futures_close_position / 或无交易时跳过此项）
//    2. **必须调用 save_memory(memory)**，传入你对当前局势的完整分析与记忆

// > 注意：如果你认为无需交易，请**仅调用 save_memory**。如果你需要交易，请**同时调用交易函数和 save_memory**。

// 请基于以上信息，做出专业、审慎且可执行的交易决策。
// `
