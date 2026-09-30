#!/bin/bash -ex

DIR=$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )
cd ${DIR}

SNAP_DIR=${DIR}/../build/snap

go test ./...

for hook in install configure pre-refresh post-refresh; do
  CGO_ENABLED=0 go build -o ${SNAP_DIR}/meta/hooks/${hook} ./cmd/${hook}
done

CGO_ENABLED=0 go build -o ${SNAP_DIR}/bin/cli ./cmd/cli
