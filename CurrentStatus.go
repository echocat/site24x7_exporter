package main

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
	Id             string                  `json:"monitor_id"`
	Name           string                  `json:"name"`
	Type           string                  `json:"monitor_type"`
	Status         int                     `json:"status"`
	AttributeKey   string                  `json:"attribute_key"`
	AttributeValue int                     `json:"attribute_value"`
	Locations      []CurrentStatusLocation `json:"locations"`
}

type CurrentStatusLocation struct {
	Name   string `json:"location_name"`
	Status int    `json:"status"`
}
