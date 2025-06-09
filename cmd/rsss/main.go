package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"rsss/pkg/config"
	"rsss/pkg/rss"
	"rsss/pkg/tui"
)

func main() {
	if len(os.Args) == 1 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "--tui", "--menu":
		url := ""
		if len(os.Args) >= 3 {
			url = os.Args[2]
		}
		if err := runTUI(url); err != nil {
			fmt.Printf("Error running TUI: %v\n", err)
			os.Exit(1)
		}
	default:
		url := os.Args[1]
		if err := runCLI(url); err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
	}
}

func printUsage() {
	fmt.Println("Usage: rsss <RSS_URL>")
	fmt.Println("       rsss --tui [RSS_URL]")
	fmt.Println("       rsss --menu")
	fmt.Println("Example: rsss https://feeds.feedburner.com/oreilly/radar")
	fmt.Println("         rsss --tui https://feeds.bbci.co.uk/news/rss.xml")
	fmt.Println("         rsss --menu")
}

func runTUI(url string) error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	feeds, err := config.LoadFeeds(cfg.FeedsFile)
	if err != nil {
		return err
	}

	if url != "" {
		feeds.Feeds = []rss.FeedInfo{{Name: "Command Line Feed", URL: url}}
	}

	rssClient := rss.NewClient(10 * time.Second)
	model := tui.NewModel(cfg, feeds, rssClient)

	p := tea.NewProgram(model)

	// Run the TUI
	_, err = p.Run()
	if err != nil {
		return err
	}

	return nil
}

func runCLI(url string) error {
	fmt.Printf("Fetching RSS feed from: %s\n\n", url)

	client := rss.NewClient(10 * time.Second)
	feed, err := client.FetchFeed(url)
	if err != nil {
		return err
	}

	displayFeed(feed)
	return nil
}

func displayFeed(feed *rss.RSS) {
	fmt.Printf("Feed: %s\n", feed.Channel.Title)
	fmt.Printf("Description: %s\n", feed.Channel.Description)
	fmt.Printf("Link: %s\n\n", feed.Channel.Link)

	for i, item := range feed.Channel.Items {
		if i >= 10 {
			break
		}
		fmt.Printf("Title: %s\n", item.Title)
		fmt.Printf("Link: %s\n", item.Link)
		fmt.Printf("Date: %s\n", item.PubDate)
		fmt.Printf("Description: %s\n\n", item.Description)
	}
}
