package main

import (
	"log"

	"github.com/Katoga/librespeed_exporter/internal/server"

	"github.com/alecthomas/kingpin/v2"
)

func main() {
	listenAddress := kingpin.Flag("web.listen-address", "Address to listen on").Default(":51423").String()
	telemetryPath := kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	enableCollectorGo := kingpin.Flag("collectors.go", "Enable GoCollector").Bool()
	enableCollectorProcess := kingpin.Flag("collectors.process", "Enable ProcessCollector").Bool()
	dataRetrieverCommand := kingpin.Flag("data-retriever-command", "Command to call to get speed data").Default("speedtest-cli").String()

	kingpin.Parse()

	log.Fatal(server.NewServer().Serve(listenAddress, telemetryPath, enableCollectorGo, enableCollectorProcess, dataRetrieverCommand))
}
