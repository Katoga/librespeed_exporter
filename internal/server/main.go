package server

import (
	"net"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
)

type server struct {
	reg *prometheus.Registry
	log zerolog.Logger
}

func NewServer(log zerolog.Logger, registry *prometheus.Registry) *server {
	s := &server{
		reg: registry,
		log: log.With().Str("component", "server").Logger(),
	}

	return s
}

func (s *server) Serve(
	listenAddress *net.TCPAddr,
	telemetryPath *string,
) error {
	http.Handle(*telemetryPath, promhttp.HandlerFor(s.reg, promhttp.HandlerOpts{Registry: s.reg}))

	s.log.Info().Msgf("listening on %v", listenAddress.String())
	s.log.Info().Msgf("serving metrics on %s", *telemetryPath)

	return http.ListenAndServe(listenAddress.String(), nil)
}
