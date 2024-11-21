package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	ActiveSubscribers = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "active_subscribers",
		Help: "The total number of active streams",
	})
)

func init() {
	prometheus.MustRegister(ActiveSubscribers)
}
