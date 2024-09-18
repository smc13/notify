package notify

import (
	"context"
)

type Notifiable interface {
	// NotifiableID returns the unique identifier of the notifiable, can be used to identify the notifiable for logging and debugging
	NotifiableID() string
	// ShouldSendNotification returns a boolean indicating if the notification should be sent to the given channel
	ShouldSendNotification(ctx context.Context, channel Channel, notif Notification) bool
}
