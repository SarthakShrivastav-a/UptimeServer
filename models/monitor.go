package models

type Monitor struct {
	MonitorID      string `json:"monitor_id"`
	URL            string `json:"url"`
	ErrorCondition string `json:"error_condition"`
}
