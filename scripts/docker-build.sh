#!/bin/sh

env GOOS=linux GOARCH=amd64 go build -o devctl-proxy
docker build -t devctl-proxy .