package main

import (
	"log"

	"github.com/Katoga/librespeed_exporter/internal/server"
)

func main() {
	port := uint16(51423)
	includeSystemCollectors := false

	log.Fatal(server.NewServer().Serve(port, includeSystemCollectors))
}
