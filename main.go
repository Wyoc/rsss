package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func fetchRSS(url string) (*RSS, error) {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Get(url)
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

func displayFeed(rss *RSS) {
	fmt.Printf("Feed: %s\n", rss.Channel.Title)
	fmt.Printf("Description: %s\n", rss.Channel.Description)
	fmt.Printf("Link: %s\n\n", rss.Channel.Link)

	for i, item := range rss.Channel.Items {
		if i >= 10 {
			break
		}
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Date: %s\n", item.PubDate)
		fmt.Printf("Description: %s\n\n", item.Description)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: rsss <RSS_URL>")
		fmt.Println("Example: rsss https://feeds.feedburner.com/oreilly/radar")
		os.Exit(1)
	}

	url := os.Args[1]
	
	fmt.Printf("Fetching RSS feed from: %s\n\n", url)

	rss, err := fetchRSS(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	displayFeed(rss)
}