package main

import (
	"chatapi/internal/auth"
	"chatapi/internal/chatroom"
	"chatapi/internal/config"
	"chatapi/internal/database"
	"chatapi/internal/message"
	"chatapi/internal/ratelimit"
	"chatapi/internal/server"
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/time/rate"
)

func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db := database.InitDB(conf.DatabaseURL)
	jwtKey := []byte(conf.JWTSecret)

	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtKey)
	authHandler := auth.NewHandler(authService, authRepo, jwtKey)

	chatroomRepo := chatroom.NewRepository(db)
	chatroomService := chatroom.NewService(chatroomRepo)
	chatroomHandler := chatroom.NewHandler(chatroomService)

	messageRepo := message.NewRepository(db)
	messageService := message.NewService(messageRepo, chatroomRepo)
	messageHandler := message.NewHandler(messageService)

	rl := ratelimit.New(rate.Limit(5), 10)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go rl.Cleanup(ctx, 5*time.Minute, 15*time.Minute)

	handler := server.New(authHandler, chatroomHandler, messageHandler, rl)

	srv := &http.Server{
		Addr:         ":" + conf.Port,
		Handler:      handler,
		ReadTimeout:  conf.ReadTimeout,
		WriteTimeout: conf.WriteTimeout,
		IdleTimeout:  conf.IdleTimeout,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()
	log.Println("Server started on", srv.Addr)

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	if err := db.Close(); err != nil {
		log.Println("Error closing DB:", err)
	}

	log.Println("Server exited cleanly")
}
