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
