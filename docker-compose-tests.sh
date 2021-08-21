#!/bin/bash
set -x
set -e
docker-compose up -d postgres || exit 1
docker-compose build --no-cache --force-rm
docker-compose run pgtester pgtester /etc/pgtestdata/examples || exit 2
echo "All is as expected"
