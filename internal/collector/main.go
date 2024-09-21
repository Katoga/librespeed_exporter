package collector

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
)

type responseItem struct {
	Ping     float64 `json:"ping"`
	Jitter   float64 `json:"jitter"`
	Upload   float64 `json:"upload"`
	Download float64 `json:"download"`
}

type results struct {
	Upload   float64
	Download float64
	Ping     float64
	Jitter   float64
}

type collector struct {
	dataRetrieverCommand *string
	dataRetrieverArgs    []string
}

func NewCollector(dataRetrieverCommand *string, librespeedServer *uint8) *collector {
	c := &collector{
		dataRetrieverCommand: dataRetrieverCommand,
	}

	dataRetrieverArgs := []string{
		"--json",
	}
	if librespeedServer != nil {
		c.dataRetrieverArgs = append(dataRetrieverArgs, []string{"--server", fmt.Sprintf("%d", *librespeedServer)}...)
	}

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
	content, errDownload := c.download()
	if errDownload != nil {
		log.Panicf("Getting libresped data failed: %s", errDownload)
	}

	response := []responseItem{}
	errJson := json.Unmarshal(content, &response)
	if errJson != nil {
		log.Panicf("Parsing JSON failed: %s", errJson)
	}

	res := response[0]

	return results{
		Upload:   res.Upload,
		Download: res.Download,
		Ping:     res.Ping,
		Jitter:   res.Jitter,
	}
}

func (c *collector) download() ([]byte, error) {

	cmd := exec.Command(*c.dataRetrieverCommand, c.dataRetrieverArgs...)
	output, errRun := cmd.Output()
	if errRun != nil {
		log.Panicf("Command failed: %s", errRun)
	}

	return output, nil
}
