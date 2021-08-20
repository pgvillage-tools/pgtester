#!/bin/bash
set -x
set -e
docker-compose up -d postgres || exit 1
docker-compose build --no-cache --force-rm
docker-compose run pgtester /pgtester -f tests1.yaml || exit 2
docker-compose run pgtester /pgtester -f tests2.yaml || [ $? -ne 5 ] && exit 3
echo "All is as expected"
