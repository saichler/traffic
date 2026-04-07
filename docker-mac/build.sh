#!/usr/bin/env bash
set -e

cd ../go/generator
GOARCH=arm64 GOOS=linux go build -o ../../docker-mac/generator
cd ../../docker-mac

docker build --platform=linux/arm64 -t saichler/traffic-generator:latest-arm64 .
docker push saichler/traffic-generator:latest-arm64
