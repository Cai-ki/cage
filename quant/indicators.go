package quant

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"strings"
	"time"

	"github.com/Cai-ki/cage/sugar"
	futures "github.com/adshao/go-binance/v2/futures"
)

// 技术指标配置
type IndicatorConfig struct {
	EMAs       []int `json:"emas"`       // [12, 26, 50]
	MAs        []int `json:"mas"`        // [20, 50]
	RSI        []int `json:"rsi"`        // [14]
	MACD       bool  `json:"macd"`       // 是否计算MACD
	Stochastic []int `json:"stochastic"` // [14, 3]
	ATR        []int `json:"atr"`        // [14]
	Bollinger  []int `json:"bollinger"`  // [20, 2]
}

// 单个周期技术指标
type TechnicalIndicator struct {
	Timestamp  time.Time          `json:"timestamp"`
	Symbol     string             `json:"symbol"`
	Price      float64            `json:"price"`
	Indicators map[string]float64 `json:"indicators"` // 指标值
}

// 多周期技术指标
type MultiTimeframeIndicator struct {
	Timestamp  time.Time                      `json:"timestamp"`
	Symbol     string                         `json:"symbol"`
	Timeframes map[string]*TechnicalIndicator `json:"timeframes"` // 周期 -> 指标
}

// 指标计算器
type IndicatorCalculator struct {
	config *IndicatorConfig
}

func NewIndicatorCalculator(config *IndicatorConfig) *IndicatorCalculator {
	return &IndicatorCalculator{config: config}
}

// 计算单个周期指标
func (ic *IndicatorCalculator) Calculate(symbol, timeframe string, klines []*futures.Kline) *TechnicalIndicator {
	if len(klines) == 0 {
		log.Println("Empty klines for", symbol, timeframe)
		return nil
	}

	indicator := &TechnicalIndicator{
		Timestamp:  time.Now(),
		Symbol:     symbol,
		Price:      sugar.Must(sugar.StrToT[float64](klines[len(klines)-1].Close)),
		Indicators: make(map[string]float64),
	}

	// 计算EMA
	for _, period := range ic.config.EMAs {
		key := fmt.Sprintf("ema_%d", period)
		indicator.Indicators[key] = ic.calculateEMA(klines, period)
	}

	// 计算MA
	for _, period := range ic.config.MAs {
		key := fmt.Sprintf("ma_%d", period)
		indicator.Indicators[key] = ic.calculateMA(klines, period)
	}

	// 计算RSI
	for _, period := range ic.config.RSI {
		key := fmt.Sprintf("rsi_%d", period)
		rsi := ic.calculateRSI(klines, period)
		indicator.Indicators[key] = rsi
	}

	// 计算MACD
	if ic.config.MACD {
		dif, dea, macd := ic.calculateMACD(klines)
		indicator.Indicators["macd_dif"] = dif
		indicator.Indicators["macd_dea"] = dea
		indicator.Indicators["macd"] = macd
	}

	// 计算随机指标
	if len(ic.config.Stochastic) >= 2 {
		k, d := ic.calculateStochastic(klines, ic.config.Stochastic[0], ic.config.Stochastic[1])
		indicator.Indicators["stoch_k"] = k
		indicator.Indicators["stoch_d"] = d
	}

	// 计算ATR
	for _, period := range ic.config.ATR {
		key := fmt.Sprintf("atr_%d", period)
		indicator.Indicators[key] = ic.calculateATR(klines, period)
	}

	// 计算布林带
	if len(ic.config.Bollinger) >= 2 {
		upper, middle, lower := ic.calculateBollinger(klines, ic.config.Bollinger[0], float64(ic.config.Bollinger[1]))
		indicator.Indicators["bb_upper"] = upper
		indicator.Indicators["bb_middle"] = middle
		indicator.Indicators["bb_lower"] = lower
	}

	return indicator
}

// 计算多周期指标
func (ic *IndicatorCalculator) CalculateMultiTimeframe(symbol string, timeframeData map[string][]*futures.Kline) *MultiTimeframeIndicator {
	multi := &MultiTimeframeIndicator{
		Timestamp:  time.Now(),
		Symbol:     symbol,
		Timeframes: make(map[string]*TechnicalIndicator),
	}

	// 计算每个周期
	for tf, klines := range timeframeData {
		multi.Timeframes[tf] = ic.Calculate(symbol, tf, klines)
	}

	return multi
}

// 计算EMA
func (ic *IndicatorCalculator) calculateEMA(klines []*futures.Kline, period int) float64 {
	// 在calculateEMA中
	if len(klines) < period {
		return 0
	}
	// 使用更精确的初始SMA
	sum := 0.0
	for i := 0; i < period; i++ {
		sum += sugar.Must(sugar.StrToT[float64](klines[i].Close))
	}
	ema := sum / float64(period)

	// 从第period根开始递归
	multiplier := 2.0 / (float64(period) + 1.0)
	for i := period; i < len(klines); i++ {
		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		ema = close*multiplier + ema*(1-multiplier)
	}
	return ema
}

// 计算MA
func (ic *IndicatorCalculator) calculateMA(klines []*futures.Kline, period int) float64 {
	if len(klines) < period {
		return 0
	}

	sum := 0.0
	for i := len(klines) - period; i < len(klines); i++ {
		sum += sugar.Must(sugar.StrToT[float64](klines[i].Close))
	}
	return sum / float64(period)
}

// 计算RSI
func (ic *IndicatorCalculator) calculateRSI(klines []*futures.Kline, period int) float64 {
	if len(klines) <= period {
		return 50
	}

	gains := make([]float64, 0, period)
	losses := make([]float64, 0, period)

	// 计算初始14根的变化
	for i := 1; i <= period; i++ {
		prev := sugar.Must(sugar.StrToT[float64](klines[i-1].Close))
		curr := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		change := curr - prev

		if change > 0 {
			gains = append(gains, change)
		} else {
			losses = append(losses, -change)
		}
	}

	// Wilders平滑初始值
	avgGain := 0.0
	avgLoss := 0.0
	if len(gains) > 0 {
		avgGain = sum(gains) / float64(period)
	}
	if len(losses) > 0 {
		avgLoss = sum(losses) / float64(period)
	}

	// 递归计算后续值
	for i := period + 1; i < len(klines); i++ {
		prev := sugar.Must(sugar.StrToT[float64](klines[i-1].Close))
		curr := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		change := curr - prev

		currentGain := 0.0
		currentLoss := 0.0
		if change > 0 {
			currentGain = change
		} else {
			currentLoss = -change
		}

		avgGain = (avgGain*float64(period-1) + currentGain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + currentLoss) / float64(period)
	}

	if avgLoss == 0 {
		return 100
	}
	rs := avgGain / avgLoss
	return 100 - (100 / (1 + rs))
}

// 辅助函数
func sum(vals []float64) float64 {
	s := 0.0
	for _, v := range vals {
		s += v
	}
	return s
}

// 计算MACD
func (ic *IndicatorCalculator) calculateMACD(klines []*futures.Kline) (float64, float64, float64) {
	if len(klines) < 35 { // 需要26+9=35根K线
		return 0, 0, 0
	}

	// 1. 预计算所有EMA值
	ema12Series := make([]float64, len(klines))
	ema26Series := make([]float64, len(klines))

	// 计算12日EMA序列
	multiplier12 := 2.0 / 13.0
	sum12 := 0.0
	for i := 0; i < 12; i++ {
		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		sum12 += close
	}
	ema12 := sum12 / 12.0
	ema12Series[11] = ema12
	for i := 12; i < len(klines); i++ {
		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		ema12 = close*multiplier12 + ema12*(1-multiplier12)
		ema12Series[i] = ema12
	}

	// 计算26日EMA序列
	multiplier26 := 2.0 / 27.0
	sum26 := 0.0
	for i := 0; i < 26; i++ {
		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		sum26 += close
	}
	ema26 := sum26 / 26.0
	ema26Series[25] = ema26
	for i := 26; i < len(klines); i++ {
		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		ema26 = close*multiplier26 + ema26*(1-multiplier26)
		ema26Series[i] = ema26
	}

	// 2. 计算DIF序列 (从第26根开始)
	difSeries := make([]float64, len(klines))
	for i := 26; i < len(klines); i++ {
		difSeries[i] = ema12Series[i] - ema26Series[i]
	}

	// 3. 计算DEA (DIF的9日EMA，从第34根开始)
	if len(klines) < 35 {
		return 0, 0, 0
	}
	sumDEA := 0.0
	for i := 26; i < 35; i++ {
		sumDEA += difSeries[i]
	}
	dea := sumDEA / 9.0
	multiplierDEA := 2.0 / 10.0

	for i := 35; i < len(klines); i++ {
		dea = difSeries[i]*multiplierDEA + dea*(1-multiplierDEA)
	}

	// 4. 取最后值
	lastIdx := len(klines) - 1
	dif := difSeries[lastIdx]
	macd := (dif - dea) * 2

	return dif, dea, macd
}

// 计算随机指标
func (ic *IndicatorCalculator) calculateStochastic(klines []*futures.Kline, period, smooth int) (float64, float64) {
	if len(klines) < period+smooth-1 {
		return 50, 50
	}

	// 1. 计算%K（最近period根K线）
	lookbackStart := len(klines) - period
	lowestLow := math.MaxFloat64
	highestHigh := -math.MaxFloat64

	for i := lookbackStart; i < len(klines); i++ {
		low := sugar.Must(sugar.StrToT[float64](klines[i].Low))
		high := sugar.Must(sugar.StrToT[float64](klines[i].High))
		if low < lowestLow {
			lowestLow = low
		}
		if high > highestHigh {
			highestHigh = high
		}
	}

	currentClose := sugar.Must(sugar.StrToT[float64](klines[len(klines)-1].Close))
	if highestHigh == lowestLow {
		return 50, 50
	}
	k := 100 * (currentClose - lowestLow) / (highestHigh - lowestLow)

	// 2. 计算%D（%K的smooth日SMA）
	kValues := make([]float64, 0, smooth)
	for i := len(klines) - smooth; i < len(klines); i++ {
		if i < period-1 {
			continue
		}
		windowStart := i - period + 1
		windowLow := math.MaxFloat64
		windowHigh := -math.MaxFloat64

		for j := windowStart; j <= i; j++ {
			low := sugar.Must(sugar.StrToT[float64](klines[j].Low))
			high := sugar.Must(sugar.StrToT[float64](klines[j].High))
			if low < windowLow {
				windowLow = low
			}
			if high > windowHigh {
				windowHigh = high
			}
		}

		close := sugar.Must(sugar.StrToT[float64](klines[i].Close))
		if windowHigh == windowLow {
			kValues = append(kValues, 50)
		} else {
			kVal := 100 * (close - windowLow) / (windowHigh - windowLow)
			kValues = append(kValues, kVal)
		}
	}

	d := 0.0
	if len(kValues) > 0 {
		sumK := 0.0
		for _, v := range kValues {
			sumK += v
		}
		d = sumK / float64(len(kValues))
	}

	return k, d
}

// 计算ATR
func (ic *IndicatorCalculator) calculateATR(klines []*futures.Kline, period int) float64 {
	if len(klines) < period+1 {
		return 0
	}

	// 1. 计算所有TR值（从第2根K线开始）
	trValues := make([]float64, len(klines)-1)
	for i := 1; i < len(klines); i++ {
		high := sugar.Must(sugar.StrToT[float64](klines[i].High))
		low := sugar.Must(sugar.StrToT[float64](klines[i].Low))
		prevClose := sugar.Must(sugar.StrToT[float64](klines[i-1].Close))

		tr1 := high - low
		tr2 := math.Abs(high - prevClose)
		tr3 := math.Abs(low - prevClose)
		trValues[i-1] = math.Max(tr1, math.Max(tr2, tr3)) // TR[0]对应第2根K线
	}

	// 2. 初始ATR = 前period个TR的平均值
	sumTR := 0.0
	for i := 0; i < period; i++ {
		sumTR += trValues[i]
	}
	atr := sumTR / float64(period)

	// 3. 从第(period+1)个TR开始递归更新（Wilder平滑）
	for i := period; i < len(trValues); i++ {
		atr = (atr*(float64(period)-1) + trValues[i]) / float64(period)
	}

	return atr
}

// 计算布林带
func (ic *IndicatorCalculator) calculateBollinger(klines []*futures.Kline, period int, stdDev float64) (float64, float64, float64) {
	ma := ic.calculateMA(klines, period)

	sumSquares := 0.0
	for i := len(klines) - period; i < len(klines); i++ {
		diff := sugar.Must(sugar.StrToT[float64](klines[i].Close)) - ma
		sumSquares += diff * diff
	}

	std := math.Sqrt(sumSquares / float64(period))
	upper := ma + (std * stdDev)
	lower := ma - (std * stdDev)

	return upper, ma, lower
}

// 输出原始JSON数据
func (multi *MultiTimeframeIndicator) ToJSON() string {
	data := map[string]interface{}{
		"timestamp":  multi.Timestamp,
		"symbol":     multi.Symbol,
		"timeframes": make(map[string]interface{}),
	}

	for tf, indicator := range multi.Timeframes {
		data["timeframes"].(map[string]interface{})[tf] = map[string]interface{}{
			"price":      indicator.Price,
			"indicators": indicator.Indicators,
		}
	}

	jsonData, _ := json.Marshal(data)
	return string(jsonData)
}

func Format(t time.Time) string {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	return t.In(loc).Format("2006-01-02 15:04:05")
}

// 输出清晰格式（推荐）
func (multi *MultiTimeframeIndicator) ToSimpleString() string {
	var buf strings.Builder

	// 头部信息
	buf.WriteString(fmt.Sprintf("symbol: %s\n", multi.Symbol))
	buf.WriteString(fmt.Sprintf("timestamp: %s\n", Format(multi.Timestamp)))
	buf.WriteString("---\n")

	// 各周期数据
	for tf, indicator := range multi.Timeframes {
		buf.WriteString(fmt.Sprintf("[%s周期]\n", tf))
		buf.WriteString(fmt.Sprintf("  price: %.2f\n", indicator.Price))

		// 分组输出指标
		buf.WriteString("  趋势指标: ")
		for k, v := range indicator.Indicators {
			if strings.HasPrefix(k, "ema_") || strings.HasPrefix(k, "ma_") {
				buf.WriteString(fmt.Sprintf("%s=%.2f ", k, v))
			}
		}
		buf.WriteString("\n")

		buf.WriteString("  动量指标: ")
		for k, v := range indicator.Indicators {
			if strings.Contains(k, "rsi") || strings.Contains(k, "macd") || strings.Contains(k, "stoch") {
				buf.WriteString(fmt.Sprintf("%s=%.2f ", k, v))
			}
		}
		buf.WriteString("\n")

		buf.WriteString("  波动指标: ")
		for k, v := range indicator.Indicators {
			if strings.Contains(k, "atr") || strings.Contains(k, "bb_") {
				buf.WriteString(fmt.Sprintf("%s=%.2f ", k, v))
			}
		}
		buf.WriteString("\n\n")
	}

	return buf.String()
}
