package notify

import (
	"context"
	"log/slog"
)

// BeforeHook is a function that is called before a notification is sent to all channels.
// If it returns an error, the notification process will be stopped and the error will be returned
type BeforeHook func(ctx context.Context, channel Channel, notifiable Notifiable, notif Notification) error

// AfterHook is a function that is called after a notification is sent to all channels.
// It does not return an error
type AfterHook func(ctx context.Context, channel Channel, notifiable Notifiable, notif Notification, err error)

// SlogBeforeHook returns a BeforeHook that logs the notification to a slog logger, before it is sent to a channel
// The logger will log the channel name, notification kind, and the notifiable
func SlogBeforeHook(logger *slog.Logger) BeforeHook {
	return func(ctx context.Context, channel Channel, notifiable Notifiable, notif Notification) error {
		logger.DebugContext(
			ctx,
			"sending notification",
			slog.String("channel", channel.Name()),
			slog.String("kind", notif.Kind()),
			slog.String("notifiable", notifiable.NotifiableID()),
		)

		return nil
	}
}

// SlogAfterHook returns an AfterHook that logs the notification to a slog logger, after it is sent to a channel
// The logger will log the channel name, notification kind, and the notifiable
func SlogAfterHook(logger *slog.Logger, errKey string) AfterHook {
	return func(ctx context.Context, channel Channel, notifiable Notifiable, notif Notification, err error) {
		if err == nil {
			logger.InfoContext(
				ctx,
				"notification sent",
				slog.String("channel", channel.Name()),
				slog.String("kind", notif.Kind()),
				slog.String("notifiable", notifiable.NotifiableID()),
			)

			return
		}

		logger.ErrorContext(
			ctx,
			"error while sending notification",
			slog.String("channel", channel.Name()),
			slog.String("kind", notif.Kind()),
			slog.Any("notifiable", notifiable),
			slog.String(errKey, err.Error()),
		)
	}
}
