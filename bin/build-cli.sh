#!/usr/bin/env sh

set -eu

(
    target=cli/kubebootstrapper
    tag=$(git rev-parse HEAD)
    CGO_ENABLED=0 go build -o $target -mod=readonly -ldflags "-s -w -X github.com/michalschott/kubebootstrapper/pkg/version.Version=$tag" ./cli
    echo $target
)