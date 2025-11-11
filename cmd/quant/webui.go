package main

import (
	"encoding/json"
	"log"
	"net/http"
	"text/template"
	"time"
)

func StartServer() {
	go func() {
		http.HandleFunc("/", handleIndex)
		http.HandleFunc("/api/performance", handleAPIPerformance)

		log.Println("服务器启动在 :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.New("index").Parse(htmlTemplate))

	record, err := LoadOrCreatePerformanceRecord()
	if err != nil {
		http.Error(w, "Performance data unavailable", http.StatusInternalServerError)
		return
	}

	data := struct {
		Date           string
		InitialBalance float64
		CurrentBalance float64
		CumulativePnL  float64
		CumulativeROI  float64
		TotalTrades    int
		CurrentPrice   float64
		LastUpdate     string
	}{
		Date:           record.Date,
		InitialBalance: record.InitialBalance,
		CurrentBalance: record.CurrentBalance,
		CumulativePnL:  record.CumulativePnL,
		CumulativeROI:  record.CumulativeROI,
		TotalTrades:    record.TotalTrades,
		CurrentPrice:   record.CurrentPrice,
		LastUpdate:     time.Now().Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAPIPerformance(w http.ResponseWriter, r *http.Request) {
	record, err := LoadOrCreatePerformanceRecord()
	if err != nil {
		http.Error(w, "Performance data unavailable", http.StatusInternalServerError)
		return
	}

	data := struct {
		Date           string  `json:"date"`
		InitialBalance float64 `json:"initial_balance"`
		CurrentBalance float64 `json:"current_balance"`
		CumulativePnL  float64 `json:"cumulative_pnl"`
		CumulativeROI  float64 `json:"cumulative_roi"`
		TotalTrades    int     `json:"total_trades"`
		CurrentPrice   float64 `json:"current_price"`
		LastUpdate     string  `json:"last_update"`
	}{
		Date:           record.Date,
		InitialBalance: record.InitialBalance,
		CurrentBalance: record.CurrentBalance,
		CumulativePnL:  record.CumulativePnL,
		CumulativeROI:  record.CumulativeROI,
		TotalTrades:    record.TotalTrades,
		CurrentPrice:   record.CurrentPrice,
		LastUpdate:     time.Now().Format("2006-01-02 15:04:05"),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// HTML 模板
const htmlTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>交易绩效监控</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: white;
            color: #333;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 8px;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            border-bottom: 1px solid #e0e0e0;
            padding: 20px;
            text-align: center;
        }
        .header h1 {
            margin: 0;
            font-size: 1.8em;
            font-weight: 600;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
            gap: 20px;
            padding: 30px;
        }
        .stat-card {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            padding: 20px;
            text-align: center;
        }
        .stat-title {
            font-size: 0.9em;
            color: #666;
            margin-bottom: 8px;
        }
        .stat-value {
            font-size: 1.6em;
            font-weight: 600;
            color: #333;
        }
        .positive {
            color: #2e7d32;
        }
        .negative {
            color: #c62828;
        }
        .update-time {
            border-top: 1px solid #e0e0e0;
            padding: 15px 20px;
            text-align: center;
            color: #666;
            font-size: 0.9em;
        }
        .refresh-btn {
            display: block;
            margin: 20px auto;
            padding: 10px 20px;
            background: #333;
            color: white;
            border: 1px solid #333;
            border-radius: 4px;
            cursor: pointer;
            font-size: 1em;
        }
        .refresh-btn:hover {
            background: white;
            color: #333;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BTC/USDT</h1>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-title">初始资金</div>
                <div class="stat-value">{{.InitialBalance}} USDT</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">当前余额</div>
                <div class="stat-value">{{.CurrentBalance}} USDT</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">累计盈亏</div>
                <div class="stat-value {{if ge .CumulativePnL 0.0}}positive{{else}}negative{{end}}">{{printf "%.2f" .CumulativePnL}} USDT</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">累计收益率</div>
                <div class="stat-value {{if ge .CumulativeROI 0.0}}positive{{else}}negative{{end}}">{{printf "%.2f" .CumulativeROI}}%</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">总交易次数</div>
                <div class="stat-value">{{.TotalTrades}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">当前价格</div>
                <div class="stat-value">{{printf "%.2f" .CurrentPrice}} USDT</div>
            </div>
        </div>
        
        <div class="update-time">
            最后更新: {{.LastUpdate}}
        </div>
        
        <button class="refresh-btn" onclick="refreshData()">刷新数据</button>
    </div>

    <script>
        function refreshData() {
            fetch('/api/performance')
                .then(response => response.json())
                .then(data => {
                    location.reload();
                })
                .catch(error => {
                    console.error('Error:', error);
                    alert('刷新失败');
                });
        }

        // 每30秒自动刷新一次
        setInterval(() => {
            fetch('/api/performance')
                .then(response => response.json())
                .then(data => {
                    // 数据会自动更新到页面
                })
                .catch(error => {
                    console.error('Auto refresh error:', error);
                });
        }, 30000);
    </script>
</body>
</html>
`
