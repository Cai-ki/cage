package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type Record struct {
	Timestamp string  `json:"timestamp"`
	USDT      float64 `json:"usdt"`
	BTC       float64 `json:"btc"`
	NetValue  float64 `json:"net_value"`
}

type PriceRecord struct {
	Timestamp string  `json:"timestamp"`
	Price     float64 `json:"price"`
}

type AIHistoryRecord struct {
	Timestamp   string `json:"timestamp"`
	Symbol      string `json:"symbol"`
	Action      string `json:"action"`
	Price       string `json:"price"`
	Reason      string `json:"reason"`
	AnalystView string `json:"analyst_view"`
	RiskView    string `json:"risk_view"`
}

type HTMLServer struct {
	csvPath            string
	tradingHistoryPath string
	data               []Record
	priceData          []PriceRecord
	aiHistory          []AIHistoryRecord
	mu                 sync.RWMutex
	running            bool
}

func NewHTMLServer(csvPath, tradingHistoryPath string) *HTMLServer {
	return &HTMLServer{
		csvPath:            csvPath,
		tradingHistoryPath: tradingHistoryPath,
	}
}

const indexHTML = `<!DOCTYPE html>
<html>
<head>
  <meta charset="utf-8">
  <title>Portfolio Dashboard</title>
  <script src="https://cdn.jsdelivr.net/npm/chart.js"></script>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; padding: 20px; background: #f9f9f9; }
    .container { max-width: 1200px; margin: 0 auto; }
    h2 { text-align: center; color: #333; margin-bottom: 20px; }
    .chart-container { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); height: 500px; margin-bottom: 20px; }
    .history-panel { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); max-height: 400px; overflow-y: auto; }
    .history-item { padding: 12px 0; border-bottom: 1px solid #eee; }
    .history-item:last-child { border-bottom: none; }
    .history-time { font-weight: bold; color: #555; }
    .history-action { display: inline-block; margin-left: 10px; padding: 2px 6px; border-radius: 4px; font-size: 0.85em; }
    .history-action.hold { background: #e3f2fd; color: #1976d2; }
    .history-action.buy,
    .history-action.sell { background: #e8f5e9; color: #2e7d32; }
    .history-price { color: #ff9800; font-weight: bold; }
    .history-reason { margin-top: 6px; font-size: 0.9em; color: #666; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Portfolio Dashboard</h2>
    <div class="chart-container"><canvas id="portfolioChart"></canvas></div>

    <h3>AI Decision History</h3>
    <div class="history-panel" id="historyList">
      <!-- Filled by JS -->
    </div>
  </div>

  <script>
    const ctx = document.getElementById('portfolioChart').getContext('2d');
    const chart = new Chart(ctx, {
      type: 'line',
      data: {
        labels: [],
        datasets: [
          {
            label: 'Net Value (USDT)',
            data: [],
            borderColor: '#4CAF50',
            backgroundColor: 'rgba(76, 175, 80, 0.1)',
            borderWidth: 2,
            tension: 0.3,
            fill: false
          },
          {
            label: 'USDT Balance',
            data: [],
            borderColor: '#2196F3',
            backgroundColor: 'rgba(33, 150, 243, 0.1)',
            borderWidth: 2,
            tension: 0.3,
            fill: false
          },
          {
            label: 'BTC Balance',
            data: [],
            borderColor: '#FF9800',
            backgroundColor: 'rgba(255, 152, 0, 0.1)',
            borderWidth: 2,
            tension: 0.3,
            fill: false
          },
          {
            label: 'BTC Price (USDT)',
            data: [],
            borderColor: '#9C27B0',
            backgroundColor: 'rgba(156, 39, 176, 0.1)',
            borderWidth: 2,
            tension: 0.3,
            fill: false,
            yAxisID: 'y1'
          }
        ]
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        interaction: { mode: 'index', intersect: false },
        scales: {
          x: {
            title: { display: true, text: 'Time' },
            ticks: { maxRotation: 45, autoSkip: true }
          },
          y: {
            title: { display: true, text: 'Portfolio Amount' },
            beginAtZero: false
          },
          y1: {
            position: 'right',
            title: { display: true, text: 'BTC Price (USDT)' },
            grid: { drawOnChartArea: false },
            beginAtZero: false
          }
        },
        plugins: {
          legend: { position: 'top' },
          tooltip: { mode: 'index' }
        },
        animation: { duration: 0 }
      }
    });

    async function fetchData() {
      try {
        const [portfolioRes, priceRes, historyRes] = await Promise.all([
          fetch('/data').then(r => r.ok ? r.json() : Promise.reject('portfolio')),
          fetch('/price-data').then(r => r.ok ? r.json() : Promise.reject('price')),
          fetch('/ai-history').then(r => r.ok ? r.json() : Promise.reject('history'))
        ]);

        // Update chart
        const labels = portfolioRes.map(d => d.timestamp);
        chart.data.labels = labels;
        chart.data.datasets[0].data = portfolioRes.map(d => d.net_value);
        chart.data.datasets[1].data = portfolioRes.map(d => d.usdt);
        chart.data.datasets[2].data = portfolioRes.map(d => d.btc);
        chart.data.datasets[3].data = priceRes.map(d => d.price);
        chart.update('none');

        // Update history list
        const historyList = document.getElementById('historyList');
        historyList.innerHTML = historyRes.map(function(item) {
          return '<div class="history-item">' +
            '<div>' +
              '<span class="history-time">' + new Date(item.timestamp).toLocaleString() + '</span>' +
              '<span class="history-action ' + item.action + '">' + item.action.toUpperCase() + '</span>' +
              '<span class="history-price">@ ' + item.price + ' USDT</span>' +
            '</div>' +
            '<div class="history-reason">' + (item.reason || '') + '</div>' +
            '<details style="margin-top:6px; font-size:0.85em; color:#555;">' +
              '<summary style="cursor:pointer; user-select:none;">🔍 Analyst & Risk Details</summary>' +
              '<div style="margin-top:4px; padding-left:10px; border-left:2px solid #ddd;">' +
                '<div><strong>Analyst:</strong> ' + (item.analyst_view || '') + '</div>' +
                '<div style="margin-top:6px;"><strong>Risk:</strong> ' + (item.risk_view || '') + '</div>' +
              '</div>' +
            '</details>' +
          '</div>';
        }).join('');
      } catch (e) {
        console.error('Fetch failed:', e);
      }
    }

    fetchData();
    setInterval(fetchData, 2000);
  </script>
</body>
</html>`

func (s *HTMLServer) Start(addr string) {
	if s.running {
		return
	}
	s.running = true

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(indexHTML))
		} else {
			http.NotFound(w, r)
		}
	})

	http.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.data)
	})

	http.HandleFunc("/price-data", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.priceData)
	})

	http.HandleFunc("/ai-history", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(s.aiHistory)
	})

	go func() {
		log.Printf("[HTMLServer] Serving dashboard at http://%s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTMLServer] Server error: %v", err)
		}
	}()
}

func (s *HTMLServer) Update() {
	log.Printf("[Update] Reading CSV: %s", s.csvPath)

	// --- Parse CSV ---
	var data []Record
	if file, err := os.Open(s.csvPath); err == nil {
		defer file.Close()
		reader := csv.NewReader(file)
		reader.FieldsPerRecord = -1
		records, err := reader.ReadAll()
		if err != nil {
			log.Printf("[Update] CSV parse error: %v", err)
		} else {
			if len(records) > 0 {
				records = records[1:] // skip header
			}
			for _, row := range records {
				if len(row) < 4 {
					continue
				}
				usdt, err1 := strconv.ParseFloat(row[1], 64)
				btc, err2 := strconv.ParseFloat(row[2], 64)
				net, err3 := strconv.ParseFloat(row[3], 64)
				if err1 != nil || err2 != nil || err3 != nil {
					continue
				}
				data = append(data, Record{
					Timestamp: row[0],
					USDT:      usdt,
					BTC:       btc,
					NetValue:  net,
				})
			}
		}
	}

	// --- Parse trading_history.json ---
	var priceData []PriceRecord
	var aiHistory []AIHistoryRecord

	if s.tradingHistoryPath != "" {
		if dataBytes, err := os.ReadFile(s.tradingHistoryPath); err == nil {
			var raw []map[string]interface{}
			if json.Unmarshal(dataBytes, &raw) == nil {
				for _, item := range raw {
					ts := item["timestamp"].(string)
					priceStr := item["price"].(string)
					price, _ := strconv.ParseFloat(priceStr, 64)

					priceData = append(priceData, PriceRecord{
						Timestamp: ts,
						Price:     price,
					})

					aiHistory = append(aiHistory, AIHistoryRecord{
						Timestamp:   ts,
						Symbol:      item["symbol"].(string),
						Action:      item["action"].(string),
						Price:       priceStr,
						Reason:      item["reason"].(string),
						AnalystView: item["analyst_view"].(string),
						RiskView:    item["risk_view"].(string),
					})
				}
			}
		}
	}

	// Sort CSV data by timestamp (assume RFC3339)
	sort.Slice(data, func(i, j int) bool {
		ti, _ := time.Parse(time.RFC3339Nano, data[i].Timestamp)
		tj, _ := time.Parse(time.RFC3339Nano, data[j].Timestamp)
		return ti.Before(tj)
	})

	// Reverse AI history so latest is first
	reversedAI := make([]AIHistoryRecord, len(aiHistory))
	for i, r := range aiHistory {
		reversedAI[len(reversedAI)-1-i] = r
	}

	// --- Update under lock ---
	s.mu.Lock()
	s.data = data
	s.priceData = priceData
	s.aiHistory = reversedAI
	s.mu.Unlock()

	log.Printf("[Update] CSV records: %d, Price records: %d, AI history: %d",
		len(data), len(priceData), len(aiHistory))
}
