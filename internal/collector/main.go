package collector

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

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
	dataRetrieverArgs := []string{
		"--json",
	}
	if *librespeedServer != uint8(0) {
		dataRetrieverArgs = append(dataRetrieverArgs, []string{"--server", fmt.Sprintf("%d", *librespeedServer)}...)
	}

	c := &collector{
		log:                  log.With().Str("component", "collector").Logger(),
		dataRetrieverCommand: dataRetrieverCommand,
		dataRetrieverArgs:    dataRetrieverArgs,
	}

	return c
}

func (c *collector) Collect(ch chan<- prometheus.Metric) {
	c.log.Info().Msg("collecting")

	duration, results, errResults := c.getResults()
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_duration_seconds",
			"Duration of the mearurement in seconds",
			[]string{"server"},
			nil,
		),
		prometheus.GaugeValue,
		duration.Seconds(),
		results.Server.Url,
	)
	if errResults != nil {
		c.log.Error().Err(errResults).Msg("collecting failed")
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				"librespeed_success",
				"Success of speed measuring",
				[]string{},
				nil,
			),
			prometheus.GaugeValue,
			0,
		)
		return
	}

	c.log.Info().Msg("collecting succeeded")
	ch <- prometheus.MustNewConstMetric(
		prometheus.NewDesc(
			"librespeed_success",
			"Success of speed measuring",
			[]string{},
			nil,
		),
		prometheus.GaugeValue,
		1,
	)
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

func (c *collector) getResults() (time.Duration, results, error) {
	duration, content, errDownload := c.download()
	if errDownload != nil {
		return duration, results{}, errDownload
	}

	response := []responseItem{}
	errJson := json.Unmarshal(content, &response)
	if errJson != nil {
		return duration, results{}, errJson
	}

	res := response[0]

	return duration,
		results{
			Upload:   res.Upload,
			Download: res.Download,
			Ping:     res.Ping,
			Jitter:   res.Jitter,
			Server:   res.Server,
		}, nil
}

func (c *collector) download() (time.Duration, []byte, error) {
	c.log.Info().Msg("downloading")

	cmd := exec.Command(*c.dataRetrieverCommand, c.dataRetrieverArgs...)
	start := time.Now()
	output, errRun := cmd.Output()
	duration := time.Since(start)

	if errRun != nil {
		return duration, nil, errRun
	}

	c.log.Info().Msg("downloaded")

	return duration, output, nil
}
