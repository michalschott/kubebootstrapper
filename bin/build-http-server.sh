#!/usr/bin/env sh
#
# NOTE:
# need kubectl 1.19+
#

set -eu

(
    tag=$(git rev-parse HEAD)
    docker build http -t http:$tag
    echo "docker run --rm http:$tag"
    kind load docker-image http:$tag
    kubectl delete ns kubebootstrapperserver || true
    kubectl create ns kubebootstrapperserver
    ./bin/kubectl -n kubebootstrapperserver create deployment http --image=http:$tag --port=8000
    kubectl -n kubebootstrapperserver expose deployment http --port=8000 --target-port=8000
)
