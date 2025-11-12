package quant

import (
	"fmt"
	"testing"

	futures "github.com/adshao/go-binance/v2/futures"
)

// 模拟K线数据生成
func generateMockKlines(count int, basePrice float64) []*futures.Kline {
	var klines []*futures.Kline
	price := basePrice

	for i := 0; i < count; i++ {
		// 模拟价格波动
		change := (float64(i%10) - 5) * 10 // -50 到 +50 的波动
		price += change

		kline := &futures.Kline{
			Open:   fmt.Sprint(price - 5),
			High:   fmt.Sprint(price + 10),
			Low:    fmt.Sprint(price - 10),
			Close:  fmt.Sprint(price),
			Volume: fmt.Sprint(1000 + float64(i)*100),
		}
		klines = append(klines, kline)
	}

	return klines
}

func TestTechnicalIndicators(t *testing.T) {
	// 创建配置
	config := &IndicatorConfig{
		EMAs:       []int{12, 26},
		MAs:        []int{20},
		RSI:        []int{14},
		MACD:       true,
		Stochastic: []int{14, 3},
		ATR:        []int{14},
		Bollinger:  []int{20, 2},
	}

	calculator := NewIndicatorCalculator(config)

	// 生成模拟数据
	klines5m := generateMockKlines(100, 50000)
	klines15m := generateMockKlines(80, 50000)
	klines1h := generateMockKlines(50, 50000)

	// 计算多周期指标
	timeframeData := map[string][]*futures.Kline{
		"5m":  klines5m,
		"15m": klines15m,
		"1h":  klines1h,
	}

	multiIndicator := calculator.CalculateMultiTimeframe("BTCUSDT", timeframeData)

	// 测试输出
	fmt.Println("=== JSON 输出 ===")
	fmt.Println(multiIndicator.ToJSON())

	fmt.Println("\n=== 简洁格式 ===")
	fmt.Println(multiIndicator.ToSimpleString())

	// 验证数据完整性
	if multiIndicator.Symbol != "BTCUSDT" {
		t.Errorf("Expected symbol BTCUSDT, got %s", multiIndicator.Symbol)
	}

	if len(multiIndicator.Timeframes) != 3 {
		t.Errorf("Expected 3 timeframes, got %d", len(multiIndicator.Timeframes))
	}

	// 验证5分钟周期数据
	if indicator5m, exists := multiIndicator.Timeframes["5m"]; exists {
		if indicator5m.Price <= 0 {
			t.Error("5m price should be positive")
		}

		if _, exists := indicator5m.Indicators["ema_12"]; !exists {
			t.Error("5m should have ema_12 indicator")
		}
	}

	fmt.Println("✅ 所有测试通过!")
}
