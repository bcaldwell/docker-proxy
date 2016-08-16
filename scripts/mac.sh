#!/bin/sh

echo "\nadding loopback alias"
sudo ifconfig lo0 10.0.0.100 alias

echo "\nadding routing rule"
echo "rdr pass on lo0 inet proto tcp from any to 10.0.0.100 port 0:30000 -> 127.0.0.1 port 8080" | sudo pfctl -ef -

echo "\nbuilding docker image"
sh scripts/docker-build.sh

echo "\n starting docker container"
echo "docker run -v /var/run/docker.sock:/var/run/docker.sock -it -p 8080:80 --rm devctl-proxy"
docker run -v /var/run/docker.sock:/var/run/docker.sock -it -p 8080:80 --rm devctl-proxy