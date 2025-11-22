package middleware

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	totalRequestsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_requests_received",
		Help: "The total number of requests received",
	})
	totalResponsesSent = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_responses_sent",
		Help: "The total number of responses sent",
	})
	totalProcessingTimeMicroseconds = promauto.NewCounter(prometheus.CounterOpts{
		Name: "total_processing_time_microseconds",
		Help: "The total (cumulative) time taken to process all requests in microseconds",
	})
	activeRequests = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "in_flight_requests",
		Help: "The number of 'active' in-flight requests",
	})
)

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		totalRequestsReceived.Inc()
		activeRequests.Inc()

		next.ServeHTTP(w, r)

		totalResponsesSent.Inc()
		activeRequests.Dec()

		duration := time.Since(start).Microseconds()
		totalProcessingTimeMicroseconds.Add(float64(duration))
	})
}
