package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/echocat/site24x7_exporter/utils"
	"github.com/prometheus/client_golang/prometheus"
	"log"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
	"io"
	"io/ioutil"
)

var (
	currentStatusApiUri, _ = url.Parse("https://www.site24x7.com/api/current_status")
)

type Site24x7Exporter struct {
	uri         string
	accessToken string

	mutex sync.RWMutex

	up     prometheus.Gauge
	status prometheus.Gauge

	client *http.Client
}

func NewSite24x7Exporter(accessToken string, timeout time.Duration) *Site24x7Exporter {

	return &Site24x7Exporter{
		accessToken: accessToken,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the Site24x7 instance query successful?",
		}),
		status: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "status",
			Help:      "What is the status of the target monitor?",
		}),
		client: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					RootCAs: utils.LoadInternalCaBundle(),
				},
				Dial: func(netw, addr string) (net.Conn, error) {
					c, err := net.DialTimeout(netw, addr, timeout)
					if err != nil {
						return nil, err
					}
					if err := c.SetDeadline(time.Now().Add(timeout)); err != nil {
						return nil, err
					}
					return c, nil
				},
			},
		},
	}
}

func (instance *Site24x7Exporter) createRequestFor(url *url.URL) *http.Request {
	return &http.Request{
		Method:     "GET",
		URL:        url,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Authorization": []string{fmt.Sprintf("Zoho-authtoken %s", instance.accessToken)},
		},
		Host: url.Host,
	}
}

func (instance *Site24x7Exporter) executeAndEvaluate(request *http.Request, target interface{}) error {
	response, err := instance.client.Do(request)
	defer func() {
		io.Copy(ioutil.Discard, response.Body)
		response.Body.Close()
	}()
	if err != nil {
		return fmt.Errorf("Could not execute request %v. Got: %v", request.URL, err)
	}
	if response.StatusCode < 200 || response.StatusCode >= 400 {
		return fmt.Errorf("Could not execute request %v. Got: %v - %v", request.URL, response.StatusCode, response.Status)
	}
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(target)
	if err != nil {
		return fmt.Errorf("Could not execute request %v. Could not decode response. Got: %v", request.URL, err)
	}
	return nil
}

func (instance *Site24x7Exporter) retrieveCurrentStatus() (*CurrentStatus, error) {
	request := instance.createRequestFor(currentStatusApiUri)
	restObject := CurrentStatus{}
	err := instance.executeAndEvaluate(request, &restObject)
	if err != nil {
		return nil, err
	}
	if restObject.ErrorCode != 0 {
		return nil, fmt.Errorf("Could not execute request %v. Could not decode response. Got: #%d - %s", request.URL, restObject.ErrorCode, restObject.Message)
	}
	return &restObject, nil
}

// Describe describes all the metrics ever exported by the
// exporter. It implements prometheus.Collector.
func (instance *Site24x7Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- instance.up.Desc()
	ch <- instance.status.Desc()
}

// Collect fetches the stats from configured site24x7 and
// delivers them as Prometheus metrics. It implements prometheus.Collector.
func (instance *Site24x7Exporter) Collect(ch chan<- prometheus.Metric) {
	instance.mutex.Lock() // To protect metrics from concurrent collects.
	defer instance.mutex.Unlock()

	currentStatus, err := instance.retrieveCurrentStatus()
	if err != nil {
		log.Printf("Failed to retreive current status. Cause: %v", err)
		return
	}

	NewStatusFor(currentStatus).Collect(ch)
	NewAttributesFor(currentStatus).Collect(ch)
}
