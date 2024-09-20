#!/usr/bin/env bash

set -euo pipefail

go build -o ./out/librespeed_exporter -tags 'netgo' -ldflags "-w -s -linkmode 'external' -extldflags '-static'" ./cmd/server

command -v upx > /dev/null && upx -qqq ./out/librespeed_exporter
