package metrics

import (
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type NamedHistogram struct {
	Name      string
	Histogram *prometheus.HistogramVec
}

func NewNamedHistogram(name string, buckets []float64) *NamedHistogram {

	hist := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Buckets: buckets,
		},
		[]string{"handler", "status"},
	)

	prometheus.MustRegister(hist)

	return &NamedHistogram{
		Name:      name,
		Histogram: hist,
	}

}

func (nh *NamedHistogram) Observe(name string, startTime time.Time) {

	duration := time.Since(startTime)
	nh.Histogram.WithLabelValues(name).Observe(duration.Seconds())

}

func (nh *NamedHistogram) ObserveHandler(name string, startTime time.Time, status int) {

	duration := time.Since(startTime)
	nh.Histogram.WithLabelValues(name, strconv.Itoa(status)).Observe(duration.Seconds())

}
