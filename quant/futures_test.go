package quant_test

import (
	"encoding/json"
	"math"
	"strconv"
	"testing"

	_ "github.com/Cai-ki/cage/config"
	"github.com/Cai-ki/cage/quant"
	"github.com/Cai-ki/cage/sugar"
)

func TestFuturesGetBalance(t *testing.T) {
	bal := sugar.Must(quant.FuturesGetBalance("USDT"))
	t.Log("Futures USDT balance:", bal)
}

func TestFuturesGetPosition(t *testing.T) {
	pos := sugar.Must(quant.FuturesGetPosition("BTCUSDT"))
	data := sugar.Must(json.Marshal(pos))
	t.Log("Position:", string(data))
}

func TestFuturesGetKlines(t *testing.T) {
	klines := sugar.Must(quant.FuturesGetKlines("BTCUSDT", "1h", 5))
	data := sugar.Must(json.Marshal(klines))
	t.Log("Klines:", string(data))
}

// ⚠️ DANGEROUS: Only run on futures testnet with small amounts
func TestFuturesBuySellAndClose(t *testing.T) {
	// Check balance
	bal := sugar.Must(quant.FuturesGetBalance("USDT"))
	usdt, _ := strconv.ParseFloat(bal, 64)
	if usdt < 10.0 {
		t.Skip("Insufficient USDT balance on futures testnet")
	}

	symbol := "BTCUSDT"

	// Buy (open long)
	order1 := sugar.Must(quant.FuturesBuyMarket(symbol, 0.001))
	t.Logf("Opened long: %+v", order1)

	// Check position
	pos := sugar.Must(quant.FuturesGetPosition(symbol))
	t.Logf("Position after buy: %+v", pos)

	// Close position
	closeOrder := sugar.Must(quant.FuturesClosePosition(symbol))
	t.Logf("Closed position: %+v", closeOrder)

	// Verify closed
	pos2 := sugar.Must(quant.FuturesGetPosition(symbol))
	amt, err := strconv.ParseFloat(pos2.PositionAmt, 64)
	if err != nil {
		t.Fatalf("Failed to parse position amount: %v", err)
	}

	// 判断是否接近 0（考虑浮点精度）
	if math.Abs(amt) > 1e-8 {
		t.Errorf("Position not fully closed: %s", pos2.PositionAmt)
	}
}
