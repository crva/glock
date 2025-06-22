# glock

## Roadmap

### ✅ MVP: Send HTTP Requests

- [x] Send N `GET` requests to a given URL
- [x] Display total execution time
- [x] Count successful and failed requests
- [x] Measure average, min, and max response times

---

### 🔄 Concurrency

- [x] Flag to set the number of goroutines
- [x] Distribution of requests across workers
- [x] Use of `sync.WaitGroup`
- [x] Parallel measurement of response times

---

### ⚙️ Parameterization

- [ ] `-method` flag (GET, POST, etc.)
- [ ] `-data` flag (JSON request body)
- [ ] `-headers` flag (multiple HTTP headers)
- [ ] Input validation (URL format, JSON syntax…)

---

### 📈 Performance Report

- [ ] Detailed report:
  - Total number of requests
  - Success / failure rate
  - Min / max / average response time
  - Percentiles (p95, p99)
- [ ] Export results to `.json` or `.csv`

---

### 🌐 Local HTTP Server

- [ ] Embedded HTTP server (`localhost:8080`)
- [ ] `/metrics` route exposing stats in JSON
- [ ] Live updates of data

---

### 🎨 React Frontend

- [ ] React app with visual charts (Chart.js, Recharts, etc.)
- [ ] Displays:
  - Requests per second over time
  - Response latency
  - Error rate
- [ ] Connects to Go backend via `/metrics`
