package notify_test

import (
	"fmt"
	"testing"
	"time"

	_ "github.com/Cai-ki/cage/localconfig"
	"github.com/Cai-ki/cage/notify"
)

func TestSend(t *testing.T) {
	now := time.Now().Format("2006-01-02 15:04:05 MST")
	subject := "Test Notification"
	body := fmt.Sprintf("This is a test notification from notify.Send()!\nSent at: %s", now)

	if err := notify.Send(subject, body); err != nil {
		t.Fatal("Failed to send notification:", err)
	}
}
