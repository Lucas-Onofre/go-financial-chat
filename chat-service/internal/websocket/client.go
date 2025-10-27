package websocket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	Hub      *Hub
	Conn     *websocket.Conn
	Send     chan []byte
	UserID   string
	RoomID   string
	Username string
}

type Message struct {
	Type      string `json:"type"`
	UserID    string `json:"user_id"`
	Username  string `json:"username"`
	RoomID    string `json:"room_id"`
	Content   string `json:"content"`
	Timestamp int64  `json:"timestamp"`
}

type MessageType string

const (
	MessageTypeChat       MessageType = "default"
	MessageTypeUserJoined MessageType = "user_joined"
	MessageTypeUserLeft   MessageType = "user_left"
)

func (mt MessageType) ToString() string {
	return string(mt)
}

func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, messageBytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(messageBytes, &message); err != nil {
			log.Printf("error unmarshaling message: %v", err)
			continue
		}

		message.UserID = c.UserID
		message.Username = c.Username
		message.RoomID = c.RoomID

		c.Hub.Broadcast <- message
	}
}

func (c *Client) WritePump() {
	defer c.Conn.Close()

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
				log.Printf("error writing message: %v", err)
				return
			}
		case <-ticker.C:
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("error sending ping: %v", err)
				return
			}
		}
	}
}
