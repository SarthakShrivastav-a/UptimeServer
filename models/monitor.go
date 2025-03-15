package models

type ErrorCondition struct {
	TriggerOn string `json:"triggerOn"` // enum : STATUS_NOT, RESPONSE_CONTAINS, TIMEOUT will update more baad mein
	Value     []int  `json:"value"`     // list of hhtttp status codes
}

type Monitor struct {
	MonitorID      string         `json:"monitor_id"`
	URL            string         `json:"url"`
	ErrorCondition ErrorCondition `json:"error_condition"`
}
