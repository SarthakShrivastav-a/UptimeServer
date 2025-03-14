package handlers

import (
	"Uptime/models"
	"Uptime/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
)

func AddMonitorHandler(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var monitor models.Monitor
		if err := json.NewDecoder(r.Body).Decode(&monitor); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		if err := repository.AddMonitor(db, monitor); err != nil {
			fmt.Print(err)
			http.Error(w, "Failed to add monitor", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
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
