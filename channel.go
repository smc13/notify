package notify

import (
	"context"
)

type Channel interface {
	//  Notify sends a notification to this channel
	Notify(context.Context, []string, Notification) error
}
