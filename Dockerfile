FROM scratch
MAINTAINER contact@echocat.org

COPY build/out/site24x7_exporter-linux-amd64 /usr/bin/site24x7_exporter

ENTRYPOINT ["/usr/bin/site24x7_exporter"]
