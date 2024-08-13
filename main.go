package main

import (
	"Saving-Account-Banking-System/app"
	"Saving-Account-Banking-System/repository"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/rs/cors"
	"go.uber.org/zap"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	log.Print(ctx)
	defer cancel()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("=> Starting Banking Application...")
	fmt.Println("*** WELCOME to BANKING SYSTEM !! ***")

	// To Initialize Database
	database, err := repository.InitializeDB()
	if err != nil {
		log.Fatalln(err)
	}
	defer database.Close()

	// repository.InsertSeedData()

	// Initialize Service
	services := app.NewServices(database)

	// Initialize RouterCORS middleware
	router := app.NewRouter(services)

	// CORS middleware
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:1925"},
		AllowCredentials: true,
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
	})

	server := &http.Server{
		Addr:    "localhost:1925",
		Handler: cors.Handler(router),
	}

	go func() {
		fmt.Println("Server started")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("\n Server error: %s", err)
		}
	}()

	// Listen for an interrupt signal to gracefully shutdown the server
	c := make(chan os.Signal, 1) //buffered chan
	signal.Notify(c, os.Interrupt)
	<-c

	// Create a deadline for shutting down the server
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Print("=> Shutting Down Banking Application...")

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}
}
