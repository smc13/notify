package notify

import (
	"context"
)

type Channel interface {
	// Name returns the name of this channel and should be unique. Useful for logging and debugging
	Name() string
	//  Notify sends a notification to this channel
	Notify(context.Context, Notifiable, Notification) error
}
