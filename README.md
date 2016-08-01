# site24x7 Exporter

site24x7 exporter for prometheus.io, written in go.

## Get it

Download your version from the [releases page](https://github.com/echocat/site24x7_exporter/releases/latest). For older version see [archive page](https://github.com/echocat/site24x7_exporter/releases).

Example:
```bash
sudo curl -SL https://github.com/echocat/site24x7_exporter/releases/download/v0.1.0/site24x7_exporter-linux-amd64 \
    > /usr/bin/site24x7_exporter
sudo chmod +x /usr/bin/site24x7_exporter
```

## Use it

### Usage

```
Usage: site24x7_exporter <flags>
Flags:
  -site24x7.timeout duration
        Timeout for trying to get stats from site24x7. (default 5s)
  -site24x7.token string
        Token to access the API of site24x7.
        See: https://www.site24x7.com/app/client#/admin/developer/api
  -web.listen-address string
        Address to listen on for web interface and telemetry. (default ":9112")
  -web.telemetry-path string
        Path under which to expose metrics. (default "/metrics")
  -web.tls-cert string
        Path to PEM file that conains the certificate (and optionally also the private key in PEM format).
        This should include the whole certificate chain.
        If provided: The web socket will be a HTTPS socket.
        If not provided: Only HTTP.
  -web.tls-client-ca string
        Path to PEM file that conains the CAs that are trused for client connections.
        If provided: Connecting clients should present a certificate signed by one of this CAs.
        If not provided: Every client will be accepted.
  -web.tls-private-key string
        Path to PEM file that contains the private key (if not contained in web.tls-cert file).
```

### Examples

```bash
# Simply start the exporter with your token and listen on 0.0.0.0:9112
site24x7_exporter -site24x7.token=mySecrectToken

# Start the exporter with your token and listen on 0.0.0.0:9112
# ...it also secures the connector via SSL 
site24x7_exporter -listen.address=:8443 \
    -web.tls-cert=my.server.com.pem

# Simply start the exporter with your token and listen on 0.0.0.0:9112
# ...secures the connector via SSL
# ...and requires client certificates signed by your authority
site24x7_exporter -listen.address=:8443 \
    -web.tls-cert=my.server.com.pem \
    -web.tls-client-ca=ca.pem
```

## Build it

### Precondition

For building site24x7_exporter there is only:

1. a compatible operating system (Linux, Windows or Mac OS X)
2. and a working [Java 8](http://www.oracle.com/technetwork/java/javase/downloads/index.html) installation required.

There is no need for a working and installed Go installation (or anything else). The build system will download every dependency and build it if necessary.

> **Hint:** The Go runtime build by the build system will be placed under ``~/.go/sdk``.

### Run build process

On Linux and Mac OS X:
```bash
# Build binaries (includes test)
./gradlew build

# Run tests (but do not build binaries)
./gradlew test

# Build binaries and release it on GitHub
# Environment variable GITHUB_TOKEN is required
./gradlew build githubRelease
```

On Windows:
```bash
# Build binaries (includes test)
gradlew build

# Run tests (but do not build binaries)
gradlew test

# Build binaries and release it on GitHub
# Environment variable GITHUB_TOKEN is required
gradlew build githubRelease
```

### Build artifacts

* Compiled and lined binaries can be found under ``./build/out/site24x7_exporter-*``

## Contributing

site24x7_exporter is an open source project of [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/site24x7_exporter)
by [fork us](https://github.com/echocat/site24x7_exporter/fork).

If you commit code to this project you have to accept that this code will be released under the [license](#license) of this project.


## License

See [LICENSE](LICENSE) file.
