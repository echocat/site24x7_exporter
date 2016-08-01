package main

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type Status struct {
	CurrentStatus  *CurrentStatus
}

func NewStatusFor(currentStatus *CurrentStatus) *Status {
	return &Status{
		CurrentStatus:  currentStatus,
	}
}

func (instance *Status) Describe(ch chan<- *prometheus.Desc) {
	ch <- instance.Desc()
}

func (instance *Status) Desc() *prometheus.Desc {
	return prometheus.NewDesc(
		fmt.Sprintf("%s.monitor.status", namespace),
		"Was is the status of the target monitor?",
		[]string{},
		prometheus.Labels{},
	)
}

func (instance *Status) Collect(ch chan<- prometheus.Metric) {
	for _, monitor := range (*instance).CurrentStatus.Data.Monitors {
		element := &StatusElement{
			Parent: instance,
			Monitor: monitor,
		}
		ch <- element
	}
	for _, monitorGroup := range (*instance).CurrentStatus.Data.MonitorGroups {
		for _, monitor := range monitorGroup.Monitors {
			element := &StatusElement{
				Parent: instance,
				MonitorGroup: monitorGroup,
				Monitor: monitor,
			}
			ch <- element
		}
	}
}

type StatusElement struct {
	Parent       *Status
	MonitorGroup CurrentStatusMonitorGroup
	Monitor      CurrentStatusMonitor
}

func (instance *StatusElement) Write(out *dto.Metric) error {
	out.Counter = &dto.Counter{Value: proto.Float64(float64(instance.Monitor.Status))}
	label := []*dto.LabelPair{
		labelPairFor("monitorId", instance.Monitor.Id),
		labelPairFor("monitorDisplayName", instance.Monitor.Name),
	}
	if instance.MonitorGroup.Id != "" {
		label = append(label,
			labelPairFor("monitorGroupId", instance.MonitorGroup.Id),
			labelPairFor("monitorGroupDisplayName", instance.MonitorGroup.Name),
		)
	}
	out.Label = label
	return nil
}

func (instance *StatusElement) Desc() *prometheus.Desc {
	return instance.Parent.Desc()
}

func labelPairFor(name string, value string) *dto.LabelPair {
	return &dto.LabelPair{
		Name:  &name,
		Value: &value,
	}
}
