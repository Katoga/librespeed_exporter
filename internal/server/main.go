package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct{}

func NewServer() *server {
	s := &server{}

	return s
}

func (s *server) Serve(port uint16, includeSystemCollectors bool) error {
	reg := prometheus.NewRegistry()

	if includeSystemCollectors {
		reg.MustRegister(
			collectors.NewGoCollector(),
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	http.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))

	log.Printf("listening on %d", port)

	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
