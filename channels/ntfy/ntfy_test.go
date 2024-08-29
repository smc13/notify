package ntfy

import (
	"context"
	"testing"
	"time"

	"github.com/smc13/notify"
	"github.com/stretchr/testify/assert"
)

type user struct{}

func (u *user) RouteNotificationFor(ctx context.Context, channel notify.Channel, notif notify.Notification) ([]string, error) {
	return []string{}, nil
}

type notification struct{}

func (n *notification) Kind() string    { return "notification" }
func (n *notification) Subject() string { return "Test notification" }
func (n *notification) Content() string { return "This is a test notification" }
func (n *notification) ToNtfy() (Message, error) {
	return Message{
		Title:    "Hello world!!!!",
		Message:  "This is a more detailed message",
		Tags:     []string{"test", "notification"},
		Priority: MessagePriorityHigh,
		Delay:    MessageDelay(time.Now().Add(10 * time.Second)),
	}, nil
}

func TestNtfyChannel(t *testing.T) {
	ntfy := New("http://127.0.0.1:8045", "test")
	service := notify.NewNotify(ntfy)

	err := service.Notify(context.Background(), &user{}, &notification{})
	assert.NoError(t, err)
}
