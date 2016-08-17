FROM nginx:1.10-alpine

ADD docker-proxy /docker-proxy/docker-proxy

ADD scripts/docker-init.sh /docker-proxy/docker-init.sh

RUN chmod +x /docker-proxy/*

ADD nginx-default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443

CMD docker-proxy/docker-init.sh
