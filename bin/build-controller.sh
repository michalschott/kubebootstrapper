#!/usr/bin/env sh

set -eu

(
    target=controller/controller
    tag=$(git rev-parse HEAD)
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $target -mod=readonly -ldflags "-s -w -X github.com/michalschott/kubebootstrapper/pkg/version.Version=$tag" ./controller/cmd
    echo $target
    docker build controller -t controller:$tag
    echo "docker run --rm controller:$tag"
    kind load docker-image controller:$tag
    kubectl -n kubebootstrapper delete pod --all
)
