package tui

import (
	"time"

	"rsss/pkg/rss"
)

// FetchMsg represents the result of fetching RSS feeds
type FetchMsg struct {
	Articles []rss.Article
	Err      error
}

// TickMsg represents a timer tick for auto-refresh
type TickMsg time.Time

// SaveMsg represents the result of a save operation
type SaveMsg struct {
	Success bool
	Err     error
}

// OpenURLMsg represents the result of opening a URL in browser
type OpenURLMsg struct {
	URL string
	Err error
}