package quant_test

import (
	"fmt"
	"testing"

	"github.com/Cai-ki/cage/quant"
	futures "github.com/adshao/go-binance/v2/futures"
)

func TestTechnicalIndicators(t *testing.T) {
	// 创建配置
	config := &quant.IndicatorConfig{
		EMAs:       []int{12, 26},
		MAs:        []int{20},
		RSI:        []int{14},
		MACD:       true,
		Stochastic: []int{14, 3},
		ATR:        []int{14},
		Bollinger:  []int{20, 2},
	}

	calculator := quant.NewIndicatorCalculator(config)

	// 生成数据
	k5m, _ := quant.FuturesGetKlines("BTCUSDT", "5m", 100)
	k15m, _ := quant.FuturesGetKlines("BTCUSDT", "15m", 80)
	k1h, _ := quant.FuturesGetKlines("BTCUSDT", "1h", 50)

	// 计算多周期指标
	timeframeData := map[string][]*futures.Kline{
		"5m":  k5m,
		"15m": k15m,
		"1h":  k1h,
	}

	multiIndicator := calculator.CalculateMultiTimeframe("BTCUSDT", timeframeData)

	// 测试输出
	fmt.Println("=== JSON 输出 ===")
	fmt.Println(multiIndicator.ToJSON())

	fmt.Println("\n=== 简洁格式 ===")
	timeframeOrder := []string{"1h", "15m", "5m"}
	fmt.Println(multiIndicator.ToSimpleString(timeframeOrder))

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
