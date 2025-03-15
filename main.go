package main

import (
	"Uptime/config"
	"Uptime/handlers"
	"Uptime/services"
	"log"
	"net/http"
)

func main() {
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database")
	}
	defer db.Close()

	http.HandleFunc("/add_monitor", handlers.AddMonitorHandler(db))
	http.HandleFunc("/get_monitors", handlers.GetAllMonitorsHandler(db))

	services.StartMonitoring(db)

	log.Println("Server running on port 8081")
	http.ListenAndServe(":8081", nil)
}
