#!/bin/bash

set -eu

version="$1"

# trigger docker hub
curl -X "POST" "https://registry.hub.docker.com/u/kingori/sanaa/trigger/${DOCKER_HUB_TRIGGER_TOKEN}/" \
     -d '{"source_type": "Tag", "source_name": "'${version}'"}' \
     -H "Content-Type: application/json; charset=utf-8" \
     -sS
