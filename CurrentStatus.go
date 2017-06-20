package main

import (
	"encoding/json"
	"errors"
)

var (
	// String value returned by current_status API to indicate 'no value'.
	dashValue = `"-"`

	// Error value returned by AttributeValue() in case |dashValue| is observed.
	errUndefined = errors.New("no metric value")
)

type CurrentStatus struct {
	Code      int               `json:"code"`
	ErrorCode int               `json:"error_code"`
	Message   string            `json:"message"`
	Data      CurrentStatusData `json:"data"`
}

type CurrentStatusData struct {
	Monitors      []CurrentStatusMonitor      `json:"monitors"`
	MonitorGroups []CurrentStatusMonitorGroup `json:"monitor_groups"`
}

type CurrentStatusMonitorGroup struct {
	Id       string                 `json:"group_id"`
	Name     string                 `json:"group_name"`
	Status   int                    `json:"status"`
	Monitors []CurrentStatusMonitor `json:"monitors"`
}

type CurrentStatusMonitor struct {
	Id                string                  `json:"monitor_id"`
	Name              string                  `json:"name"`
	Type              string                  `json:"monitor_type"`
	Status            int                     `json:"status"`
	AttributeKey      string                  `json:"attribute_key"`
	RawAttributeValue json.RawMessage         `json:"attribute_value"`
	Locations         []CurrentStatusLocation `json:"locations"`
}

type CurrentStatusLocation struct {
	Name   string `json:"location_name"`
	Status int    `json:"status"`
}

// If the raw |AttributeValue| equals the empty "-", coerce to -1.
// Otherwise unmarshal as an integer.
func (csm CurrentStatusMonitor) AttributeValue() (int, error) {
	if string(csm.RawAttributeValue) == dashValue {
		return 0, errUndefined
	} else {
		var i int
		var err = json.Unmarshal(csm.RawAttributeValue, &i)
		return i, err
	}
}
