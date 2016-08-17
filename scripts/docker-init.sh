#!/bin/sh

/docker-proxy/docker-proxy &

nginx -g "daemon off;"