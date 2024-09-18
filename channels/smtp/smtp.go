package smtp

import (
	"context"

	"github.com/smc13/notify"
	"github.com/wneessen/go-mail"
)

type SmtpNotification interface {
	ToSMTPMessage(notify.Notifiable) (*mail.Msg, error)
}

type SmtpNotifiable interface {
	SMTPAddress() string
}

type SmtpChannel struct {
	client *mail.Client
	from   string
}

type SmtpOption func(*SmtpChannel)

func New(client *mail.Client, opts ...SmtpOption) *SmtpChannel {
	return &SmtpChannel{client: client}
}

// WithFrom sets the from address for the email
func WithFrom(from string) SmtpOption {
	return func(c *SmtpChannel) {
		c.from = from
	}
}

func (c *SmtpChannel) Name() string { return "smtp" }

func (c *SmtpChannel) Notify(ctx context.Context, notifiable notify.Notifiable, notif notify.Notification) error {
	message, err := c.notificationToMailMessage(notifiable, notif)
	if err != nil {
		return err
	}

	if message == nil {
		return nil
	}

	if c.from != "" && len(message.GetFromString()) == 0 {
		if err := message.From(c.from); err != nil {
			return err
		}
	}

	if notifiable, ok := notifiable.(SmtpNotifiable); ok {
		if err := message.AddTo(notifiable.SMTPAddress()); err != nil {
			return err
		}
	}

	return c.client.DialAndSendWithContext(ctx, message)
}

func (c *SmtpChannel) notificationToMailMessage(notifiable notify.Notifiable, notif notify.Notification) (*mail.Msg, error) {
	if smtpNotif, ok := notif.(SmtpNotification); ok {
		return smtpNotif.ToSMTPMessage(notifiable)
	}

	msg := mail.NewMsg()
	msg.Subject(notif.Subject())
	msg.SetBodyString(mail.TypeTextPlain, notif.Content())
	msg.SetDate()

	return msg, nil
}
