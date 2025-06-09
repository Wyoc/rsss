package tui

import (
	"fmt"
	"strings"
)

// View renders the current state of the TUI
func (m *Model) View() string {
	switch m.State {
	case StateMenu:
		return m.viewMenu()
	case StateFeedView:
		return m.viewFeedView()
	case StateArticleView:
		return m.viewArticleView()
	case StateManageFeeds:
		return m.viewManageFeeds()
	case StateConfigure:
		return m.viewConfigure()
	case StateAddFeed:
		return m.viewAddFeed()
	case StateRemoveFeed:
		return m.viewRemoveFeed()
	default:
		return "Unknown state"
	}
}

// viewMenu renders the main menu
func (m *Model) viewMenu() string {
	var b strings.Builder

	b.WriteString(m.Styles.Header.Render("📡 RSS Reader"))
	b.WriteString("\n\n")

	menuItems := []string{"📰 See Feeds", "⚙️ Manage Feeds", "🎨 Configure"}

	for i, item := range menuItems {
		style := m.Styles.Normal
		if i == m.MenuSelected {
			style = m.Styles.Selected
		}
		b.WriteString(style.Render(fmt.Sprintf("  %s", item)))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter to select, q to quit"))

	return m.Styles.Menu.Render(b.String())
}

// viewFeedView renders the feed view with articles
func (m *Model) viewFeedView() string {
	var b strings.Builder

	// Note: Height calculation is now handled in getMaxVisibleArticles()
	
	// Compact header - just the essential info on one line
	headerInfo := fmt.Sprintf("📰 Latest Articles | Updated: %s", m.LastRefresh.Format("15:04:05"))
	
	if m.Err != nil {
		headerInfo += fmt.Sprintf(" | Error: %v", m.Err)
	}
	
	b.WriteString(m.Styles.Title.Render(headerInfo))
	b.WriteString("\n")

	// Handle special states
	if m.Loading {
		b.WriteString(m.Styles.Normal.Render("Loading feeds..."))
		b.WriteString("\n")
		b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter to read, 'r' to refresh, Esc to menu"))
		return b.String()
	}

	if len(m.Feeds.Feeds) == 0 {
		b.WriteString(m.Styles.Error.Render("No feeds configured! Go to 'Manage Feeds' to add RSS feeds first."))
		b.WriteString("\n")
		b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter to read, 'r' to refresh, Esc to menu"))
		return b.String()
	}

	if len(m.Articles) == 0 {
		b.WriteString(m.Styles.Normal.Render("No articles found. Press 'r' to refresh."))
		b.WriteString("\n")
		b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter to read, 'r' to refresh, Esc to menu"))
		return b.String()
	}

	// Calculate terminal width for responsive layout
	terminalWidth := m.Width
	if terminalWidth == 0 {
		terminalWidth = 80 // Fallback
	}

	// Display articles using available height with scrolling
	maxArticles := m.getMaxVisibleArticles()
	
	// Calculate visible range with scrolling
	startIdx := m.ViewportTop
	endIdx := min(len(m.Articles), startIdx+maxArticles)
	
	for i := range endIdx - startIdx {
		articleIdx := startIdx + i
		if articleIdx >= len(m.Articles) {
			break
		}
		article := m.Articles[articleIdx]
		
		style := m.Styles.Normal
		if articleIdx == m.Selected {
			style = m.Styles.Selected
		}

		// Format time and feed name with responsive width
		timeStr := article.PubDate.Format("15:04")
		feedName := article.FeedName
		
		// Adjust feed name width based on terminal size
		feedNameWidth := min(15, max(8, terminalWidth/6)) // 8-15 chars based on width
		if len(feedName) > feedNameWidth {
			feedName = feedName[:feedNameWidth-3] + "..."
		}
		
		// Create aligned columns: [TIME] [FEED_NAME] TITLE
		prefix := fmt.Sprintf("%s %-*s ", timeStr, feedNameWidth, feedName)
		
		// Calculate available space for title using actual terminal width
		titleMaxWidth := max(20, terminalWidth-len(prefix)-2) // 2 for margins
		
		title := article.Title
		if len(title) > titleMaxWidth {
			title = title[:titleMaxWidth-3] + "..."
		}

		line := fmt.Sprintf("%s%s", prefix, title)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	// Help text
	b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter to read, 'r' to refresh, Esc to menu"))

	return b.String()
}

// viewManageFeeds renders the feed management view
func (m *Model) viewManageFeeds() string {
	var b strings.Builder

	b.WriteString(m.Styles.Title.Render("⚙️ Manage Feeds"))
	b.WriteString("\n\n")

	if len(m.Feeds.Feeds) == 0 {
		b.WriteString(m.Styles.Normal.Render("No feeds configured."))
	} else {
		b.WriteString(m.Styles.Accent.Render("Current Feeds:"))
		b.WriteString("\n")
		for i, feed := range m.Feeds.Feeds {
			style := m.Styles.Normal
			if i == m.Selected {
				style = m.Styles.Selected
			}
			b.WriteString(style.Render(fmt.Sprintf("%d. %s (%s)", i+1, feed.Name, feed.URL)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter/d to delete selected, 'a' to add, Esc to menu"))

	return b.String()
}

// viewConfigure renders the configuration view
func (m *Model) viewConfigure() string {
	var b strings.Builder

	b.WriteString(m.Styles.Title.Render("🎨 Configuration"))
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Accent.Render("Settings:"))
	b.WriteString("\n")

	// Refresh Rate option
	style := m.Styles.Normal
	if m.Selected == 0 {
		style = m.Styles.Selected
	}
	b.WriteString(style.Render(fmt.Sprintf("  Refresh Rate: %v", m.Config.RefreshRate)))
	b.WriteString("\n")

	// Color Theme option
	style = m.Styles.Normal
	if m.Selected == 1 {
		style = m.Styles.Selected
	}
	b.WriteString(style.Render(fmt.Sprintf("  Color Theme: %s", m.Config.ColorTheme)))
	b.WriteString("\n")

	b.WriteString(m.Styles.Normal.Render(fmt.Sprintf("  Feeds File: %s", m.Config.FeedsFile)))
	b.WriteString("\n\n")
	b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to navigate, Enter/Space to change, Esc to menu"))

	return b.String()
}

// viewAddFeed renders the add feed view
func (m *Model) viewAddFeed() string {
	var b strings.Builder

	b.WriteString(m.Styles.Title.Render("➕ Add Feed"))
	b.WriteString("\n\n")

	b.WriteString(m.Styles.Normal.Render("Enter feed name|URL (or just URL):"))
	b.WriteString("\n")
	b.WriteString(m.Styles.Selected.Render(m.Input + "█"))
	b.WriteString("\n\n")
	b.WriteString(m.Styles.Normal.Render("Examples:"))
	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render("  BBC News|https://feeds.bbci.co.uk/news/rss.xml"))
	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render("  https://feeds.bbci.co.uk/news/rss.xml"))
	b.WriteString("\n\n")
	b.WriteString(m.Styles.Normal.Render("Press Enter to save, Esc to cancel"))

	return b.String()
}

// viewRemoveFeed renders the remove feed view
func (m *Model) viewRemoveFeed() string {
	var b strings.Builder

	b.WriteString(m.Styles.Title.Render("🗑️ Remove Feed"))
	b.WriteString("\n\n")

	if len(m.Feeds.Feeds) == 0 {
		b.WriteString(m.Styles.Normal.Render("No feeds to remove."))
	} else {
		b.WriteString(m.Styles.Normal.Render("Select feed to remove:"))
		b.WriteString("\n\n")

		for i, feed := range m.Feeds.Feeds {
			style := m.Styles.Normal
			if i == m.Selected {
				style = m.Styles.Selected
			}
			b.WriteString(style.Render(fmt.Sprintf("%d. %s (%s)", i+1, feed.Name, feed.URL)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render("Use ↑/↓ to select, Enter to remove, Esc to cancel"))

	return b.String()
}

// viewArticleView renders the article content view
func (m *Model) viewArticleView() string {
	var b strings.Builder

	if len(m.Articles) == 0 || m.Selected >= len(m.Articles) {
		b.WriteString(m.Styles.Error.Render("No article selected"))
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Normal.Render("Press Esc to return to feed list"))
		return b.String()
	}

	article := m.Articles[m.Selected]
	
	// Get terminal width for responsive layout
	terminalWidth := m.Width
	if terminalWidth == 0 {
		terminalWidth = 80 // Fallback
	}

	// Article header with title (word-wrapped if needed)
	title := article.Title
	if len(title) > terminalWidth-4 { // Account for emoji and padding
		title = m.wrapText(title, terminalWidth-4)
	}
	b.WriteString(m.Styles.Title.Render("📖 " + title))
	b.WriteString("\n\n")

	// Article metadata - compact for mobile, expanded for wider screens
	timeStr := article.PubDate.Format("15:04 on 2006-01-02")
	if terminalWidth < 60 {
		// Compact layout for narrow screens
		b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("🕒 %s", timeStr)))
		b.WriteString("\n")
		b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("📰 %s", article.FeedName)))
	} else {
		// Full layout for wider screens
		b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("🕒 %s | 📰 %s", timeStr, article.FeedName)))
	}
	b.WriteString("\n")
	
	// URL - wrap if too long
	url := article.Link
	if len(url) > terminalWidth-4 {
		url = url[:terminalWidth-7] + "..."
	}
	b.WriteString(m.Styles.Normal.Render(fmt.Sprintf("🔗 %s", url)))
	b.WriteString("\n\n")

	// Article content
	if article.Description != "" {
		// Clean up HTML tags and decode entities for better readability
		content := article.Description
		
		// Basic HTML tag removal (simple approach)
		content = strings.ReplaceAll(content, "<br>", "\n")
		content = strings.ReplaceAll(content, "<br/>", "\n")
		content = strings.ReplaceAll(content, "<br />", "\n")
		content = strings.ReplaceAll(content, "<p>", "\n")
		content = strings.ReplaceAll(content, "</p>", "\n")
		
		// Remove remaining HTML tags (basic regex replacement)
		for strings.Contains(content, "<") && strings.Contains(content, ">") {
			start := strings.Index(content, "<")
			end := strings.Index(content[start:], ">")
			if end != -1 {
				content = content[:start] + content[start+end+1:]
			} else {
				break
			}
		}
		
		// Basic HTML entity decoding
		content = strings.ReplaceAll(content, "&#8217;", "'")
		content = strings.ReplaceAll(content, "&#8220;", "\"")
		content = strings.ReplaceAll(content, "&#8221;", "\"")
		content = strings.ReplaceAll(content, "&amp;", "&")
		content = strings.ReplaceAll(content, "&lt;", "<")
		content = strings.ReplaceAll(content, "&gt;", ">")
		content = strings.ReplaceAll(content, "&quot;", "\"")
		content = strings.ReplaceAll(content, "&#160;", " ")
		content = strings.ReplaceAll(content, "&nbsp;", " ")
		
		// Clean up excessive whitespace
		content = strings.TrimSpace(content)
		lines := strings.Split(content, "\n")
		var cleanLines []string
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				cleanLines = append(cleanLines, line)
			}
		}
		content = strings.Join(cleanLines, "\n\n")
		
		// Wrap content to terminal width
		wrappedContent := m.wrapText(content, terminalWidth-4) // Leave margin
		b.WriteString(m.Styles.Normal.Render(wrappedContent))
	} else {
		b.WriteString(m.Styles.Normal.Render("No content available for this article."))
	}

	b.WriteString("\n\n")
	b.WriteString(m.Styles.Normal.Render("Press 'o' to open in browser, Esc to return to feed list"))

	return b.String()
}