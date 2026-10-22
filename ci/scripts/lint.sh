#!/bin/bash -eux

cwd=$(pwd)

pushd $cwd/dp-cantabular-filter-flex-api
  make lint
popd
