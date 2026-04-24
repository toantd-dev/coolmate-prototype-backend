package handlers

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

type MetricsHandler struct {
	startTime time.Time
}

func NewMetricsHandler() *MetricsHandler {
	return &MetricsHandler{
		startTime: time.Now(),
	}
}

func (m *MetricsHandler) ServeMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	uptime := time.Since(m.startTime).Seconds()

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	metrics := fmt.Sprintf(
		"# HELP coolmate_uptime_seconds Service uptime in seconds\n"+
			"# TYPE coolmate_uptime_seconds gauge\n"+
			"coolmate_uptime_seconds %.0f\n"+
			"# HELP process_resident_memory_bytes Resident memory in bytes\n"+
			"# TYPE process_resident_memory_bytes gauge\n"+
			"process_resident_memory_bytes %d\n"+
			"# HELP go_goroutines Number of goroutines\n"+
			"# TYPE go_goroutines gauge\n"+
			"go_goroutines %d\n",
		uptime,
		memStats.Alloc,
		runtime.NumGoroutine(),
	)

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(metrics))
}
