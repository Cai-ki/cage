// Package quant provides a minimal AI-friendly interface to crypto exchanges.
// Currently supports Binance (Spot) via environment variables.
package quant

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	binance "github.com/adshao/go-binance/v2"
)

var (
	client *binance.Client
)

func init() {
	exchange := os.Getenv("EXCHANGE_NAME")
	if exchange == "" {
		exchange = "binance"
	}

	if strings.ToLower(exchange) != "binance" {
		return
	}

	apiKey := os.Getenv("EXCHANGE_API_KEY")
	secretKey := os.Getenv("EXCHANGE_API_SECRET")
	if apiKey == "" || secretKey == "" {
		return
	}

	// Enable Testnet if needed
	if os.Getenv("BINANCE_TESTNET") == "true" {
		binance.UseTestnet = true
	}

	client = binance.NewClient(apiKey, secretKey)
}

// standardizeSymbol converts "BTC/USDT" → "BTCUSDT"
func standardizeSymbol(symbol string) string {
	s := strings.ToUpper(symbol)
	s = strings.ReplaceAll(s, "/", "")
	s = strings.ReplaceAll(s, "-", "")
	s = strings.ReplaceAll(s, "_", "")
	if len(s) > 20 {
		s = s[:20]
	}
	return s
}

// FetchLatestPrice returns the latest price as a binance.SymbolPrice.
func FetchLatestPrice(symbol string) (*binance.SymbolPrice, error) {
	symbol = standardizeSymbol(symbol)
	prices, err := binance.NewClient("", "").NewListPricesService().Symbol(symbol).Do(context.Background())
	if err != nil {
		return nil, fmt.Errorf("fetch ticker: %w", err)
	}
	if len(prices) == 0 {
		return nil, fmt.Errorf("no price found for %s", symbol)
	}
	return prices[0], nil
}

// FetchHistoricalPrices returns historical klines.
func FetchHistoricalPrices(symbol, interval string, limit int) ([]*binance.Kline, error) {
	symbol = standardizeSymbol(symbol)
	return binance.NewClient("", "").NewKlinesService().
		Symbol(symbol).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
}

// BuyMarket buys using quote currency (e.g., USDT).
func BuyMarket(symbol string, quoteAmount float64) (*binance.CreateOrderResponse, error) {
	if client == nil {
		return nil, errors.New("client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	symbol = standardizeSymbol(symbol)
	qStr := fmt.Sprintf("%.8f", quoteAmount)
	return client.NewCreateOrderService().
		Symbol(symbol).
		Type(binance.OrderTypeMarket).
		Side(binance.SideTypeBuy).
		QuoteOrderQty(qStr).
		Do(context.Background())
}

// SellMarket sells base asset.
func SellMarket(symbol string, baseAmount float64) (*binance.CreateOrderResponse, error) {
	if client == nil {
		return nil, errors.New("client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	symbol = standardizeSymbol(symbol)
	qStr := fmt.Sprintf("%.8f", baseAmount)
	return client.NewCreateOrderService().
		Symbol(symbol).
		Type(binance.OrderTypeMarket).
		Side(binance.SideTypeSell).
		Quantity(qStr).
		Do(context.Background())
}

// GetAccountBalance returns the free balance of a currency.
func GetAccountBalance(currency string) (string, error) {
	if client == nil {
		return "", errors.New("client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	account, err := client.NewGetAccountService().Do(context.Background())
	if err != nil {
		return "", err
	}
	for _, b := range account.Balances {
		if strings.EqualFold(b.Asset, currency) {
			return b.Free, nil
		}
	}
	return "0", nil
}

// ListOpenOrders returns open orders as []*binance.Order.
func ListOpenOrders(symbol string) ([]*binance.Order, error) {
	if client == nil {
		return nil, errors.New("client not initialized (missing EXCHANGE_API_KEY/SECRET)")
	}
	symbol = standardizeSymbol(symbol)
	return client.NewListOpenOrdersService().Symbol(symbol).Do(context.Background())
}
