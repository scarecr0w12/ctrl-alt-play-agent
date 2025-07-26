package api

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModManager handles mod installation and management that the panel expects
type ModManager struct {
	baseDir string
}

// NewModManager creates a new mod manager
func NewModManager(baseDir string) *ModManager {
	if baseDir == "" {
		baseDir = "/opt/gameservers"
	}
	return &ModManager{
		baseDir: baseDir,
	}
}

// ModInfo represents information about an installed mod
type ModInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}

// Mod management operations that the panel expects
func (s *Server) handleInstallMod(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	modID, ok := data["modId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid modId",
		}
	}

	modURL, _ := data["modUrl"].(string)
	modVersion, _ := data["version"].(string)

	// For now, simulate mod installation
	// In a real implementation, this would download and install the mod
	serverDir := filepath.Join("/opt/gameservers", serverID, "mods")
	if err := os.MkdirAll(serverDir, 0755); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create mods directory: %v", err),
		}
	}

	// Create a simple mod info file to track installation
	modInfoPath := filepath.Join(serverDir, modID+".mod")
	modContent := fmt.Sprintf("id=%s\nversion=%s\nurl=%s\ninstalled=true\n", modID, modVersion, modURL)

	if err := os.WriteFile(modInfoPath, []byte(modContent), 0644); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to install mod: %v", err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"modId":    modID,
			"version":  modVersion,
			"message":  "Mod installed successfully",
		},
	}
}

func (s *Server) handleUninstallMod(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	modID, ok := data["modId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid modId",
		}
	}

	// Remove the mod info file
	modInfoPath := filepath.Join("/opt/gameservers", serverID, "mods", modID+".mod")
	if err := os.Remove(modInfoPath); err != nil && !os.IsNotExist(err) {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to uninstall mod: %v", err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"modId":    modID,
			"message":  "Mod uninstalled successfully",
		},
	}
}

func (s *Server) handleListMods(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	modsDir := filepath.Join("/opt/gameservers", serverID, "mods")

	// Check if mods directory exists
	if _, err := os.Stat(modsDir); os.IsNotExist(err) {
		return CommandResponse{
			Success: true,
			Data: map[string]interface{}{
				"serverId": serverID,
				"mods":     []ModInfo{},
				"count":    0,
			},
		}
	}

	// List all .mod files
	entries, err := os.ReadDir(modsDir)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to list mods: %v", err),
		}
	}

	var mods []ModInfo
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), ".mod") {
			modID := strings.TrimSuffix(entry.Name(), ".mod")

			// Read mod info from file
			modInfoPath := filepath.Join(modsDir, entry.Name())
			content, err := os.ReadFile(modInfoPath)
			if err != nil {
				continue
			}

			// Parse simple key=value format
			modInfo := ModInfo{
				ID:      modID,
				Name:    modID, // Use ID as name for now
				Version: "unknown",
				Enabled: true,
			}

			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				if strings.Contains(line, "=") {
					parts := strings.SplitN(line, "=", 2)
					if len(parts) == 2 {
						key, value := strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1])
						switch key {
						case "version":
							modInfo.Version = value
						case "name":
							modInfo.Name = value
						case "description":
							modInfo.Description = value
						}
					}
				}
			}

			mods = append(mods, modInfo)
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"mods":     mods,
			"count":    len(mods),
		},
	}
}
