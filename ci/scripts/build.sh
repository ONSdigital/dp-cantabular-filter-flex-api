#!/bin/bash -eux

pushd dp-cantabular-filter-flex-api
  make build
  cp build/dp-cantabular-filter-flex-api Dockerfile.concourse ../build
popd
