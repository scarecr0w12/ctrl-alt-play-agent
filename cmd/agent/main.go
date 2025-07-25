package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/api"
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

	// Initialize API server
	apiServer := api.NewServer(cfg, dockerManager)

	// Start combined API/Health server in background
	go func() {
		log.Printf("Starting combined API/Health server on :%s", cfg.HealthPort)
		if err := apiServer.StartServer(cfg.HealthPort, healthServer); err != nil {
			log.Printf("API/Health server error: %v", err)
		}
	}()

	// Initialize WebSocket client
	wsClient := client.NewClient(cfg, dockerManager)

	// Try to connect to panel (but don't fail if it's not available)
	if err := wsClient.Connect(); err != nil {
		log.Printf("Warning: Could not connect to panel: %v", err)
		log.Printf("Agent will continue running with HTTP API only")
		healthServer.SetConnectionStatus(false)
	} else {
		// Update health status to connected
		healthServer.SetConnectionStatus(true)

		// Start client
		if err := wsClient.Start(); err != nil {
			log.Printf("Warning: Error starting WebSocket client: %v", err)
			healthServer.SetConnectionStatus(false)
		} else {
			log.Printf("Agent connected to panel at %s as node %s", cfg.PanelURL, cfg.NodeID)
		}
	}

	log.Printf("Health check available at http://localhost:%s/health", cfg.HealthPort)
	log.Printf("API commands available at http://localhost:%s/api/command", cfg.HealthPort)

	// Wait for interrupt signal
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	log.Println("Received interrupt signal, shutting down...")

	// Update health status to disconnected
	healthServer.SetConnectionStatus(false)

	// Graceful shutdown
	if wsClient != nil {
		wsClient.Stop()
	}
	log.Println("Agent stopped")
}
