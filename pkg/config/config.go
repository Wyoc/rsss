package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"

	"rsss/pkg/rss"
)

// Config represents the application configuration
type Config struct {
	RefreshRate time.Duration `json:"refresh_rate"`
	FeedsFile   string        `json:"feeds_file"`
	ColorTheme  string        `json:"color_theme"`
	ConfigFile  string        `json:"-"`
}

// FeedConfig represents the feeds configuration
type FeedConfig struct {
	Feeds []rss.FeedInfo `json:"feeds"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	configDir := filepath.Join(homeDir, ".config", "rsss")

	return &Config{
		RefreshRate: 5 * time.Minute,
		FeedsFile:   filepath.Join(configDir, "feeds.json"),
		ColorTheme:  "default",
		ConfigFile:  filepath.Join(configDir, "config.json"),
	}
}

// Load loads configuration from file or returns default if file doesn't exist
func Load() (*Config, error) {
	config := DefaultConfig()

	if _, err := os.Stat(config.ConfigFile); os.IsNotExist(err) {
		return config, nil
	}

	data, err := os.ReadFile(config.ConfigFile)
	if err != nil {
		return config, err
	}

	if err := json.Unmarshal(data, config); err != nil {
		return config, err
	}

	return config, nil
}

// Save saves the configuration to file
func (c *Config) Save() error {
	dir := filepath.Dir(c.ConfigFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(c.ConfigFile, data, 0644)
}

// LoadFeeds loads feed configuration from file
func LoadFeeds(filename string) (*FeedConfig, error) {
	feeds := &FeedConfig{Feeds: []rss.FeedInfo{}}

	if _, err := os.Stat(filename); os.IsNotExist(err) {
		// Add default feeds for new users
		feeds.Feeds = []rss.FeedInfo{
			{Name: "BBC News", URL: "https://feeds.bbci.co.uk/news/rss.xml"},
			{Name: "TechCrunch", URL: "https://techcrunch.com/feed/"},
		}
		// Save the default feeds
		feeds.Save(filename)
		return feeds, nil
	}

	data, err := os.ReadFile(filename)
	if err != nil {
		return feeds, err
	}

	if err := json.Unmarshal(data, feeds); err != nil {
		return feeds, err
	}

	return feeds, nil
}

// Save saves the feed configuration to file
func (f *FeedConfig) Save(filename string) error {
	dir := filepath.Dir(filename)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}