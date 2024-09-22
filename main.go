package main

import (
	"os"

	"github.com/Katoga/librespeed_exporter/cmd/librespeed_exporter"

	"github.com/rs/zerolog"
)

func main() {
	log := zerolog.New(os.Stderr).With().Timestamp().Logger()

	error := librespeed_exporter.NewLibrespeedExporter(log).Run()

	log.Fatal().Err(error)
}
