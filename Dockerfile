FROM nginx:1.10-alpine

ADD devctl-proxy /devctl-proxy/devctl-proxy

ADD scripts/docker-init.sh /devctl-proxy/docker-init.sh

RUN chmod +x /devctl-proxy/*

ADD nginx-default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443

CMD devctl-proxy/docker-init.sh
