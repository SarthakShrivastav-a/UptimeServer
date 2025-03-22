package services

import (
	"Uptime/config"
	"Uptime/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func post(monitor models.Monitor, status string, responseTime time.Duration) {
	alertURL := config.GetSpringBootURL()

	alertData, _ := json.Marshal(map[string]interface{}{
		"monitorId":     monitor.MonitorID,
		"status":        status,
		"triggerReason": monitor.ErrorCondition.TriggerOn,
		"checkedAt":     time.Now().Format(time.RFC3339),
		"responseTime":  responseTime.Milliseconds(),
	})
	resp, err := http.Post(alertURL, "application/json", bytes.NewBuffer(alertData))
	if err != nil {
		fmt.Println("Failed to send status update:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Status update sent to Spring Boot for %s: %s, response status: %s\n",
		monitor.URL, status, resp.Status)
}
