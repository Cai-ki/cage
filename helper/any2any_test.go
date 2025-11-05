package helper_test

import (
	"fmt"
	"testing"

	"github.com/Cai-ki/cage/helper"
)

// 测试 JSON → SQL
func TestJsonToSql(t *testing.T) {
	jsonStr := `{
		"symbol": "BTCUSDT",
		"price": 67234.5,
		"volume": 123.456,
		"is_margin": true,
		"timestamp": "2025-11-05T12:34:56Z",
		"user": {
			"id": "user123",
			"balance": 10000
		},
		"tags": ["crypto", "long"]
	}`
	sql, err := helper.JsonToSql(jsonStr)
	if err != nil {
		t.Fatalf("JsonToSql error: %v", err)
	}
	fmt.Println("=== JSON → SQL ===")
	fmt.Println(sql)
}

// 测试 JSON → Go Struct
func TestJsonToGoStruct(t *testing.T) {
	jsonStr := `{
		"symbol": "ETHUSDT",
		"open": 3200.1,
		"high": 3250.0,
		"low": 3180.5,
		"close": 3220.75,
		"volume": 500.23,
		"timestamp": "2025-11-05T12:00:00Z",
		"final": true
	}`
	structCode, err := helper.JsonToGoStruct(jsonStr)
	if err != nil {
		t.Fatalf("JsonToGoStruct error: %v", err)
	}
	fmt.Println("\n=== JSON → Go Struct ===")
	fmt.Println(structCode)
}

// 测试 SQL → JSON Schema
func TestSqlToJSONSchema(t *testing.T) {
	sql := `CREATE TABLE klines (
		symbol TEXT NOT NULL,
		open REAL,
		high REAL,
		low REAL,
		close REAL NOT NULL,
		volume REAL,
		timestamp DATETIME,
		final BOOLEAN
	);`
	schema, err := helper.SqlToJSONSchema(sql)
	if err != nil {
		t.Fatalf("SqlToJSONSchema error: %v", err)
	}
	fmt.Println("\n=== SQL → JSON Schema ===")
	fmt.Println(schema)
}

// 测试 CSV → SQL
func TestCsvToSql(t *testing.T) {
	csvSample := `timestamp,symbol,price,volume,buyer_is_maker
2025-11-05 12:34:56,BTCUSDT,67234.5,0.123,true
2025-11-05 12:35:00,ETHUSDT,3220.75,1.456,false`
	sql, err := helper.CsvToSql(csvSample)
	if err != nil {
		t.Fatalf("CsvToSql error: %v", err)
	}
	fmt.Println("\n=== CSV → SQL ===")
	fmt.Println(sql)
}

// 测试 自然语言描述 → SQL
func TestDescriptionToSql(t *testing.T) {
	desc := "一张表用于记录加密货币4小时K线数据，包含交易对名称（字符串）、开盘价（小数）、最高价、最低价、收盘价、成交量（小数）、时间戳（时间）以及是否为最终K线（布尔值）"
	sql, err := helper.DescriptionToSql(desc)
	if err != nil {
		t.Fatalf("DescriptionToSql error: %v", err)
	}
	fmt.Println("\n=== Description → SQL ===")
	fmt.Println(sql)
}

// 测试 JSON → Protobuf
func TestJsonToProto(t *testing.T) {
	jsonStr := `{"symbol": "BTCUSDT", "price": 67234.5, "final": true}`
	proto, err := helper.JsonToProto(jsonStr)
	if err != nil {
		t.Fatalf("JsonToProto error: %v", err)
	}
	fmt.Println("\n=== JSON → Protobuf ===")
	fmt.Println(proto)
}
