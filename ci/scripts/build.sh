#!/bin/bash -eux

pushd dp-frontend-interactives-controller
  make build
  cp build/dp-frontend-interactives-controller Dockerfile.concourse ../build
popd
