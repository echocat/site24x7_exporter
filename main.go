package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"os"
	"strings"
	"time"
)

const (
	namespace = "site24x7"
)

var (
	name     = "site24x7_exporter"
	version  = "development"
	revision = "development"
	compiled = time.Now().Format("2006-01-02T15:04:05Z")

	listenAddress = flag.String("web.listen-address", ":9112", "Address to listen on for web interface and telemetry.")
	metricsPath   = flag.String("web.telemetry-path", "/metrics", "Path under which to expose metrics.")
	tlsCert       = flag.String("web.tls-cert", "", "Path to PEM file that conains the certificate (and optionally also the private key in PEM format).\n"+
		"\tThis should include the whole certificate chain.\n"+
		"\tIf provided: The web socket will be a HTTPS socket.\n"+
		"\tIf not provided: Only HTTP.")
	tlsPrivateKey = flag.String("web.tls-private-key", "", "Path to PEM file that contains the private key (if not contained in web.tls-cert file).")
	tlsClientCa   = flag.String("web.tls-client-ca", "", "Path to PEM file that conains the CAs that are trused for client connections.\n"+
		"\tIf provided: Connecting clients should present a certificate signed by one of this CAs.\n"+
		"\tIf not provided: Every client will be accepted.")
	site24x7Token = flag.String("site24x7.token", "", "Token to access the API of site24x7.\n"+
		"\tSee: https://www.site24x7.com/app/client#/admin/developer/api")
	site24x7Timeout = flag.Duration("site24x7.timeout", 5*time.Second, "Timeout for trying to get stats from site24x7.")

	flagsBuffer = &bytes.Buffer{}
)

func main() {
	parseUsage()

	exporter := NewSite24x7Exporter(*site24x7Token, *site24x7Timeout)
	prometheus.MustRegister(exporter)

	err := startServer(*metricsPath, *listenAddress, *tlsCert, *tlsPrivateKey, *tlsClientCa)
	if err != nil {
		log.Fatalf("Could not start server. Cause: %v", err)
	}
}

func parseUsage() {
	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	flags.SetOutput(flagsBuffer)
	flags.ErrorHandling()
	flags.Usage = func() {
		errorString := flagsBuffer.String()
		if len(errorString) > 0 {
			printUsage(strings.TrimSpace(errorString))
		} else {
			printUsage(nil)
		}
	}
	if err := flags.Parse(os.Args[1:]); err == flag.ErrHelp {
		os.Exit(0)
	} else if err != nil {
		os.Exit(1)
	}
	assertUsage()
}

func assertUsage() {
	if len(strings.TrimSpace(*listenAddress)) == 0 {
		fail("Missing -web.listen-address")
	}
	if len(strings.TrimSpace(*site24x7Token)) == 0 {
		fail("Missing -site24x7.token")
	}
}

func fail(err interface{}) {
	printUsage(err)
	os.Exit(1)
}

func printUsage(err interface{}) {
	fmt.Fprintf(os.Stderr, "%v (version: %v, revision: %v, build: %v)\n", name, version, revision, compiled)
	fmt.Fprint(os.Stderr, "URL: https://github.com/echocat/site24x7_exporter\n")
	fmt.Fprint(os.Stderr, "Author(s): Gregor Noczinski (gregor@noczinski.eu)\n")
	fmt.Fprint(os.Stderr, "\n")

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n\n", err)
	}

	fmt.Fprintf(os.Stderr, "Usage: %v <flags>\n", os.Args[0])
	fmt.Fprint(os.Stderr, "Flags:\n")
	flag.CommandLine.SetOutput(os.Stderr)
	flag.CommandLine.PrintDefaults()
	flag.CommandLine.SetOutput(flagsBuffer)
}
