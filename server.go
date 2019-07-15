package main

import (
	"crypto/tls"
	"github.com/echocat/site24x7_exporter/utils"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

type bufferedLogWriter struct {
	buf []byte
}

func (w *bufferedLogWriter) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func createHttpServerLogWrapper() *log.Logger {
	return log.New(&bufferedLogWriter{}, "", 0)
}

func startServer(metricsPath, listenAddress, tlsCert, tlsPrivateKey, tlsClientCa string) error {
	server := &http.Server{
		Addr:     listenAddress,
		ErrorLog: createHttpServerLogWrapper(),
	}
	http.Handle(metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Site24x7 Exporter</title></head>
             <body>
             <h1>Site24x7 Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	if len(tlsCert) > 0 {
		clientValidation := "no"
		if len(tlsClientCa) > 0 && len(tlsCert) > 0 {
			certificates, err := utils.LoadCertificatesFrom(tlsClientCa)
			if err != nil {
				log.Fatalf("Couldn't load client CAs from %s. Got: %s", tlsClientCa, err)
			}
			server.TLSConfig = &tls.Config{
				ClientCAs:  certificates,
				ClientAuth: tls.RequireAndVerifyClientCert,
			}
			clientValidation = "yes"
		}
		targetTlsPrivateKey := tlsPrivateKey
		if len(targetTlsPrivateKey) <= 0 {
			targetTlsPrivateKey = tlsCert
		}
		log.Printf("Listening on %s (scheme=HTTPS, secured=TLS, clientValidation=%s)", listenAddress, clientValidation)
		return server.ListenAndServeTLS(tlsCert, targetTlsPrivateKey)
	} else {
		log.Printf("Listening on %s (scheme=HTTP, secured=no, clientValidation=no)", server.Addr)
		return server.ListenAndServe()
	}
}
