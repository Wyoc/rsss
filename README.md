# RSS Reader

A modern RSS client written in Go with both CLI and TUI modes.

## Features

- 📰 **Dual Interface**: Command-line and interactive terminal UI
- 🔄 **Auto-refresh**: Configurable refresh intervals (1, 5, 15 minutes)
- 🎨 **Themes**: Multiple color themes (default, dark, ocean)
- 📱 **Feed Management**: Add, remove, and organize RSS feeds
- ⚡ **Fast**: Concurrent feed fetching with proper error handling
- 💾 **Persistent**: Configuration and feeds saved as JSON

## Installation

```bash
# Clone the repository
git clone https://github.com/wyoc/rsss.git
cd rsss

# Build the binary
make build

# Or install directly
make install
```

## Usage

### CLI Mode (Quick feed reading)

```bash
# Read a single RSS feed
./build/rsss https://feeds.bbci.co.uk/news/rss.xml
```

### TUI Mode (Interactive interface)

```bash
# Launch interactive menu
./build/rsss --menu

# Or directly with a feed
./build/rsss --tui https://feeds.bbci.co.uk/news/rss.xml
```

### TUI Navigation

- **Main Menu**: Use ↑/↓ to navigate, Enter to select
- **Feed View**: Navigate articles with ↑/↓, Enter to view link, 'r' to refresh
- **Manage Feeds**: 'a' to add, 'd' to delete feeds
- **Configure**: Use ↑/↓ to select setting, Enter/Space to change
- **Universal**: Esc to go back, 'q' to quit

## Development

```bash
# Run tests
make test

# Run with coverage
make test-coverage

# Format code
make fmt

# Run linter
make lint

# Clean build artifacts
make clean
```

## Project Structure

```
├── cmd/rsss/           # Main application
├── pkg/
│   ├── config/         # Configuration management
│   ├── rss/           # RSS parsing and fetching
│   └── tui/           # Terminal user interface
├── build/             # Build artifacts
├── Makefile          # Build automation
└── README.md         # This file
```

## Configuration

Configuration files are automatically created in `~/.config/rsss/`:

- `config.json` - Application settings
- `feeds.json` - RSS feed list

## Default Feeds

On first run, the application creates default feeds:
- BBC News
- TechCrunch

## License

MIT License