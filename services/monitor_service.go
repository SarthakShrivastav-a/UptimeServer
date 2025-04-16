package services

import (
	"Uptime/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	workerCount   = 5
	checkInterval = 30 * time.Second
)

func StartMonitoring(db *sql.DB) { //init the workerr pools
	monitorChan := make(chan models.Monitor, 10)

	for i := 0; i < workerCount; i++ {
		go monitorWorker(monitorChan)
	}

	go func() { // fetch monitors every 30 seconds
		for {
			monitors, err := fetchMonitors(db)
			if err != nil {
				fmt.Println("Error fetching monitors:", err)
				time.Sleep(10 * time.Second)
				continue
			}
			//add monitors to the channel
			for _, monitor := range monitors {
				monitorChan <- monitor
			}

			time.Sleep(checkInterval)
		}
	}()
}

func monitorWorker(monitorChan <-chan models.Monitor) { //processe thee urlls and checks if they match error conditions
	for monitor := range monitorChan {
		timeout := 10 // default time of 10sec
		if monitor.ErrorCondition.TriggerOn == "TIMEOUT" && len(monitor.ErrorCondition.Value) > 0 {
			timeout = monitor.ErrorCondition.Value[0]
		}

		statusCode, responseTime, err := checkWebsiteStatus(monitor.URL, timeout)
		if shouldTriggerAlert(monitor.ErrorCondition, statusCode, responseTime, err) {
			fmt.Printf("Alert triggered for %s (Status: %d, Time: %v)\n", monitor.URL, statusCode, responseTime)
			post(monitor, "DOWN", responseTime)
		} else {
			post(monitor, "UP", responseTime)
			fmt.Printf("%v Working Fine\n", monitor.URL)
		}
	}
}

func checkWebsiteStatus(url string, timeoutSeconds int) (int, time.Duration, error) { //sends a req returns the status code response time and error agar mile
	client := http.Client{
		Timeout: time.Duration(timeoutSeconds) * time.Second,
	}
	start := time.Now()
	resp, err := client.Get(url)
	duration := time.Since(start)

	if err != nil {
		return 0, duration, err // timeout or request failure
	}
	defer resp.Body.Close()

	return resp.StatusCode, duration, nil
}

// determines if the error condition is met
func shouldTriggerAlert(cond models.ErrorCondition, status int, responseTime time.Duration, err error) bool {
	switch cond.TriggerOn {
	case "STATUS_NOT":
		// ttrigger alert if sttus is not in the allowed list
		for _, v := range cond.Value {
			if status == v {
				return false //allowed status no alert
			}
		}
		return true // sttaus not allowed trigger a nuke lmao

	case "RESPONSE_CONTAINS":
		// trigger if matches in the list
		for _, v := range cond.Value {
			if status == v {
				return true // found a matching error status trigger alert!!!!!!PANICCCCCCCCC OMFG
			}
		}
		return false // its fineee no need for alerrts

	case "TIMEOUT":
		// request timed out, trigger alert
		if err != nil {
			return true
		}
		// checkfor if response time exceeds the allowed limit
		if len(cond.Value) > 0 && responseTime.Seconds() > float64(cond.Value[0]) {
			return true
		}
		return false
	}

	return false //the ddeafualt is that no alert
}

func fetchMonitors(db *sql.DB) ([]models.Monitor, error) {
	rows, err := db.Query("SELECT monitor_id, url, error_condition FROM monitors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []models.Monitor
	for rows.Next() {
		var m models.Monitor
		var errorConditionJSON string

		// Scan monitor ID, URL, and error_condition (stored as JSON)
		if err := rows.Scan(&m.MonitorID, &m.URL, &errorConditionJSON); err != nil {
			return nil, err
		}

		// Parse JSON into ErrorCondition struct
		if err := json.Unmarshal([]byte(errorConditionJSON), &m.ErrorCondition); err != nil {
			fmt.Println("Error parsing error_condition JSON:", err)
			continue // Skip this monitor if parsing fails
		}

		monitors = append(monitors, m)
	}
	return monitors, nil
}

// func fetchMonitors(db *sql.DB) ([]models.Monitor, error) {
// 	// Mocking a monitor for testing with the Node.js server
// 	monitors := []models.Monitor{
// 		{
// 			MonitorID: "test123",
// 			URL:       "http://localhost:3000/test",
// 			ErrorCondition: models.ErrorCondition{
// 				TriggerOn: "STATUS_NOT",
// 				Value:     []int{200}, // Should trigger alert if != 200
// 			},
// 		},
// 	}
// 	return monitors, nil
// }

// // parseStatusValues converts a string "200,201" into []int{200, 201}
// func parseStatusValues(valueStr string) []int {
// 	var values []int
// 	fmt.Sscanf(valueStr, "%d,%d", &values)
// 	return values
// }
