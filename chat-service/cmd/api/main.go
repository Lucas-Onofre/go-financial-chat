package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/websocket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/auth/jwt"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/broker"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/dao"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/handler"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/repository"
	"github.com/Lucas-Onofre/financial-chat/chat-service/internal/user/service"
)

func main() {
	mux := http.NewServeMux()

	// Postgres
	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(
		&dao.User{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	// RabbitMQ
	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"))

	var rb *broker.RabbitMQBroker
	var retryErr error
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		rb, retryErr = broker.NewRabbitMQBroker(rabbitmqURL)
		if retryErr == nil {
			break
		}
		log.Printf("Failed to connect to RabbitMQ (attempt %d/%d): %v", i+1, maxRetries, err)
		time.Sleep(5 * time.Second)
	}
	if retryErr != nil {
		log.Fatal("Failed to connect to RabbitMQ after all retries:", err)
	}
	defer rb.Close()

	// User auth
	jwtService := jwt.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	userRepo := repository.NewRepository(db)
	userService := service.New(userRepo, jwtService)
	userHandler := handler.New(*userService)

	mux.HandleFunc("/register", handleMethod(http.MethodPost, userHandler.Register))
	mux.HandleFunc("/login", handleMethod(http.MethodPost, userHandler.Login))

	// Websocket Hub
	hub := websocket.NewHub()
	go hub.Run()

	// Websocket
	mux.HandleFunc("/ws", websocket.WsHandler(hub, jwtService))

	// Health check
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Starting server on port 8081")
	if err := http.ListenAndServe(":8081", corsMiddleware(mux)); err != nil {
		log.Fatal(err)
	}
}

func handleMethod(method string, handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		handlerFunc(w, r)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
