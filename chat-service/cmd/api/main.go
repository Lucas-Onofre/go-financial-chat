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
	rb, err := broker.NewRabbitMQBroker(rabbitmqURL)
	if err != nil {
		log.Fatal(err)
	}
	defer rb.Close()

	jwtService := jwt.NewJWTService(os.Getenv("SECRET_KEY"), 24*time.Hour)
	userRepo := repository.NewRepository(db)
	userService := service.New(userRepo, jwtService)
	userHandler := handler.New(*userService)

	mux.HandleFunc("/register", userHandler.Register)
	mux.HandleFunc("/login", userHandler.Login)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	if err := http.ListenAndServe(":8081", nil); err != nil {
		log.Fatal(err)
	}
}
