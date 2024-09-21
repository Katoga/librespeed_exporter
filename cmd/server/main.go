package main

import (
	"os"

	"github.com/Katoga/librespeed_exporter/internal/collector"
	"github.com/Katoga/librespeed_exporter/internal/server"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	listenAddress := kingpin.Flag("web.listen-address", "Address to listen on").Default(":51423").TCP()
	telemetryPath := kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	enableCollectorGo := kingpin.Flag("collectors.go", "Enable GoCollector").Bool()
	enableCollectorProcess := kingpin.Flag("collectors.process", "Enable ProcessCollector").Bool()
	dataRetrieverCommand := kingpin.Flag("data-retriever-command", "Command to call to get speed data").Default("speedtest-cli").ExistingFile()
	librespeedServer := kingpin.Flag("librespeed.server", "Librespeed server to get speed data (zero means 'none specified')").Default("0").Uint8()

	kingpin.Parse()

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		collector.NewCollector(log, dataRetrieverCommand, librespeedServer),
	)
	if *enableCollectorGo {
		registry.MustRegister(
			collectors.NewGoCollector(),
		)
	}
	if *enableCollectorProcess {
		registry.MustRegister(
			collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		)
	}

	log.Info().Msg("serving")

	error := server.NewServer(log, registry).Serve(*listenAddress, telemetryPath)
	log.Fatal().Msg(error.Error())
}
