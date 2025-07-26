package api

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// FileManager handles file operations that the panel expects
type FileManager struct {
	baseDir string
}

// NewFileManager creates a new file manager
func NewFileManager(baseDir string) *FileManager {
	if baseDir == "" {
		baseDir = "/opt/gameservers"
	}
	return &FileManager{
		baseDir: baseDir,
	}
}

// File operations that the panel expects
func (s *Server) handleListFiles(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	pathStr, _ := data["path"].(string)
	if pathStr == "" {
		pathStr = "/"
	}

	// Build the full path (serverId as subdirectory)
	serverDir := filepath.Join("/opt/gameservers", serverID)
	fullPath := filepath.Join(serverDir, pathStr)

	// Ensure we don't go outside the server directory
	if !strings.HasPrefix(fullPath, serverDir) {
		return CommandResponse{
			Success: false,
			Error:   "Invalid path - cannot access files outside server directory",
		}
	}

	// Check if directory exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Path does not exist: %s", pathStr),
		}
	}

	// List directory contents
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to list directory: %v", err),
		}
	}

	var files []map[string]interface{}
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		fileInfo := map[string]interface{}{
			"name":     entry.Name(),
			"type":     "file",
			"size":     info.Size(),
			"modified": info.ModTime(),
		}

		if entry.IsDir() {
			fileInfo["type"] = "directory"
		}

		files = append(files, fileInfo)
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"path":     pathStr,
			"files":    files,
		},
	}
}

func (s *Server) handleReadFile(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	filePath, ok := data["path"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file path",
		}
	}

	// Build the full path
	serverDir := filepath.Join("/opt/gameservers", serverID)
	fullPath := filepath.Join(serverDir, filePath)

	// Ensure we don't go outside the server directory
	if !strings.HasPrefix(fullPath, serverDir) {
		return CommandResponse{
			Success: false,
			Error:   "Invalid path - cannot access files outside server directory",
		}
	}

	// Read the file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to read file: %v", err),
		}
	}

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get file info: %v", err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"path":     filePath,
			"content":  string(content),
			"size":     info.Size(),
			"modified": info.ModTime(),
		},
	}
}

func (s *Server) handleWriteFile(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	filePath, ok := data["path"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file path",
		}
	}

	content, ok := data["content"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file content",
		}
	}

	// Build the full path
	serverDir := filepath.Join("/opt/gameservers", serverID)
	fullPath := filepath.Join(serverDir, filePath)

	// Ensure we don't go outside the server directory
	if !strings.HasPrefix(fullPath, serverDir) {
		return CommandResponse{
			Success: false,
			Error:   "Invalid path - cannot access files outside server directory",
		}
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create directory: %v", err),
		}
	}

	// Write the file
	if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to write file: %v", err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"path":     filePath,
			"message":  "File written successfully",
		},
	}
}

func (s *Server) handleUploadFile(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	filePath, ok := data["path"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file path",
		}
	}

	// Expect base64 encoded content
	contentB64, ok := data["content"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file content",
		}
	}

	// Decode base64 content
	content, err := base64.StdEncoding.DecodeString(contentB64)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Invalid base64 content: %v", err),
		}
	}

	// Build the full path
	serverDir := filepath.Join("/opt/gameservers", serverID)
	fullPath := filepath.Join(serverDir, filePath)

	// Ensure we don't go outside the server directory
	if !strings.HasPrefix(fullPath, serverDir) {
		return CommandResponse{
			Success: false,
			Error:   "Invalid path - cannot access files outside server directory",
		}
	}

	// Create directory if it doesn't exist
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to create directory: %v", err),
		}
	}

	// Write the file
	if err := os.WriteFile(fullPath, content, 0644); err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to upload file: %v", err),
		}
	}

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"path":     filePath,
			"size":     len(content),
			"message":  "File uploaded successfully",
		},
	}
}

func (s *Server) handleDownloadFile(data map[string]interface{}) CommandResponse {
	serverID, ok := data["serverId"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid serverId",
		}
	}

	filePath, ok := data["path"].(string)
	if !ok {
		return CommandResponse{
			Success: false,
			Error:   "Missing or invalid file path",
		}
	}

	// Build the full path
	serverDir := filepath.Join("/opt/gameservers", serverID)
	fullPath := filepath.Join(serverDir, filePath)

	// Ensure we don't go outside the server directory
	if !strings.HasPrefix(fullPath, serverDir) {
		return CommandResponse{
			Success: false,
			Error:   "Invalid path - cannot access files outside server directory",
		}
	}

	// Read the file
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to read file: %v", err),
		}
	}

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		return CommandResponse{
			Success: false,
			Error:   fmt.Sprintf("Failed to get file info: %v", err),
		}
	}

	// Encode content as base64 for safe transport
	contentB64 := base64.StdEncoding.EncodeToString(content)

	return CommandResponse{
		Success: true,
		Data: map[string]interface{}{
			"serverId": serverID,
			"path":     filePath,
			"content":  contentB64,
			"size":     info.Size(),
			"modified": info.ModTime(),
			"encoding": "base64",
		},
	}
}
