package notify

import (
	"context"
)

type Notifiable interface {
	// ShouldSendNotification returns a boolean indicating if the notification should be sent to the given channel
	ShouldSendNotification(ctx context.Context, channel Channel, notif Notification) bool
}
