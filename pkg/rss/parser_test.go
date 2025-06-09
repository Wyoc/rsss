package rss

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	timeout := 5 * time.Second
	client := NewClient(timeout)

	if client == nil {
		t.Fatal("NewClient returned nil")
	}

	if client.httpClient.Timeout != timeout {
		t.Errorf("Expected timeout %v, got %v", timeout, client.httpClient.Timeout)
	}
}

func TestFetchFeed(t *testing.T) {
	// Create a test RSS feed
	testRSS := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Test Feed</title>
		<link>https://example.com</link>
		<description>A test RSS feed</description>
		<item>
			<title>Test Article</title>
			<link>https://example.com/article1</link>
			<description>This is a test article</description>
			<pubDate>Mon, 01 Jan 2024 12:00:00 GMT</pubDate>
		</item>
	</channel>
</rss>`

	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testRSS))
	}))
	defer server.Close()

	client := NewClient(5 * time.Second)
	feed, err := client.FetchFeed(server.URL)

	if err != nil {
		t.Fatalf("FetchFeed returned error: %v", err)
	}

	if feed.Channel.Title != "Test Feed" {
		t.Errorf("Expected title 'Test Feed', got '%s'", feed.Channel.Title)
	}

	if len(feed.Channel.Items) != 1 {
		t.Errorf("Expected 1 item, got %d", len(feed.Channel.Items))
	}

	item := feed.Channel.Items[0]
	if item.Title != "Test Article" {
		t.Errorf("Expected item title 'Test Article', got '%s'", item.Title)
	}
}

func TestFetchFeedHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := NewClient(5 * time.Second)
	_, err := client.FetchFeed(server.URL)

	if err == nil {
		t.Fatal("Expected error for HTTP 404, got nil")
	}
}

func TestFetchMultipleFeeds(t *testing.T) {
	testRSS1 := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Feed 1</title>
		<link>https://example1.com</link>
		<description>First test feed</description>
		<item>
			<title>Article 1</title>
			<link>https://example1.com/article1</link>
			<description>First article</description>
			<pubDate>Mon, 01 Jan 2024 12:00:00 GMT</pubDate>
		</item>
	</channel>
</rss>`

	testRSS2 := `<?xml version="1.0" encoding="UTF-8"?>
<rss version="2.0">
	<channel>
		<title>Feed 2</title>
		<link>https://example2.com</link>
		<description>Second test feed</description>
		<item>
			<title>Article 2</title>
			<link>https://example2.com/article1</link>
			<description>Second article</description>
			<pubDate>Tue, 02 Jan 2024 12:00:00 GMT</pubDate>
		</item>
	</channel>
</rss>`

	server1 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testRSS1))
	}))
	defer server1.Close()

	server2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(testRSS2))
	}))
	defer server2.Close()

	feeds := []FeedInfo{
		{Name: "Test Feed 1", URL: server1.URL},
		{Name: "Test Feed 2", URL: server2.URL},
	}

	client := NewClient(5 * time.Second)
	articles, err := client.FetchMultipleFeeds(feeds)

	if err != nil {
		t.Fatalf("FetchMultipleFeeds returned error: %v", err)
	}

	if len(articles) != 2 {
		t.Errorf("Expected 2 articles, got %d", len(articles))
	}

	// Check if articles are sorted by date (newest first)
	if !articles[0].PubDate.After(articles[1].PubDate) {
		t.Error("Articles are not sorted by date (newest first)")
	}
}

func TestParseTime(t *testing.T) {
	testCases := []struct {
		input    string
		expected bool // whether parsing should succeed
	}{
		{"Mon, 01 Jan 2024 12:00:00 GMT", true},
		{"Mon, 01 Jan 2024 12:00:00 -0700", true},
		{"2024-01-01T12:00:00Z", true},
		{"invalid date", false}, // This will return time.Now(), so we can't test exact value
	}

	for _, tc := range testCases {
		result := parseTime(tc.input)
		if tc.expected && result.IsZero() {
			t.Errorf("parseTime(%s) returned zero time, expected valid time", tc.input)
		}
	}
}