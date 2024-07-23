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

	"go.uber.org/zap"
)

func main() {
	fmt.Println("Hello")
	ctx, cancel := context.WithCancel(context.Background())
	log.Print(ctx)
	defer cancel()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	logger.Info("Starting Banking Application...")
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

	// Initialize Router
	router := app.NewRouter(services)

	server := &http.Server{
		Addr:    "localhost:1925",
		Handler: router,
	}

	go func() {
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

	log.Print("Shutting Down Banking Application...")

	// Attempt to gracefully shut down the server
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server shutdown error: %s\n", err)
	}

	log.Print("Banking Application has been gracefully shut down.")
}
