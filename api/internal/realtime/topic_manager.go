package realtime

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
	"github.com/raghavyuva/nixopus-api/internal/types"
)

// SubscribeToTopic adds a connection to the specified topic.
//
// The function takes a topic string and a connection as parameters and
// stores the connection in the topic map. It is safe to call this function
// concurrently from multiple goroutines.
func (s *SocketServer) SubscribeToTopic(topic topics, conn *websocket.Conn) {
	s.topicsMu.Lock()
	defer s.topicsMu.Unlock()

	if _, exists := s.topics[topic]; !exists {
		s.topics[topic] = make(map[*websocket.Conn]bool)
	}
	s.topics[topic][conn] = true

	// Confirm subscription to the client
	conn.WriteJSON(types.Payload{
		Action: "subscribed",
		Topic:  string(topic),
		Data:   nil,
	})

	log.Printf("Client %s subscribed to topic %s", conn.RemoteAddr(), topic)
}

// SubscribeToResourceTopic subscribes a connection to a resource-specific topic.
//
// The function takes a topic string, a resource ID string, and a connection as parameters and
// subscribes the connection to the resource-specific topic. It is safe to call this function
// concurrently from multiple goroutines.
func (s *SocketServer) SubscribeToResourceTopic(topic topics, resourceID string, conn *websocket.Conn) {
	s.topicsMu.Lock()
	defer s.topicsMu.Unlock()

	resourceTopicKey := fmt.Sprintf("%s:%s", string(topic), resourceID)

	if _, exists := s.resourceTopics[resourceTopicKey]; !exists {
		s.resourceTopics[resourceTopicKey] = make(map[*websocket.Conn]bool)
	}
	s.resourceTopics[resourceTopicKey][conn] = true

	conn.WriteJSON(types.Payload{
		Action: "subscribed",
		Topic:  string(topic),
		Data: map[string]interface{}{
			"resourceId": resourceID,
		},
	})

	log.Printf("Client %s subscribed to resource topic %s with ID %s",
		conn.RemoteAddr(), topic, resourceID)
}

// UnsubscribeFromTopic removes a connection from the specified topic.
//
// The function takes a topic string and a connection as parameters and
// removes the connection from the topic map. It is safe to call this function
// concurrently from multiple goroutines.
func (s *SocketServer) UnsubscribeFromTopic(topic topics, conn *websocket.Conn) {
	s.topicsMu.Lock()
	defer s.topicsMu.Unlock()

	if connections, exists := s.topics[topic]; exists {
		delete(connections, conn)

		if len(connections) == 0 {
			delete(s.topics, topic)
		}

		conn.WriteJSON(types.Payload{
			Action: "unsubscribed",
			Topic:  string(topic),
			Data:   nil,
		})

		log.Printf("Client %s unsubscribed from topic %s", conn.RemoteAddr(), topic)
	}
}

// UnsubscribeFromResourceTopic removes a connection from a resource-specific topic.
//
// The function takes a topic string, a resource ID string, and a connection as parameters
// and removes the connection from the resource-specific topic map. It is safe to call this
// function concurrently from multiple goroutines.
func (s *SocketServer) UnsubscribeFromResourceTopic(topic topics, resourceID string, conn *websocket.Conn) {
	s.topicsMu.Lock()
	defer s.topicsMu.Unlock()

	resourceTopicKey := fmt.Sprintf("%s:%s", string(topic), resourceID)

	if connections, exists := s.resourceTopics[resourceTopicKey]; exists {
		delete(connections, conn)

		if len(connections) == 0 {
			delete(s.resourceTopics, resourceTopicKey)
		}

		conn.WriteJSON(types.Payload{
			Action: "unsubscribed",
			Topic:  string(topic),
			Data: map[string]interface{}{
				"resourceId": resourceID,
			},
		})

		log.Printf("Client %s unsubscribed from resource topic %s with ID %s",
			conn.RemoteAddr(), topic, resourceID)
	}
}

// BroadcastToTopic sends a message to all connections subscribed to the specified topic.
//
// The function takes a topic string and a payload as parameters and
// sends the payload to all connections subscribed to the topic.
func (s *SocketServer) BroadcastToTopic(topic topics, payload interface{}) {
	s.topicsMu.RLock()
	defer s.topicsMu.RUnlock()

	if connections, exists := s.topics[topic]; exists {
		for conn := range connections {
			err := conn.WriteJSON(types.Payload{
				Action: "message",
				Topic:  string(topic),
				Data:   payload,
			})

			if err != nil {
				log.Printf("Error broadcasting to client %s: %v", conn.RemoteAddr(), err)
				go func(c *websocket.Conn) {
					s.UnsubscribeFromTopic(topic, c)
				}(conn)
			}
		}
		log.Printf("Broadcast message to %d clients on topic %s", len(connections), topic)
	}
}

// BroadcastToResourceTopic sends a message to all connections subscribed to the specific resource topic.
//
// The function takes a topic string, a resource ID string, and a payload as parameters and
// sends the payload to all connections subscribed to the resource-specific topic.
func (s *SocketServer) BroadcastToResourceTopic(topic topics, resourceID string, payload interface{}) {
	s.topicsMu.RLock()
	defer s.topicsMu.RUnlock()

	resourceTopicKey := fmt.Sprintf("%s:%s", string(topic), resourceID)

	if connections, exists := s.resourceTopics[resourceTopicKey]; exists {
		for conn := range connections {
			err := conn.WriteJSON(types.Payload{
				Action: "message",
				Topic:  string(topic),
				Data: map[string]interface{}{
					"resourceId": resourceID,
					"payload":    payload,
				},
			})

			if err != nil {
				log.Printf("Error broadcasting to client %s: %v", conn.RemoteAddr(), err)
				go func(c *websocket.Conn) {
					s.UnsubscribeFromResourceTopic(topic, resourceID, c)
				}(conn)
			}
		}
		log.Printf("Broadcast message to %d clients on resource topic %s with ID %s",
			len(connections), topic, resourceID)
	}
}
