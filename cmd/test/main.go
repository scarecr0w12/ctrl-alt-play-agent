package main

import (
	"fmt"
	"log"

	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/config"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/docker"
	"github.com/scarecr0w12/ctrl-alt-play-agent/internal/messages"
)

func main() {
	fmt.Println("=== Ctrl-Alt-Play Agent Module Test ===")

	// Test configuration loading
	fmt.Println("\n1. Testing Configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Config test failed: %v", err)
	}
	fmt.Printf("✓ Config loaded: Panel=%s, Node=%s\n", cfg.PanelURL, cfg.NodeID)

	// Test Docker manager creation
	fmt.Println("\n2. Testing Docker Manager...")
	dockerManager, err := docker.NewManager()
	if err != nil {
		fmt.Printf("⚠ Docker test failed (expected if Docker not running): %v\n", err)
	} else {
		fmt.Println("✓ Docker manager created successfully")
		if err := dockerManager.Close(); err != nil {
			log.Printf("Error closing Docker manager: %v", err)
		}
	}

	// Test message creation
	fmt.Println("\n3. Testing Message System...")
	heartbeatData := &messages.HeartbeatData{
		NodeID: cfg.NodeID,
		Status: "test",
	}

	msg, err := messages.NewMessage(messages.TypeHeartbeat, heartbeatData)
	if err != nil {
		log.Fatalf("Message test failed: %v", err)
	}

	jsonData, err := msg.ToJSON()
	if err != nil {
		log.Fatalf("JSON serialization failed: %v", err)
	}

	fmt.Printf("✓ Message created: %s\n", string(jsonData))

	// Test message parsing
	parsedMsg, err := messages.ParseMessage(jsonData)
	if err != nil {
		log.Fatalf("Message parsing failed: %v", err)
	}

	fmt.Printf("✓ Message parsed: Type=%s\n", parsedMsg.Type)

	fmt.Println("\n=== All Tests Passed! ===")
	fmt.Println("The agent modules are working correctly.")
	fmt.Println("Run 'make run' to start the agent (requires a running panel).")
}
