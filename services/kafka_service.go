package services

import (
	"Uptime/models"
	"Uptime/repository"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"

	"github.com/segmentio/kafka-go"
)

type monitorLifecycleEvent struct {
	EventType      string                `json:"eventType"`
	MonitorID      string                `json:"monitorId"`
	UserID         string                `json:"userId"`
	URL            string                `json:"url"`
	ErrorCondition models.ErrorCondition `json:"errorCondition"`
}

func kafkaBrokers() []string {
	raw := os.Getenv("KAFKA_BOOTSTRAP_SERVERS")
	if raw == "" {
		return nil
	}
	return strings.Split(raw, ",")
}

func StartMonitorLifecycleConsumer(db *sql.DB) {
	brokers := kafkaBrokers()
	if len(brokers) == 0 {
		log.Println("Kafka disabled: KAFKA_BOOTSTRAP_SERVERS not set")
		return
	}

	topic := os.Getenv("MONITOR_LIFECYCLE_TOPIC")
	if topic == "" {
		topic = "sentinel.monitor.lifecycle"
	}

	groupID := os.Getenv("KAFKA_CONSUMER_GROUP")
	if groupID == "" {
		groupID = "sentinel-go-probe"
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})

	go func() {
		for {
			message, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Printf("Kafka lifecycle read failed: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			var event monitorLifecycleEvent
			if err := json.Unmarshal(message.Value, &event); err != nil {
				log.Printf("Invalid monitor lifecycle event: %v", err)
				continue
			}

			monitor := models.Monitor{
				MonitorID:      event.MonitorID,
				URL:            event.URL,
				ErrorCondition: event.ErrorCondition,
			}

			switch event.EventType {
			case "monitor.created", "monitor.updated":
				if err := repository.UpsertMonitor(db, monitor); err != nil {
					log.Printf("Failed to upsert monitor from Kafka: %v", err)
				}
			case "monitor.deleted":
				if err := repository.DeleteMonitor(db, event.MonitorID); err != nil {
					log.Printf("Failed to delete monitor from Kafka: %v", err)
				}
			default:
				log.Printf("Ignoring unknown monitor lifecycle event type: %s", event.EventType)
			}
		}
	}()
}

func publishMonitorCheckToKafka(monitor models.Monitor, status string, responseTime time.Duration, statusCode int) bool {
	brokers := kafkaBrokers()
	if len(brokers) == 0 {
		return false
	}

	topic := os.Getenv("MONITOR_CHECKS_TOPIC")
	if topic == "" {
		topic = "sentinel.monitor.checks"
	}

	probeID := os.Getenv("PROBE_ID")
	if probeID == "" {
		probeID = "local-probe"
	}

	payload, err := json.Marshal(map[string]interface{}{
		"eventId":       time.Now().Format(time.RFC3339Nano) + "-" + monitor.MonitorID,
		"eventType":     "monitor.check.completed",
		"monitorId":     monitor.MonitorID,
		"probeId":       probeID,
		"status":        status,
		"triggerReason": monitor.ErrorCondition.TriggerOn,
		"checkedAt":     time.Now().Format(time.RFC3339Nano),
		"responseTime":  responseTime.Milliseconds(),
		"statusCode":    statusCode,
	})
	if err != nil {
		log.Printf("Failed to serialize monitor check event: %v", err)
		return false
	}

	writer := kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.Hash{},
	}
	defer writer.Close()

	err = writer.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(monitor.MonitorID),
		Value: payload,
	})
	if err != nil {
		log.Printf("Failed to publish monitor check event: %v", err)
		return false
	}

	return true
}
