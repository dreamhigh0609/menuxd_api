package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	_ "github.com/joho/godotenv/autoload"

	"github.com/super-radmir/menuxd_api/internal/server"
	"github.com/super-radmir/menuxd_api/internal/storage"
)

func main() {
	debug := flag.Bool("debug", false, "Debug mode activation")
	flag.Parse()

	var port string
	if port = os.Getenv("PORT"); port == "" {
		port = "1323"
	}

	server.SetConfig(port, *debug)
	srv, err := server.New()
	if err != nil {
		log.Fatal(err)
	}

	// Start the server.
	go func() {
		log.Printf("Server running on http://localhost:%s", port)
		log.Fatal(srv.ListenAndServe())
	}()

	// Wait for an interrupt.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c

	// Attempt a graceful shutdown.
	storage.Close()
	server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	srv.Shutdown(ctx)
}
