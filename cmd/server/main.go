package main

import (
	"log"

	"github.com/Katoga/librespeed_exporter/internal/server"

	"github.com/alecthomas/kingpin/v2"
)

func main() {
	listenAddress := kingpin.Flag("web.listen-address", "Address to listen on").Default(":51423").String()
	includeSystemCollectors := false

	kingpin.Parse()

	log.Fatal(server.NewServer().Serve(listenAddress, includeSystemCollectors))
}
