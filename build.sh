#!/usr/bin/env bash

set -euo pipefail

go build -o ./out/ -tags 'netgo' -ldflags "-X main.version=$(date -Iseconds) -w -s -linkmode 'external' -extldflags '-static'"

command -v upx > /dev/null && upx -qqq ./out/librespeed_exporter
