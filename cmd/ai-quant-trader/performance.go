package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/Cai-ki/cage/quant"
)

type NetValuePoint struct {
	Timestamp time.Time
	USDT      float64
	BTC       float64
	NetValue  float64 // USDT + BTC * current_price_at_that_time
}

func calculateNetValueCurve(symbol string) []NetValuePoint {
	history := loadAllHistory(symbol)
	if len(history) == 0 {
		// 这理论上不会发生，因为 ensureInitRecord 会创建
		usdtStr, _ := quant.GetAccountBalance("USDT")
		btcStr, _ := quant.GetAccountBalance("BTC")
		usdt, _ := strconv.ParseFloat(usdtStr, 64)
		btc, _ := strconv.ParseFloat(btcStr, 64)
		latest, _ := quant.FetchLatestPrice(symbol)
		price, _ := strconv.ParseFloat(latest.Price, 64)
		net := usdt + btc*price
		return []NetValuePoint{{
			Timestamp: time.Now().UTC().Add(-time.Hour),
			USDT:      usdt,
			BTC:       btc,
			NetValue:  net,
		}}
	}

	// Find the first "init" record
	var initRecord *TradeRecord
	for i, r := range history {
		if r.Action == "init" {
			initRecord = &history[i]
			break
		}
	}

	if initRecord == nil {
		// Fallback: treat first record as init (should not happen)
		initRecord = &history[0]
	}

	// Start from init state
	usdt := initRecord.AmountUSDT
	btc := initRecord.AmountBTC

	// Sort full history by time
	sort.Slice(history, func(i, j int) bool {
		return history[i].Timestamp.Before(history[j].Timestamp)
	})

	var curve []NetValuePoint

	// Add init point
	priceAtInit, _ := strconv.ParseFloat(initRecord.Price, 64)
	if priceAtInit == 0 {
		priceAtInit = 1.0 // fallback
	}
	curve = append(curve, NetValuePoint{
		Timestamp: initRecord.Timestamp,
		USDT:      usdt,
		BTC:       btc,
		NetValue:  usdt + btc*priceAtInit,
	})

	// Replay all non-init records
	for _, r := range history {
		if r.Action == "init" {
			continue // already processed
		}
		price, _ := strconv.ParseFloat(r.Price, 64)
		if price == 0 {
			continue
		}
		switch r.Action {
		case "buy":
			buyUSDT := r.AmountUSDT
			if buyUSDT > 0 && buyUSDT <= usdt {
				btcBought := buyUSDT / price
				usdt -= buyUSDT
				btc += btcBought
			}
		case "sell":
			sellBTC := r.AmountBTC
			if sellBTC > 0 && sellBTC <= btc {
				usdt += sellBTC * price
				btc -= sellBTC
			}
		case "hold":
			// no change
		}
		net := usdt + btc*price
		curve = append(curve, NetValuePoint{
			Timestamp: r.Timestamp,
			USDT:      usdt,
			BTC:       btc,
			NetValue:  net,
		})
	}

	return curve
}

// Export to CSV for plotting
func exportNetValueCSV(symbol string) error {
	curve := calculateNetValueCurve(symbol)
	file, err := os.Create("testdata/net_value.csv")
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	writer.Write([]string{"timestamp", "usdt", "btc", "net_value"})
	for _, p := range curve {
		writer.Write([]string{
			p.Timestamp.Format("2006-01-02 15:04:05"),
			fmt.Sprintf("%.4f", p.USDT),
			fmt.Sprintf("%.6f", p.BTC),
			fmt.Sprintf("%.4f", p.NetValue),
		})
	}
	return nil
}
