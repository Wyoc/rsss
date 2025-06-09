# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a full-featured RSS client written in Go that supports both CLI and TUI modes. The application includes feed management, configuration, and a modern terminal user interface built with Bubble Tea.

## Project Structure

- `cmd/rsss/` - Main application entry point
- `pkg/rss/` - RSS feed parsing and fetching logic
- `pkg/config/` - Configuration management
- `pkg/tui/` - Terminal user interface components
- `build/` - Build artifacts (ignored by git)

## Common Commands

**Using Makefile (recommended):**
```bash
make build          # Build the binary
make test           # Run tests
make test-coverage  # Run tests with coverage
make run-cli        # Run CLI mode with BBC News
make run-tui        # Run TUI menu mode
make clean          # Clean build artifacts
make help           # Show all available targets
```

**Manual commands:**
```bash
# Build
go build -o build/rsss ./cmd/rsss

# Test
go test -v ./...

# Run CLI mode
./build/rsss <RSS_URL>

# Run TUI mode
./build/rsss --menu
./build/rsss --tui <RSS_URL>
```

**Test with a known working feed:**
```bash
./build/rsss https://feeds.bbci.co.uk/news/rss.xml
```

## Architecture

**Package Structure:**

1. **`pkg/rss`** - RSS feed handling
   - `Client` - HTTP client with configurable timeout
   - `RSS`, `Channel`, `Item` - XML unmarshaling structs
   - `Article` - Processed article with parsed date
   - `FetchFeed()` - Single feed fetching
   - `FetchMultipleFeeds()` - Multi-feed fetching with sorting

2. **`pkg/config`** - Configuration management
   - `Config` - Application settings (refresh rate, theme, file paths)
   - `FeedConfig` - Feed list management
   - Load/Save functionality with JSON persistence
   - Default configuration creation

3. **`pkg/tui`** - Terminal user interface
   - `Model` - Bubble Tea application state
   - `Styles` - Lip Gloss styling with theme support
   - State machine with menu navigation
   - Async command handling for feed fetching

4. **`cmd/rsss`** - Main application
   - CLI argument parsing
   - Mode selection (CLI vs TUI)
   - Application initialization

**Application Flow:**

**CLI Mode:**
1. Parse URL from command line
2. Create RSS client with timeout
3. Fetch and parse single feed
4. Display articles in simple text format

**TUI Mode:**
1. Load configuration and feeds from JSON files
2. Initialize Bubble Tea model with RSS client
3. Start with main menu (See Feeds, Manage Feeds, Configure)
4. Handle navigation and state transitions
5. Auto-refresh feeds based on configured interval

**Key Features:**
- Configurable refresh intervals (1, 5, 15 minutes)
- Multiple color themes (default, dark, ocean)
- Persistent feed management with JSON storage
- Chronological article sorting across all feeds
- Error handling for network and parsing failures

## Project Steps

### Phase 1 - Implement TUI
[x] use bubble tea to create a nice TUI (https://github.com/charmbracelet/bubbletea)
[x] setup a TUI with 3 elements in the menu: see feed, manage feed and configure
[x] the see feed will check all the feeds every N minutes, and display them chronologically. We should be able to read the content by pressing Enter of a title
[x] The manage feed should allow us to add or remove feeds (the feeds should be saved in a json file)
[x] The configure should allow us to select the refresh rate, the location of the json file and the color theme



## Reminder

When you execute a step, mark it as complete