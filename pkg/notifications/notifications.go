package notifications

import (
	"fmt"
	"os/exec"
	"runtime"
)

// Notifier interface for different notification systems
type Notifier interface {
	Send(title, message string) error
}

// LinuxNotifier uses notify-send for Linux systems
type LinuxNotifier struct{}

// Send sends a notification using notify-send
func (n *LinuxNotifier) Send(title, message string) error {
	cmd := exec.Command("notify-send", 
		"--app-name=RSSS",
		"--icon=rss",
		"--urgency=normal",
		title, 
		message)
	return cmd.Run()
}

// NoOpNotifier does nothing (fallback for unsupported systems)
type NoOpNotifier struct{}

func (n *NoOpNotifier) Send(title, message string) error {
	return nil
}

// NewNotifier creates a new notifier based on the operating system
func NewNotifier() Notifier {
	switch runtime.GOOS {
	case "linux":
		// Check if notify-send is available
		if _, err := exec.LookPath("notify-send"); err == nil {
			return &LinuxNotifier{}
		}
		fallthrough
	default:
		return &NoOpNotifier{}
	}
}

// SendArticleNotification sends a notification about new articles
func SendArticleNotification(notifier Notifier, count int) error {
	var title, message string
	
	if count == 1 {
		title = "RSSS - New Article"
		message = "1 new article available!"
	} else {
		title = "RSSS - New Articles"
		message = fmt.Sprintf("%d new articles available!", count)
	}
	
	return notifier.Send(title, message)
}