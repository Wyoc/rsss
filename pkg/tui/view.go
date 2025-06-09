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

	b.WriteString(m.Styles.Title.Render("📰 Latest Articles"))
	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render(fmt.Sprintf("Last updated: %s", m.LastRefresh.Format("15:04:05"))))
	b.WriteString("\n\n")

	if m.Loading {
		b.WriteString(m.Styles.Normal.Render("Loading feeds..."))
		return b.String()
	}

	if m.Err != nil {
		b.WriteString(m.Styles.Error.Render(fmt.Sprintf("Error: %v", m.Err)))
		b.WriteString("\n")
	}

	if len(m.Feeds.Feeds) == 0 {
		b.WriteString(m.Styles.Error.Render("No feeds configured!"))
		b.WriteString("\n")
		b.WriteString(m.Styles.Normal.Render("Go to 'Manage Feeds' to add some RSS feeds first."))
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Normal.Render("Press Esc to return to menu"))
		return b.String()
	}

	if len(m.Articles) == 0 {
		b.WriteString(m.Styles.Normal.Render("No articles found. Press 'r' to refresh."))
		return b.String()
	}

	for i, article := range m.Articles {
		if i >= 20 {
			break
		}

		style := m.Styles.Normal
		if i == m.Selected {
			style = m.Styles.Selected
		}

		// Format time and feed name with consistent width
		timeStr := article.PubDate.Format("15:04")
		feedName := article.FeedName
		if len(feedName) > 15 {
			feedName = feedName[:12] + "..."
		}
		
		// Create aligned columns: [TIME] [FEED_NAME] TITLE
		// Time: 5 chars, Feed: 15 chars padded
		prefix := fmt.Sprintf("%s %-15s ", timeStr, feedName)
		
		// Calculate available space for title (assuming 80 char terminal width)
		titleMaxWidth := max(20, 80-len(prefix)-2) // 2 for margins
		
		title := article.Title
		if len(title) > titleMaxWidth {
			title = title[:titleMaxWidth-3] + "..."
		}

		line := fmt.Sprintf("%s%s", prefix, title)
		b.WriteString(style.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
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

	// Article header with title
	b.WriteString(m.Styles.Title.Render("📖 " + article.Title))
	b.WriteString("\n\n")

	// Article metadata
	timeStr := article.PubDate.Format("15:04 on 2006-01-02")
	b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("🕒 %s | 📰 %s", timeStr, article.FeedName)))
	b.WriteString("\n")
	b.WriteString(m.Styles.Normal.Render(fmt.Sprintf("🔗 %s", article.Link)))
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
		
		b.WriteString(m.Styles.Normal.Render(content))
	} else {
		b.WriteString(m.Styles.Normal.Render("No content available for this article."))
	}

	b.WriteString("\n\n")
	b.WriteString(m.Styles.Normal.Render("Press 'o' to open in browser, Esc to return to feed list"))

	return b.String()
}