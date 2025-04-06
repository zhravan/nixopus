package realtime

import (
	"encoding/json"
	"log"
)

// handleNotifications processes incoming PostgreSQL notifications from the notification channel.
//
// This method listens on the provided channel for notifications related to application changes.
// Upon receiving a notification, it checks if the notification is from the "application_changes"
// channel. If it is, the method attempts to parse the JSON payload to extract the table, action,
// and ID details. If parsing is successful, it constructs a message containing the parsed data
// and broadcasts it to the appropriate topic using the BroadcastToTopic method.
//
// Parameters:
//
//	notificationChan - a channel that receives *PostgresNotification objects.
//
// Errors and any issues parsing the notification payload are logged, but do not stop the
// processing of further notifications.
func (s *SocketServer) handleNotifications(notificationChan <-chan *PostgresNotification) {
	for notification := range notificationChan {
		// fmt.Printf("Received notification on channel %s: %s\n",
		// 	notification.Channel, notification.Payload)

		if notification.Channel == "application_changes" {
			var parsedPayload struct {
				Table         string                 `json:"table"`
				Action        string                 `json:"action"`
				ApplicationID string                 `json:"application_id"`
				Data          map[string]interface{} `json:"data"`
			}

			if err := json.Unmarshal([]byte(notification.Payload), &parsedPayload); err != nil {
				log.Printf("Error parsing notification payload: %v", err)
				continue
			}

			resourceID := parsedPayload.ApplicationID

			messageData := map[string]interface{}{
				"table":          parsedPayload.Table,
				"action":         parsedPayload.Action,
				"application_id": parsedPayload.ApplicationID,
				"data":           parsedPayload.Data,
			}

			s.BroadcastToTopic(MonitorApplicationDeployment, resourceID, messageData)
		}
	}
}