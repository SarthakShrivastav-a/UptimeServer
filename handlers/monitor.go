package handlers

import (
	"Uptime/models"
	"Uptime/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func AddMonitorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("AddMonitor function called")

		// Read the raw request body for debugging
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		fmt.Printf("Received Raw JSON: %s\n", string(body))

		var monitor models.Monitor
		err = json.Unmarshal(body, &monitor)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		fmt.Printf("Decoded Monitor struct: %+v\n", monitor)

		if err := repository.AddMonitor(db, monitor); err != nil {
			log.Printf("Failed to add monitor to DB: %v", err)
			http.Error(w, "Failed to add monitor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintln(w, "Created a new entry in the database")
		log.Println("Created a new entry in DB")
	}
}

func DeleteMonitorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		fmt.Println("DeleteMonitor function called")

		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		fmt.Printf("Received Raw JSON: %s\n", string(body))

		var request struct {
			ID string `json:"id"`
		}

		err = json.Unmarshal(body, &request)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		fmt.Printf("Decoded Delete Request: %+v\n", request)

		if request.ID == "" {
			log.Printf("Invalid monitor ID: %s", request.ID)
			http.Error(w, "Invalid monitor ID", http.StatusBadRequest)
			return
		}

		if err := repository.DeleteMonitor(db, request.ID); err != nil {
			log.Printf("Failed to delete monitor from DB: %v", err)
			http.Error(w, "Failed to delete monitor", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Successfully deleted monitor from database")
		log.Println("Deleted monitor from DB")
	}
}
func UpdateMonitorErrorConditionHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("UpdateMonitorErrorCondition function called")

		// Read the request body
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Error reading request body: %v", err)
			http.Error(w, "Failed to read request body", http.StatusBadRequest)
			return
		}
		log.Printf("Received Raw JSON: %s\n", string(body))

		// Define a struct to parse just the fields we need
		type ErrorConditionUpdate struct {
			MonitorID      string                `json:"monitor_id"`
			ErrorCondition models.ErrorCondition `json:"error_condition"`
		}

		var update ErrorConditionUpdate
		err = json.Unmarshal(body, &update)
		if err != nil {
			log.Printf("Error decoding JSON: %v", err)
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}

		if update.MonitorID == "" {
			log.Printf("Invalid monitor ID: %v", update.MonitorID)
			http.Error(w, "Invalid or missing monitor ID", http.StatusBadRequest)
			return
		}

		// Validate error condition
		if update.ErrorCondition.TriggerOn == "" {
			log.Printf("Invalid error condition: TriggerOn is empty")
			http.Error(w, "Invalid error condition: TriggerOn is required", http.StatusBadRequest)
			return
		}

		log.Printf("Processing error condition update for monitor ID: %s", update.MonitorID)

		if err := repository.UpdateMonitorErrorCondition(db, update.MonitorID, update.ErrorCondition); err != nil {
			log.Printf("Failed to update monitor error condition: %v", err)
			if strings.Contains(err.Error(), "no monitor found") {
				http.Error(w, err.Error(), http.StatusNotFound)
			} else {
				http.Error(w, "Failed to update monitor error condition", http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "Successfully updated monitor error condition")
		log.Println("Updated monitor error condition in DB")
	}
}

// func UpdateMonitorHandler(db *sql.DB) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Println("UpdateMonitor function called")

// 		// Read the raw request body for debugging
// 		body, err := io.ReadAll(r.Body)
// 		if err != nil {
// 			log.Printf("Error reading request body: %v", err)
// 			http.Error(w, "Failed to read request body", http.StatusBadRequest)
// 			return
// 		}
// 		fmt.Printf("Received Raw JSON: %s\n", string(body))

// 		var monitor models.Monitor
// 		err = json.Unmarshal(body, &monitor)
// 		if err != nil {
// 			log.Printf("Error decoding JSON: %v", err)
// 			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
// 			return
// 		}

// 		if monitor.MonitorID == "" {
// 			log.Printf("Invalid monitor ID: %v", monitor.MonitorID)
// 			http.Error(w, "Invalid or missing monitor ID", http.StatusBadRequest)
// 			return
// 		}

// 		fmt.Printf("Decoded Monitor struct for update: %+v\n", monitor)

// 		if err := repository.UpdateMonitor(db, monitor); err != nil {
// 			log.Printf("Failed to update monitor in DB: %v", err)
// 			if strings.Contains(err.Error(), "no monitor found") {
// 				http.Error(w, err.Error(), http.StatusNotFound)
// 			} else {
// 				http.Error(w, "Failed to update monitor", http.StatusInternalServerError)
// 			}
// 			return
// 		}

// 		w.WriteHeader(http.StatusOK)
// 		fmt.Fprintln(w, "Successfully updated monitor in the database")
// 		log.Println("Updated monitor in DB")
// 	}
// }

func GetAllMonitorsHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		monitors, err := repository.GetAllMonitors(db)
		if err != nil {
			http.Error(w, "Failed to fetch monitors", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(monitors)
	}
}
