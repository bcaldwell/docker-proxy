#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -o docker-proxy
docker build -t benjamincaldwell/docker-proxy .