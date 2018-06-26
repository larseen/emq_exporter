FROM quay.io/prometheus/busybox:latest
MAINTAINER Kristoffer Larsen <kristoffer@larsen.so>

COPY emq_exporter /bin/emq_exporter

EXPOSE      9444
ENTRYPOINT  [ "/bin/emq_exporter" ]
