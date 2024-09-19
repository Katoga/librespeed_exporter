package server

import (
	"log"
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type server struct {
	reg *prometheus.Registry
}

func NewServer(registry *prometheus.Registry) *server {
	s := &server{
		reg: registry,
	}

	return s
}

func (s *server) Serve(
	listenAddress *net.TCPAddr,
	telemetryPath *string,
) error {
	http.Handle(*telemetryPath, promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{Registry: s.reg}))

	log.Printf("listening on %v", listenAddress.String())
	log.Printf("serving metrics on %s", *telemetryPath)

	return http.ListenAndServe(listenAddress.String(), nil)
}
