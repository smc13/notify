package ably

import (
	"context"

	ably_go "github.com/ably/ably-go/ably"
	"github.com/smc13/notify"
)

type AblyChannel struct {
	client *ably_go.REST
}

type AblyNotification interface {
	ToAblyData() (any, error)
}

func New(client *ably_go.REST) *AblyChannel {
	return &AblyChannel{client: client}
}

func (a *AblyChannel) Notify(ctx context.Context, channelIds []string, notif notify.Notification) error {
	message, err := a.notificationToAblyMessage(notif)
	if err != nil {
		return err
	}

	for _, channelID := range channelIds {
		if err := a.client.Channels.Get(channelID).Publish(ctx, notif.Kind(), message); err != nil {
			return err
		}
	}

	return nil
}

type simpleAlbyData struct {
	Subject      string `json:"subject"`
	Content      string `json:"content"`
	Notification any    `json:"notification"`
}

func (a *AblyChannel) notificationToAblyMessage(notif notify.Notification) (any, error) {
	if ablyNotif, ok := notif.(AblyNotification); ok {
		return ablyNotif.ToAblyData()
	}

	return &simpleAlbyData{
		Subject:      notif.Subject(),
		Content:      notif.Content(),
		Notification: notif,
	}, nil
}
