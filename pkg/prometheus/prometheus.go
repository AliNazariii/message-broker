package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
)

func StartPrometheusServer(port string) {
	http.Handle("/metrics", promhttp.Handler())

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		logrus.Errorln(err)
	}
}
