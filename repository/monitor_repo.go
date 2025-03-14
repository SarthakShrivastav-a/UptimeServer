package repository

import (
	"Uptime/models"
	"database/sql"
)

func AddMonitor(db *sql.DB, monitor models.Monitor) error {
	_, err := db.Exec("INSERT INTO monitors (monitor_id, url, error_condition) VALUES (?, ?, ?)",
		monitor.MonitorID, monitor.URL, monitor.ErrorCondition)
	return err
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
		err = rows.Scan(&monitor.MonitorID, &monitor.URL, &monitor.ErrorCondition)
		if err != nil {
			return nil, err
		}
		monitors = append(monitors, monitor)
	}
	return monitors, nil
}
