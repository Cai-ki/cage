// Package notify provides a simple way to send notifications.
package notify

import (
	"fmt"
	"sync"
)

// Notifier is the interface that wraps the basic Send method.
type Notifier interface {
	Send(subject, message string) error
}

var (
	defaultNotifier Notifier
	once            sync.Once
	initErr         error
)

// initDefault initializes the default notifier (currently only email).
func initDefault() {
	defaultNotifier, initErr = NewEmailNotifier()
}

// Send sends a notification using the default channel (e.g., email).
// It reads configuration from environment variables.
func Send(subject, body string) error {
	once.Do(initDefault)
	if initErr != nil {
		return fmt.Errorf("failed to initialize notifier: %w", initErr)
	}
	return defaultNotifier.Send(subject, body)
}
