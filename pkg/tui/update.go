package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"rsss/pkg/rss"
)

// Init initializes the TUI model
func (m *Model) Init() tea.Cmd {
	var cmds []tea.Cmd

	if len(m.Feeds.Feeds) > 0 {
		cmds = append(cmds, FetchAllFeedsCmd(m.RSSClient, m.Feeds))
	}

	cmds = append(cmds, TickCmd(m.Config.RefreshRate))

	return tea.Batch(cmds...)
}

// Update handles all TUI messages
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		
	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case FetchMsg:
		m.Loading = false
		
		// Check for new articles before updating the article list
		if msg.Err == nil && len(msg.Articles) > 0 {
			m.checkForNewArticles(msg.Articles)
		}
		
		m.Articles = msg.Articles
		m.Err = msg.Err
		m.LastRefresh = time.Now()
		m.Selected = 0
		m.ViewportTop = 0 // Reset viewport when new articles load

	case TickMsg:
		if time.Since(m.LastRefresh) >= m.Config.RefreshRate {
			return m, FetchAllFeedsCmd(m.RSSClient, m.Feeds)
		}
		return m, TickCmd(m.Config.RefreshRate)

	case SaveMsg:
		if msg.Success {
			m.Err = nil
		} else {
			m.Err = msg.Err
		}

	case OpenURLMsg:
		m.State = StateFeedView // Return to feed view after opening URL
		if msg.Err != nil {
			m.Err = msg.Err
		} else {
			m.Err = nil // Clear any previous errors
		}
	}

	return m, nil
}

// handleKeyPress routes key presses to the appropriate state handler
func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Global notification dismissal
	if m.ShowNotification && (msg.String() == "space" || msg.String() == "enter" || msg.String() == "n") {
		m.dismissNotification()
		return m, nil
	}
	
	switch m.State {
	case StateMenu:
		return m.updateMenu(msg)
	case StateFeedView:
		return m.updateFeedView(msg)
	case StateArticleView:
		return m.updateArticleView(msg)
	case StateManageFeeds:
		return m.updateManageFeeds(msg)
	case StateConfigure:
		return m.updateConfigure(msg)
	case StateAddFeed:
		return m.updateAddFeed(msg)
	case StateRemoveFeed:
		return m.updateRemoveFeed(msg)
	}
	return m, nil
}

// updateMenu handles menu navigation
func (m *Model) updateMenu(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q":
		return m, tea.Quit
	case "up", "k":
		if m.MenuSelected > 0 {
			m.MenuSelected--
		}
	case "down", "j":
		if m.MenuSelected < 2 {
			m.MenuSelected++
		}
	case "enter":
		switch m.MenuSelected {
		case 0:
			m.State = StateFeedView
			m.Selected = 0
		case 1:
			m.State = StateManageFeeds
			m.Selected = 0
		case 2:
			m.State = StateConfigure
			m.Selected = 0
		}
	}
	return m, nil
}

// updateFeedView handles feed view navigation
func (m *Model) updateFeedView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.State = StateMenu
		return m, nil
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
			// Scroll viewport up if needed
			if m.Selected < m.ViewportTop {
				m.ViewportTop = m.Selected
			}
		}
	case "down", "j":
		if m.Selected < len(m.Articles)-1 {
			m.Selected++
			// Scroll viewport down if needed
			maxVisible := m.getMaxVisibleArticles()
			if m.Selected >= m.ViewportTop+maxVisible {
				m.ViewportTop = m.Selected - maxVisible + 1
			}
		}
	case "enter":
		if len(m.Articles) > 0 && m.Selected < len(m.Articles) {
			// Switch to article view to show content
			m.State = StateArticleView
			return m, nil
		}
	case "r":
		m.Loading = true
		return m, FetchAllFeedsCmd(m.RSSClient, m.Feeds)
	}
	return m, nil
}

// updateManageFeeds handles feed management
func (m *Model) updateManageFeeds(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.State = StateMenu
		return m, nil
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < len(m.Feeds.Feeds)-1 {
			m.Selected++
		}
	case "a":
		m.State = StateAddFeed
		m.Input = ""
		return m, nil
	case "d":
		if len(m.Feeds.Feeds) > 0 {
			m.State = StateRemoveFeed
			m.Selected = 0 // Reset selection for remove view
			return m, nil
		}
	case "enter":
		if len(m.Feeds.Feeds) > 0 && m.Selected < len(m.Feeds.Feeds) {
			// Enter can also be used to delete the selected feed
			m.State = StateRemoveFeed
			return m, nil
		}
	}
	return m, nil
}

// updateConfigure handles configuration changes
func (m *Model) updateConfigure(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.State = StateMenu
		return m, nil
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < 2 {
			m.Selected++
		}
	case "enter", "space":
		switch m.Selected {
		case 0:
			if m.Config.RefreshRate == 1*time.Minute {
				m.Config.RefreshRate = 5 * time.Minute
			} else if m.Config.RefreshRate == 5*time.Minute {
				m.Config.RefreshRate = 15 * time.Minute
			} else {
				m.Config.RefreshRate = 1 * time.Minute
			}
			m.Config.Save()
		case 1:
			themes := []string{"default", "dark", "ocean"}
			for i, theme := range themes {
				if theme == m.Config.ColorTheme {
					m.Config.ColorTheme = themes[(i+1)%len(themes)]
					break
				}
			}
			m.Styles = NewStyles(m.Config.ColorTheme)
			m.Config.Save()
		case 2:
			m.Config.EnableNotifications = !m.Config.EnableNotifications
			m.Config.Save()
		}
	}
	return m, nil
}

// updateAddFeed handles feed addition
func (m *Model) updateAddFeed(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.State = StateManageFeeds
		return m, nil
	case "enter":
		if m.Input != "" {
			parts := strings.SplitN(m.Input, "|", 2)
			name := strings.TrimSpace(parts[0])
			url := name
			if len(parts) == 2 {
				url = strings.TrimSpace(parts[1])
			}

			m.Feeds.Feeds = append(m.Feeds.Feeds, rss.FeedInfo{Name: name, URL: url})
			if err := m.Feeds.Save(m.Config.FeedsFile); err != nil {
				m.Err = err
			}
			m.State = StateManageFeeds
			return m, FetchAllFeedsCmd(m.RSSClient, m.Feeds)
		}
	case "backspace":
		if len(m.Input) > 0 {
			m.Input = m.Input[:len(m.Input)-1]
		}
	default:
		m.Input += msg.String()
	}
	return m, nil
}

// updateRemoveFeed handles feed removal
func (m *Model) updateRemoveFeed(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.State = StateManageFeeds
		// Ensure selected index is valid
		if m.Selected >= len(m.Feeds.Feeds) && len(m.Feeds.Feeds) > 0 {
			m.Selected = len(m.Feeds.Feeds) - 1
		} else if len(m.Feeds.Feeds) == 0 {
			m.Selected = 0
		}
		return m, nil
	case "up", "k":
		if m.Selected > 0 {
			m.Selected--
		}
	case "down", "j":
		if m.Selected < len(m.Feeds.Feeds)-1 {
			m.Selected++
		}
	case "enter":
		if m.Selected < len(m.Feeds.Feeds) {
			// Remove feed using slices operations for better performance
			copy(m.Feeds.Feeds[m.Selected:], m.Feeds.Feeds[m.Selected+1:])
			m.Feeds.Feeds = m.Feeds.Feeds[:len(m.Feeds.Feeds)-1]
			
			if err := m.Feeds.Save(m.Config.FeedsFile); err != nil {
				m.Err = err
			}
			m.State = StateManageFeeds
			// Adjust selection after removal
			if m.Selected >= len(m.Feeds.Feeds) && len(m.Feeds.Feeds) > 0 {
				m.Selected = len(m.Feeds.Feeds) - 1
			} else if len(m.Feeds.Feeds) == 0 {
				m.Selected = 0
			}
			return m, FetchAllFeedsCmd(m.RSSClient, m.Feeds)
		}
	}
	return m, nil
}

// updateArticleView handles article content view
func (m *Model) updateArticleView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "q", "esc":
		m.State = StateFeedView
		return m, nil
	case "o":
		// 'o' for "open" - open URL in browser and return to feed view
		if len(m.Articles) > 0 && m.Selected < len(m.Articles) {
			return m, OpenURLCmd(m.Articles[m.Selected].Link)
		}
	}
	return m, nil
}

// GetSelectedArticle returns the currently selected article
func (m *Model) GetSelectedArticle() *rss.Article {
	if len(m.Articles) > 0 && m.Selected < len(m.Articles) {
		return &m.Articles[m.Selected]
	}
	return nil
}