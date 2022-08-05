package prometheus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"therealbroker/pkg/log"
)

type APIMetrics struct {
	MethodCount       *prometheus.CounterVec
	MethodDuration    *prometheus.SummaryVec
	ActiveSubscribers prometheus.Gauge
}

func NewPrometheusServer(log *log.Logger, port string) *APIMetrics {
	metrics := &APIMetrics{
		MethodCount: promauto.NewCounterVec(prometheus.CounterOpts{
			Name: "method_count",
			Help: "The total number of rpc calls",
		}, []string{"method", "status"}),
		MethodDuration: prometheus.NewSummaryVec(
			prometheus.SummaryOpts{
				Name:       "method_durations_nanoseconds",
				Help:       "Total Rpc latency.",
				Objectives: map[float64]float64{0.5: 0.05, 0.95: 0.005, 0.99: 0.001},
			}, []string{"method", "status"}),
		ActiveSubscribers: promauto.NewGauge(prometheus.GaugeOpts{
			Name: "active_subscribers",
			Help: "The total number of active streams",
		}),
	}
	prometheus.MustRegister(metrics.MethodDuration)
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+port, nil)
		if err != nil {
			log.Errorln(err)
			return
		}
	}()
	return metrics
}
