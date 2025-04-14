package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func GetHandler(registry *prometheus.Registry) http.Handler {

	return promhttp.HandlerFor(registry, promhttp.HandlerOpts{})

}
