package tui

import (
	"time"

	"rsss/pkg/config"
	"rsss/pkg/rss"
)

// AppState represents the current state of the application
type AppState int

const (
	StateMenu AppState = iota
	StateFeedView
	StateManageFeeds
	StateConfigure
	StateAddFeed
	StateRemoveFeed
)

// Model represents the TUI application model
type Model struct {
	State        AppState
	MenuSelected int
	Selected     int
	Feeds        *config.FeedConfig
	Articles     []rss.Article
	Config       *config.Config
	Styles       *Styles
	Loading      bool
	Err          error
	Input        string
	LastRefresh  time.Time
	RSSClient    *rss.Client
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, feeds *config.FeedConfig, rssClient *rss.Client) *Model {
	return &Model{
		State:        StateMenu,
		MenuSelected: 0,
		Selected:     0,
		Feeds:        feeds,
		Config:       cfg,
		Styles:       NewStyles(cfg.ColorTheme),
		Loading:      true,
		LastRefresh:  time.Now(),
		RSSClient:    rssClient,
	}
}