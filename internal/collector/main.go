package collector

import (
	"math/rand/v2"

	"github.com/prometheus/client_golang/prometheus"
)

type results struct {
	Upload   float64
	Download float64
	Ping     float64
	Jitter   float64
}

type collector struct {
}

func NewCollector() *collector {
	c := &collector{}

	return c
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	results := c.getResults()

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_upload_bps",
			"Upload speed in bits per second",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		results.Upload,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_download_bps",
			"Download speed in bits per second",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		results.Download,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_ping_seconds",
			"Ping in seconds",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		results.Ping,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_jitter_seconds",
			"Jitter in seconds",
			nil,
			nil,
		),
		prometheus.GaugeValue,
		results.Jitter,
	)
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *collector) getResults() results {
	return results{
		Upload:   rand.Float64() * 100,
		Download: rand.Float64() * 100,
		Ping:     rand.Float64() * 10,
		Jitter:   rand.Float64() * 10,
	}
}
