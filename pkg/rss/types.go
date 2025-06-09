package rss

import (
	"encoding/xml"
	"time"
)

// RSS represents the root RSS element
type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

// Channel represents an RSS channel
type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

// Item represents an RSS item/article
type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

// Article represents a processed RSS article with parsed date
type Article struct {
	Title       string
	Link        string
	Description string
	PubDate     time.Time
	FeedName    string
}