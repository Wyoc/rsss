package tui

import (
	"testing"
	"time"

	"rsss/pkg/config"
	"rsss/pkg/rss"
)

func TestNewModel(t *testing.T) {
	cfg := config.DefaultConfig()
	feeds := &config.FeedConfig{
		Feeds: []rss.FeedInfo{
			{Name: "Test Feed", URL: "https://example.com/feed.xml"},
		},
	}
	rssClient := rss.NewClient(5 * time.Second)

	model := NewModel(cfg, feeds, rssClient)

	if model.State != StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", model.State)
	}

	if model.MenuSelected != 0 {
		t.Errorf("Expected MenuSelected to be 0, got %d", model.MenuSelected)
	}

	if model.Selected != 0 {
		t.Errorf("Expected Selected to be 0, got %d", model.Selected)
	}

	if model.Feeds != feeds {
		t.Error("Expected feeds to be set correctly")
	}

	if model.Config != cfg {
		t.Error("Expected config to be set correctly")
	}

	if model.RSSClient != rssClient {
		t.Error("Expected RSS client to be set correctly")
	}

	if !model.Loading {
		t.Error("Expected model to be in loading state initially")
	}
}

func TestGetSelectedArticle(t *testing.T) {
	cfg := config.DefaultConfig()
	feeds := &config.FeedConfig{}
	rssClient := rss.NewClient(5 * time.Second)
	model := NewModel(cfg, feeds, rssClient)

	// Test with no articles
	article := model.GetSelectedArticle()
	if article != nil {
		t.Error("Expected nil article when no articles are present")
	}

	// Add some test articles
	model.Articles = []rss.Article{
		{
			Title:    "Test Article 1",
			Link:     "https://example.com/1",
			FeedName: "Test Feed",
			PubDate:  time.Now(),
		},
		{
			Title:    "Test Article 2",
			Link:     "https://example.com/2",
			FeedName: "Test Feed",
			PubDate:  time.Now(),
		},
	}

	// Test with first article selected
	model.Selected = 0
	article = model.GetSelectedArticle()
	if article == nil {
		t.Fatal("Expected article to be returned")
	}
	if article.Title != "Test Article 1" {
		t.Errorf("Expected 'Test Article 1', got '%s'", article.Title)
	}

	// Test with second article selected
	model.Selected = 1
	article = model.GetSelectedArticle()
	if article == nil {
		t.Fatal("Expected article to be returned")
	}
	if article.Title != "Test Article 2" {
		t.Errorf("Expected 'Test Article 2', got '%s'", article.Title)
	}

	// Test with out-of-bounds selection
	model.Selected = 10
	article = model.GetSelectedArticle()
	if article != nil {
		t.Error("Expected nil article for out-of-bounds selection")
	}
}

func TestStateTransitions(t *testing.T) {
	cfg := config.DefaultConfig()
	feeds := &config.FeedConfig{}
	rssClient := rss.NewClient(5 * time.Second)
	model := NewModel(cfg, feeds, rssClient)

	// Test initial state
	if model.State != StateMenu {
		t.Errorf("Expected initial state StateMenu, got %v", model.State)
	}

	// Test state transitions
	testCases := []struct {
		from AppState
		to   AppState
	}{
		{StateMenu, StateFeedView},
		{StateMenu, StateManageFeeds},
		{StateMenu, StateConfigure},
		{StateManageFeeds, StateAddFeed},
		{StateManageFeeds, StateRemoveFeed},
	}

	for _, tc := range testCases {
		model.State = tc.from
		model.State = tc.to
		if model.State != tc.to {
			t.Errorf("State transition from %v to %v failed", tc.from, tc.to)
		}
	}
}