FROM ververica/flink-statefun:2.2.0

RUN mkdir -p /opt/statefun/modules/functions
ADD module.yaml /opt/statefun/modules/functions