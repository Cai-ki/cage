package main

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"time"
)

func StartServer() {
	go func() {
		http.HandleFunc("/", handleIndex)
		http.HandleFunc("/api/performance", handleAPIPerformance)
		http.HandleFunc("/api/performance/history", handleAPIPerformanceHistory)

		log.Println("服务器启动在 :8080")
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	funcMap := template.FuncMap{
		"reverse": reverse,
		"sub": func(a, b int) int {
			return a - b
		},
	}
	tmpl := template.Must(template.New("index").Funcs(funcMap).Parse(htmlTemplate))

	records := GetAllPerformanceRecords()

	if len(records) == 0 {
		var err error
		records, err = LoadOrCreatePerformanceRecord()
		if err != nil {
			http.Error(w, "Performance data unavailable", http.StatusInternalServerError)
			return
		}
	}

	data := struct {
		Records []PerformanceRecord
	}{
		Records: records,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleAPIPerformance(w http.ResponseWriter, r *http.Request) {
	records := GetAllPerformanceRecords()

	var currentRecord *PerformanceRecord
	if len(records) > 0 {
		currentRecord = &records[len(records)-1]
	} else {
		var err error
		records, err = LoadOrCreatePerformanceRecord()
		if err != nil {
			http.Error(w, "Performance data unavailable", http.StatusInternalServerError)
			return
		}
		if len(records) > 0 {
			currentRecord = &records[len(records)-1]
		}
	}

	if currentRecord != nil {
		response := struct {
			Date           string  `json:"date"`
			InitialBalance float64 `json:"initial_balance"`
			CurrentBalance float64 `json:"current_balance"`
			CumulativePnL  float64 `json:"cumulative_pnl"`
			CumulativeROI  float64 `json:"cumulative_roi"`
			TotalTrades    int     `json:"total_trades"`
			CurrentPrice   float64 `json:"current_price"`
			LastUpdate     string  `json:"last_update"`
		}{
			Date:           currentRecord.Date,
			InitialBalance: currentRecord.InitialBalance,
			CurrentBalance: currentRecord.CurrentBalance,
			CumulativePnL:  currentRecord.CumulativePnL,
			CumulativeROI:  currentRecord.CumulativeROI,
			TotalTrades:    currentRecord.TotalTrades,
			CurrentPrice:   currentRecord.CurrentPrice,
			LastUpdate:     time.Now().Format("2006-01-02 15:04:05"),
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "No performance data", http.StatusNotFound)
	}
}

func handleAPIPerformanceHistory(w http.ResponseWriter, r *http.Request) {
	records := GetAllPerformanceRecords()

	history := make([]map[string]interface{}, len(records))
	for i, record := range records {
		history[i] = map[string]interface{}{
			"date":            record.Date,
			"initial_balance": record.InitialBalance,
			"current_balance": record.CurrentBalance,
			"cumulative_pnl":  record.CumulativePnL,
			"cumulative_roi":  record.CumulativeROI,
			"total_trades":    record.TotalTrades,
			"current_price":   record.CurrentPrice,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(history)
}

// 反转切片的辅助函数
func reverse(records []PerformanceRecord) []PerformanceRecord {
	result := make([]PerformanceRecord, len(records))
	for i, j := 0, len(records)-1; i < len(records); i, j = i+1, j-1 {
		result[i] = records[j]
	}
	return result
}

// HTML 模板
const htmlTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>BTC/USDT AI交易监控</title>
    <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
    <script src="https://cdn.jsdelivr.net/npm/marked/marked.min.js"></script>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            margin: 0;
            padding: 20px;
            background: white;
            color: #333;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        .header {
            border-bottom: 1px solid #e0e0e0;
            padding: 20px;
            text-align: center;
            margin-bottom: 20px;
        }
        .header h1 {
            margin: 0;
            font-size: 1.8em;
            font-weight: 600;
        }
        .stats-grid {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
            gap: 15px;
            margin-bottom: 30px;
        }
        .stat-card {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            padding: 15px;
            text-align: center;
        }
        .stat-title {
            font-size: 0.85em;
            color: #666;
            margin-bottom: 8px;
        }
        .stat-value {
            font-size: 1.4em;
            font-weight: 600;
            color: #333;
        }
        .positive {
            color: #2e7d32;
        }
        .negative {
            color: #c62828;
        }
        .chart-container {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            padding: 20px;
            margin-bottom: 30px;
        }
        .chart-row {
            display: grid;
            grid-template-columns: 1fr;
            gap: 20px;
            margin-bottom: 20px;
        }
        @media (min-width: 768px) {
            .chart-row {
                grid-template-columns: 1fr 1fr;
            }
        }
        .chart-wrapper {
            height: 300px;
            position: relative;
        }
        .details-section {
            background: white;
            border: 1px solid #e0e0e0;
            border-radius: 6px;
            margin-bottom: 20px;
        }
        .details-header {
            background: #f5f5f5;
            padding: 15px;
            border-bottom: 1px solid #e0e0e0;
            cursor: pointer;
            display: flex;
            justify-content: space-between;
            align-items: center;
        }
        .details-header:hover {
            background: #eee;
        }
        .details-content {
            padding: 15px;
            display: none;
        }
        .details-content.expanded {
            display: block;
        }
        .toggle-btn {
            background: #333;
            color: white;
            border: none;
            padding: 5px 10px;
            border-radius: 4px;
            cursor: pointer;
            font-size: 0.8em;
        }
        .toggle-btn:hover {
            background: #555;
        }
        .detail-item {
            margin-bottom: 10px;
            padding: 8px 0;
            border-bottom: 1px solid #f0f0f0;
        }
        .detail-label {
            font-weight: 600;
            color: #555;
            margin-bottom: 4px;
        }
        .detail-value {
            color: #333;
            word-break: break-word;
        }
        .markdown-content {
            line-height: 1.6;
            font-size: 0.9em;
        }
        .markdown-content h1, .markdown-content h2, .markdown-content h3 {
            margin-top: 16px;
            margin-bottom: 8px;
        }
        .markdown-content p {
            margin: 8px 0;
        }
        .markdown-content code {
            background: #f5f5f5;
            padding: 2px 4px;
            border-radius: 3px;
            font-family: monospace;
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
        .no-data {
            text-align: center;
            padding: 40px;
            color: #666;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>BTC/USDT AI交易监控</h1>
        </div>
        
        <div class="stats-grid">
            <div class="stat-card">
                <div class="stat-title">初始资金</div>
                <div class="stat-value">{{if .Records}}{{(index .Records (sub (len .Records) 1)).InitialBalance}} USDT{{else}}0.00 USDT{{end}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">当前余额</div>
                <div class="stat-value">{{if .Records}}{{(index .Records (sub (len .Records) 1)).CurrentBalance}} USDT{{else}}0.00 USDT{{end}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">累计盈亏</div>
                <div class="stat-value {{if and .Records (ge (index .Records (sub (len .Records) 1)).CumulativePnL 0.0)}}positive{{else}}negative{{end}}">{{if .Records}}{{printf "%.2f" (index .Records (sub (len .Records) 1)).CumulativePnL}} USDT{{else}}0.00 USDT{{end}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">累计收益率</div>
                <div class="stat-value {{if and .Records (ge (index .Records (sub (len .Records) 1)).CumulativeROI 0.0)}}positive{{else}}negative{{end}}">{{if .Records}}{{printf "%.2f" (index .Records (sub (len .Records) 1)).CumulativeROI}}%{{else}}0.00%{{end}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">总交易次数</div>
                <div class="stat-value">{{if .Records}}{{(index .Records (sub (len .Records) 1)).TotalTrades}}{{else}}0{{end}}</div>
            </div>
            
            <div class="stat-card">
                <div class="stat-title">当前价格</div>
                <div class="stat-value">{{if .Records}}{{printf "%.2f" (index .Records (sub (len .Records) 1)).CurrentPrice}} USDT{{else}}0.00 USDT{{end}}</div>
            </div>
        </div>
        
        <div class="chart-container">
            <h3>绩效变化趋势</h3>
            <div class="chart-row">
                <div class="chart-wrapper">
                    <canvas id="balanceChart"></canvas>
                </div>
                <div class="chart-wrapper">
                    <canvas id="pnlChart"></canvas>
                </div>
            </div>
            <div class="chart-row">
                <div class="chart-wrapper">
                    <canvas id="roiChart"></canvas>
                </div>
                <div class="chart-wrapper">
                    <canvas id="priceChart"></canvas>
                </div>
            </div>
        </div>
        
        <h3>交易记录详情</h3>
        {{if .Records}}
            {{range $index, $record := (reverse .Records)}}
            <div class="details-section">
                <div class="details-header" onclick="toggleDetails({{$index}})">
                    <span>{{$record.Date}} - 余额: {{printf "%.2f" $record.CurrentBalance}} USDT</span>
                    <button class="toggle-btn">展开/收起</button>
                </div>
                <div class="details-content" id="details-{{$index}}">
                    <div class="detail-item">
                        <div class="detail-label">初始资金</div>
                        <div class="detail-value">{{$record.InitialBalance}} USDT</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">当前余额</div>
                        <div class="detail-value">{{$record.CurrentBalance}} USDT</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">累计盈亏</div>
                        <div class="detail-value {{if ge $record.CumulativePnL 0.0}}positive{{else}}negative{{end}}">{{printf "%.2f" $record.CumulativePnL}} USDT</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">累计收益率</div>
                        <div class="detail-value {{if ge $record.CumulativeROI 0.0}}positive{{else}}negative{{end}}">{{printf "%.2f" $record.CumulativeROI}}%</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">总交易次数</div>
                        <div class="detail-value">{{$record.TotalTrades}}</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">当前价格</div>
                        <div class="detail-value">{{printf "%.2f" $record.CurrentPrice}} USDT</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">Prompt</div>
                        <div class="detail-value markdown-content" id="prompt-{{$index}}">{{$record.Prompt}}</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">Decision</div>
                        <div class="detail-value markdown-content" id="decision-{{$index}}">{{$record.Decision}}</div>
                    </div>
                    <div class="detail-item">
                        <div class="detail-label">Tool Calls</div>
                        <div class="detail-value">{{$record.ToolCalls}}</div>
                    </div>
                </div>
            </div>
            {{end}}
        {{else}}
            <div class="no-data">暂无交易记录</div>
        {{end}}
        
        <button class="refresh-btn" onclick="refreshData()">刷新数据</button>
    </div>

    <script>
        // 获取历史数据并绘制图表
        async function loadChartData() {
            try {
                const response = await fetch('/api/performance/history');
                const history = await response.json();
                
                if (history.length === 0) {
                    return;
                }

                // 使用所有记录，而不是按日期分组
                const labels = history.map(item => item.date);
                const balanceData = history.map(item => item.current_balance);
                const pnlData = history.map(item => item.cumulative_pnl);
                const roiData = history.map(item => item.cumulative_roi);
                const priceData = history.map(item => item.current_price);

                // 创建余额图表
                const balanceCtx = document.getElementById('balanceChart').getContext('2d');
                new Chart(balanceCtx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: '余额 (USDT)',
                            data: balanceData,
                            borderColor: 'rgb(75, 192, 192)',
                            backgroundColor: 'rgba(75, 192, 192, 0.1)',
                            tension: 0.1,
                            fill: false
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                            y: {
                                beginAtZero: false
                            }
                        }
                    }
                });

                // 创建盈亏图表
                const pnlCtx = document.getElementById('pnlChart').getContext('2d');
                new Chart(pnlCtx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: '累计盈亏 (USDT)',
                            data: pnlData,
                            borderColor: 'rgb(255, 99, 132)',
                            backgroundColor: 'rgba(255, 99, 132, 0.1)',
                            tension: 0.1,
                            fill: false
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                            y: {
                                beginAtZero: false
                            }
                        }
                    }
                });

                // 创建收益率图表
                const roiCtx = document.getElementById('roiChart').getContext('2d');
                new Chart(roiCtx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: '累计收益率 (%)',
                            data: roiData,
                            borderColor: 'rgb(54, 162, 235)',
                            backgroundColor: 'rgba(54, 162, 235, 0.1)',
                            tension: 0.1,
                            fill: false
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                            y: {
                                beginAtZero: false
                            }
                        }
                    }
                });

                // 创建价格图表
                const priceCtx = document.getElementById('priceChart').getContext('2d');
                new Chart(priceCtx, {
                    type: 'line',
                    data: {
                        labels: labels,
                        datasets: [{
                            label: '价格 (USDT)',
                            data: priceData,
                            borderColor: 'rgb(153, 102, 255)',
                            backgroundColor: 'rgba(153, 102, 255, 0.1)',
                            tension: 0.1,
                            fill: false
                        }]
                    },
                    options: {
                        responsive: true,
                        maintainAspectRatio: false,
                        scales: {
                            y: {
                                beginAtZero: false
                            }
                        }
                    }
                });

            } catch (error) {
                console.error('Error loading chart data:', error);
            }
        }

        function toggleDetails(index) {
            const content = document.getElementById('details-' + index);
            if (content) {
                content.classList.toggle('expanded');
                
                // 渲染Markdown内容（只在第一次展开时渲染）
                if (content.classList.contains('expanded')) {
                    const promptEl = document.getElementById('prompt-' + index);
                    const decisionEl = document.getElementById('decision-' + index);
                    
                    if (promptEl && promptEl.innerHTML.trim() && !promptEl.dataset.rendered) {
                        promptEl.innerHTML = marked.parse(promptEl.textContent);
                        promptEl.dataset.rendered = 'true';
                    }
                    
                    if (decisionEl && decisionEl.innerHTML.trim() && !decisionEl.dataset.rendered) {
                        decisionEl.innerHTML = marked.parse(decisionEl.textContent);
                        decisionEl.dataset.rendered = 'true';
                    }
                }
            }
        }

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

        // 页面加载完成后加载图表
        document.addEventListener('DOMContentLoaded', function() {
            loadChartData();
        });
    </script>
</body>
</html>
`
