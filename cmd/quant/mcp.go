package main

import (
	_ "embed"
	"encoding/json"

	"github.com/Cai-ki/cage/llm/mcp"
	"github.com/Cai-ki/cage/quant"
)

//go:embed mcp.json
var mcpString string

var globalMemory string

func init() {
	type futures_buy_market_args struct {
		Symbol   string  `json:"symbol"`
		Quantity float64 `json:"quantity"`
	}
	futures_buy_market := func(args futures_buy_market_args) (interface{}, error) {
		rsp, err := quant.FuturesBuyMarket(args.Symbol, args.Quantity)
		resultBytes, err := json.Marshal(rsp)
		return map[string]interface{}{"result": string(resultBytes)}, err
	}
	mcp.RegisterTool("futures_buy_market", futures_buy_market, futures_buy_market_args{})

	type futures_sell_market_args struct {
		Symbol   string  `json:"symbol"`
		Quantity float64 `json:"quantity"`
	}
	futures_sell_market := func(args futures_sell_market_args) (interface{}, error) {
		rsp, err := quant.FuturesSellMarket(args.Symbol, args.Quantity)
		resultBytes, err := json.Marshal(rsp)
		return map[string]interface{}{"result": string(resultBytes)}, err
	}
	mcp.RegisterTool("futures_sell_market", futures_sell_market, futures_sell_market_args{})

	type futures_close_position_args struct {
		Symbol string `json:"symbol"`
	}
	futures_close_position := func(args futures_close_position_args) (interface{}, error) {
		rsp, err := quant.FuturesClosePosition(args.Symbol)
		resultBytes, err := json.Marshal(rsp)
		return map[string]interface{}{"result": string(resultBytes)}, err
	}
	mcp.RegisterTool("futures_close_position", futures_close_position, futures_close_position_args{})

	type save_memory_args struct {
		Memory string `json:"memory"`
	}
	save_memory := func(args save_memory_args) (interface{}, error) {
		globalMemory = args.Memory
		return map[string]interface{}{"result": args.Memory}, nil
	}
	mcp.RegisterTool("save_memory", save_memory, save_memory_args{})
}
