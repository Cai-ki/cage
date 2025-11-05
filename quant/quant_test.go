package quant_test

import (
	"encoding/json"
	"strconv"
	"testing"

	_ "github.com/Cai-ki/cage/config"
	"github.com/Cai-ki/cage/quant"
	"github.com/Cai-ki/cage/sugar"
)

func TestFetchLatestPrice(t *testing.T) {
	price := sugar.Must(quant.FetchLatestPrice("BTC/USDT"))

	// Verify JSON serializable
	data := sugar.Must(json.Marshal(price))
	t.Log(string(data))
}

func TestFetchHistoricalPrices(t *testing.T) {
	klines := sugar.Must(quant.FetchHistoricalPrices("BTC/USDT", "1h", 3))

	data := sugar.Must(json.Marshal(klines))
	t.Log(string(data))
}

func TestGetAccountBalance(t *testing.T) {
	bal := sugar.Must(quant.GetAccountBalance("USDT"))
	t.Log(bal)
	bal = sugar.Must(quant.GetAccountBalance("BTC"))
	t.Log(bal)
}

func TestListOpenOrders(t *testing.T) {
	orders := sugar.Must(quant.ListOpenOrders("BTC/USDT"))

	data := sugar.Must(json.Marshal(orders))
	t.Log(string(data))
}

// ⚠️ Uncomment ONLY if you want to test real orders on Binance Testnet
func TestBuyAndSellMarket(t *testing.T) {
	bal := sugar.Must(quant.GetAccountBalance("USDT"))
	usdt, _ := strconv.ParseFloat(bal, 64)
	if usdt < 15 {
		t.Skip("Insufficient USDT on testnet")
	}

	// Buy
	order := sugar.Must(quant.BuyMarket("BTC/USDT", 15.0))
	t.Logf("Buy order: %+v", order)

	// Sell ALL the BTC you just bought
	btcBal := sugar.Must(quant.GetAccountBalance("BTC"))
	btcFloat, _ := strconv.ParseFloat(btcBal, 64)

	// Due to precision, use a tiny epsilon
	if btcFloat > 0.000001 {
		sellOrder := sugar.Must(quant.SellMarket("BTC/USDT", btcFloat))
		t.Logf("Sell order: %+v", sellOrder)
	}
}
