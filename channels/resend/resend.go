package resend

import (
	"context"

	"github.com/resend/resend-go/v2"
	"github.com/smc13/notify"
)

type ResendNotification interface {
	ToResendMessage(notify.Notifiable) (*resend.SendEmailRequest, error)
}

type ResendNotifiable interface {
	ResendAddress() string
}

type ResendChannel struct {
	client *resend.Client
	from   string
}

type ResendOption func(*ResendChannel)

func New(apiKey string, opts ...ResendOption) *ResendChannel {
	return NewFromClient(resend.NewClient(apiKey), opts...)
}

func NewFromClient(client *resend.Client, opts ...ResendOption) *ResendChannel {
	c := &ResendChannel{client: client}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func WithFrom(from string) ResendOption {
	return func(c *ResendChannel) {
		c.from = from
	}
}

func (c *ResendChannel) Name() string { return "resend" }

func (c *ResendChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notif notify.Notification) error {
	message, err := c.notificationToResendMessage(notifiable, notif)
	if err != nil {
		return err
	}

	if message == nil {
		return nil
	}

	if c.from != "" && message.From == "" {
		message.From = c.from
	}

	if notifiable, ok := notifiable.(ResendNotifiable); ok {
		message.To = append(message.To, notifiable.ResendAddress())
	}

	_, err = c.client.Emails.SendWithContext(ctx, message)
	return err
}

func (c *ResendChannel) notificationToResendMessage(notifiable notify.Notifiable, notif notify.Notification) (*resend.SendEmailRequest, error) {
	if resendNotif, ok := notif.(ResendNotification); ok {
		return resendNotif.ToResendMessage(notifiable)
	}

	return &resend.SendEmailRequest{
		Subject: notif.Subject(),
		Text:    notif.Content(),
	}, nil
}
