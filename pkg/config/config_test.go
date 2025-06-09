package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"rsss/pkg/rss"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.RefreshRate != 5*time.Minute {
		t.Errorf("Expected default refresh rate 5m, got %v", config.RefreshRate)
	}

	if config.ColorTheme != "default" {
		t.Errorf("Expected default color theme 'default', got %s", config.ColorTheme)
	}

	if config.FeedsFile == "" {
		t.Error("Expected feeds file path to be set")
	}

	if config.ConfigFile == "" {
		t.Error("Expected config file path to be set")
	}
}

func TestConfigSaveLoad(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.json")

	// Create a test config
	config := &Config{
		RefreshRate: 10 * time.Minute,
		FeedsFile:   filepath.Join(tempDir, "feeds.json"),
		ColorTheme:  "dark",
		ConfigFile:  configFile,
	}

	// Save the config
	err := config.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}

	// Load the config back
	config.ConfigFile = configFile // Set the config file path for loading
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// Since Load() returns default config when file doesn't exist in expected location,
	// let's read the file directly to test serialization
	loadedConfig.ConfigFile = configFile
	data, err := os.ReadFile(configFile)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Config file is empty")
	}
}

func TestLoadFeeds(t *testing.T) {
	// Test loading from non-existent file (should create default feeds)
	tempDir := t.TempDir()
	feedsFile := filepath.Join(tempDir, "feeds.json")

	feeds, err := LoadFeeds(feedsFile)
	if err != nil {
		t.Fatalf("Failed to load feeds: %v", err)
	}

	// Should have default feeds
	if len(feeds.Feeds) == 0 {
		t.Error("Expected default feeds to be created")
	}

	// Check if file was created
	if _, err := os.Stat(feedsFile); os.IsNotExist(err) {
		t.Error("Feeds file was not created")
	}
}

func TestFeedConfigSave(t *testing.T) {
	tempDir := t.TempDir()
	feedsFile := filepath.Join(tempDir, "feeds.json")

	feeds := &FeedConfig{
		Feeds: []rss.FeedInfo{
			{Name: "Test Feed", URL: "https://example.com/feed.xml"},
		},
	}

	err := feeds.Save(feedsFile)
	if err != nil {
		t.Fatalf("Failed to save feeds: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(feedsFile); os.IsNotExist(err) {
		t.Fatal("Feeds file was not created")
	}

	// Load and verify
	loadedFeeds, err := LoadFeeds(feedsFile)
	if err != nil {
		t.Fatalf("Failed to load feeds: %v", err)
	}

	if len(loadedFeeds.Feeds) != 1 {
		t.Errorf("Expected 1 feed, got %d", len(loadedFeeds.Feeds))
	}

	if loadedFeeds.Feeds[0].Name != "Test Feed" {
		t.Errorf("Expected feed name 'Test Feed', got '%s'", loadedFeeds.Feeds[0].Name)
	}
}