#!/usr/bin/env bash
set -e

cd ../go/generator
GOARCH=amd64 GOOS=linux go build -o ../../docker/generator
cd ../../docker

docker build --platform=linux/amd64 -t saichler/traffic-generator:latest .
docker push saichler/traffic-generator:latest