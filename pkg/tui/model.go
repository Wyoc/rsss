package tui

import (
	"strings"
	"time"

	"rsss/pkg/config"
	"rsss/pkg/rss"
)

// AppState represents the current state of the application
type AppState int

const (
	StateMenu AppState = iota
	StateFeedView
	StateArticleView
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
	Width        int
	Height       int
	ViewportTop  int // For scrolling in feed view
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

// getMaxVisibleArticles calculates how many articles can fit on screen
func (m *Model) getMaxVisibleArticles() int {
	availableHeight := m.Height
	if availableHeight == 0 {
		availableHeight = 24 // Fallback
	}
	// Reserve space for header and help text (3 lines total)
	contentHeight := availableHeight - 3
	return max(1, contentHeight-1)
}

// wrapText wraps text to fit within the specified width
func (m *Model) wrapText(text string, width int) string {
	if len(text) <= width {
		return text
	}
	
	words := strings.Fields(text)
	if len(words) == 0 {
		return text
	}
	
	var lines []string
	var currentLine strings.Builder
	
	for _, word := range words {
		// If adding this word would exceed width, start a new line
		if currentLine.Len() > 0 && currentLine.Len()+1+len(word) > width {
			lines = append(lines, currentLine.String())
			currentLine.Reset()
		}
		
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}
	
	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}
	
	return strings.Join(lines, "\n")
}