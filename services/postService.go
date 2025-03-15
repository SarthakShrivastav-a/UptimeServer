package services

import (
	"Uptime/models"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func post(monitor models.Monitor) {
	alertURL := "http://your-spring-boot-server.com/api/alert" // Replace with actual Spring Boot endpoint

	// Construct JSON payload
	alertData, _ := json.Marshal(map[string]interface{}{
		"monitor_id":     monitor.MonitorID,
		"url":            monitor.URL,
		"trigger_reason": monitor.ErrorCondition.TriggerOn,
		"timestamp":      time.Now().Format(time.RFC3339),
	})
	resp, err := http.Post(alertURL, "application/json", bytes.NewBuffer(alertData))
	if err != nil {
		fmt.Println("Failed to send alert:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("✅ Alert sent to Spring Boot, response status:", resp.Status)
}
