package collector

import "github.com/prometheus/client_golang/prometheus"

type collector struct {
}

func NewCollector() *collector {
	c := &collector{}

	return c
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {

}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {

}
