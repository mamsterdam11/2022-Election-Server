package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/mamsterdam11/Snowflake-News/news"
	"github.com/mamsterdam11/Snowflake-News/server"
)

// Snowflake-News server
func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)

	// gracefully handle Interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		receivedSignal := <-c
		log.Printf("Received %s, shutting down server", receivedSignal)
		cancel()
	}()

	// Create and start NewsCollector
	newsCollector := news.NewNewsCollector()
	newsCollector.Start(ctx)

	// Create and start WebServer
	webServer := server.NewWebServer(newsCollector)
	webServer.Serve(ctx)
}
