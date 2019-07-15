# site24x7 Exporter

site24x7 exporter for prometheus.io, written in go.

## Get it

### Binary distribution

Download your version from the [releases page](https://github.com/echocat/site24x7_exporter/releases/latest). For older version see [archive page](https://github.com/echocat/site24x7_exporter/releases).

Example:
```bash
sudo curl -SL https://github.com/echocat/site24x7_exporter/releases/download/v0.1.5/site24x7_exporter-linux-amd64 \
    > /usr/bin/site24x7_exporter
sudo chmod +x /usr/bin/site24x7_exporter
```

### Docker image

Image: ``docker pull echocat/site24x7_exporter``

You can go to [Docker Hub Tags page](https://hub.docker.com/r/echocat/site24x7_exporter/tags/) to see all available tags or you can simply use ``latest``.

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

#### Binary distribution

```bash
# Simply start the exporter with your token and listen on 0.0.0.0:9112
site24x7_exporter \
    -site24x7.token=mySecrectToken

# Start the exporter with your token and listen on 0.0.0.0:9112
# ...it also secures the connector via SSL 
site24x7_exporter \
    -listen.address=:8443 \
    -web.tls-cert=my.server.com.pem

# Simply start the exporter with your token and listen on 0.0.0.0:9112
# ...secures the connector via SSL
# ...and requires client certificates signed by your authority
site24x7_exporter \
    -listen.address=:8443 \
    -web.tls-cert=my.server.com.pem \
    -web.tls-client-ca=ca.pem
```

#### Docker image

```bash
# Simply start the exporter with your token and listen on 0.0.0.0:9112
docker run -p9112:9112 echocat/site24x7_exporter \
    -site24x7.token=mySecrectToken

# Start the exporter with your token and listen on 0.0.0.0:9112
# ...it also secures the connector via SSL 
docker run -p9112:9112 -v/etc/certs:/etc/certs:ro echocat/site24x7_exporter \
    -listen.address=:8443 \
    -web.tls-cert=/etc/certs/my.server.com.pem

# Simply start the exporter with your token and listen on 0.0.0.0:9112
# ...secures the connector via SSL
# ...and requires client certificates signed by your authority
docker run -p9112:9112 -v/etc/certs:/etc/certs:ro echocat/site24x7_exporter \
    -listen.address=:8443 \
    -web.tls-cert=my.server.com.pem \
    -web.tls-client-ca=ca.pem
```

## Metrics

### ``site24x7_monitor_status``

**Type**: Counter

#### Labels

| Name | Example | Description |
| ---- | ------- | ----------- |
| ``monitorId`` | ``123456789012345678`` | Internal ID assigned by site24x7 of monitor |
| ``monitorDisplayName`` | ``My service`` | Display name assigned by you of monitor |
| ``monitorGroupId`` | ``123456789012345678`` | Internal ID assigned by site24x7 of monitor group (optional) |
| ``monitorGroupDisplayName`` | ``My data center`` | Display name assigned by you of monitor group  (optional) |

#### Possible values

| Value | Description |
| ----- | ----------- |
| ``0`` | Down |
| ``1`` | Up |
| ``2`` | Trouble |
| ``5`` | Suspended |
| ``7`` | Maintenance |
| ``9`` | Discovery |
| ``10`` | Discovery Error |

## Contributing

site24x7_exporter is an open source project of [echocat](https://echocat.org).
So if you want to make this project even better, you can contribute to this project on [Github](https://github.com/echocat/site24x7_exporter)
by [fork us](https://github.com/echocat/site24x7_exporter/fork).

If you commit code to this project you have to accept that this code will be released under the [license](#license) of this project.


## License

See [LICENSE](LICENSE) file.
