package rss

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"sort"
	"time"
)

// Client handles RSS feed fetching and parsing
type Client struct {
	httpClient *http.Client
}

// NewClient creates a new RSS client with the specified timeout
func NewClient(timeout time.Duration) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// FetchFeed fetches and parses an RSS feed from the given URL
func (c *Client) FetchFeed(url string) (*RSS, error) {
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch RSS feed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var rss RSS
	if err := xml.Unmarshal(body, &rss); err != nil {
		return nil, fmt.Errorf("failed to parse RSS: %w", err)
	}

	return &rss, nil
}

// FetchMultipleFeeds fetches multiple RSS feeds and returns all articles sorted by date
func (c *Client) FetchMultipleFeeds(feeds []FeedInfo) ([]Article, error) {
	var allArticles []Article

	for _, feed := range feeds {
		rss, err := c.FetchFeed(feed.URL)
		if err != nil {
			// Continue with other feeds if one fails
			continue
		}

		for _, item := range rss.Channel.Items {
			allArticles = append(allArticles, Article{
				Title:       item.Title,
				Link:        item.Link,
				Description: item.Description,
				PubDate:     parseTime(item.PubDate),
				FeedName:    feed.Name,
			})
		}
	}

	// Sort by publication date (newest first)
	sort.Slice(allArticles, func(i, j int) bool {
		return allArticles[i].PubDate.After(allArticles[j].PubDate)
	})

	return allArticles, nil
}

// FeedInfo represents RSS feed configuration
type FeedInfo struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// parseTime attempts to parse various date formats commonly used in RSS feeds
func parseTime(pubDate string) time.Time {
	layouts := []string{
		time.RFC1123Z,
		time.RFC1123,
		"Mon, 2 Jan 2006 15:04:05 -0700",
		"Mon, 2 Jan 2006 15:04:05 MST",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05Z",
	}

	for _, layout := range layouts {
		if t, err := time.Parse(layout, pubDate); err == nil {
			return t
		}
	}
	return time.Now()
}