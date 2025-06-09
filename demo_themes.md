# RSS Reader UI Improvements & Notifications

## Enhanced Features

### 🔔 **New Article Notifications**
- **Smart Detection**: Automatically detects when new articles are published
- **Persistent Tracking**: Remembers which articles you've seen across app restarts  
- **Configurable**: Enable/disable notifications in the settings
- **Dual Notifications**: Shows both in-app banners AND system notifications
- **Cross-Platform**: Uses `notify-send` on Linux, falls back gracefully on other systems
- **Easy Dismissal**: Press Space, Enter, or 'n' to dismiss in-app notifications

### 🎨 New Color Themes
- **default**: Modern indigo and pink theme with clean styling
- **dark**: Blue and pink theme optimized for dark environments  
- **ocean**: Sky blue and teal theme with oceanic colors
- **sunset**: Warm amber and orange theme with sunset colors
- **forest**: Green and lime theme with natural earth tones

### 🔧 Visual Improvements

#### Enhanced Menu
- Header with subtitle showing current theme and feed count
- Descriptive menu items with context on selection
- Improved visual hierarchy with borders and spacing

#### Article Feed View
- Enhanced status bar with refresh info and health status
- Color-coded time stamps and feed names
- Visual selection indicators (▶)
- Responsive layout adapting to terminal width
- Scroll indicators for long lists

#### Feed Management
- Empty state guidance for new users  
- Enhanced feed display with URLs on selection
- Action buttons with clear labels
- Visual indicators for selected items

#### Configuration View
- Theme preview when selecting color themes
- Visual checkmarks for current selections
- Enhanced option descriptions
- **Notification toggle** for enabling/disabling alerts
- Responsive settings layout

#### Article Reading View
- Improved header with visual separation
- Enhanced metadata display with responsive layout
- Better content formatting and HTML cleanup
- Visual separators and improved spacing
- Action buttons with clear navigation

### 🎯 UX Enhancements
- Consistent visual indicators throughout the app
- Responsive design that adapts to terminal size
- Enhanced help text with contextual information
- Better error states and empty state handling
- Improved spacing and visual hierarchy
- More intuitive navigation cues

## Usage

The improved UI automatically applies when running the TUI mode:

```bash
# Build the enhanced version
make build

# Run in TUI mode to see improvements
make run-tui

# Or run directly
./build/rsss --menu
```

### 📁 **File Structure**
The notification system creates these files in `~/.config/rsss/`:
- `config.json` - Application settings including notification preferences
- `feeds.json` - Your RSS feed subscriptions  
- `seen.json` - Tracks which articles you've already seen

### 🔔 **How Notifications Work**
1. **First Run**: All current articles are marked as "seen" without notification
2. **Subsequent Runs**: New articles trigger both system and in-app notifications
3. **System Notifications**: Appears in your desktop notification center (Linux with `notify-send`)
4. **In-App Notifications**: Shows a banner within the TUI interface
5. **Persistence**: Article tracking survives app restarts
6. **Smart Refresh**: Only shows notifications for truly new content
7. **Requirements**: On Linux, ensure `libnotify` is installed (`pacman -S libnotify` on Arch)

All existing functionality remains the same, but with a much more polished interface and intelligent new article notifications!