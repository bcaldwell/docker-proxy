FROM nginx:1.10-alpine

COPY docker-proxy /docker-proxy/docker-proxy

COPY scripts/docker-init.sh /docker-proxy/docker-init.sh

RUN chmod +x /docker-proxy/*

COPY nginx-default.conf /etc/nginx/conf.d/default.conf

EXPOSE 80 443

CMD docker-proxy/docker-init.sh
