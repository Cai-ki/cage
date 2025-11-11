package quant

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	futures "github.com/adshao/go-binance/v2/futures"
)

var (
	futuresClient *futures.Client
)

func init() {
	exchange := os.Getenv("EXCHANGE_NAME")
	if exchange == "" {
		exchange = "binance"
	}

	if strings.ToLower(exchange) != "binance" {
		return
	}

	apiKey := os.Getenv("EXCHANGE_FUTURES_API_KEY")
	secretKey := os.Getenv("EXCHANGE_FUTURES_API_SECRET")
	if apiKey == "" || secretKey == "" {
		return
	}

	// Enable Testnet if needed
	if os.Getenv("BINANCE_TESTNET") == "true" {
		futures.UseTestnet = true
	}

	futuresClient = futures.NewClient(apiKey, secretKey)
}

// standardizeSymbolFutures converts "BTC/USDT" → "BTCUSDT"
// Same as spot, but kept separate for clarity
func standardizeSymbolFutures(symbol string) string {
	return standardizeSymbol(symbol)
}

// FuturesGetBalance returns USDT-margined futures account balance
func FuturesGetBalance(asset string) (string, error) {
	if futuresClient == nil {
		return "", errors.New("futures client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	resp, err := futuresClient.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", err
	}
	for _, b := range resp.Assets {
		if strings.EqualFold(b.Asset, asset) {
			return b.WalletBalance, nil
		}
	}
	return "0", nil
}

// FuturesGetPosition returns position info for a symbol
func FuturesGetPosition(symbol string) (*futures.PositionRisk, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	positions, err := futuresClient.NewGetPositionRiskService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, err
	}
	if len(positions) == 0 {
		return nil, fmt.Errorf("no position found for %s", symbol)
	}
	return positions[0], nil
}

// FuturesBuyMarket opens a long position with market order
func FuturesBuyMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	qStr := fmt.Sprintf("%.8f", quantity)
	return futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Type(futures.OrderTypeMarket).
		Side(futures.SideTypeBuy).
		Quantity(qStr).
		Do(context.Background())
}

// FuturesSellMarket opens a short position or closes long
func FuturesSellMarket(symbol string, quantity float64) (*futures.CreateOrderResponse, error) {
	if futuresClient == nil {
		return nil, errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	qStr := fmt.Sprintf("%.8f", quantity)
	return futuresClient.NewCreateOrderService().
		Symbol(symbol).
		Type(futures.OrderTypeMarket).
		Side(futures.SideTypeSell).
		Quantity(qStr).
		Do(context.Background())
}

// FuturesClosePosition closes current position (auto detect direction)
func FuturesClosePosition(symbol string) (*futures.CreateOrderResponse, error) {
	pos, err := FuturesGetPosition(symbol)
	if err != nil {
		return nil, err
	}
	qty := pos.PositionAmt
	if qty == "0" {
		return nil, fmt.Errorf("no open position for %s", symbol)
	}

	// Determine close side
	var side futures.SideType
	qtyAbs := qty
	if strings.HasPrefix(qty, "-") {
		side = futures.SideTypeBuy // short → buy to close
		qtyAbs = qty[1:]
	} else {
		side = futures.SideTypeSell // long → sell to close
	}

	return futuresClient.NewCreateOrderService().
		Symbol(standardizeSymbolFutures(symbol)).
		Type(futures.OrderTypeMarket).
		Side(side).
		Quantity(qtyAbs).
		ReduceOnly(true). // important: only reduce
		Do(context.Background())
}

// FuturesGetKlines returns historical futures klines
func FuturesGetKlines(symbol, interval string, limit int) ([]*futures.Kline, error) {
	symbol = standardizeSymbolFutures(symbol)
	return futuresClient.NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
}

func FuturesGetTickerPrice(symbol string) (string, error) {
	if futuresClient == nil {
		return "", errors.New("futures client not initialized")
	}
	symbol = standardizeSymbolFutures(symbol)
	resp, err := futuresClient.NewPremiumIndexService().
		Symbol(symbol).
		Do(context.Background())
	if err != nil {
		return "", err
	}
	if len(resp) > 0 {
		return resp[0].MarkPrice, nil // 或者 resp[0].LastPrice
	}
	return "", errors.New("no price data")
}
