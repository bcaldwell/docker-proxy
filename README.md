# docker-proxy
proxy http requests to docker containers by dns

docker-proxy redirects traffic to a specfic docker container as spesified by a DNS record. The request is proxied by NGINX running in the docker-proxy container.

## Getting started
### Mac
`git clone git@github.com:benjamincaldwell/docker-proxy.git`
`cd docker-proxy`

`sh scripts/mac.sh`

To add a container to be tracked by docker-proxy, /etc/hosts needs to be modified. For conveince a script is provided:
`source scripts/docker-proxy.sh`

`
`docker-proxy docker-args -l proxy-hostname=example.com image`

Example:
`docker-proxy -l proxy-hostname=example.com nginx`


## Limitations
Current implimentation only supports http call due to the nature of how DNS hows. :(
