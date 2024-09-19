package main

import (
	"log"

	"github.com/Katoga/librespeed_exporter/internal/collector"
	"github.com/Katoga/librespeed_exporter/internal/server"

	"github.com/alecthomas/kingpin/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

func main() {
	listenAddress := kingpin.Flag("web.listen-address", "Address to listen on").Default(":51423").String()
	telemetryPath := kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	enableCollectorGo := kingpin.Flag("collectors.go", "Enable GoCollector").Bool()
	enableCollectorProcess := kingpin.Flag("collectors.process", "Enable ProcessCollector").Bool()
	dataRetrieverCommand := kingpin.Flag("data-retriever-command", "Command to call to get speed data").Default("speedtest-cli").String()

	kingpin.Parse()

	registry := prometheus.NewPedanticRegistry()
	registry.MustRegister(
		collector.NewCollector(*dataRetrieverCommand),
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

	log.Fatal(server.NewServer(registry).Serve(listenAddress, telemetryPath))
}
