package quant

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Cai-ki/cage/sugar"
	futures "github.com/adshao/go-binance/v2/futures"
	"github.com/markcheno/go-talib"
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

// 专业指标计算器（使用go-talib）
type ProfessionalIndicatorCalculator struct {
	config *IndicatorConfig
}

func NewIndicatorCalculator(config *IndicatorConfig) *ProfessionalIndicatorCalculator {
	return &ProfessionalIndicatorCalculator{config: config}
}

// 转换K线数据为价格序列
func (pic *ProfessionalIndicatorCalculator) convertKlinesToPriceData(klines []*futures.Kline) ([]float64, []float64, []float64) {
	closes := make([]float64, len(klines))
	highs := make([]float64, len(klines))
	lows := make([]float64, len(klines))

	for i, kline := range klines {
		closes[i] = sugar.Must(sugar.StrToT[float64](kline.Close))
		highs[i] = sugar.Must(sugar.StrToT[float64](kline.High))
		lows[i] = sugar.Must(sugar.StrToT[float64](kline.Low))
	}

	return closes, highs, lows
}

// 计算单个周期指标
// 计算单个周期指标
func (pic *ProfessionalIndicatorCalculator) Calculate(symbol, timeframe string, klines []*futures.Kline) *TechnicalIndicator {
	if len(klines) == 0 {
		return nil
	}

	closes, highs, lows := pic.convertKlinesToPriceData(klines)

	indicator := &TechnicalIndicator{
		Timestamp:  time.Now(),
		Symbol:     symbol,
		Price:      closes[len(closes)-1],
		Indicators: make(map[string]float64),
	}

	// 计算EMA - 需要至少period根K线
	for _, period := range pic.config.EMAs {
		if len(closes) < period {
			continue // 跳过，数据不足
		}
		key := fmt.Sprintf("ema_%d", period)
		ema := talib.Ema(closes, period)
		if len(ema) > 0 {
			indicator.Indicators[key] = ema[len(ema)-1]
		}
	}

	// 计算MA (SMA) - 需要至少period根K线
	for _, period := range pic.config.MAs {
		if len(closes) < period {
			continue // 跳过，数据不足
		}
		key := fmt.Sprintf("ma_%d", period)
		ma := talib.Sma(closes, period)
		if len(ma) > 0 {
			indicator.Indicators[key] = ma[len(ma)-1]
		}
	}

	// 计算RSI - 需要至少period+1根K线
	for _, period := range pic.config.RSI {
		if len(closes) < period+1 {
			continue // 跳过，数据不足
		}
		key := fmt.Sprintf("rsi_%d", period)
		rsi := talib.Rsi(closes, period)
		if len(rsi) > 0 {
			indicator.Indicators[key] = rsi[len(rsi)-1]
		}
	}

	// 计算MACD - 需要至少35根K线（26+9）
	if pic.config.MACD && len(closes) >= 35 {
		macd, macdSignal, macdHist := talib.Macd(closes, 12, 26, 9)
		if len(macd) > 0 {
			indicator.Indicators["macd_dif"] = macd[len(macd)-1]
			indicator.Indicators["macd_dea"] = macdSignal[len(macdSignal)-1]
			indicator.Indicators["macd"] = macdHist[len(macdHist)-1]
		}
	}

	// 计算随机指标 - 需要至少max(周期)+3根K线
	if len(pic.config.Stochastic) >= 2 {
		period := max(pic.config.Stochastic[0], pic.config.Stochastic[1])
		if len(closes) >= period+3 {
			fastK, fastD := talib.Stoch(highs, lows, closes,
				pic.config.Stochastic[0], // %K周期
				3,                        // %K平滑周期
				0,                        // %K移动平均类型 (0=SMA)
				pic.config.Stochastic[1], // %D周期
				0)                        // %D移动平均类型 (0=SMA)
			if len(fastK) > 0 {
				indicator.Indicators["stoch_k"] = fastK[len(fastK)-1]
				indicator.Indicators["stoch_d"] = fastD[len(fastD)-1]
			}
		}
	}

	// 计算ATR - 需要至少period+1根K线
	for _, period := range pic.config.ATR {
		if len(closes) < period+1 {
			continue // 跳过，数据不足
		}
		key := fmt.Sprintf("atr_%d", period)
		atr := talib.Atr(highs, lows, closes, period)
		if len(atr) > 0 {
			indicator.Indicators[key] = atr[len(atr)-1]
		}
	}

	// 计算布林带 - 需要至少period根K线
	if len(pic.config.Bollinger) >= 2 {
		period := pic.config.Bollinger[0]
		if len(closes) >= period {
			upper, middle, lower := talib.BBands(closes,
				period,                           // 周期
				float64(pic.config.Bollinger[1]), // 上轨标准差倍数
				float64(pic.config.Bollinger[1]), // 下轨标准差倍数
				0)                                // MA类型 (0=SMA)
			if len(upper) > 0 {
				indicator.Indicators["bb_upper"] = upper[len(upper)-1]
				indicator.Indicators["bb_middle"] = middle[len(middle)-1]
				indicator.Indicators["bb_lower"] = lower[len(lower)-1]
			}
		}
	}

	return indicator
}

// 计算多周期指标
func (pic *ProfessionalIndicatorCalculator) CalculateMultiTimeframe(symbol string, timeframeData map[string][]*futures.Kline) *MultiTimeframeIndicator {
	multi := &MultiTimeframeIndicator{
		Timestamp:  time.Now(),
		Symbol:     symbol,
		Timeframes: make(map[string]*TechnicalIndicator),
	}

	// 计算每个周期
	for tf, klines := range timeframeData {
		multi.Timeframes[tf] = pic.Calculate(symbol, tf, klines)
	}

	return multi
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

// 时间格式化
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
