package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type Attributes struct {
	CurrentStatus *CurrentStatus
}

func NewAttributesFor(currentStatus *CurrentStatus) *Attributes {
	return &Attributes{
		CurrentStatus: currentStatus,
	}
}

func (instance *Attributes) Describe(ch chan<- *prometheus.Desc) {
	ch <- instance.Desc()
}

func (instance *Attributes) Desc() *prometheus.Desc {
	return prometheus.NewDesc(
		fmt.Sprintf("%s_monitor_attribute", namespace),
		"Attributes of the target monitor",
		[]string{},
		prometheus.Labels{},
	)
}

func (instance *Attributes) Collect(ch chan<- prometheus.Metric) {
	for _, monitor := range (*instance).CurrentStatus.Data.Monitors {
		element := &AttributesElement{
			Parent:  instance,
			Monitor: monitor,
		}
		ch <- element
	}
	for _, monitorGroup := range (*instance).CurrentStatus.Data.MonitorGroups {
		for _, monitor := range monitorGroup.Monitors {
			element := &AttributesElement{
				Parent:       instance,
				MonitorGroup: monitorGroup,
				Monitor:      monitor,
			}
			ch <- element
		}
	}
}

type AttributesElement struct {
	Parent       *Attributes
	MonitorGroup CurrentStatusMonitorGroup
	Monitor      CurrentStatusMonitor
}

func (instance *AttributesElement) Write(out *dto.Metric) error {
	out.Gauge = &dto.Gauge{Value: proto.Float64(float64(instance.Monitor.AttributeValue))}
	label := []*dto.LabelPair{
		labelPairFor("attributeKey", instance.Monitor.AttributeKey),
		labelPairFor("monitorId", instance.Monitor.Id),
		labelPairFor("monitorDisplayName", instance.Monitor.Name),
		labelPairFor("monitorGroupId", instance.MonitorGroup.Id),
		labelPairFor("monitorGroupDisplayName", instance.MonitorGroup.Name),
	}
	out.Label = label
	return nil
}

func (instance *AttributesElement) Desc() *prometheus.Desc {
	return instance.Parent.Desc()
}
