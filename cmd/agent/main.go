package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/client"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/health"
)

func main() {
	log.Println("Starting Ctrl-Alt-Play Agent...")

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Initialize health check server
	healthServer := health.NewServer(cfg.NodeID, "1.0.0")

	// Start health check server in background
	go func() {
		log.Printf("Starting health check server on :%s", cfg.HealthPort)
		if err := healthServer.StartServer(cfg.HealthPort); err != nil {
			log.Printf("Health check server error: %v", err)
		}
	}()

	// Initialize Docker manager
	dockerManager, err := docker.NewManager()
	if err != nil {
		log.Fatalf("Error initializing Docker manager: %v", err)
	}
	defer func() {
		if err := dockerManager.Close(); err != nil {
			log.Printf("Error closing Docker manager: %v", err)
		}
	}()

	// Initialize WebSocket client
	wsClient := client.NewClient(cfg, dockerManager)

	// Connect to panel
	if err := wsClient.Connect(); err != nil {
		log.Fatalf("Error connecting to panel: %v", err)
	}

	// Update health status to connected
	healthServer.SetConnectionStatus(true)

	// Start client
	if err := wsClient.Start(); err != nil {
		log.Fatalf("Error starting client: %v", err)
	}

	log.Printf("Agent connected to panel at %s as node %s", cfg.PanelURL, cfg.NodeID)
	log.Printf("Health check available at http://localhost:%s/health", cfg.HealthPort)

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Println("Received interrupt signal, shutting down...")

	// Update health status to disconnected
	healthServer.SetConnectionStatus(false)

	// Graceful shutdown
	wsClient.Stop()
	log.Println("Agent stopped")
}
