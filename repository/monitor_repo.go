package repository

import (
	"Uptime/models"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
)

func AddMonitor(db *sql.DB, monitor models.Monitor) error {

	errorConditionJSON, err := json.Marshal(monitor.ErrorCondition) // struct to json string
	if err != nil {
		log.Println("Error serializing ErrorCondition:", err)
		return err
	}

	_, err = db.Exec("INSERT INTO monitors (monitor_id, url, error_condition) VALUES (?, ?, ?)",
		monitor.MonitorID, monitor.URL, string(errorConditionJSON))
	return err
}
func DeleteMonitor(db *sql.DB, monitorID int) error {
	result, err := db.Exec("DELETE FROM monitors WHERE monitor_id = ?", monitorID)
	if err != nil {
		log.Println("Error deleting monitor:", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		log.Println("Error getting rows affected:", err)
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no monitor found with ID %d", monitorID)
	}

	log.Printf("Successfully deleted monitor with ID %d", monitorID)
	return nil
}

func GetAllMonitors(db *sql.DB) ([]models.Monitor, error) {
	rows, err := db.Query("SELECT monitor_id, url, error_condition FROM monitors")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var monitors []models.Monitor

	for rows.Next() {
		var monitor models.Monitor
		var errorConditionJSON string

		if err := rows.Scan(&monitor.MonitorID, &monitor.URL, &errorConditionJSON); err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(errorConditionJSON), &monitor.ErrorCondition); err != nil { //jsson string back to struct
			log.Println("Error deserializing ErrorCondition:", err)
			return nil, err
		}

		monitors = append(monitors, monitor)
	}

	return monitors, nil
}
