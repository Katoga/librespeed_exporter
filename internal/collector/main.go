package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rs/zerolog"
)

type responseServer struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}
type responseItem struct {
	Ping     float64        `json:"ping"`
	Jitter   float64        `json:"jitter"`
	Upload   float64        `json:"upload"`
	Download float64        `json:"download"`
	Server   responseServer `json:"server"`
}

type results struct {
	Upload   float64
	Download float64
	Ping     float64
	Jitter   float64
	Server   responseServer
}

type collector struct {
	log                  zerolog.Logger
	dataRetrieverCommand *string
	dataRetrieverArgs    []string
}

func NewCollector(log zerolog.Logger, dataRetrieverCommand *string, librespeedServer *uint8) *collector {
	c := &collector{
		log:                  log.With().Str("component", "collector").Logger(),
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
	c.log.Info().Msg("collecting")
	results := c.getResults()

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_upload_bps",
			"Upload speed in bits per second",
			[]string{"server"},
			nil,
		),
		prometheus.GaugeValue,
		results.Upload*1000000,
		results.Server.Url,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_download_bps",
			"Download speed in bits per second",
			[]string{"server"},
			nil,
		),
		prometheus.GaugeValue,
		results.Download*1000000,
		results.Server.Url,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_ping_seconds",
			"Ping in seconds",
			[]string{"server"},
			nil,
		),
		prometheus.GaugeValue,
		results.Ping/1000,
		results.Server.Url,
	)

	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_jitter_seconds",
			"Jitter in seconds",
			[]string{"server"},
			nil,
		),
		prometheus.GaugeValue,
		results.Jitter/1000,
		results.Server.Url,
	)
}

func (c *collector) Describe(ch chan<- *prometheus.Desc) {
	prometheus.DescribeByCollect(c, ch)
}

func (c *collector) getResults() results {
	content, errDownload := c.download()
	if errDownload != nil {
		c.log.Panic().Err(errDownload).Msg("Getting libresped data failed")
	}

	response := []responseItem{}
	errJson := json.Unmarshal(content, &response)
	if errJson != nil {
		c.log.Panic().Err(errJson).Msg("Parsing JSON failed")
	}

	res := response[0]

	return results{
		Upload:   res.Upload,
		Download: res.Download,
		Ping:     res.Ping,
		Jitter:   res.Jitter,
		Server:   res.Server,
	}
}

func (c *collector) download() ([]byte, error) {
	c.log.Info().Msg("downloading")

	cmd := exec.Command(*c.dataRetrieverCommand, c.dataRetrieverArgs...)
	output, errRun := cmd.Output()
	if errRun != nil {
		c.log.Panic().Err(errRun).Msg("Command failed")
	}

	c.log.Info().Msg("downloaded")

	return output, nil
}
