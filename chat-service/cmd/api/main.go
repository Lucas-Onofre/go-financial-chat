package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

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

	db, err := gorm.Open(postgres.Open(os.Getenv("DATABASE_URL")), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database:", err)
	}

	if err := db.AutoMigrate(
		&dao.User{},
	); err != nil {
		log.Fatal("failed to migrate database:", err)
	}

	rabbitmqURL := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		os.Getenv("RABBITMQ_USER"),
		os.Getenv("RABBITMQ_PASSWORD"),
		os.Getenv("RABBITMQ_HOST"),
		os.Getenv("RABBITMQ_PORT"))

	// Retry logic for RabbitMQ connection
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

	jwtService := jwt.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	userRepo := repository.NewRepository(db)
	userService := service.New(userRepo, jwtService)
	userHandler := handler.New(*userService)

	mux.HandleFunc("/register", handleMethod(http.MethodPost, userHandler.Register))
	mux.HandleFunc("/login", handleMethod(http.MethodPost, userHandler.Login))
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	log.Println("Starting server on port 8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
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
