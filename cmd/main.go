package main

import (
	"encoding/json"
	"fmt"
	"log"

	_ "github.com/Cai-ki/cage/localconfig"
	"github.com/Cai-ki/cage/quant"
)

func main() {
	// 1. 获取价格（直接给 LLM）
	price, err := quant.FetchLatestPrice("BTC/USDT")
	if err != nil {
		log.Fatal(err)
	}
	priceJSON, _ := json.Marshal(price)
	fmt.Printf("Price for LLM: %s\n", priceJSON)

	// 2. 获取余额
	bal, err := quant.GetAccountBalance("USDT")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("USDT Balance: %s\n", bal)

	// 3. 下单（谨慎！）
	order, err := quant.BuyMarket("BTC/USDT", 10.0)
	if err != nil {
		log.Fatal(err)
	}
	orderJSON, _ := json.Marshal(order)
	fmt.Printf("Order result: %s\n", orderJSON)
}
