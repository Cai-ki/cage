package main

import (
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type Record struct {
	Timestamp string  `json:"timestamp"`
	USDT      float64 `json:"usdt"`
	BTC       float64 `json:"btc"`
	NetValue  float64 `json:"net_value"`
}

type HTMLServer struct {
	csvPath string
	data    []Record
	mu      sync.RWMutex
	running bool
}

func NewHTMLServer(csvPath string) *HTMLServer {
	return &HTMLServer{
		csvPath: csvPath,
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
    .container { max-width: 1000px; margin: 0 auto; }
    h2 { text-align: center; color: #333; }
    .chart-container { background: white; padding: 15px; border-radius: 8px; box-shadow: 0 2px 6px rgba(0,0,0,0.1); height: 500px; }
  </style>
</head>
<body>
  <div class="container">
    <h2>Portfolio Overview</h2>
    <div class="chart-container"><canvas id="portfolioChart"></canvas></div>
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
            title: { display: true, text: 'Amount' },
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

    function fetchData() {
      fetch('/data')
        .then(r => r.ok ? r.json() : Promise.reject('Network error'))
        .then(data => {
          const labels = data.map(d => d.timestamp);
          chart.data.labels = labels;
          chart.data.datasets[0].data = data.map(d => d.net_value);
          chart.data.datasets[1].data = data.map(d => d.usdt);
          chart.data.datasets[2].data = data.map(d => d.btc);
          chart.update('none');
        })
        .catch(e => console.error('Fetch failed:', e));
    }

    fetchData();
    setInterval(fetchData, 1000);
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

	go func() {
		log.Printf("[HTMLServer] Serving dashboard at http://%s", addr)
		if err := http.ListenAndServe(addr, nil); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[HTMLServer] Server error: %v", err)
		}
	}()
}

func (s *HTMLServer) Update() {
	log.Printf("[Update] Attempting to read CSV: %s", s.csvPath)

	file, err := os.Open(s.csvPath)
	if err != nil {
		log.Printf("[Update] FAILED to open CSV: %v", err)
		return
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	records, err := reader.ReadAll()
	if err != nil {
		log.Printf("[Update] FAILED to parse CSV: %v", err)
		return
	}

	log.Printf("[Update] Read %d raw rows", len(records))

	if len(records) > 0 {
		records = records[1:]
	}

	var data []Record
	for i, row := range records {
		if len(row) < 4 {
			log.Printf("[Update] Skip invalid row %d: %v", i, row)
			continue
		}
		usdt, err1 := strconv.ParseFloat(row[1], 64)
		btc, err2 := strconv.ParseFloat(row[2], 64)
		net, err3 := strconv.ParseFloat(row[3], 64)
		if err1 != nil || err2 != nil || err3 != nil {
			log.Printf("[Update] Parse error on row %d: %v (usdt=%v, btc=%v, net=%v)", i, row, err1, err2, err3)
			continue
		}
		data = append(data, Record{
			Timestamp: row[0],
			USDT:      usdt,
			BTC:       btc,
			NetValue:  net,
		})
	}

	s.mu.Lock()
	s.data = data
	s.mu.Unlock()

	log.Printf("[Update] SUCCESS: loaded %d valid records", len(data))
}
