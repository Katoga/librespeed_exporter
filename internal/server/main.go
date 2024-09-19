package server

import (
	"log"
	"net/http"

	"github.com/Katoga/librespeed_exporter/internal/collector"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct{}

func NewServer() *server {
	s := &server{}

	return s
}

func (s *server) Serve(listenAddress *string, telemetryPath *string, includeSystemCollectors bool) error {
	reg := prometheus.NewPedanticRegistry()

	reg.MustRegister(
		collector.NewCollector(),
	)

	if includeSystemCollectors {
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	http.Handle(*telemetryPath, promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	log.Printf("listening on %s", *listenAddress)
	log.Printf("serving metrics on %s", *telemetryPath)

	return http.ListenAndServe(*listenAddress, nil)
}
