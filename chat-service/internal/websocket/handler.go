package websocket

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"

	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func WsHandler(hub *Hub, jwtService *jwt.JWTService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "token required", http.StatusUnauthorized)
			return
		}

		claims, err := jwtService.ValidateToken(token)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		username := claims.Username
		roomID := r.URL.Query().Get("room")
		if roomID == "" {
			roomID = "general"
		}

		fmt.Printf("websocket connection for user %s (%s) joining room %s\n", claims.UserID, username, roomID)

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "could not open websocket connection", http.StatusBadRequest)
			return
		}

		client := &Client{
			Hub:      hub,
			Conn:     conn,
			Send:     make(chan []byte, 256),
			UserID:   claims.UserID,
			RoomID:   roomID,
			Username: username,
		}

		client.Hub.Register <- client

		go client.WritePump()
		go client.ReadPump()
	}
}
