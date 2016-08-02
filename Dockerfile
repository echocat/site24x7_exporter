FROM java:8-jdk-alpine
MAINTAINER coontact@echocat.org

# Prepare build environment
COPY . /tmp/buildroot/site24x7_exporter
RUN apk update && \
	apk add bash go

# Build exporter
RUN cd /tmp/buildroot/site24x7_exporter && \
	export GOROOT=/usr/lib/go && \
	chmod +x /tmp/buildroot/site24x7_exporter/gradlew && \
	./gradlew --stacktrace --info -Dplatforms=linux-amd64 build && \
	cp /tmp/buildroot/site24x7_exporter/build/out/site24x7_exporter-linux-amd64 /usr/bin/site24x7_exporter && \
	chmod +x /usr/bin/site24x7_exporter

# Clean up build environment
RUN apk del bash go readline ncurses-libs ncurses-terminfo ncurses-terminfo-base && \
	rm -rf /var/cache/* && \
	rm -rf /root/.go && \
	rm -rf /root/.gradle && \
	rm -rf /media && \
	rm -rf /tmp/*

ENTRYPOINT ["/usr/bin/site24x7_exporter"]
