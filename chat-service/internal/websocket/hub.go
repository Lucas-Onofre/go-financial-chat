package websocket

import (
	"encoding/json"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/broker"
	"log"
	"time"
)

// Hub maintains the set of active clients and broadcasts messages to the clients;
// Rooms: Map of RoomID to set of Clients;
// Broadcast: Channel responsible for broadcasting messages to rooms;
// Register: Responsible for registering new clients and adding them to rooms;
// Unregister: Responsible for unregistering clients and removing them from rooms.
type Hub struct {
	Rooms      map[string]map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Broker     broker.Producer
}

func NewHub(rb broker.Producer) *Hub {
	return &Hub{
		Rooms:      make(map[string]map[*Client]bool),
		Broadcast:  make(chan Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broker:     rb,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if h.Rooms[client.RoomID] == nil {
				h.Rooms[client.RoomID] = make(map[*Client]bool)
			}

			// Close any unwanted existing connections for the same user in the room
			for c := range h.Rooms[client.RoomID] {
				if c.UserID == client.UserID {
					log.Printf("user %s already in room %s, closing old connection", c.UserID, client.RoomID)
					c.Conn.Close()
				}
			}

			h.Rooms[client.RoomID][client] = true
			log.Printf("client %s joined room %s", client.UserID, client.RoomID)

			joinMessage := Message{
				Type:      MessageTypeUserJoined.ToString(),
				UserID:    client.UserID,
				Username:  client.Username,
				RoomID:    client.RoomID,
				Content:   client.Username + " joined the room",
				Timestamp: time.Now().Unix(),
			}
			h.broadcastToRoom(client.RoomID, joinMessage)

		case client := <-h.Unregister:
			if room, exists := h.Rooms[client.RoomID]; exists {
				if _, ok := room[client]; ok {
					delete(room, client)
					close(client.Send)

					log.Printf("client %s left room %s", client.UserID, client.RoomID)

					leaveMessage := Message{
						Type:      MessageTypeUserLeft.ToString(),
						UserID:    client.UserID,
						Username:  client.Username,
						RoomID:    client.RoomID,
						Content:   client.Username + " left the room",
						Timestamp: time.Now().Unix(),
					}
					h.broadcastToRoom(client.RoomID, leaveMessage)

					if len(room) == 0 {
						delete(h.Rooms, client.RoomID)
					}
				}
			}

		case message := <-h.Broadcast:
			log.Printf("Broadcast to room %s, clients: %d", message.RoomID, len(h.Rooms[message.RoomID]))
			h.broadcastToRoom(message.RoomID, message)
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, message Message) {
	if room, exists := h.Rooms[roomID]; exists {
		messageBytes, err := json.Marshal(message)
		if err != nil {
			log.Printf("error marshaling message: %v", err)
			return
		}

		for client := range room {
			select {
			case client.Send <- messageBytes:
			default:
				close(client.Send)
				delete(room, client)
			}
		}
	}
}
