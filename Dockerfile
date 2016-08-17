FROM        quay.io/prometheus/busybox:latest
MAINTAINER  Meiqia Developers<dev@meiqia.com>

COPY elasticsearch_exporter /bin/elasticsearch_exporter

EXPOSE      9108
ENTRYPOINT  [ "/bin/elasticsearch_exporter" ]
