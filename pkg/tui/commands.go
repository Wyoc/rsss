package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"rsss/pkg/browser"
	"rsss/pkg/config"
	"rsss/pkg/rss"
)

// FetchAllFeedsCmd fetches all configured RSS feeds
func FetchAllFeedsCmd(client *rss.Client, feeds *config.FeedConfig) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		articles, err := client.FetchMultipleFeeds(feeds.Feeds)
		return FetchMsg{Articles: articles, Err: err}
	})
}

// TickCmd creates a ticker command for auto-refresh
func TickCmd(interval time.Duration) tea.Cmd {
	return tea.Tick(interval, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}

// OpenURLCmd opens a URL in the default browser
func OpenURLCmd(url string) tea.Cmd {
	return tea.Cmd(func() tea.Msg {
		err := browser.Open(url)
		return OpenURLMsg{URL: url, Err: err}
	})
}