#!/bin/bash -eux

export cwd=$(pwd)

pushd $cwd/dp-cantabular-filter-flex-api
  make audit
popd
