package notify

import (
	"context"
	"errors"
)

type Notifiable interface {
	// RouteNotificationFor returns an slice of ids that the notification should be routed to (if any)
	// The format of the ids will depend on the channel (eg: discord channel ids, ntfy topic names, etc)
	// If the notification should not be sent to the channel, it should return ErrSkipNotification
	RouteNotificationFor(ctx context.Context, channel Channel, notif Notification) ([]string, error)
}

type GroupedNotifiable struct {
	notifiables []Notifiable
}

func NewGroupedNotifiable(notifiables ...Notifiable) *GroupedNotifiable {
	return &GroupedNotifiable{notifiables: notifiables}
}

func (g *GroupedNotifiable) AddNotifiable(notifiable Notifiable) {
	g.notifiables = append(g.notifiables, notifiable)
}

func (g *GroupedNotifiable) RouteNotificationFor(ctx context.Context, channel Channel, notif Notification) ([]string, error) {
	var ids []string
	for _, notifiable := range g.notifiables {
		routedIds, err := notifiable.RouteNotificationFor(ctx, channel, notif)
		if err != nil {
			if errors.Is(err, ErrSkipNotification) {
				continue
			}

			return nil, err
		}

		ids = append(ids, routedIds...)
	}

	return ids, nil
}
