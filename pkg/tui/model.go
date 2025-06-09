package tui

import (
	"fmt"
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
	
	// Notification system
	SeenArticles    map[string]bool // Track seen article URLs
	NewArticleCount int             // Count of new articles since last check
	ShowNotification bool           // Whether to show notification
	NotificationMsg  string         // Notification message to display
}

// NewModel creates a new TUI model
func NewModel(cfg *config.Config, feeds *config.FeedConfig, rssClient *rss.Client) *Model {
	// Load seen articles from file
	seenArticles := make(map[string]bool)
	if seen, err := config.LoadSeenArticles(cfg.SeenArticlesFile); err == nil {
		seenArticles = seen.Articles
	}
	
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
		
		// Initialize notification system with loaded data
		SeenArticles:     seenArticles,
		NewArticleCount:  0,
		ShowNotification: false,
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

// checkForNewArticles compares current articles with seen articles and updates notification state
func (m *Model) checkForNewArticles(articles []rss.Article) {
	if len(m.SeenArticles) == 0 {
		// First run - mark all current articles as seen without notification
		for _, article := range articles {
			m.SeenArticles[article.Link] = true
		}
		m.saveSeenArticles()
		return
	}
	
	newCount := 0
	for _, article := range articles {
		if !m.SeenArticles[article.Link] {
			newCount++
			m.SeenArticles[article.Link] = true
		}
	}
	
	if newCount > 0 && m.Config.EnableNotifications {
		m.NewArticleCount = newCount
		m.ShowNotification = true
		if newCount == 1 {
			m.NotificationMsg = "🔔 1 new article available!"
		} else {
			m.NotificationMsg = fmt.Sprintf("🔔 %d new articles available!", newCount)
		}
		m.saveSeenArticles()
	} else if newCount > 0 {
		// Still save seen articles even if notifications are disabled
		m.saveSeenArticles()
	}
}

// dismissNotification clears the current notification
func (m *Model) dismissNotification() {
	m.ShowNotification = false
	m.NotificationMsg = ""
	m.NewArticleCount = 0
}

// saveSeenArticles saves the current seen articles to file
func (m *Model) saveSeenArticles() {
	seen := &config.SeenArticles{Articles: m.SeenArticles}
	seen.Save(m.Config.SeenArticlesFile)
}